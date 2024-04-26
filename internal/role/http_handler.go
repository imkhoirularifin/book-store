package role

import (
	"errors"
	"gramedia-service/internal/domain"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type HttpRoleHandler struct {
	roleSvc domain.RoleService
}

func NewHttpHandler(r fiber.Router, roleSvc domain.RoleService) {
	handler := &HttpRoleHandler{roleSvc: roleSvc}

	r.Get("/", handler.Fetch)
	r.Get("/:id", handler.GetById)
}

// Fetch used to get list of role
//
//	@Summary		Get list of role
//	@Description	Get list of roles
//	@Tags			roles
//	@Accept			json
//	@Produce		json
//	@Param			page	query		int				false	"Page number (default 1)"
//	@Param			size	query		int				false	"Size of page (default 10)"
//	@Param			q		query		string			false	"Search query"
//	@Header			200		{string}	X-Cursor		"Next page"
//	@Header			200		{string}	X-Total-Count	"Total item"
//	@Header			200		{string}	X-Max-Page		"Max page"
//	@Success		200		{array}		domain.Success	"List of roles"
//	@Failure		400		{object}	domain.Error	"Bad Request"
//	@Failure		500		{object}	domain.Error	"Internal Server Error"
//	@Router			/roles [get]
func (h *HttpRoleHandler) Fetch(c *fiber.Ctx) error {
	page, size, query := c.QueryInt("page", 1), c.QueryInt("size", 10), c.Query("q")
	if page <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(domain.Error{
			Code:    fiber.StatusBadRequest,
			Message: "page must be a positive integer",
		})
	}
	if size <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(domain.Error{
			Code:    fiber.StatusBadRequest,
			Message: "size must be a positive integer",
		})
	}

	filter := &domain.Role{Name: query}

	roles, nextPage, err := h.roleSvc.Fetch(page, size, filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.Error{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	if roles == nil {
		return c.Status(fiber.StatusNotFound).JSON(domain.Error{
			Code:    fiber.StatusNotFound,
			Message: "roles not found",
		})
	}

	totalItem, err := h.roleSvc.Count(filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.Error{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	maxPage := int(totalItem) / size

	if nextPage > 0 && nextPage <= maxPage {
		c.Set("X-Cursor", strconv.Itoa(nextPage))
	}
	c.Set("X-Total-Count", strconv.Itoa(int(totalItem)))
	c.Set("X-Max-Page", strconv.Itoa(maxPage))

	return c.Status(fiber.StatusOK).JSON(domain.Success{
		Code:    fiber.StatusOK,
		Message: "success",
		Data:    roles,
	})
}

// GetByID used to get role by id
//
//	@Summary		Get role by id
//	@Description	Get role by id
//	@Tags			roles
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int				true	"role ID"
//	@Success		200	{object}	domain.Success	"role detail"
//	@Failure		400	{object}	domain.Error	"Bad Request"
//	@Failure		404	{object}	domain.Error	"Not Found"
//	@Failure		500	{object}	domain.Error	"Internal Server Error"
//	@Router			/roles/{id} [get]
func (h *HttpRoleHandler) GetById(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.Error{
			Code:    fiber.StatusBadRequest,
			Message: "invalid role id",
		})
	}

	role, err := h.roleSvc.GetById(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(domain.Error{
				Code:    fiber.StatusNotFound,
				Message: err.Error(),
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
		Data:    role,
	})
}
