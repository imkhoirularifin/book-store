package auth

import (
	"errors"
	"gramedia-service/internal/domain"
	"gramedia-service/internal/middleware/validation"
	"gramedia-service/internal/utilities"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type HttpAuthHandler struct {
	authSvc domain.AuthService
}

func NewHttpHandler(r fiber.Router, authSvc domain.AuthService) {
	handler := &HttpAuthHandler{authSvc: authSvc}

	r.Post("/token", validation.New[domain.AuthRequest](), handler.GetToken)
}

// GetToken used to get JWT Token
//
//	@Summary		Get JWT Token
//	@Description	Get JWT Token
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			auth	body		domain.AuthRequest	true	"user credential"
//	@Success		200	{object}	domain.Success	"token detail"
//	@Failure		400	{object}	domain.Error	"Bad Request"
//	@Failure		404	{object}	domain.Error	"Not Found"
//	@Failure		500	{object}	domain.Error	"Internal Server Error"
//	@Router			/auth/token [post]
func (h *HttpAuthHandler) GetToken(c *fiber.Ctx) error {
	authReq := utilities.ExtractStructFromValidator[domain.AuthRequest](c)

	token, err := h.authSvc.GetToken(authReq)
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return c.Status(fiber.StatusUnauthorized).JSON(domain.Error{
				Code:    fiber.StatusUnauthorized,
				Message: "invalid email or password",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(domain.Error{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(domain.Success{
		Code:    fiber.StatusOK,
		Message: "success",
		Data:    token,
	})
}
