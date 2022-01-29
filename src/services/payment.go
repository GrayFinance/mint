package services

import (
	"fmt"

	"github.com/GrayFinance/mint/src/models"
	"github.com/GrayFinance/mint/src/storage"
)

type Payment struct {
	WalletID string
	UserID   string
}

func (p *Payment) CreatePayment(payment models.Payment) (models.Payment, error) {
	wl := Wallet{UserID: p.UserID}

	wallet, err := wl.GetWallet(p.WalletID)
	if err != nil {
		return payment, err
	}

	balance := wallet.Balance
	if (payment.Category == "withdraw") && (payment.Value+payment.Fee > balance) {
		err := fmt.Errorf("Payment value is larger than the balance available on the Wallet.")
		return payment, err
	}

	if storage.DB.Create(&payment).Error != nil {
		err := fmt.Errorf("Not possible create payment.")
		return payment, err
	}

	if payment.Category == "withdraw" {
		balance = balance - (payment.Value + payment.Fee)
	}

	if payment.Category == "deposit" {
		balance = balance + payment.Value
	}

	if storage.DB.Model(&wallet).Update("balance", balance).Error != nil {
		err := fmt.Errorf("It was not possible to update the wallet.")
		return payment, err
	}
	return payment, nil
}
