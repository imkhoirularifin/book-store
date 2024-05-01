package book

import (
	"errors"
	"book-store/internal/domain"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type bookService struct {
	bookRepo domain.BookRepository
}

// Count
func (b *bookService) Count(filter *domain.Book) (int64, error) {
	count, err := b.bookRepo.Count(filter)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// Delete
func (b *bookService) Delete(id uint) error {
	return b.bookRepo.Delete(id)
}

// Fetch
func (b *bookService) Fetch(page int, size int, filter *domain.Book) ([]*domain.Book, int, error) {
	books, nextCursor, err := b.bookRepo.Fetch(page, size, filter)
	if err != nil {
		return nil, 0, err
	}

	return books, nextCursor, nil
}

// GetById
func (b *bookService) GetById(id uint) (*domain.Book, error) {
	book, err := b.bookRepo.GetById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.ErrNotFound
		}
		return nil, err
	}

	return book, nil
}

// Store
func (b *bookService) Store(book *domain.Book) error {
	return b.bookRepo.Store(book)
}

// Update
func (b *bookService) Update(book *domain.Book) error {
	return b.bookRepo.Update(book)
}

func NewBookService(bookRepo domain.BookRepository) domain.BookService {
	return &bookService{
		bookRepo: bookRepo,
	}
}
