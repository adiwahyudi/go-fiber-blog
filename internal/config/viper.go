package config

import (
	"fmt"
	"github.com/spf13/viper"
)

// NewViper is a function to load config from .env
// You can change the implementation, for example load from env file, consul, etcd, etc
func NewViper() *viper.Viper {
	config := viper.New()
	config.AddConfigPath(".")
	config.SetConfigType("env")
	config.SetConfigName(".env")
	config.AutomaticEnv()
	err := config.ReadInConfig()

	if err != nil {
		fmt.Printf("Fatal error .env file: %w \n", err)
	}

	return config
}
