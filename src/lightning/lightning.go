package lightning

import (
	"encoding/json"
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
		return
	}

	for {
		data, err := sub.ReadBytes('\n')
		if err != nil {
			log.Fatal(2, err)
		}

		payload_invoice := gjson.ParseBytes(data)
		if payload_invoice.Type != gjson.JSON {
			continue
		}

		if payload_invoice.Get("result").Get("settled").Bool() == false {
			continue
		}

		decode_invoice, err := Lightning.DecodeInvoice(payload_invoice.Get("result").Get("payment_request").String())
		if err != nil {
			continue
		}

		payment_hash := decode_invoice.Get("payment_hash").String()
		get_payment, err := storage.REDIS.Get(payment_hash).Result()
		if err != nil {
			continue
		}

		var payment models.Payment
		if err := json.Unmarshal([]byte(get_payment), &payment); err != nil {
			continue
		}
		payment.Pending = true

		var wallet models.Wallet
		if storage.DB.Model(wallet).Where("wallet_id = ?", payment.WalletID).First(&wallet).Error != nil {
			continue
		}

		if storage.DB.Create(&payment).Error != nil {
			continue
		}

		balance := wallet.Balance + payment.Value
		if storage.DB.Model(&wallet).Update("balance", balance).Error != nil {
			continue
		}

		pipe := storage.REDIS.Pipeline()
		pipe.Del(payment_hash)
		pipe.Exec()
	}
}
