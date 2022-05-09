package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Entry struct {
	Name       string
	Surname    string
	Tel        string
	LastAccess string
}

// CSV file resides in the current directory
var CSV_FILE = "csv.data"

type PhoneBook []Entry

var data = PhoneBook{}
var index map[string]int

func readCSVFile(filepath string) error {
	_, err := os.Stat(filepath)
	if err != nil {
		return err
	}

	f, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer f.Close()

	// CSV file read all at once
	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return err
	}

	for _, line := range lines {
		temp := Entry{
			Name:       line[0],
			Surname:    line[1],
			Tel:        line[2],
			LastAccess: line[3],
		}
		// Storing to global variable
		data = append(data, temp)
	}

	return nil
}

func saveCSVFile(filepath string) error {
	csvFile, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer csvFile.Close()

	csvWriter := csv.NewWriter(csvFile)
	for _, row := range data {
		temp := []string{row.Name, row.Surname, row.Tel, row.LastAccess}
		_ = csvWriter.Write(temp)
	}
	csvWriter.Flush()
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

func insert(pS *Entry) error {
	// If it already exists, do not add it
	_, ok := index[(*pS).Tel]
	if ok {
		return fmt.Errorf("%s already exists", pS.Tel)
	}
	data = append(data, *pS)

	// Update the index
	_ = createIndex()

	err := saveCSVFile(CSV_FILE)
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

	err := saveCSVFile(CSV_FILE)
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

func list() {
	sort.Sort(data)
	for _, v := range data {
		fmt.Println(v)
	}
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

// Setup CSV file
func setCSVFile() error {
	filepath := os.Getenv("PHONEBOOK")
	if filepath != "" {
		CSV_FILE = filepath
	}

	_, err := os.Stat(CSV_FILE)
	if err != nil {
		fmt.Println("Creating", CSV_FILE)
		f, err := os.Create(CSV_FILE)
		if err != nil {
			f.Close()
			return err
		}
		f.Close()
	}

	fileInfo, err := os.Stat(CSV_FILE)
	mode := fileInfo.Mode()
	if !mode.IsRegular() {
		return fmt.Errorf("%s not a regular file", CSV_FILE)
	}
	return nil
}

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		exe := path.Base(arguments[0])
		fmt.Printf("Usage: %s insert|delete|search|list <arguments>\n", exe)
		return
	}

	if err := setCSVFile(); err != nil {
		fmt.Println(err)
		return
	}

	if err := readCSVFile(CSV_FILE); err != nil {
		fmt.Println(err)
		return
	}

	if err := createIndex(); err != nil {
		fmt.Println("Cannot create index.")
		return
	}

	// Differentiate between the commands
	switch arguments[1] {
	case "insert":
		if len(arguments) != 5 {
			exe := path.Base(arguments[0])
			fmt.Printf("Usage: %s insert Name Surname Telephone\n", exe)
			return
		}
		t := strings.ReplaceAll(arguments[4], "-", "")
		if !matchTel(t) {
			fmt.Println("Not a valid telephone number:", t)
			return
		}
		temp := initS(arguments[2], arguments[3], t)
		// If it is nil, there was an error
		if temp != nil {
			err := insert(temp)
			if err != nil {
				fmt.Println(err)
				return
			}
		}

	case "delete":
		if len(arguments) != 3 {
			exe := path.Base(arguments[0])
			fmt.Printf("Usage: %s delete Telephone\n", exe)
			return
		}
		t := strings.ReplaceAll(arguments[2], "-", "")
		if !matchTel(t) {
			fmt.Println("Not a valid telephone number:", t)
			return
		}
		err := deleteEntry(t)
		if err != nil {
			fmt.Println(err)
		}

	case "search":
		if len(arguments) != 3 {
			exe := path.Base(arguments[0])
			fmt.Printf("Usage: %s search Telephone\n", exe)
			return
		}
		t := strings.ReplaceAll(arguments[2], "-", "")
		if !matchTel(t) {
			fmt.Println("Not a valid telephone number:", t)
			return
		}
		result := search(t)
		if result == nil {
			fmt.Println("Entry not found:", t)
			return
		}
		fmt.Println(*result)

	case "list":
		list()

	default:
		fmt.Println("Not a valid option")
	}
}
