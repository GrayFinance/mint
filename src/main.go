package main

import (
	"log"

	"github.com/GrayFinance/mint/src/api"
	"github.com/GrayFinance/mint/src/bitcoin"
	"github.com/GrayFinance/mint/src/config"
	"github.com/GrayFinance/mint/src/lightning"
	"github.com/GrayFinance/mint/src/storage"
)

func init() {
	if err := config.Loads(); err != nil {
		log.Fatal(err)
	}

	if err := storage.Connect(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	go bitcoin.Start()
	go lightning.Start()

	api.Start()
}
