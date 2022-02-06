package services

import (
	"fmt"
	"math"

	"github.com/GrayFinance/mint/src/bitcoin"
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
	var (
		wallet  models.Wallet
		payment models.Payment
	)
	if storage.DB.Model(wallet).Where("user_id = ? AND wallet_id = ?", w.UserID, wallet_id).First(&wallet).Error != nil {
		err := fmt.Errorf("Wallet not found.")
		return payment, err
	}

	if value == 0 {
		err := fmt.Errorf("Value should be greater than zero to be able to transfer.")
		return payment, err
	}

	if wallet.Balance < value {
		err := fmt.Errorf("Available balance is less than the value is transferred.")
		return payment, err
	}

	var dest_user models.User
	if storage.DB.Model(dest_user).Where("tag_name = ?", destination).First(&dest_user).Error != nil {
		err := fmt.Errorf("Tag name not found.")
		return payment, err
	}

	var dest_wallet models.Wallet
	if storage.DB.Model(dest_wallet).Where("user_id = ? AND label = ?", dest_user.UserID, "default").First(&dest_wallet).Error != nil {
		err := fmt.Errorf("Wallet not found.")
		return payment, err
	}

	payment = models.Payment{
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
		return payment, err
	}

	payment.ID = payment.ID + 1
	payment.UserID = wallet.UserID
	payment.Category = "withdraw"
	payment.WalletID = wallet.WalletID
	if storage.DB.Model(payment).Create(&payment).Error != nil {
		err := fmt.Errorf("It was not possible to create the payment.")
		return payment, err
	}

	if storage.DB.Model(&wallet).Update("balance", wallet.Balance-value).Error != nil {
		err := fmt.Errorf("Unable to update wallet balance.")
		return payment, err
	}

	if storage.DB.Model(&dest_wallet).Update("balance", dest_wallet.Balance+value).Error != nil {
		err := fmt.Errorf("Unable to update wallet balance.")
		return payment, err
	}
	return payment, nil
}

func (w *Wallet) BitcoinWithdraw(wallet_id string, address string, value uint64, feerate uint64, description string) (models.Payment, error) {
	var wallet models.Wallet
	if storage.DB.Model(wallet).Where("wallet_id = ? AND user_id = ?", wallet_id, w.UserID).First(&wallet).Error != nil {
		err := fmt.Errorf("Wallet not found.")
		return models.Payment{}, err
	}

	if value >= wallet.Balance {
		err := fmt.Errorf("Value is greater or equal to the wallet balance.")
		return models.Payment{}, err
	}

	if feerate <= 0 {
		err := fmt.Errorf("Feerate can not be 0.")
		return models.Payment{}, err
	}

	get_address_info, err := bitcoin.Bitcoin.GetAddressInfo(address)
	if err != nil || get_address_info.Get("ismine").Bool() == true {
		err := fmt.Errorf("It is not possible to send funds to an internal address.")
		return models.Payment{}, err
	}

	validate_address, err := bitcoin.Bitcoin.ValidateAddress(address)
	if err != nil {
		err := fmt.Errorf("It was not possible to validate the address.")
		return models.Payment{}, err
	}

	if validate_address.Get("isvalid").Bool() == false {
		err := fmt.Errorf("Invalid address.")
		return models.Payment{}, err
	}

	tx, err := bitcoin.Bitcoin.CreateRawTransaction(
		[]map[string]interface{}{},
		[]map[string]interface{}{{address: float64(value) / math.Pow(10, 8)}},
	)
	if err != nil {
		err := fmt.Errorf("Could not create transaction.")
		return models.Payment{}, err
	}

	tx, err = bitcoin.Bitcoin.FundRawTransaction(tx.String(), feerate)
	if err != nil {
		err := fmt.Errorf("Unable to fund the transaction.")
		return models.Payment{}, err
	}

	total_value := uint64(tx.Get("fee").Float()*math.Pow(10, 8)) + value
	if total_value > wallet.Balance {
		err := fmt.Errorf("Value is greater than the funds available in your wallet.")
		return models.Payment{}, err
	}

	get_all_balance, err := bitcoin.Bitcoin.GetBalance()
	if err != nil || uint64(get_all_balance.Float()*math.Pow(10, 8)) < total_value {
		err := fmt.Errorf("Your transaction was not processed as there is no onchain liquidity from our service at the moment. Please wait until later.")
		return models.Payment{}, err
	}

	signtx, err := bitcoin.Bitcoin.SignRawTransactionWithWallet(tx.Get("hex").String())
	if err != nil || signtx.Get("hex").String() == "" {
		err := fmt.Errorf("Unable to sign transaction.")
		return models.Payment{}, err
	}

	dec_tx, err := bitcoin.Bitcoin.DecodeRawTransaction(signtx.Get("hex").String())
	if err != nil || dec_tx.Get("txid").String() == "" {
		err := fmt.Errorf("Unable to decode the transaction.")
		return models.Payment{}, err
	}

	payment := models.Payment{
		Pending:     true,
		AssetID:     "bitcoin",
		AssetName:   "bitcoin",
		Value:       value,
		Fee:         total_value - value,
		Description: description,
		HashID:      dec_tx.Get("txid").String(),
		Category:    "withdraw",
		Network:     "bitcoin",
		UserID:      w.UserID,
		WalletID:    wallet.WalletID,
	}
	if storage.DB.Create(&payment).Error != nil {
		err := fmt.Errorf("Unable to create transaction.")
		return payment, err
	}

	if storage.DB.Model(&wallet).Update("balance", wallet.Balance-total_value).Error != nil {
		err := fmt.Errorf("Unable to update wallet balance.")
		return payment, err
	}

	send_raw_tx, err := bitcoin.Bitcoin.SendRawTransaction(signtx.Get("hex").String())
	if err != nil || send_raw_tx.String() == "" {
		err := fmt.Errorf("Cannot propagate transaction.")
		return payment, err
	}
	return payment, nil
}
