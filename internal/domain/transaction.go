package domain

import (
	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	UserId             uint                 `json:"user_id" gorm:"not null" validate:"required"`
	User               *User                `json:"user,omitempty" gorm:"foreignKey:UserId"`
	CustomerId         uint                 `json:"customer_id" gorm:"not null" validate:"required"`
	Customer           *Customer            `json:"customer,omitempty" gorm:"foreignKey:CustomerId"`
	TotalPrice         int                  `json:"total_price" gorm:"not null" validate:"required"`
	TransactionDetails []*TransactionDetail `json:"transaction_details,omitempty"`
}

type TransactionStoreRequest struct {
	UserId             uint                             `json:"user_id" validate:"required"`
	CustomerId         uint                             `json:"customer_id" validate:"required"`
	TransactionDetails []*TransactionDetailStoreRequest `json:"transaction_details" validate:"required"`
}

type TransactionUpdateRequest struct {
	UserId             uint                 `json:"user_id"`
	CustomerId         uint                 `json:"customer_id"`
	TransactionDetails []*TransactionDetail `json:"transaction_details"`
}

type TransactionRepository interface {
	Fetch(page int, size int, filter *Transaction) ([]*Transaction, int, error)
	GetById(id uint) (*Transaction, error)
	Count(filter *Transaction) (int64, error)
	Store(transaction *Transaction) error
	Update(transaction *Transaction) error
	Delete(id uint) error
}

type TransactionService interface {
	Fetch(page int, size int, filter *Transaction) ([]*Transaction, int, error)
	GetById(id uint) (*Transaction, error)
	Count(filter *Transaction) (int64, error)
	Store(transaction *TransactionStoreRequest) error
	Update(transaction *Transaction) error
	Delete(id uint) error
}
