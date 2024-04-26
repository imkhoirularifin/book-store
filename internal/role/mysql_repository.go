package role

import (
	"gramedia-service/internal/domain"

	"gorm.io/gorm"
)

type mysqlRoleRepository struct {
	db *gorm.DB
}

// Count
func (m *mysqlRoleRepository) Count(filter *domain.Role) (int64, error) {
	var count int64
	query := m.db.Model(&domain.Role{})

	if filter.Name != "" {
		query = query.Where("name LIKE ?", "%"+filter.Name+"%")
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

// Fetch
func (m *mysqlRoleRepository) Fetch(page int, size int, filter *domain.Role) ([]*domain.Role, int, error) {
	var roles []*domain.Role

	offset := (page - 1) * size
	query := m.db

	if filter.Name != "" {
		query = query.Where("name LIKE ?", "%"+filter.Name+"%")
	}

	if err := query.Order("created_at DESC").Offset(offset).Limit(size).Find(&roles).Error; err != nil {
		return nil, 0, err
	}

	var nextCursor int
	if len(roles) > 0 {
		nextCursor = page + 1 // Next page
	}

	return roles, nextCursor, nil
}

// GetById
func (m *mysqlRoleRepository) GetById(id uint) (*domain.Role, error) {
	var role *domain.Role

	if err := m.db.First(&role, id).Error; err != nil {
		return nil, err
	}

	return role, nil
}

func NewMysqlRoleRepository(db *gorm.DB) domain.RoleRepository {
	return &mysqlRoleRepository{db: db}
}
