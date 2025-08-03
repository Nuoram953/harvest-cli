package flagutils

import (
	"fmt"
	"harvest-cli/internal/api"
	"harvest-cli/internal/ui"
	"regexp"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

func DoIfProjectFlagMissing(cmd *cobra.Command, client *api.Client) int64 {
	selectedProject, _ := ui.SelectProjectInteractively(client)
	ui.WriteTextStep("Project", selectedProject.Name)
	return selectedProject.ID
}

func DoIfTaskFlagMissing(cmd *cobra.Command, client *api.Client, projectId int64) int64 {
	selectedTask, _ := ui.SelectTaskInteractively(client, projectId)
	ui.WriteTextStep("Task", selectedTask.Name)
	return selectedTask.ID
}

func DoIfDateFlagMissing(cmd *cobra.Command, client *api.Client) string {
	date, _ := ui.TextInputDate("When was the entry made?")
	ui.WriteTextStep("Date", date)
	return date
}

func DoIfTimeFlagMissing(cmd *cobra.Command, client *api.Client) float64 {
	input, err := ui.SimpleTextInput("What was the duration?", "(ex. 60m / 1h / 1h30m)")
	if err != nil {
		fmt.Println("Error reading input:", err)
		return 0
	}

	input = strings.TrimSpace(input)
	re := regexp.MustCompile(`(?i)^\s*(?:(\d+(?:\.\d+)?)\s*h)?\s*(?:(\d+)\s*m)?\s*$`)
	matches := re.FindStringSubmatch(input)

	if matches == nil {
		fmt.Println("Invalid duration format. Please use '60m', '1h', or '1h30m'.")
		return 0
	}

	hoursStr := matches[1]
	minutesStr := matches[2]

	var hours float64
	var minutes int

	if hoursStr != "" {
		h, err := strconv.ParseFloat(hoursStr, 64)
		if err != nil {
			fmt.Println("Invalid hours format.")
			return 0
		}
		hours += h
	}

	if minutesStr != "" {
		m, err := strconv.Atoi(minutesStr)
		if err != nil {
			fmt.Println("Invalid minutes format.")
			return 0
		}
		hours += float64(m) / 60.0
		minutes += m
	}

	ui.WriteTextStep("Duration", fmt.Sprintf("%.0fh%dm", hours, minutes))
	return hours
}
func DoIfConfirmFlagMissing(cmd *cobra.Command, client *api.Client) bool {
	confirm, err := ui.Confirm("Create entry", "Are you sure you want to create this entry?")
	if err != nil || !confirm {
		return false
	}

	return true
}
