package domain

import "github.com/golang-jwt/jwt/v5"

type JwtTokenClaims struct {
	jwt.RegisteredClaims
	// list yang dibuat di payload
	UserName string `json:"user_name"`
	RoleName string `json:"role_name"`
}

type Token struct {
	Token string `json:"token"`
}

type AuthRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthService interface {
	GetToken(userCredential *AuthRequest) (Token, error)
}
