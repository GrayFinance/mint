package models

import "gorm.io/gorm"

type User struct {
	gorm.Model

	UserID       string `gorm:"not null;unique" json:"user_id"`
	TagName      string `gorm:"not null;unique" json:"tag_name"`
	Username     string `gorm:"not null;unique" json:"username"`
	Password     string `gorm:"not null" json:"password"`
	MasterAPIKey string `gorm:"not null;unique" json:"master_api_key"`
}

type Wallet struct {
	gorm.Model

	UserID         string `gorm:"index;not null" json:"user_id"`
	Label          string `gorm:"not null" json:"label"`
	Balance        uint64 `gorm:"default 0" json:"balance"`
	WalletID       string `gorm:"not null" json:"wallet_id"`
	WalletReadKey  string `gorm:"not null" json:"wallet_read_key"`
	WalletAdminKey string `gorm:"not null" json:"wallet_admin_key"`
}

type Address struct {
	gorm.Model

	Address  string `gorm:"not null" json:"address"`
	Network  string `gorm:"not null;default bitcoin" json:"network"`
	UserID   string `gorm:"not null" json:"user_id"`
	WalletID string `gorm:"not null" json:"wallet_id"`
}

type Payment struct {
	gorm.Model

	Pending     bool   `gorm:"not null;default true" json:"pending"`
	AssetID     string `gorm:"not null" json:"asset_id"`
	AssetName   string `gorm:"not null;default bitcoin" json:"asset_name"`
	Value       uint64 `gorm:"not null" json:"value"`
	Fee         uint64 `gorm:"not null;default 0" json:"fee"`
	Description string `json:"description"`
	HashID      string `gorm:"not null" json:"hash_id"`
	Preimage    string `json:"preimage"`
	Invoice     string `json:"invoice"`
	Category    string `gorm:"not null" json:"category"`
	Network     string `gorm:"not null" json:"network"`
	UserID      string `gorm:"index;not null" json:"user_id"`
	WalletID    string `gorm:"index;not null" json:"wallet_id"`
}
