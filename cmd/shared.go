package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"harvest-cli/internal/api"
	"harvest-cli/internal/config"
)

// Global flags shared across commands
var (
	apiURL       string
	apiKey       string
	timeout      int
	verbose      bool
	outputFormat string
	limit        int
	offset       int
	filter       string
	force        bool
)

func addGlobalFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&apiURL, "api-url", "", "API base URL")
	cmd.PersistentFlags().StringVar(&apiKey, "api-key", "", "API key")
	cmd.PersistentFlags().IntVar(&timeout, "timeout", 30, "Request timeout")
	cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
}

func addViewFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&outputFormat, "format", "f", "table", "Output format (table, json)")
	cmd.Flags().IntVarP(&limit, "limit", "l", 10, "Limit results")
	cmd.Flags().IntVarP(&offset, "offset", "o", 0, "Offset for pagination")
	cmd.Flags().StringVar(&filter, "filter", "", "Filter results")
}

func addDeleteFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&force, "force", false, "Skip confirmation")
}

func createAPIClient() (*api.Client, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Override config with flags if provided
	url := cfg.APIURL
	if apiURL != "" {
		url = apiURL
	}

	key := cfg.APIKey
	if apiKey != "" {
		key = apiKey
	}

	timeoutVal := cfg.Timeout
	if timeout != 30 {
		timeoutVal = timeout
	}

	return api.NewClient(url, key, timeoutVal)
}

func buildListParams() api.ListParams {
	return api.ListParams{
		Limit:  limit,
		Offset: offset,
		Filter: filter,
	}
}

func confirmDelete(resourceType, id string) bool {
	if force {
		return true
	}

	fmt.Printf("Are you sure you want to delete %s %s? [y/N]: ", resourceType, id)
	var response string
	fmt.Scanln(&response)
	return strings.ToLower(response) == "y" || strings.ToLower(response) == "yes"
}
