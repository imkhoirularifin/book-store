package user

import (
	"book-store/internal/domain"
	"book-store/internal/utilities"
	"errors"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type userService struct {
	userRepo domain.UserRepository
}

// Count
func (u *userService) Count(filter *domain.User) (int64, error) {
	count, err := u.userRepo.Count(filter)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// Delete
func (u *userService) Delete(id uint) error {
	return u.userRepo.Delete(id)
}

// Fetch
func (u *userService) Fetch(page int, size int, filter *domain.User) ([]*domain.User, int, error) {
	users, nextCursor, err := u.userRepo.Fetch(page, size, filter)
	if err != nil {
		return nil, 0, err
	}

	return users, nextCursor, nil
}

// GetById
func (u *userService) GetById(id uint) (*domain.User, error) {
	user, err := u.userRepo.GetById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.ErrNotFound
		}
		return nil, err
	}

	return user, nil
}

// Store
func (u *userService) Store(user *domain.User) error {
	// Hash password
	hashedPassword, err := utilities.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	return u.userRepo.Store(user)
}

// Update
func (u *userService) Update(user *domain.User) error {
	return u.userRepo.Update(user)
}

func NewUserService(userRepo domain.UserRepository) domain.UserService {
	return &userService{
		userRepo: userRepo,
	}
}
