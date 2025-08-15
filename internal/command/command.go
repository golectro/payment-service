package command

import (
	"fmt"
	"golectro-payment/internal/migrations"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type CommandExecutor struct {
	DB    *gorm.DB
	Viper *viper.Viper
}

func NewCommandExecutor(viper *viper.Viper, db *gorm.DB) *CommandExecutor {
	return &CommandExecutor{
		DB:    db,
		Viper: viper,
	}
}

func (ce *CommandExecutor) Execute(logger *logrus.Logger) bool {
	args := os.Args[1:]
	if len(args) == 0 {
		return true
	}

	run := false
	for _, arg := range args {
		switch arg {
		case "--migrate":
			ce.handleMigrate(logger)
		case "--seed":
			ce.handleSeed(logger)
		case "--create-db":
			ce.handleCreateDB(logger)
		case "--drop-db":
			ce.handleDropDB(logger)
		case "--drop-table":
			ce.handleDropTable(logger)
		case "--run":
			run = true
		}
	}

	return run
}

func (ce *CommandExecutor) handleMigrate(logger *logrus.Logger) {
	if err := migrations.Migrate(ce.DB); err != nil {
		logger.Fatalf("❌ Migration failed: %v", err)
	}
	logger.Println("✅ Migration completed")
}

func (ce *CommandExecutor) handleSeed(logger *logrus.Logger) {
	if err := migrations.Seeder(ce.DB, logger); err != nil {
		logger.Fatalf("❌ Seeder failed: %v", err)
	}
	logger.Println("✅ Seeder completed")
}

func (ce *CommandExecutor) handleCreateDB(logger *logrus.Logger) {
	dbName := ce.Viper.GetString("DB_NAME")
	if dbName == "" {
		logger.Fatal("❌ DB_NAME is not set in env")
	}

	sql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", dbName)
	if err := ce.DB.Exec(sql).Error; err != nil {
		logger.Fatalf("❌ Failed to create database: %v", err)
	}
	logger.Printf("✅ Database '%s' created\n", dbName)
}

func (ce *CommandExecutor) handleDropDB(logger *logrus.Logger) {
	dbName := ce.Viper.GetString("DB_NAME")
	if dbName == "" {
		logger.Fatal("❌ DB_NAME is not set in env")
	}

	sql := fmt.Sprintf("DROP DATABASE IF EXISTS `%s`", dbName)
	if err := ce.DB.Exec(sql).Error; err != nil {
		logger.Fatalf("❌ Failed to drop database: %v", err)
	}
	logger.Printf("✅ Database '%s' dropped\n", dbName)
}

func (ce *CommandExecutor) handleDropTable(logger *logrus.Logger) {
	tables := ce.Viper.GetString("DROP_TABLE_NAMES")
	if tables == "" {
		logger.Fatal("❌ DROP_TABLE_NAMES is not set in env")
	}

	tableList := strings.SplitSeq(tables, ",")
	for table := range tableList {
		table = strings.TrimSpace(table)
		if table == "" {
			continue
		}

		sql := fmt.Sprintf("DROP TABLE IF EXISTS `%s`", table)
		if err := ce.DB.Exec(sql).Error; err != nil {
			logger.Fatalf("❌ Failed to drop table '%s': %v", table, err)
		}
		logger.Printf("✅ Table '%s' dropped\n", table)
	}
}
