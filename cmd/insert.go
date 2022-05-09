/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// insertCmd represents the insert command
var insertCmd = &cobra.Command{
	Use:   "insert",
	Short: "Insert new data",
	Long:  `This command inserts new data into the phonebook application.`,
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		surname, _ := cmd.Flags().GetString("surname")
		tel, _ := cmd.Flags().GetString("tel")

		t := strings.ReplaceAll(tel, "-", "")
		if !validData(name, surname, t) {
			return
		}

		entry := initEntry(name, surname, t)
		if entry == nil {
			fmt.Println("Not a valid record:", entry)
			return
		}

		// insert data
		err := insert(entry)
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(insertCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// insertCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// insertCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	insertCmd.Flags().StringP("name", "n", "", "firstname")
	insertCmd.Flags().StringP("surname", "s", "", "surname")
	insertCmd.Flags().StringP("tel", "t", "", "telephone")
}

func validData(name, surname, tel string) bool {
	if name == "" {
		fmt.Println("Not a valid name!")
		return false
	}

	if surname == "" {
		fmt.Println("Not a valid surname!")
		return false
	}

	if tel == "" {
		fmt.Println("Not a valid telephone number!")
		return false
	}

	if !matchTel(tel) {
		fmt.Println("Not a valid telephone number:", tel)
		return false
	}

	return true
}

func insert(entry *Entry) error {
	// if it already exists, do not add it
	_, ok := index[(*entry).Tel]
	if ok {
		return fmt.Errorf("%s already exists", entry.Tel)
	}
	data = append(data, *entry)

	// save the data
	err := saveJSONFile(JSONFILE)
	if err != nil {
		return err
	}
	return nil
}
