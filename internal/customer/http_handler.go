package customer

import (
	"book-store/internal/domain"
	"book-store/internal/middleware/jwt"
	"book-store/internal/middleware/validation"
	"book-store/internal/utilities"
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type HttpCustomerHandler struct {
	customerSvc    domain.CustomerService
	authMiddleware jwt.AuthMiddleware
}

func NewHttpHandler(r fiber.Router, customerSvc domain.CustomerService, authMiddleware jwt.AuthMiddleware) {
	handler := &HttpCustomerHandler{
		customerSvc:    customerSvc,
		authMiddleware: authMiddleware,
	}

	r.Get("/", handler.Fetch)
	r.Get("/:id", handler.GetById)
	r.Post("/", authMiddleware.RequireRole("admin", "employee"), validation.New[domain.CustomerStoreRequest](), handler.Store)
	r.Put("/:id", authMiddleware.RequireRole("admin", "employee"), validation.New[domain.CustomerUpdateRequest](), handler.Update)
	r.Delete("/:id", authMiddleware.RequireRole("admin"), handler.Delete)
}

// Fetch used to get list of customer
//
//	@Summary		Get list of customer
//	@Description	Get list of customers
//	@Tags			customers
//	@Accept			json
//	@Produce		json
//	@Param			page	query		int				false	"Page number (default 1)"
//	@Param			size	query		int				false	"Size of page (default 10)"
//	@Param			q		query		string			false	"Search query"
//	@Header			200		{string}	X-Cursor		"Next page"
//	@Header			200		{string}	X-Total-Count	"Total item"
//	@Header			200		{string}	X-Max-Page		"Max page"
//	@Success		200		{array}		domain.Success	"List of customers"
//	@Failure		400		{object}	domain.Error	"Bad Request"
//	@Failure		500		{object}	domain.Error	"Internal Server Error"
//	@Router			/customers [get]
func (h *HttpCustomerHandler) Fetch(c *fiber.Ctx) error {
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

	filter := &domain.Customer{Name: query}
	customers, nextPage, err := h.customerSvc.Fetch(page, size, filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.Error{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	if customers == nil {
		return c.Status(fiber.StatusNotFound).JSON(domain.Error{
			Code:    fiber.StatusNotFound,
			Message: "customers not found",
		})
	}

	totalItem, err := h.customerSvc.Count(filter)
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
		Data:    customers,
	})
}

// GetByID used to get customer by id
//
//	@Summary		Get customer by id
//	@Description	Get customer by id
//	@Tags			customers
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int				true	"customer ID"
//	@Success		200	{object}	domain.Success	"customer detail"
//	@Failure		400	{object}	domain.Error	"Bad Request"
//	@Failure		404	{object}	domain.Error	"Not Found"
//	@Failure		500	{object}	domain.Error	"Internal Server Error"
//	@Router			/customers/{id} [get]
func (h *HttpCustomerHandler) GetById(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.Error{
			Code:    fiber.StatusBadRequest,
			Message: "invalid customer id",
		})
	}

	customer, err := h.customerSvc.GetById(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(domain.Error{
				Code:    fiber.StatusNotFound,
				Message: "customer not found",
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
		Data:    customer,
	})
}

// Store used to store customer
//
//	@Summary		Store customer
//	@Description	Store customer
//	@Tags			customers
//	@Accept			json
//	@Produce		json
//	@Param			customer	body		domain.CustomerStoreRequest	true	"customer data"
//	@Success		201		{object}	domain.Success				"customer detail"
//	@Failure		400		{object}	domain.Error				"Bad Request"
//	@Failure		500		{object}	domain.Error				"Internal Server Error"
//	@Router			/customers [post]
//
// @Security Bearer
func (h *HttpCustomerHandler) Store(c *fiber.Ctx) error {
	customerReq := utilities.ExtractStructFromValidator[domain.CustomerStoreRequest](c)

	customer := &domain.Customer{
		Name:        customerReq.Name,
		Email:       customerReq.Email,
		PhoneNumber: customerReq.PhoneNumber,
	}

	if err := h.customerSvc.Store(customer); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.Error{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(domain.Success{
		Code:    fiber.StatusCreated,
		Message: "success",
		Data:    customer,
	})
}

// Update used to update customer
//
//	@Summary		Update customer
//	@Description	Update customer
//	@Tags			customers
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int							true	"Customer ID"
//	@Param			customer	body		domain.CustomerUpdateRequest	true	"Customer data"
//	@Success		200		{object}	domain.Success				"Customer detail"
//	@Failure		400		{object}	domain.Error				"Bad Request"
//	@Failure		404		{object}	domain.Error				"Not Found"
//	@Failure		500		{object}	domain.Error				"Internal Server Error"
//	@Router			/customers/{id} [put]
//
// @Security Bearer
func (h *HttpCustomerHandler) Update(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.Error{
			Code:    fiber.StatusBadRequest,
			Message: "invalid customer id",
		})
	}

	customerReq := utilities.ExtractStructFromValidator[domain.CustomerUpdateRequest](c)

	customer := &domain.Customer{
		Model:       gorm.Model{ID: uint(id)},
		Name:        customerReq.Name,
		Email:       customerReq.Email,
		PhoneNumber: customerReq.PhoneNumber,
	}

	if err := h.customerSvc.Update(customer); err != nil {
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
		Data:    customer,
	})
}

// Delete used to delete customer
//
//	@Summary		Delete customer
//	@Description	Delete customer
//	@Tags			customers
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int				true	"Customer ID"
//	@Success		200	{object}	domain.Success	"Success delete customer"
//	@Failure		400	{object}	domain.Error	"Bad Request"
//	@Failure		500	{object}	domain.Error	"Internal Server Error"
//	@Router			/customers/{id} [delete]
//
// @Security Bearer
func (h *HttpCustomerHandler) Delete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.Error{
			Code:    fiber.StatusBadRequest,
			Message: "invalid customer id",
		})
	}

	if err := h.customerSvc.Delete(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.Error{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(domain.Success{
		Code:    fiber.StatusOK,
		Message: "success",
	})
}
