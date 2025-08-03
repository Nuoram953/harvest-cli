package cmd

import (
	"fmt"

	"harvest-cli/internal/api"
	"harvest-cli/internal/flagutils"
	"harvest-cli/internal/ui"

	"github.com/spf13/cobra"
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
	entryTime      float64
)

func init() {
	entryCmd.AddCommand(entryCreateCmd)

	addGlobalFlags(entryCmd)

	entryCreateCmd.Flags().Int64VarP(&entryProjectId, "project", "p", 0, "Project ID")
	entryCreateCmd.Flags().Int64VarP(&entryTaskId, "task", "t", 0, "Task ID")
	entryCreateCmd.Flags().StringVarP(&entryDate, "date", "d", "", "Date for the entry (YYYY-MM-DD)")
	entryCreateCmd.Flags().Float64VarP(&entryTime, "time", "", 0, "Duration in hour or minutes (e.g., 1h or 60m)")
}

func runEntryCreate(cmd *cobra.Command, args []string) error {
	client, err := createAPIClient()
	if err != nil {
		return err
	}

	if !cmd.Flags().Changed("project") {
		entryProjectId = flagutils.DoIfProjectFlagMissing(cmd, client)
	}

	if !cmd.Flags().Changed("task") {
		entryTaskId = flagutils.DoIfTaskFlagMissing(cmd, client, entryProjectId)
	}

	if !cmd.Flags().Changed("date") {
		entryDate = flagutils.DoIfDateFlagMissing(cmd, client)
	}

	if !cmd.Flags().Changed("time") {
		entryTime = flagutils.DoIfTimeFlagMissing(cmd, client)
	}

	if !cmd.Flags().Changed("noconfirm") {
		if !flagutils.DoIfConfirmFlagMissing(cmd, client) {
			ui.WriteTextError("Entry creation cancelled.")
			return nil
		}
	}

	entry := api.CreateEntryRequest{
		ProjectId: entryProjectId,
		TaskId:    entryTaskId,
		Date:      entryDate,
		Hours:     entryTime,
	}

	_, err = client.CreateEntry(entry)
	if err != nil {
		return fmt.Errorf("Failed to create entry: %w", err)
	}

	ui.WriteTextSuccess("Entry created successfully!")

	return nil
}
