package infrastructure

import (
	"errors"
	"gramedia-service/internal/domain"

	"github.com/gofiber/fiber/v2"
)

func defaultErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	msg := fiber.ErrInternalServerError.Error()

	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
		msg = e.Message
	}

	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
	return c.Status(code).JSON(&domain.Error{
		Code:    code,
		Message: msg,
	})
}
