/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an entry",
	Long:  `Delete an entry from phonebook application`,
	Run: func(cmd *cobra.Command, args []string) {
		// get key
		key, _ := cmd.Flags().GetString("key")
		if key == "" {
			fmt.Println("Not a valid key!")
			return
		}

		// remove data
		err := deleteEntry(key)
		if err != nil {
			fmt.Println(err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	deleteCmd.Flags().StringP("key", "k", "", "Key to delete")
}

func deleteEntry(key string) error {
	i, ok := index[key]
	if !ok {
		return fmt.Errorf("%s cannot be found.\n", key)
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
