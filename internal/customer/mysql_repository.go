package customer

import (
	"book-store/internal/domain"

	"gorm.io/gorm"
)

type mysqlCustomerRepository struct {
	db *gorm.DB
}

// Count
func (m *mysqlCustomerRepository) Count(filter *domain.Customer) (int64, error) {
	var count int64
	query := m.db.Model(&domain.Customer{})

	if filter.Name != "" {
		query = query.Where("name LIKE ?", "%"+filter.Name+"%")
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

// Fetch
func (m *mysqlCustomerRepository) Fetch(page int, size int, filter *domain.Customer) ([]*domain.Customer, int, error) {
	var customers []*domain.Customer

	offset := (page - 1) * size
	query := m.db

	if filter.Name != "" {
		query = query.Where("name LIKE ?", "%"+filter.Name+"%")
	}

	if err := query.Order("created_at DESC").Offset(offset).Limit(size).Find(&customers).Error; err != nil {
		return nil, 0, err
	}

	var nextCursor int
	if len(customers) > 0 {
		nextCursor = page + 1 // Next page
	}

	return customers, nextCursor, nil
}

// GetById
func (m *mysqlCustomerRepository) GetById(id uint) (*domain.Customer, error) {
	var customer *domain.Customer

	if err := m.db.First(&customer, id).Error; err != nil {
		return nil, err
	}

	return customer, nil
}

// Store
func (m *mysqlCustomerRepository) Store(customer *domain.Customer) error {
	return m.db.Create(customer).Error
}

// Update
func (m *mysqlCustomerRepository) Update(customer *domain.Customer) error {
	return m.db.Updates(customer).Error
}

// Delete
func (m *mysqlCustomerRepository) Delete(id uint) error {
	return m.db.Delete(&domain.Customer{}, id).Error
}

func NewMysqlCustomerRepository(db *gorm.DB) domain.CustomerRepository {
	return &mysqlCustomerRepository{db: db}
}
