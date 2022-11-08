package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	HUE HueConfig
}

type HueConfig struct {
	Host string
	User string
}

// Load LoadConfig reads configuration from file or environment variables.
func Load(path string) (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	config := &Config{}
	if err := viper.Unmarshal(config); err != nil {
		return nil, err
	}

	return config, nil
}
