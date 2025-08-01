package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"harvest-cli/internal/api"
)

var entryCmd = &cobra.Command{
	Use:   "entry",
	Short: "Manage entries via API",
	Long:  `Create, view, edit, and delete entries through the API`,
}

var entryCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new time entry",
	RunE:  runEntryCreate,
}

var entryViewCmd = &cobra.Command{
	Use:   "view [id]",
	Short: "View a time entry",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runEntryView,
}

var entryEditCmd = &cobra.Command{
	Use:   "edit [id]",
	Short: "Edit a time entry",
	Args:  cobra.ExactArgs(1),
	RunE:  runEntryEdit,
}

var entryDeleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Delete a time entry",
	Args:  cobra.ExactArgs(1),
	RunE:  runEntryDelete,
}

var (
	entryTitle   string
	entryContent string
	entryTags    []string
	entryStatus  string
	entryPrivate bool
)

func init() {
	entryCmd.AddCommand(entryCreateCmd)
	entryCmd.AddCommand(entryViewCmd)
	entryCmd.AddCommand(entryEditCmd)
	entryCmd.AddCommand(entryDeleteCmd)

	addGlobalFlags(entryCmd)

	entryCreateCmd.Flags().StringVarP(&entryTitle, "title", "t", "", "Entry title (required)")
	entryCreateCmd.Flags().StringVarP(&entryContent, "content", "c", "", "Entry content")
	entryCreateCmd.Flags().StringSliceVar(&entryTags, "tags", []string{}, "Tags (comma-separated)")
	entryCreateCmd.Flags().StringVarP(&entryStatus, "status", "s", "draft", "Entry status")
	entryCreateCmd.Flags().BoolVar(&entryPrivate, "private", false, "Make entry private")
	entryCreateCmd.MarkFlagRequired("title")

	entryEditCmd.Flags().StringVarP(&entryTitle, "title", "t", "", "New entry title")
	entryEditCmd.Flags().StringVarP(&entryContent, "content", "c", "", "New entry content")
	entryEditCmd.Flags().StringSliceVar(&entryTags, "tags", []string{}, "New tags")
	entryEditCmd.Flags().StringVarP(&entryStatus, "status", "s", "", "New entry status")
	entryEditCmd.Flags().BoolVar(&entryPrivate, "private", false, "Set entry privacy")

	addViewFlags(entryViewCmd)
	addDeleteFlags(entryDeleteCmd)
}

func runEntryCreate(cmd *cobra.Command, args []string) error {
	client, err := createAPIClient()
	if err != nil {
		return err
	}

	entry := api.CreateEntryRequest{
		Title:   entryTitle,
		Content: entryContent,
		Tags:    entryTags,
		Status:  entryStatus,
		Private: entryPrivate,
	}

	result, err := client.CreateEntry(entry)
	if err != nil {
		return fmt.Errorf("failed to create entry: %w", err)
	}

	fmt.Printf("Entry created successfully!\nID: %s\nTitle: %s\n", result.ID, result.Title)
	return nil
}

func runEntryView(cmd *cobra.Command, args []string) error {
	client, err := createAPIClient()
	if err != nil {
		return err
	}

	if len(args) == 1 {
		entry, err := client.GetEntry(args[0])
		if err != nil {
			return fmt.Errorf("failed to get entry: %w", err)
		}
		return displayEntry(entry)
	}

	entries, err := client.ListEntries(buildListParams())
	if err != nil {
		return fmt.Errorf("failed to list entries: %w", err)
	}

	return displayEntries(entries)
}

func runEntryEdit(cmd *cobra.Command, args []string) error {
	client, err := createAPIClient()
	if err != nil {
		return err
	}

	updates := api.UpdateEntryRequest{}
	if cmd.Flags().Changed("title") {
		updates.Title = &entryTitle
	}
	if cmd.Flags().Changed("content") {
		updates.Content = &entryContent
	}
	if cmd.Flags().Changed("tags") {
		updates.Tags = &entryTags
	}
	if cmd.Flags().Changed("status") {
		updates.Status = &entryStatus
	}
	if cmd.Flags().Changed("private") {
		updates.Private = &entryPrivate
	}

	result, err := client.UpdateEntry(args[0], updates)
	if err != nil {
		return fmt.Errorf("failed to update entry: %w", err)
	}

	fmt.Println("Entry updated successfully!")
	return displayEntry(result)
}

func runEntryDelete(cmd *cobra.Command, args []string) error {
	if !confirmDelete("entry", args[0]) {
		return nil
	}

	client, err := createAPIClient()
	if err != nil {
		return err
	}

	err = client.DeleteEntry(args[0])
	if err != nil {
		return fmt.Errorf("failed to delete entry: %w", err)
	}

	fmt.Printf("Entry %s deleted successfully!\n", args[0])
	return nil
}

func displayEntry(entry *api.Entry) error {
	if outputFormat == "json" {
		data, _ := json.MarshalIndent(entry, "", "  ")
		fmt.Println(string(data))
		return nil
	}

	fmt.Printf("ID:      %s\n", entry.ID)
	fmt.Printf("Title:   %s\n", entry.Title)
	fmt.Printf("Content: %s\n", entry.Content)
	fmt.Printf("Tags:    %s\n", strings.Join(entry.Tags, ", "))
	fmt.Printf("Status:  %s\n", entry.Status)
	fmt.Printf("Private: %t\n", entry.Private)
	return nil
}

func displayEntries(entries []*api.Entry) error {
	if outputFormat == "json" {
		data, _ := json.MarshalIndent(entries, "", "  ")
		fmt.Println(string(data))
		return nil
	}

	fmt.Printf("%-10s %-30s %-15s %-10s\n", "ID", "Title", "Status", "Private")
	fmt.Println(strings.Repeat("-", 70))
	for _, entry := range entries {
		title := entry.Title
		if len(title) > 30 {
			title = title[:27] + "..."
		}
		fmt.Printf("%-10s %-30s %-15s %-10t\n",
			entry.ID, title, entry.Status, entry.Private)
	}
	return nil
}
