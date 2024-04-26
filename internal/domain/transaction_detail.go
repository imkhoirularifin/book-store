package domain

import (
	"gorm.io/gorm"
)

type TransactionDetail struct {
	gorm.Model
	TransactionId uint  `json:"transaction_id" gorm:"not null" validate:"required"`
	BookId        uint  `json:"book_id" gorm:"not null" validate:"required"`
	Book          *Book `json:"book,omitempty" gorm:"foreignKey:BookId"`
	Quantity      int   `json:"quantity" gorm:"not null" validate:"required"`
	SubTotal      int   `json:"sub_total" gorm:"not null" validate:"required"`
}

type TransactionDetailStoreRequest struct {
	BookId   uint `json:"book_id" validate:"required"`
	Quantity int  `json:"quantity" validate:"required"`
}

type TransactionDetailRepository interface {
	Store(transactionDetail *TransactionDetail) error
	Update(transactionDetail *TransactionDetail) error
	Delete(id uint) error
}

type TransactionDetailService interface {
	Store(transactionDetail *TransactionDetail) error
	Update(transactionDetail *TransactionDetail) error
	Delete(id uint) error
}
