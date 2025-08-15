package config

import (
	"github.com/hashicorp/vault/api"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewVaultClient(viper *viper.Viper, log *logrus.Logger) *api.Client {
	config := api.DefaultConfig()

	vaultAddr := viper.GetString("VAULT_ADDR")
	if vaultAddr == "" {
		log.Fatal("VAULT_ADDR is not set in configuration")
	}
	config.Address = vaultAddr

	client, err := api.NewClient(config)
	if err != nil {
		log.Fatalf("failed to create Vault client: %v", err)
	}

	token := viper.GetString("VAULT_TOKEN")
	if token == "" {
		log.Fatal("VAULT_TOKEN is not set in configuration")
	}
	client.SetToken(token)

	log.Info("Vault client initialized successfully")
	return client
}
