package book

import (
	"gramedia-service/internal/domain"

	"gorm.io/gorm"
)

type mysqlBookRepository struct {
	db *gorm.DB
}

// Count
func (m *mysqlBookRepository) Count(filter *domain.Book) (int64, error) {
	var count int64
	query := m.db.Model(&domain.Book{})

	if filter.Title != "" {
		query = query.Where("title LIKE ?", "%"+filter.Title+"%")
	}

	if filter.Author != "" {
		query = query.Where("author LIKE ?", "%"+filter.Author+"%")
	}

	if filter.Price > 0 {
		query = query.Where("price = ?", filter.Price)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

// Delete
func (m *mysqlBookRepository) Delete(id uint) error {
	return m.db.Delete(&domain.Book{}, id).Error
}

// Fetch
func (m *mysqlBookRepository) Fetch(page int, size int, filter *domain.Book) ([]*domain.Book, int, error) {
	var books []*domain.Book

	offset := (page - 1) * size
	query := m.db

	// Can filter by title/author/price
	if filter.Title != "" {
		query = query.Where("title LIKE ?", "%"+filter.Title+"%")
	}

	if filter.Author != "" {
		query = query.Where("author LIKE ?", "%"+filter.Author+"%")
	}

	if filter.Price > 0 {
		query = query.Where("price = ?", filter.Price)
	}

	if err := query.Order("created_at DESC").Offset(offset).Limit(size).Find(&books).Error; err != nil {
		return nil, 0, err
	}

	var nextCursor int
	if len(books) > 0 {
		nextCursor = page + 1 // Next page
	}

	return books, nextCursor, nil
}

// GetById
func (m *mysqlBookRepository) GetById(id uint) (*domain.Book, error) {
	var book *domain.Book

	if err := m.db.First(&book, id).Error; err != nil {
		return nil, err
	}

	return book, nil
}

// Store
func (m *mysqlBookRepository) Store(book *domain.Book) error {
	return m.db.Create(book).Error
}

// Update
func (m *mysqlBookRepository) Update(book *domain.Book) error {
	return m.db.Updates(book).Error
}

func NewMysqlBookRepository(db *gorm.DB) domain.BookRepository {
	return &mysqlBookRepository{db: db}
}
