package auth

import (
	"book-store/internal/domain"
	"book-store/internal/utilities"
	"errors"
)

type authService struct {
	userRepo   domain.UserRepository
	jwtService utilities.JwtTokenService
}

// GetToken
func (a *authService) GetToken(userCredential *domain.AuthRequest) (domain.Token, error) {
	user, err := a.userRepo.GetByEmail(userCredential.Email)
	if err != nil {
		return domain.Token{}, err
	}

	// Check password
	isMatch, err := utilities.ComparePassword(userCredential.Password, user.Password)
	if err != nil {
		return domain.Token{}, err
	}

	if !isMatch {
		return domain.Token{}, errors.New("invalid password")
	}

	// Generate token
	token, err := a.jwtService.GenerateToken(user)
	if err != nil {
		return domain.Token{}, err
	}

	return token, nil
}

func NewAuthService(userRepo domain.UserRepository, jwtService utilities.JwtTokenService) domain.AuthService {
	return &authService{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}
