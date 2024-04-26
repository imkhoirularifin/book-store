package customer

import (
	"errors"
	"gramedia-service/internal/domain"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type customerService struct {
	customerRepo domain.CustomerRepository
}

// Count
func (c *customerService) Count(filter *domain.Customer) (int64, error) {
	count, err := c.customerRepo.Count(filter)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// Delete
func (c *customerService) Delete(id uint) error {
	return c.customerRepo.Delete(id)
}

// Fetch
func (c *customerService) Fetch(page int, size int, filter *domain.Customer) ([]*domain.Customer, int, error) {
	customers, nextCursor, err := c.customerRepo.Fetch(page, size, filter)
	if err != nil {
		return nil, 0, err
	}

	return customers, nextCursor, nil
}

// GetById
func (c *customerService) GetById(id uint) (*domain.Customer, error) {
	customer, err := c.customerRepo.GetById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.ErrNotFound
		}
		return nil, err
	}

	return customer, nil
}

// Store
func (c *customerService) Store(customer *domain.Customer) error {
	return c.customerRepo.Store(customer)
}

// Update
func (c *customerService) Update(customer *domain.Customer) error {
	return c.customerRepo.Update(customer)
}

func NewCustomerService(customerRepo domain.CustomerRepository) domain.CustomerService {
	return &customerService{customerRepo: customerRepo}
}
