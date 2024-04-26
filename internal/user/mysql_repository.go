package user

import (
	"gramedia-service/internal/domain"

	"gorm.io/gorm"
)

type mysqlUserRepository struct {
	db *gorm.DB
}

// GetByEmail
func (m *mysqlUserRepository) GetByEmail(email string) (*domain.User, error) {
	var user *domain.User

	if err := m.db.Where("email = ?", email).Preload("Role").First(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// Count
func (m *mysqlUserRepository) Count(filter *domain.User) (int64, error) {
	var count int64
	query := m.db.Model(&domain.User{})

	if filter.Name != "" {
		query = query.Where("name LIKE ?", "%"+filter.Name+"%")
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

// Delete
func (m *mysqlUserRepository) Delete(id uint) error {
	return m.db.Delete(&domain.User{}, id).Error
}

// Fetch
func (m *mysqlUserRepository) Fetch(page int, size int, filter *domain.User) ([]*domain.User, int, error) {
	var users []*domain.User

	offset := (page - 1) * size
	query := m.db

	if filter.Name != "" {
		query = query.Where("name LIKE ?", "%"+filter.Name+"%")
	}

	if err := query.Preload("Role").Order("created_at DESC").Offset(offset).Limit(size).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	var nextCursor int
	if len(users) > 0 {
		nextCursor = page + 1 // Next page
	}

	return users, nextCursor, nil
}

// GetById
func (m *mysqlUserRepository) GetById(id uint) (*domain.User, error) {
	var user *domain.User

	if err := m.db.Preload("Role").First(&user, id).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// Store
func (m *mysqlUserRepository) Store(user *domain.User) error {
	return m.db.Create(user).Error
}

// Update
func (m *mysqlUserRepository) Update(user *domain.User) error {
	return m.db.Updates(user).Error
}

func NewMysqlUserRepository(db *gorm.DB) domain.UserRepository {
	return &mysqlUserRepository{db: db}
}
