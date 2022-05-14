package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"
)

type Entry struct {
	Name       string
	Surname    string
	Tel        string
	LastAccess string
}

// JSONFILE resides in the current directory
var JSONFILE = "./data.json"

type PhoneBook []Entry

var data = PhoneBook{}
var index map[string]int

func readJSONFile(filepath string) error {
	_, err := os.Stat(filepath)
	if err != nil {
		return err
	}

	f, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer f.Close()

	err = json.NewDecoder(f).Decode(&data)
	if err != nil {
		return err
	}

	return nil
}

func saveJSONFile(filepath string) error {
	f, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer f.Close()

	err = json.NewEncoder(f).Encode(&data)
	if err != nil {
		return err
	}

	return nil
}

func createIndex() error {
	index = make(map[string]int)
	for i, k := range data {
		key := k.Tel
		index[key] = i
	}
	return nil
}

// Initialized by the user - returns a pointer
// If it returns nil, there was an error
func initS(N, S, T string) *Entry {
	// Both of them should have a value
	if T == "" || S == "" {
		return nil
	}
	// Give LastAccess a value
	LastAccess := strconv.FormatInt(time.Now().Unix(), 10)
	return &Entry{Name: N, Surname: S, Tel: T, LastAccess: LastAccess}
}

func insert(entry *Entry) error {
	// If it already exists, do not add it
	_, ok := index[(*entry).Tel]
	if ok {
		return fmt.Errorf("%s already exists", entry.Tel)
	}
	data = append(data, *entry)

	// Update the index
	_ = createIndex()

	err := saveJSONFile(JSONFILE)
	if err != nil {
		return err
	}
	return nil
}

func deleteEntry(key string) error {
	i, ok := index[key]
	if !ok {
		return fmt.Errorf("%s cannot be found!", key)
	}
	data = append(data[:i], data[i+1:]...)

	// Update the index - key does not exist anymore
	delete(index, key)

	err := saveJSONFile(JSONFILE)
	if err != nil {
		return err
	}
	return nil
}

func search(key string) *Entry {
	i, ok := index[key]
	if !ok {
		return nil
	}
	data[i].LastAccess = strconv.FormatInt(time.Now().Unix(), 10)
	return &data[i]
}

func list() string {
	// sort.Sort(data)
	var all string
	for _, v := range data {
		all = all + v.Name + " " + v.Surname + " " + v.Tel + "\n"
	}
	all = all + fmt.Sprintf("%d records in total.", len(data))
	return all
}

func matchTel(s string) bool {
	t := []byte(s)
	re := regexp.MustCompile(`\d+$`)
	return re.Match(t)
}

// Implement sort.Interface
func (a PhoneBook) Len() int {
	return len(a)
}

// First based on surname. If they have the same surname
// take into account the name.
func (a PhoneBook) Less(i, j int) bool {
	if a[i].Surname == a[j].Surname {
		return a[i].Name < a[j].Name
	}
	return a[i].Surname < a[j].Surname
}

func (a PhoneBook) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// Setup JSON file
func setJSONFile() error {
	filepath := os.Getenv("PHONEBOOK")
	if filepath != "" {
		JSONFILE = filepath
	}

	_, err := os.Stat(JSONFILE)
	if err != nil {
		fmt.Println("Creating", JSONFILE)
		f, err := os.Create(JSONFILE)
		if err != nil {
			f.Close()
			return err
		}
		f.Close()
	}

	fileInfo, err := os.Stat(JSONFILE)
	mode := fileInfo.Mode()
	if !mode.IsRegular() {
		return fmt.Errorf("%s not a regular file", JSONFILE)
	}

	return nil
}

func main() {
	err := setJSONFile()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = readJSONFile(JSONFILE)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = createIndex()
	if err != nil {
		fmt.Println("Cannot create index.")
		return
	}

	mux := http.NewServeMux()
	s := &http.Server{
		Addr:         PORT,
		Handler:      mux,
		IdleTimeout:  10 * time.Second,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
	}

	mux.Handle("/list", http.HandlerFunc(listHandler))

	mux.Handle("/insert/", http.HandlerFunc(insertHandler))
	mux.Handle("/insert", http.HandlerFunc(insertHandler))
	mux.Handle("/search/", http.HandlerFunc(searchHandler))
	mux.Handle("/search", http.HandlerFunc(searchHandler))
	mux.Handle("/delete/", http.HandlerFunc(deleteHandler))
	mux.Handle("/", http.HandlerFunc(defaultHandler))

	fmt.Println("Ready to serve at", PORT)
	err = s.ListenAndServe()
	if err != nil {
		fmt.Println(err)
		return
	}
}
