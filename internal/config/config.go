package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	APIURL  string `mapstructure:"api_url"`
	APIKey  string `mapstructure:"api_key"`
	Timeout int    `mapstructure:"timeout"`
}

func Load() (*Config, error) {
	viper.SetConfigName("harvest-cli")
	viper.SetConfigType("yaml")

	home, _ := os.UserHomeDir()
	viper.AddConfigPath(filepath.Join(home, ".config", "harvest-cli"))
	viper.AddConfigPath(".")

	// Environment variables
	viper.SetEnvPrefix("HARVESTCLI")
	viper.AutomaticEnv()

	// Defaults
	viper.SetDefault("timeout", 30)

	var config Config
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
