package domain

import (
	"gorm.io/gorm"
)

type Customer struct {
	gorm.Model
	Name        string `json:"name" gorm:"not null"`
	Email       string `json:"email" gorm:"not null;unique"`
	PhoneNumber string `json:"phone_number" gorm:"not null"`
}

type CustomerStoreRequest struct {
	Name        string `json:"name" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	PhoneNumber string `json:"phone_number" validate:"required"`
}

type CustomerUpdateRequest struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}

type CustomerRepository interface {
	Fetch(page int, size int, filter *Customer) ([]*Customer, int, error)
	GetById(id uint) (*Customer, error)
	Count(filter *Customer) (int64, error)
	Store(customer *Customer) error
	Update(customer *Customer) error
	Delete(id uint) error
}

type CustomerService interface {
	Fetch(page int, size int, filter *Customer) ([]*Customer, int, error)
	GetById(id uint) (*Customer, error)
	Count(filter *Customer) (int64, error)
	Store(customer *Customer) error
	Update(customer *Customer) error
	Delete(id uint) error
}
