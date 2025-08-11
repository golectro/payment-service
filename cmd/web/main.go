package main

import (
	"fmt"
	"golectro-payment/internal/command"
	"golectro-payment/internal/config"
)

func main() {
	viper := config.NewViper()
	log := config.NewLogger(viper)
	db := config.NewDatabase(viper, log)
	mongo := config.NewMongoDB(viper, log)
	validate := config.NewValidator(viper)
	redis := config.NewRedis(viper, log)
	vault := config.NewVaultClient(viper, log)
	app := config.NewGin(viper, log, mongo, redis)
	executor := command.NewCommandExecutor(viper, db)

	config.Bootstrap(&config.BootstrapConfig{
		Viper:    viper,
		Log:      log,
		DB:       db,
		Mongo:    mongo,
		Validate: validate,
		App:      app,
		Redis:    redis,
		Vault:    vault,
	})

	if !executor.Execute(log) {
		return
	}

	webPort := viper.GetInt("PORT")
	err := app.Run(fmt.Sprintf(":%d", webPort))
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
