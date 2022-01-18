package storage

import (
	"strings"

	"github.com/GrayFinance/mint/src/models"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB
var err error

func Connect(uri string) error {
	options := &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	}
	if strings.HasPrefix(uri, "postgres") {
		DB, err = gorm.Open(postgres.Open(uri), options)
	} else {
		DB, err = gorm.Open(sqlite.Open(uri), options)
	}

	if err != nil {
		return err
	}

	if err := DB.AutoMigrate(
		&models.User{},
		&models.Wallet{},
		&models.Payment{},
	); err != nil {
		return err
	}
	return nil
}
