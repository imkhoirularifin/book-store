package utilities

import (
	"encoding/base64"
	"errors"
	"fmt"
	"book-store/internal/config"
	"book-store/internal/domain"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtTokenService interface {
	GenerateToken(payload *domain.User) (domain.Token, error)
	VerifyToken(tokenString string) (jwt.MapClaims, error)
}

type jwtTokenService struct {
	cfg config.Config
}

// GenerateToken
func (j *jwtTokenService) GenerateToken(payload *domain.User) (domain.Token, error) {
	// Decode private key
	decodedPrivateKey, err := base64.StdEncoding.DecodeString(j.cfg.JwtConfig.PrivateKey)
	if err != nil {
		return domain.Token{}, fmt.Errorf("could not decode key: %w", err)
	}

	// Parse decoded private key
	key, err := jwt.ParseRSAPrivateKeyFromPEM(decodedPrivateKey)
	if err != nil {
		return domain.Token{}, fmt.Errorf("failed to parse decoded private key: %w", err)
	}

	// Claims is jwt payload
	claims := domain.JwtTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    payload.Name,
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Duration(j.cfg.JwtConfig.ExpiresIn))),
		},
		UserName: payload.Name,
		RoleName: payload.Role.Name,
	}

	// Sign token
	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)
	if err != nil {
		return domain.Token{}, fmt.Errorf("failed to sign token: %w", err)
	}

	return domain.Token{Token: token}, nil
}

// VerifyToken
func (j *jwtTokenService) VerifyToken(tokenString string) (jwt.MapClaims, error) {
	// Decode public key
	decodedPublicKey, err := base64.StdEncoding.DecodeString(j.cfg.JwtConfig.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("could not decode public key %w", err)
	}

	// Parse public key
	key, err := jwt.ParseRSAPublicKeyFromPEM(decodedPublicKey)
	if err != nil {
		return nil, fmt.Errorf("validate: parse key: %w", err)
	}

	// Parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return key, nil
	})

	if err != nil {
		return nil, fmt.Errorf("error parsing token: %w", err)
	}

	// Validate token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func NewJwtTokenService(cfg config.Config) JwtTokenService {
	return &jwtTokenService{
		cfg: cfg,
	}
}
