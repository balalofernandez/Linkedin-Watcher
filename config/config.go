package config

import (
	"linkedin-watcher/infra/logger"

	"github.com/spf13/viper"
)

type Configuration struct {
	Server   ServerConfiguration
	Database DatabaseConfiguration
}

// SetupConfig configuration
func SetupConfig() error {
	var configuration *Configuration

	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		logger.Errorf("Error to reading config file, %s", err)
		return err
	}

	// Enable environment variable reading
	viper.AutomaticEnv()

	err := viper.Unmarshal(&configuration)
	if err != nil {
		logger.Errorf("error to decode, %v", err)
		return err
	}

	return nil
}
