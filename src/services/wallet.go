package services

import (
	"fmt"

	"github.com/GrayFinance/mint/src/models"
	"github.com/GrayFinance/mint/src/storage"
	"github.com/GrayFinance/mint/src/utils"
	"github.com/google/uuid"
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
	if storage.DB.Model(wallet).Where("user_id = ? AND wallet_id = ?", w.UserID, wallet_id).First(&wallet).Error != nil {
		err := fmt.Errorf("Wallet not found.")
		return wallet, err
	}

	if wallet.Label == "default" {
		err := fmt.Errorf("Wallet not possible delete.")
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
	if storage.DB.Model(wallet).Where("user_id = ? AND wallet_id = ?", w.UserID, wallet_id).First(&wallet).Error != nil {
		err := fmt.Errorf("Wallet not found.")
		return wallet, err
	}

	if wallet.Label == "default" {
		err := fmt.Errorf("Wallet not possible renamed.")
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

func (w *Wallet) Transfer(wallet_id string, destination string, value uint64, description string) (models.Payment, error) {
	var wallet models.Wallet
	if storage.DB.Model(wallet).Where("user_id = ? AND wallet_id = ?", w.UserID, wallet_id).First(&wallet).Error != nil {
		err := fmt.Errorf("Wallet not found.")
		return models.Payment{}, err
	}

	if value == 0 {
		err := fmt.Errorf("Value should be greater than zero to be able to transfer.")
		return models.Payment{}, err
	}

	if wallet.Balance < value {
		err := fmt.Errorf("Available balance is less than the value is transferred.")
		return models.Payment{}, err
	}

	var dest_user models.User
	if storage.DB.Model(dest_user).Where("tag_name = ?", destination).First(&dest_user).Error != nil {
		err := fmt.Errorf("Tag name not found.")
		return models.Payment{}, err
	}

	var dest_wallet models.Wallet
	if storage.DB.Model(dest_wallet).Where("user_id = ? AND label = ?", dest_user.UserID, "default").First(&dest_wallet).Error != nil {
		err := fmt.Errorf("Wallet not found.")
		return models.Payment{}, err
	}

	payment := models.Payment{
		Pending:     false,
		AssetID:     "bitcoin",
		AssetName:   "bitcoin",
		Value:       value,
		Description: description,
		HashID:      uuid.New().String(),
		Category:    "deposit",
		UserID:      dest_user.UserID,
		WalletID:    dest_wallet.WalletID,
	}
	if storage.DB.Model(payment).Create(&payment).Error != nil {
		err := fmt.Errorf("It was not possible to create the payment.")
		return models.Payment{}, err
	}

	payment.ID = payment.ID + 1
	payment.Category = "withdraw"
	payment.UserID = wallet.UserID
	payment.WalletID = wallet.WalletID
	if storage.DB.Model(payment).Create(&payment).Error != nil {
		err := fmt.Errorf("It was not possible to create the payment.")
		return models.Payment{}, err
	}

	if storage.DB.Model(&wallet).Update("balance", wallet.Balance-value).Error != nil {
		err := fmt.Errorf("It was not possible to update the wallet swing.")
		return models.Payment{}, err
	}

	if storage.DB.Model(&dest_wallet).Update("balance", dest_wallet.Balance+value).Error != nil {
		err := fmt.Errorf("It was not possible to update the wallet swing.")
		return models.Payment{}, err
	}
	return payment, nil
}
