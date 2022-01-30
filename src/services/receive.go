package services

import (
	"fmt"

	"github.com/GrayFinance/mint/src/bitcoin"
	"github.com/GrayFinance/mint/src/lightning"
	"github.com/GrayFinance/mint/src/models"
	"github.com/GrayFinance/mint/src/storage"
)

type Receive struct {
	UserID   string
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

func (r *Receive) CreateInvoice(value int, memo string) (models.Payment, error) {
	invoice, err := lightning.Lightning.CreateInvoice(value, memo)
	if err != nil {
		return models.Payment{}, nil
	}

	decode_invoice, err := lightning.Lightning.DecodeInvoice(invoice.Get("payment_request").String())
	if err != nil {
		return models.Payment{}, nil
	}

	payment := models.Payment{
		Pending:     true,
		AssetID:     "bitcoin",
		AssetName:   "bitcoin",
		Value:       uint64(value),
		Description: memo,
		HashID:      decode_invoice.Get("payment_hash").String(),
		Invoice:     invoice.Get("payment_request").String(),
		Category:    "deposit",
		Network:     "lightning",
		UserID:      r.UserID,
		WalletID:    r.WalletID,
	}
	if storage.DB.Create(&payment).Error != nil {
		err := fmt.Errorf("It was not possible to generate invoice.")
		return models.Payment{}, err
	}
	return payment, nil
}
