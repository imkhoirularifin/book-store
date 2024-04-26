package role

import (
	"book-store/internal/domain"
	"errors"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type roleService struct {
	roleRepo domain.RoleRepository
}

// Count
func (r *roleService) Count(filter *domain.Role) (int64, error) {
	count, err := r.roleRepo.Count(filter)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// Fetch
func (r *roleService) Fetch(page int, size int, filter *domain.Role) ([]*domain.Role, int, error) {
	roles, nextCursor, err := r.roleRepo.Fetch(page, size, filter)
	if err != nil {
		return nil, 0, err
	}

	return roles, nextCursor, nil
}

// GetById
func (r *roleService) GetById(id uint) (*domain.Role, error) {
	role, err := r.roleRepo.GetById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.ErrNotFound
		}
		return nil, err
	}

	return role, nil
}

func NewRoleService(roleRepo domain.RoleRepository) domain.RoleService {
	return &roleService{
		roleRepo: roleRepo,
	}
}
