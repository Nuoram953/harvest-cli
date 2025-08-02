package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"harvest-cli/internal/api"
	"harvest-cli/internal/ui"
	"strconv"
	"strings"
)

var entryCmd = &cobra.Command{
	Use:   "entry",
	Short: "Manage time entries via API",
	Long:  `create, view, edit, and delete entries through the API`,
}

var entryCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new time entry",
	RunE:  runEntryCreate,
}

var (
	entryProjectId int64
	entryTaskId    int64
	entryDate      string
	entryMinutes   float64
)

func init() {
	entryCmd.AddCommand(entryCreateCmd)

	addGlobalFlags(entryCmd)

	entryCreateCmd.Flags().Int64VarP(&entryProjectId, "project", "p", 0, "Project ID")
	entryCreateCmd.Flags().Int64VarP(&entryTaskId, "task", "t", 0, "Task ID")
	entryCreateCmd.Flags().StringVarP(&entryDate, "date", "d", "", "Date for the entry (YYYY-MM-DD)")
	entryCreateCmd.Flags().Float64VarP(&entryMinutes, "minute", "m", 0, "Duration in minutes")
}

func runEntryCreate(cmd *cobra.Command, args []string) error {
	client, err := createAPIClient()
	if err != nil {
		return err
	}

	if entryProjectId == 0 {
		selectedProject, _ := ui.SelectProjectInteractively(client)
		entryProjectId = selectedProject.ID
	}

	if entryTaskId == 0 {
		selectedTask, _ := ui.SelectTaskInteractively(client, entryProjectId)
		entryTaskId = selectedTask.ID
	}

	if entryDate == "" {
		date, _ := ui.TextInputDate("When was the entry made?")
		entryDate = date
	}

	if entryMinutes == 0 {
		input, err := ui.SimpleTextInput("What was the duration?", "(ex. 60m / 1h)")
		if err != nil {
			fmt.Println("Error reading input:", err)
		}

		input = strings.TrimSpace(input)
		if strings.HasSuffix(input, "m") {
			minStr := strings.TrimSuffix(input, "m")
			minutes, err := strconv.Atoi(minStr)
			if err != nil {
				fmt.Println("Invalid minutes format.")
			}
			entryMinutes = float64(minutes) / 60.0
		} else if strings.HasSuffix(input, "h") {
			hourStr := strings.TrimSuffix(input, "h")
			entryMinutes, err = strconv.ParseFloat(hourStr, 64)
			if err != nil {
				fmt.Println("Invalid hours format.")
			}
		} else {
			fmt.Println("Invalid duration format. Please use '60m' or '1h'.")
		}
	}
	if !cmd.Flags().Changed("noconfirm") {
		confirm, err := ui.Confirm("Create entry", "Are you sure you want to create this entry?")
		if err != nil {
			return fmt.Errorf("Failed to confirm entry creation: %w", err)
		}

		if !confirm {
			fmt.Println("Entry creation cancelled.")
			return nil
		}
	}

	entry := api.CreateEntryRequest{
		ProjectId: entryProjectId,
		TaskId:    entryTaskId,
		Date:      entryDate,
		Hours:     entryMinutes,
	}

	_, err = client.CreateEntry(entry)
	if err != nil {
		return fmt.Errorf("Failed to create entry: %w", err)
	}

	fmt.Printf("Entry created successfully!")

	return nil
}
