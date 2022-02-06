package config

import "github.com/kelseyhightower/envconfig"

// Loads the variables of environments specified
// in the .env file of the current directory.
var Config struct {
	API_HOST        string `envconfig:"API_HOST"`
	API_PORT        string `envconfig:"API_PORT"`
	DATABASE        string `envconfig:"DATABASE"`
	REDIS_HOST      string `envconfig:"REDIS_HOST"`
	REDIS_PASSWORD  string `envconfig:"REDIS_PASSWORD"`
	SIGN_KEY        string `envconfig:"SIGN_KEY"`
	BTC_HOST        string `envconfig:"BTC_HOST"`
	BTC_USER        string `envconfig:"BTC_USER"`
	BTC_PASS        string `envconfig:"BTC_PASS"`
	BTC_ZMQ_HASH_TX string `envconfig:"BTC_ZMQ_HASH_TX"`
	LND_HOST        string `envconfig:"LND_HOST"`
	LND_MACAROON    string `envconfig:"LND_MACAROON"`
	LND_TLS_CERT    string `envconfig:"LND_TLS_CERT"`
}

func Loads() error {
	if err := envconfig.Process("", &Config); err != nil {
		return err
	}
	return nil
}
