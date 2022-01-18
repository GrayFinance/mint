package config

import (
	"github.com/kelseyhightower/envconfig"
)

// Loads the variables of environments specified
// in the .env file of the current directory.
var Config struct {
	DB_URI         string `envconfig:"DB_URI"`
	PASS_SALT      string `envconfig:"PASS_SALT"`
	SIGN_KEY       string `envconfig:"SIGN_KEY"`
	ADMIN_USERNAME string `envconfig:"ADMIN_USERNAME"`
	ADMIN_PASSWORD string `envconfig:"ADMIN_PASSWORD"`
}

func Loads() error {
	if err := envconfig.Process("", &Config); err != nil {
		return err
	}
	return nil
}
