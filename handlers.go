package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

const PORT = ":1234"

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving:", r.URL.Path, "from", r.Host)
	w.WriteHeader(http.StatusOK)
	body := "Thanks for visiting!\n"
	fmt.Fprintf(w, "%s", body)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	// Get telephone
	paramStr := strings.Split(r.URL.Path, "/")
	fmt.Println("Path:", paramStr)
	if len(paramStr) < 3 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "Not found: "+r.URL.Path)
		return
	}

	log.Println("Serving:", r.URL.Path, "from", r.Host)
	telephone := paramStr[2]

	err := deleteEntry(telephone)
	if err != nil {
		fmt.Println(err)
		body := err.Error() + "\n"
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "%s", body)
		return
	}

	body := telephone + " deleted!\n"
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", body)
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving:", r.URL.Path, "from", r.Host)
	w.WriteHeader(http.StatusOK)
	body := list()
	fmt.Fprintf(w, "%s\n", body)
}

func insertHandler(w http.ResponseWriter, r *http.Request) {
	// Split URL
	paramStr := strings.Split(r.URL.Path, "/")
	fmt.Println("Path:", paramStr)
	if len(paramStr) < 5 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "Not enough arguments: "+r.URL.Path)
		return
	}

	name := paramStr[2]
	surname := paramStr[3]
	tel := paramStr[4]

	t := strings.ReplaceAll(tel, "-", "")
	if !matchTel(t) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Not a valid telephone number: "+tel)
		return
	}

	entry := initS(name, surname, tel)
	err := insert(entry)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Failed to add record!")
		log.Println("Failed to add record:", err)
	} else {
		log.Println("Serving:", r.URL.Path, "from", r.Host)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "New record added successfully")
	}

	log.Println("Serving:", r.URL.Path, "from", r.Host)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	// Get search value from URL
	paramStr := strings.Split(r.URL.Path, "/")
	fmt.Println("Path:", paramStr)
	if len(paramStr) < 3 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "Not found:"+r.URL.Path)
		return
	}

	var body string
	telephone := paramStr[2]
	t := search(telephone)
	if t == nil {
		w.WriteHeader(http.StatusNotFound)
		body = "Could not be found:" + telephone + "\n"
	} else {
		w.WriteHeader(http.StatusOK)
		body = t.Name + " " + t.Surname + " " + t.Tel + "\n"
	}

	fmt.Println("Serving:", r.URL.Path, "from", r.Host)
	fmt.Fprintf(w, "%s", body)
}
