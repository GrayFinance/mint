package lightning

import (
	"log"

	"github.com/GrayFinance/go-lnd"
	"github.com/GrayFinance/mint/src/config"
	"github.com/GrayFinance/mint/src/models"
	"github.com/GrayFinance/mint/src/storage"
	"github.com/tidwall/gjson"
)

var Lightning *lnd.Lnd

func Start() {
	Lightning = lnd.Connect(
		config.Config.LND_HOST,
		config.Config.LND_TLS_CERT,
		config.Config.LND_MACAROON,
	)

	log.Println("LND Connect RPC: ", config.Config.LND_HOST)

	sub, err := Lightning.InvoicesSubscribe()
	if err != nil {
		log.Fatal(err)
	}

	for {
		data, err := sub.ReadBytes('\n')
		if err != nil {
			log.Fatal(2, err)
		}

		invoice := gjson.ParseBytes(data)
		if invoice.Type != gjson.JSON {
			continue
		}

		if invoice.Get("result").Get("settled").Bool() != false {
			continue
		}

		var payment models.Payment
		if storage.DB.Model(payment).Where("invoice = ? AND pending = true", invoice.Get("result").Get("payment_request").String()).First(&payment).Error != nil {
			continue
		}

		var wallet models.Wallet
		if storage.DB.Model(wallet).Where("wallet_id = ?", payment.WalletID).First(&wallet).Error != nil {
			continue
		}

		balance := wallet.Balance + payment.Value
		if storage.DB.Model(&payment).Update("pending", false).Error != nil {
			continue
		}

		if storage.DB.Model(&wallet).Update("balance", balance).Error != nil {
			continue
		}
	}
}
