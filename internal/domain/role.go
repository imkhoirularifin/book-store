package domain

import (
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	Name string `json:"name" gorm:"not null"`
}

type RoleRepository interface {
	Fetch(page int, size int, filter *Role) ([]*Role, int, error)
	GetById(id uint) (*Role, error)
	Count(filter *Role) (int64, error)
}

type RoleService interface {
	Fetch(page int, size int, filter *Role) ([]*Role, int, error)
	GetById(id uint) (*Role, error)
	Count(filter *Role) (int64, error)
}
