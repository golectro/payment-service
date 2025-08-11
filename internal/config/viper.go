package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

func NewViper() *viper.Viper {
	config := viper.New()

	config.SetConfigFile(".env")

	if err := config.ReadInConfig(); err != nil {
		if os.IsNotExist(err) {
			log.Println("No .env file found in root directory")
		} else {
			log.Printf("Error reading .env file: %v", err)
		}
	} else {
		log.Println("Successfully loaded configuration from .env")
	}

	config.AutomaticEnv()
	return config
}
