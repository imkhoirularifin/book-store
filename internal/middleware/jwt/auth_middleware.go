package jwt

import (
	"book-store/internal/domain"
	"book-store/internal/utilities"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type AuthMiddleware interface {
	RequireRole(roles ...string) fiber.Handler
}

type authMiddleware struct {
	jwtService utilities.JwtTokenService
}

// RequireRole
func (a *authMiddleware) RequireRole(roles ...string) func(*fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		authHeader := ctx.Get("Authorization")

		tokenString := strings.Replace(authHeader, "Bearer ", "", -1)
		if tokenString == "" {
			return ctx.Status(fiber.StatusUnauthorized).JSON(domain.Error{
				Code:    fiber.StatusUnauthorized,
				Message: "unauthorized",
			})
		}

		claims, err := a.jwtService.VerifyToken(tokenString)
		if err != nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(domain.Error{
				Code:    fiber.StatusUnauthorized,
				Message: "unauthorized",
			})
		}

		userName, ok := claims["user_name"].(string)
		if !ok {
			return ctx.Status(fiber.StatusUnauthorized).JSON(domain.Error{
				Code:    fiber.StatusUnauthorized,
				Message: "unauthorized",
			})
		}

		ctx.Set("user", userName)

		var validRole bool
		for _, role := range roles {
			if role == claims["role_name"] {
				validRole = true
				break
			}
		}

		if !validRole {
			return ctx.Status(fiber.StatusForbidden).JSON(domain.Error{
				Code:    fiber.StatusForbidden,
				Message: "this role is not allowed to access this resource",
			})
		}

		return ctx.Next()
	}
}

func NewAuthMiddleware(jwtService utilities.JwtTokenService) AuthMiddleware {
	return &authMiddleware{
		jwtService: jwtService,
	}
}
