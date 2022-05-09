/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

type Entry struct {
	Name       string
	Surname    string
	Tel        string
	LastAccess string
}

// JSON file resides in the current directory
var JSONFILE = "data.json"

type PhoneBook []Entry

var data = PhoneBook{}
var index map[string]int

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

func readJOSNFile(filepath string) error {
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

func createIndex() {
	index = make(map[string]int)
	for i, k := range data {
		key := k.Tel
		index[key] = i
	}
}

// Initialized by the user - returns a pointer
// If it returns nil, there was an error
func initEntry(N, S, T string) *Entry {
	lastAccess := strconv.FormatInt(time.Now().Unix(), 10)
	return &Entry{Name: N, Surname: S, Tel: T, LastAccess: lastAccess}
}

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

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "phonebook",
	Short: "A phonebook application",
	Long:  `This is a phonebook application that uses JSON records.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := setJSONFile()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = readJOSNFile(JSONFILE)
	// io.EOF is fine because it means the file is empty
	if err != nil && err != io.EOF {
		return
	}
	createIndex()

	err = rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.go-phonebook.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
