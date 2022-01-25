package bitcoin

import (
	"encoding/hex"
	"log"
	"math"

	"github.com/GrayFinance/go-bitcoin"
	"github.com/GrayFinance/mint/src/config"
	"github.com/GrayFinance/mint/src/models"
	"github.com/GrayFinance/mint/src/storage"
	"github.com/pebbe/zmq4"
)

var Bitcoin *bitcoin.Bitcoin

func Start() {
	Bitcoin = bitcoin.Connect(
		config.Config.BTC_HOST,
		config.Config.BTC_USER,
		config.Config.BTC_PASS,
	)

	sub, err := zmq4.NewSocket(zmq4.SUB)
	if err != nil {
		log.Fatal(err)
		return
	}

	sub.Connect(config.Config.BTC_ZMQ_HASH_TX)
	sub.SetSubscribe("hashtx")

	log.Println(config.Config.BTC_ZMQ_HASH_TX)

	for {
		data, err := sub.RecvMessageBytes(0)
		if err != nil {
			log.Println(err)
			continue
		}

		topic := string(data[0])
		if topic == "hashtx" {
			txid := hex.EncodeToString(data[1])

			tx, err := Bitcoin.GetTransaction(txid)
			if err != nil {
				log.Println(err)
				continue
			}

			if tx.Get("confirmations").Uint() >= 1 {
				for _, detail := range tx.Get("details").Array() {
					if detail.Get("category").String() == "receive" {
						value := int64(detail.Get("amount").Float() * math.Pow(10, 8))
						if value >= 10000 {
							var address models.Address
							if storage.DB.Model(address).Where("address = ? AND network = ?", detail.Get("address").String(), "bitcoin").First(&address).Error != nil {
								continue
							}

							var wallet models.Wallet
							if storage.DB.Model(wallet).Where("user_id = ? AND wallet_id = ? ", address.UserID, address.WalletID).First(&wallet).Error != nil {
								continue
							}

							payment := models.Payment{
								Pending:   false,
								Value:     value,
								AssetID:   "bitcoin",
								AssetName: "bitcoin",
								HashID:    txid,
								Network:   "bitcoin",
								Category:  "deposit",
								UserID:    address.UserID,
								WalletID:  address.WalletID,
							}
							if storage.DB.Create(&payment).Error != nil {
								continue
							}

							balance := wallet.Balance + payment.Value
							if storage.DB.Model(&wallet).Update("balance", balance).Error != nil {
								continue
							}
						}
					}
				}
			}
		}
	}
}
