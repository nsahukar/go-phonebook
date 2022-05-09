/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for the number",
	Long: `Search whether a telephone number exists in the 
	phonebook application or not.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get key
		searchKey, _ := cmd.Flags().GetString("key")

		t := strings.ReplaceAll(searchKey, "-", "")
		if !validKey(t) {
			return
		}

		// Search for it
		entry := search(t)
		if entry == nil {
			fmt.Println("Number not found:", t)
			return
		}
		fmt.Println(*entry)
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// searchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// searchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	searchCmd.Flags().StringP("key", "k", "", "Key to search")
}

func validKey(k string) bool {
	if k == "" {
		fmt.Println("Not a valid key!")
		return false
	}

	if !matchTel(k) {
		fmt.Println("Not a valid telephone number:", k)
		return false
	}

	return true
}

func matchTel(s string) bool {
	t := []byte(s)
	re := regexp.MustCompile(`\d+$`)
	return re.Match(t)
}

func search(key string) *Entry {
	i, ok := index[key]
	if !ok {
		return nil
	}
	return &data[i]
}
