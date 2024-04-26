package domain

import (
	"time"

	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Title       string    `json:"title" gorm:"not null"`
	Author      string    `json:"author" gorm:"not null"`
	Price       int       `json:"price" gorm:"not null"`
	Description string    `json:"description" gorm:"not null"`
	Pages       int       `json:"pages" gorm:"not null"`
	Isbn        string    `json:"isbn" gorm:"not null"`
	Language    string    `json:"language" gorm:"not null"`
	Stock       int       `json:"stock" gorm:"not null"`
	PublishedAt time.Time `json:"published_at" gorm:"not null"`
}

type BookStoreRequest struct {
	Title       string `json:"title" validate:"required"`
	Author      string `json:"author" validate:"required"`
	Price       int    `json:"price" validate:"required"`
	Description string `json:"description" validate:"required"`
	Pages       int    `json:"pages" validate:"required"`
	Isbn        string `json:"isbn" validate:"required"`
	Language    string `json:"language" validate:"required"`
	Stock       int    `json:"stock" validate:"required"`
	PublishedAt string `json:"published_at" validate:"required"`
}

type BookUpdateRequest struct {
	Title       string `json:"title"`
	Author      string `json:"author"`
	Price       int    `json:"price"`
	Description string `json:"description"`
	Pages       int    `json:"pages"`
	Isbn        string `json:"isbn"`
	Language    string `json:"language"`
	Stock       int    `json:"stock"`
	PublishedAt string `json:"published_at"`
}

type BookService interface {
	Fetch(page int, size int, filter *Book) ([]*Book, int, error)
	GetById(id uint) (*Book, error)
	Count(filter *Book) (int64, error)
	Store(book *Book) error
	Update(book *Book) error
	Delete(id uint) error
}

type BookRepository interface {
	Fetch(page int, size int, filter *Book) ([]*Book, int, error)
	GetById(id uint) (*Book, error)
	Count(filter *Book) (int64, error)
	Store(book *Book) error
	Update(book *Book) error
	Delete(id uint) error
}
