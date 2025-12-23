/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// noteCmd represents the note command
var noteCmd = &cobra.Command{
	Use:   "note <comment...>",
	Short: "Add a note to the current engagement log",
	Long:  "Adds a note entry to your engagement state (scaffolding; functionality not yet implemented).",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("note called")
	},
	Args: cobra.MinimumNArgs(1),
}

func init() {
	rootCmd.AddCommand(noteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// noteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// noteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
