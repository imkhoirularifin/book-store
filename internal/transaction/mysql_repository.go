package transaction

import (
	"book-store/internal/domain"

	"gorm.io/gorm"
)

type mysqlTransactionRepository struct {
	db *gorm.DB
}

// Count implements domain.TransactionRepository.
func (m *mysqlTransactionRepository) Count(filter *domain.Transaction) (int64, error) {
	var count int64
	query := m.db.Model(&domain.Transaction{})

	if filter.CustomerId > 0 {
		query = query.Where("customer_id = ?", filter.CustomerId)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

// Delete
func (m *mysqlTransactionRepository) Delete(id uint) error {
	return m.db.Delete(&domain.Transaction{}, id).Error
}

// Fetch
func (m *mysqlTransactionRepository) Fetch(page int, size int, filter *domain.Transaction) ([]*domain.Transaction, int, error) {
	var transactions []*domain.Transaction

	offset := (page - 1) * size
	query := m.db.Preload("TransactionDetails")

	if filter.CustomerId > 0 {
		query = query.Where("customer_id = ?", filter.CustomerId)
	}

	if err := query.Order("created_at DESC").Offset(offset).Limit(size).Find(&transactions).Error; err != nil {
		return nil, 0, err
	}

	var nextCursor int
	if len(transactions) > 0 {
		nextCursor = page + 1 // Next page
	}

	return transactions, nextCursor, nil
}

// GetById
func (m *mysqlTransactionRepository) GetById(id uint) (*domain.Transaction, error) {
	var transaction *domain.Transaction

	if err := m.db.Preload("TransactionDetails").Preload("User").Preload("Customer").First(&transaction, id).Error; err != nil {
		return nil, err
	}

	// get book details
	transactionDetails := make([]*domain.TransactionDetail, len(transaction.TransactionDetails))
	for _, detail := range transaction.TransactionDetails {
		var book *domain.Book
		if err := m.db.First(&book, detail.BookId).Error; err != nil {
			return nil, err
		}
		detail.Book = book

		// set transactionDetails
		transactionDetails = append(transactionDetails, detail)
	}

	transaction.TransactionDetails = transactionDetails
	return transaction, nil
}

// Store
func (m *mysqlTransactionRepository) Store(transaction *domain.Transaction) error {
	return m.db.Create(transaction).Error
}

// Update
func (m *mysqlTransactionRepository) Update(transaction *domain.Transaction) error {
	return m.db.Updates(transaction).Error
}

func NewMysqlTransactionRepository(db *gorm.DB) domain.TransactionRepository {
	return &mysqlTransactionRepository{db: db}
}
