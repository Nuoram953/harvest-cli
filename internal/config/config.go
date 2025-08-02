package config

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Token  string `mapstructure:"token"`
	AccountId  string `mapstructure:"account_id"`
}

func Load() (*Config, error) {

	if err := godotenv.Load(); err != nil {
		godotenv.Load(".env.local")
	}

	viper.SetConfigName("harvest-cli")
	viper.SetConfigType("yaml")

	home, _ := os.UserHomeDir()
	viper.AddConfigPath(filepath.Join(home, ".config", "harvest-cli"))
	viper.AddConfigPath(".")

	viper.BindEnv("token", "HARVEST_TOKEN")
	viper.BindEnv("account_id", "HARVEST_ACCOUNT_ID")

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
