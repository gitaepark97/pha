// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0

package repository

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type ProductSize string

const (
	ProductSizeSmall ProductSize = "small"
	ProductSizeLarge ProductSize = "large"
)

func (e *ProductSize) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = ProductSize(s)
	case string:
		*e = ProductSize(s)
	default:
		return fmt.Errorf("unsupported scan type for ProductSize: %T", src)
	}
	return nil
}

type NullProductSize struct {
	ProductSize ProductSize
	Valid       bool // Valid is true if ProductSize is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullProductSize) Scan(value interface{}) error {
	if value == nil {
		ns.ProductSize, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.ProductSize.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullProductSize) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.ProductSize), nil
}

type Product struct {
	ID             int64       `json:"id"`
	UserID         int64       `json:"user_id"`
	Category       string      `json:"category"`
	Price          int32       `json:"price"`
	Cost           int32       `json:"cost"`
	Name           string      `json:"name"`
	Description    string      `json:"description"`
	Barcode        string      `json:"barcode"`
	ExpirationDate time.Time   `json:"expiration_date"`
	Size           ProductSize `json:"size"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
}

type Session struct {
	ID           string    `json:"id"`
	UserID       int64     `json:"user_id"`
	RefreshToken string    `json:"refresh_token"`
	UserAgent    string    `json:"user_agent"`
	ClientIp     string    `json:"client_ip"`
	IsBlocked    bool      `json:"is_blocked"`
	ExpiredAt    time.Time `json:"expired_at"`
	CreatedAt    time.Time `json:"created_at"`
}

type User struct {
	ID             int64     `json:"id"`
	PhoneNumber    string    `json:"phone_number"`
	HashedPassword string    `json:"hashed_password"`
	CreatedAt      time.Time `json:"created_at"`
}
