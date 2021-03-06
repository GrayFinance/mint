package services

import (
	"encoding/json"
	"fmt"

	"github.com/GrayFinance/mint/src/bitcoin"
	"github.com/GrayFinance/mint/src/lightning"
	"github.com/GrayFinance/mint/src/models"
	"github.com/GrayFinance/mint/src/storage"
)

type Receive struct {
	UserID   string
	UserTag  string
	WalletID string
}

func (r *Receive) GenerateAddress(network string) (models.Address, error) {
	address := models.Address{
		UserID:   r.UserID,
		WalletID: r.WalletID,
		Network:  network,
	}
	if network == "bitcoin" {
		addr, err := bitcoin.Bitcoin.GetNewAddress(r.WalletID)
		if err != nil {
			return address, err
		}
		address.Address = addr.String()
	}

	if storage.DB.Create(&address).Error != nil {
		err := fmt.Errorf("It was not possible to generate any address.")
		return address, err
	}
	return address, nil
}

func (r *Receive) GetAddress(network string) (models.Address, error) {
	var address models.Address

	if storage.DB.Model(address).Where("user_id = ? AND wallet_id = ? AND network = ?", r.UserID, r.WalletID, network).First(&address).Error != nil {
		err := fmt.Errorf("There is no address.")
		return address, err
	}
	return address, nil
}

func (r *Receive) CreateInvoice(value int64, description string) (models.Payment, error) {
	invoice, err := lightning.Lightning.CreateInvoice(int(value), description)
	if err != nil {
		err := fmt.Errorf("It was not possible to generate invoice.")
		return models.Payment{}, err
	}

	decode_invoice, err := lightning.Lightning.DecodeInvoice(invoice.Get("payment_request").String())
	if err != nil {
		err := fmt.Errorf("It was not possible decode invoice.")
		return models.Payment{}, err
	}

	payment := models.Payment{
		Pending:     true,
		AssetID:     "bitcoin",
		AssetName:   "bitcoin",
		Value:       value,
		Description: description,
		HashID:      decode_invoice.Get("payment_hash").String(),
		Invoice:     invoice.Get("payment_request").String(),
		Category:    "deposit",
		Network:     "lightning",
	}
	if r.UserID != "" && r.WalletID != "" {
		payment.UserID = r.UserID
		payment.WalletID = r.WalletID
	} else {
		var user models.User
		if storage.DB.Model(user).Where("user_tag = ?", r.UserTag).First(&user).Error != nil {
			err := fmt.Errorf("Not found user.")
			return models.Payment{}, err
		}

		var wallet models.Wallet
		if storage.DB.Model(wallet).Where("label = ? AND user_id = ?", "default", user.UserID).First(&wallet).Error != nil {
			err := fmt.Errorf("Not found wallet.")
			return models.Payment{}, err
		}
		payment.UserID = wallet.UserID
		payment.WalletID = wallet.WalletID
	}

	data, err := json.Marshal(payment)
	if err != nil {
		return models.Payment{}, err
	}

	if err = storage.REDIS.Set(payment.HashID, data, 0).Err(); err != nil {
		err := fmt.Errorf("It was not possible to generate invoice.")
		return models.Payment{}, err
	}
	return payment, nil
}
