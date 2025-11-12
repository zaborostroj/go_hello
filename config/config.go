package config

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	DB struct {
		Prefix   string
		Username string
		Password string
		Host     string
		Port     string
		Dbname   string
	}
}

func LoadConfig(path string) *Config {
	env := os.Getenv("APP_ENV")
	configName := "application"
	if env != "" {
		configName = fmt.Sprintf("application-%s", env)
	}
	viper.SetConfigName(configName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config: %s", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("Error parsing config: %s", err)
	}

	return &cfg
}
