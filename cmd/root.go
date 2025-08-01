package cmd

import (
    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "my-app",
    Short: "API wrapper for managing entries and tasks",
    Long:  `A command-line interface for interacting with the API`,
}

func Execute() error {
    return rootCmd.Execute()
}

func init() {
    rootCmd.AddCommand(entryCmd)
}
