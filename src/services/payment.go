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
	wallet := Wallet{UserID: p.UserID}

	get_wallet, err := wallet.GetWallet(p.WalletID)
	if err != nil {
		return payment, err
	}

	balance := get_wallet.Balance
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

	if storage.DB.Model(&get_wallet).Update("balance", balance).Error != nil {
		err := fmt.Errorf("It was not possible to update the wallet.")
		return payment, err
	}
	return payment, nil
}

func (p *Payment) ListPayments(offset int) ([]models.Payment, error) {
	var payments []models.Payment

	if storage.DB.Model(models.Payment{}).Where("user_id = ? AND wallet_id = ?", p.UserID, p.WalletID).Limit(10).Offset(offset).Find(&payments).Error != nil {
		err := fmt.Errorf("No transaction found.")
		return payments, err
	}
	return payments, nil
}

func (p *Payment) GetPayment(hash_id string) (models.Payment, error) {
	var payment models.Payment
	if storage.DB.Model(payment).Where("hash_id = ? AND wallet_id = ? AND user_id = ?", hash_id, p.WalletID, p.UserID).First(&payment).Error != nil {
		err := fmt.Errorf("No transaction found.")
		return payment, err
	}
	return payment, nil
}
