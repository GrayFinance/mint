package storage

import (
	"strings"

	"github.com/GrayFinance/mint/src/config"
	"github.com/GrayFinance/mint/src/models"
	"github.com/go-redis/redis"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	REDIS *redis.Client
	DB    *gorm.DB
	err   error
)

func Connect() error {
	// Connect database with ORM.
	if strings.HasPrefix(config.Config.DATABASE, "postgres") {
		DB, err = gorm.Open(postgres.Open(config.Config.DATABASE))
	} else {
		DB, err = gorm.Open(sqlite.Open(config.Config.DATABASE))
	}

	if err != nil {
		return err
	}

	// Migration database.
	err = DB.AutoMigrate(
		&models.User{},
		&models.Wallet{},
		&models.Payment{},
		&models.Address{},
	)
	if err != nil {
		return err
	}

	// Redis Connection
	REDIS = redis.NewClient(&redis.Options{
		Addr:     config.Config.REDIS_HOST,
		Password: config.Config.REDIS_PASSWORD,
		DB:       0,
	})

	if _, err := REDIS.Ping().Result(); err != nil {
		return err
	}
	return nil
}
