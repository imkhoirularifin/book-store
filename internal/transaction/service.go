package transaction

import (
	"book-store/internal/domain"
	"errors"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type transactionService struct {
	transactionRepo domain.TransactionRepository
	bookRepo        domain.BookRepository
}

// Count implements domain.TransactionService.
func (t *transactionService) Count(filter *domain.Transaction) (int64, error) {
	count, err := t.transactionRepo.Count(filter)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// Delete
func (t *transactionService) Delete(id uint) error {
	return t.transactionRepo.Delete(id)
}

// Fetch
func (t *transactionService) Fetch(page int, size int, filter *domain.Transaction) ([]*domain.Transaction, int, error) {
	transaction, nextCursor, err := t.transactionRepo.Fetch(page, size, filter)
	if err != nil {
		return nil, 0, err
	}

	return transaction, nextCursor, err
}

// GetById
func (t *transactionService) GetById(id uint) (*domain.Transaction, error) {
	transaction, err := t.transactionRepo.GetById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.ErrNotFound
		}
		return nil, err
	}

	return transaction, nil
}

// Store
func (t *transactionService) Store(transactionReq *domain.TransactionStoreRequest) error {
	var totalPrice int
	transactionDetails := make([]*domain.TransactionDetail, len(transactionReq.TransactionDetails))

	for i, detail := range transactionReq.TransactionDetails {
		// get book information
		book, err := t.bookRepo.GetById(detail.BookId)
		if err != nil {
			return err
		}

		// check stock
		if book.Stock < detail.Quantity {
			return errors.New("stock not enough")
		}

		// update stock
		book.Stock -= detail.Quantity
		if err := t.bookRepo.Update(book); err != nil {
			return err
		}

		// set transactionDetails
		transactionDetails[i] = &domain.TransactionDetail{
			BookId:   detail.BookId,
			Quantity: detail.Quantity,
			SubTotal: book.Price * detail.Quantity,
		}

		totalPrice += book.Price * detail.Quantity
	}

	transaction := &domain.Transaction{
		UserId:             transactionReq.UserId,
		CustomerId:         transactionReq.CustomerId,
		TotalPrice:         totalPrice,
		TransactionDetails: transactionDetails,
	}
	return t.transactionRepo.Store(transaction)
}

// Update
func (t *transactionService) Update(transaction *domain.Transaction) error {
	return t.transactionRepo.Update(transaction)
}

func NewTransactionService(transactionRepo domain.TransactionRepository, bookRepo domain.BookRepository) domain.TransactionService {
	return &transactionService{
		transactionRepo: transactionRepo,
		bookRepo:        bookRepo,
	}
}
