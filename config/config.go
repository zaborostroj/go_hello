package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

func LoadConfig[T any]() T {
	env := os.Getenv("APP_ENV")
	configName := "application"
	if env != "" {
		configName = fmt.Sprintf("application-%s", env)
	}

	v := viper.New()
	v.SetConfigName(configName)
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config: %s", err)
	}

	var config T
	if err := v.Unmarshal(&config); err != nil {
		log.Fatalf("Error parsing config: %s", err)
	}
	log.Printf("Application config: %+v", config)

	return config
}
