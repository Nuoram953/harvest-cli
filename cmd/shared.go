package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"harvest-cli/internal/api"
	"harvest-cli/internal/config"
)

var (
	token        string
	accountId    string
	timeout      int
	verbose      bool
	outputFormat string
	limit        int
	offset       int
	filter       string
	force        bool
)

func addGlobalFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().IntVar(&timeout, "timeout", 30, "Request timeout")
	cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
	cmd.PersistentFlags().BoolVar(&force, "noconfirm", false, "Skip confirmation")
}

func createAPIClient() (*api.Client, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return api.NewClient(cfg.Token, cfg.AccountId)
}
