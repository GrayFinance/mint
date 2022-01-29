package services

import (
	"fmt"

	"github.com/GrayFinance/mint/src/models"
	"github.com/GrayFinance/mint/src/storage"
	"github.com/GrayFinance/mint/src/utils"
)

type Wallet struct {
	UserID string `json:"userid"`
}

func (w *Wallet) CreateWallet(label string) (models.Wallet, error) {
	wallet := models.Wallet{
		Label:          label,
		UserID:         w.UserID,
		WalletID:       utils.RandomHex(16),
		WalletReadKey:  utils.RandomHex(16),
		WalletAdminKey: utils.RandomHex(16),
	}

	if storage.DB.Create(&wallet).Error != nil {
		err := fmt.Errorf("It was not able to create the wallet.")
		return models.Wallet{}, err
	}
	return wallet, nil
}

func (w *Wallet) GetWallet(wallet_id string) (models.Wallet, error) {
	var wallet models.Wallet

	if storage.DB.Where("user_id = ? AND wallet_id = ? ", w.UserID, wallet_id).First(&wallet).Error != nil {
		err := fmt.Errorf("Wallet not found.")
		return wallet, err
	}
	return wallet, nil
}

func (w *Wallet) DeleteWallet(wallet_id string) (models.Wallet, error) {
	var wallet models.Wallet

	if storage.DB.Model(wallet).Where("user_id = ? and wallet_id = ?", w.UserID, wallet_id).First(&wallet).Error != nil {
		err := fmt.Errorf("Wallet not found.")
		return wallet, err
	}

	if wallet.Balance > 0 {
		err := fmt.Errorf("Your wallet can not be deleted, the balance must be empty.")
		return wallet, err
	}

	if storage.DB.Unscoped().Delete(&wallet).Error != nil {
		err := fmt.Errorf("It is not possible to delete the wallet.")
		return wallet, err
	}
	return wallet, nil
}

func (w *Wallet) RenameWallet(wallet_id string, label string) (models.Wallet, error) {
	var wallet models.Wallet

	if storage.DB.Model(wallet).Where("user_id = ? and wallet_id = ?", w.UserID, wallet_id).First(&wallet).Error != nil {
		err := fmt.Errorf("Wallet not found.")
		return wallet, err
	}

	wallet.Label = label
	if storage.DB.Save(&wallet).Error != nil {
		err := fmt.Errorf("It was not possible to rename the wallet.")
		return wallet, err
	}
	return wallet, nil
}

func (w *Wallet) ListWallets() ([]models.Wallet, error) {
	var wallets []models.Wallet

	if storage.DB.Model(models.Wallet{}).Where("user_id = ?", w.UserID).Find(&wallets).Error != nil {
		err := fmt.Errorf("Wallets not found.")
		return wallets, err
	}
	return wallets, nil
}
