package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserID   string `gorm:"not null;unique" json:"user_id"`
	Username string `gorm:"not null;unique" json:"username"`
	Password string `gorm:"not null" json:"password"`
}

type Wallet struct {
	gorm.Model
	UserID           string `gorm:"index;not null" json:"user_id"`
	Label            string `gorm:"not null" json:"label"`
	Balance          int64  `gorm:"default 0" json:"balance"`
	WalletID         string `gorm:"not null" json:"wallet_id"`
	WalletReadKey    string `gorm:"not null" json:"wallet_read_key"`
	WalletAdminKey   string `gorm:"not null" json:"wallet_admin_key"`
	WalletInvoiceKey string `gorm:"not null" json:"wallet_invoice_key"`
}

type Address struct {
	gorm.Model
	WalletID string `gorm:"not null" json:"wallet_id"`
	Bitcoin  string `gorm:"not null;unique" json:"bitcoin"`
}

type Payment struct {
	gorm.Model
	Pending     bool   `gorm:"not null;default true" json:"pending"`
	Amount      int64  `gorm:"not null" json:"amount"`
	Fee         int64  `gorm:"not null;default 0" json:"fee"`
	Description string `json:"description"`
	Preimage    string `json:"preimage"`
	Hash        string `gorm:"not null;unique" json:"hash"`
	Tag         string `json:"tag"`
	Bolt11      string `json:"bolt11"`
	Category    string `gorm:"not null" json:"Category"`
	Network     string `json:"network"`
}
