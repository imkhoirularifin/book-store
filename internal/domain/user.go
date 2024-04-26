package domain

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `json:"name" gorm:"not null"`
	Email    string `json:"email" gorm:"not null;unique"`
	Password string `json:"-" gorm:"not null"`
	RoleId   uint   `json:"role_id" gorm:"not null"`
	Role     *Role  `json:"role,omitempty" gorm:"foreignKey:RoleId"`
}

type UserStoreRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	RoleId   uint   `json:"role_id" validate:"required"`
}

type UserUpdateRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	RoleId   uint   `json:"role_id"`
}

type UserRepository interface {
	Fetch(page int, size int, filter *User) ([]*User, int, error)
	GetById(id uint) (*User, error)
	GetByEmail(email string) (*User, error)
	Count(filter *User) (int64, error)
	Store(user *User) error
	Update(user *User) error
	Delete(id uint) error
}

type UserService interface {
	Fetch(page int, size int, filter *User) ([]*User, int, error)
	GetById(id uint) (*User, error)
	Count(filter *User) (int64, error)
	Store(user *User) error
	Update(user *User) error
	Delete(id uint) error
}
