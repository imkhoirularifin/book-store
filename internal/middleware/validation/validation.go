package validation

import (
	"gramedia-service/internal/domain"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func New[V any]() fiber.Handler {
	validate := validator.New(validator.WithRequiredStructEnabled())
	return func(c *fiber.Ctx) error {
		var v V
		if err := c.BodyParser(&v); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		if err := validate.Struct(v); err != nil {
			var errors []string
			for _, err := range err.(validator.ValidationErrors) {
				message := err.Field() + " is " + err.Tag()
				if err.Param() != "" {
					message += " " + err.Param()
				}
				errors = append(errors, message)
			}
			return c.Status(fiber.StatusBadRequest).JSON(domain.Error{
				Code:    fiber.StatusBadRequest,
				Errors:  errors,
				Message: "validation error",
			})
		}
		c.Locals("parser", &v)
		return c.Next()
	}
}
