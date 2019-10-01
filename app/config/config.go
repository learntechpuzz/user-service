package config

import (
	"errors"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var requiredConfig = map[string]string{
	"nats.server": "Missing NATS server",
}

// LoadConfiguration based on env
func LoadConfiguration(env string) error {
	// Configure logger
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	log.Info("Reading in config files...")
	log.Info("Env is " + env)

	// Get Configs
	viper.SetConfigName("config." + env)
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")

	// Read in the config file
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("fatal error config file: %s", err)
	}

	// Make sure all key config elements are set
	if shouldExit := false; true {
		for k, v := range requiredConfig {
			if !viper.IsSet(k) {
				log.Error(v)
				shouldExit = true
			}
		}

		if shouldExit {
			return errors.New("fatal configuration errors")
		}
	}

	// Log all config
	log.Infof("Configuration settings: %v", viper.AllSettings())

	return nil
}
