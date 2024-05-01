package transaction

import (
	"errors"
	"book-store/internal/domain"
	"book-store/internal/middleware/jwt"
	"book-store/internal/middleware/validation"
	"book-store/internal/utilities"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type HttpTransactionHandler struct {
	transactionSvc domain.TransactionService
	authMiddleware jwt.AuthMiddleware
}

func NewHttpHandler(r fiber.Router, transactionSvc domain.TransactionService, authMiddleware jwt.AuthMiddleware) {
	handler := &HttpTransactionHandler{
		transactionSvc: transactionSvc,
		authMiddleware: authMiddleware,
	}

	r.Get("/", handler.Fetch)
	r.Get("/:id", handler.GetById)
	r.Post("/", handler.authMiddleware.RequireRole("admin", "employee"), validation.New[domain.TransactionStoreRequest](), handler.Store)
	r.Put("/:id", handler.authMiddleware.RequireRole("admin", "employee"), validation.New[domain.TransactionUpdateRequest](), handler.Update)
	r.Delete("/:id", handler.authMiddleware.RequireRole("admin"), handler.Delete)
}

// Fetch used to get list of transaction
//
//	@Summary		Get list of transaction
//	@Description	Get list of transactions
//	@Tags			transactions
//	@Accept			json
//	@Produce		json
//	@Param			page	query		int				false	"Page number (default 1)"
//	@Param			size	query		int				false	"Size of page (default 10)"
//	@Param			q		query		string			false	"Customer Id"
//	@Header			200		{string}	X-Cursor		"Next page"
//	@Header			200		{string}	X-Total-Count	"Total item"
//	@Header			200		{string}	X-Max-Page		"Max page"
//	@Success		200		{array}		domain.Success	"List of transactions"
//	@Failure		400		{object}	domain.Error	"Bad Request"
//	@Failure		500		{object}	domain.Error	"Internal Server Error"
//	@Router			/transactions [get]
func (h *HttpTransactionHandler) Fetch(c *fiber.Ctx) error {
	page, size, query := c.QueryInt("page", 1), c.QueryInt("size", 10), c.QueryInt("q")
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

	filter := &domain.Transaction{CustomerId: uint(query)}
	transactions, nextPage, err := h.transactionSvc.Fetch(page, size, filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.Error{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	if transactions == nil {
		return c.Status(fiber.StatusNotFound).JSON(domain.Error{
			Code:    fiber.StatusNotFound,
			Message: "transactions not found",
		})
	}

	totalItem, err := h.transactionSvc.Count(filter)
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
		Data:    transactions,
	})
}

// GetByID used to get transaction by id
//
//	@Summary		Get transaction by id
//	@Description	Get transaction by id
//	@Tags			transactions
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int				true	"transaction ID"
//	@Success		200	{object}	domain.Success	"transaction detail"
//	@Failure		400	{object}	domain.Error	"Bad Request"
//	@Failure		404	{object}	domain.Error	"Not Found"
//	@Failure		500	{object}	domain.Error	"Internal Server Error"
//	@Router			/transactions/{id} [get]
func (h *HttpTransactionHandler) GetById(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.Error{
			Code:    fiber.StatusBadRequest,
			Message: "invalid transaction id",
		})
	}

	transaction, err := h.transactionSvc.GetById(uint(id))
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
		Data:    transaction,
	})
}

// Store used to store transaction
//
//	@Summary		Store transaction
//	@Description	Store transaction
//	@Tags			transactions
//	@Accept			json
//	@Produce		json
//	@Param			transaction	body		domain.TransactionStoreRequest	true	"transaction data"
//	@Success		201		{object}	domain.Success				"transaction detail"
//	@Failure		400		{object}	domain.Error				"Bad Request"
//	@Failure		500		{object}	domain.Error				"Internal Server Error"
//	@Router			/transactions [post]
//
// @Security Bearer
func (h *HttpTransactionHandler) Store(c *fiber.Ctx) error {
	transactionReq := utilities.ExtractStructFromValidator[domain.TransactionStoreRequest](c)

	transaction := &domain.TransactionStoreRequest{
		UserId:             transactionReq.UserId,
		CustomerId:         transactionReq.CustomerId,
		TransactionDetails: transactionReq.TransactionDetails,
	}

	if err := h.transactionSvc.Store(transaction); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.Error{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(domain.Success{
		Code:    fiber.StatusCreated,
		Message: "success",
		Data:    transaction,
	})
}

// Update used to update transaction
//
//	@Summary		Update transaction
//	@Description	Update transaction
//	@Tags			transactions
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int							true	"transaction ID"
//	@Param			transaction	body		domain.TransactionUpdateRequest	true	"transaction data"
//	@Success		200		{object}	domain.Success				"transaction detail"
//	@Failure		400		{object}	domain.Error				"Bad Request"
//	@Failure		404		{object}	domain.Error				"Not Found"
//	@Failure		500		{object}	domain.Error				"Internal Server Error"
//	@Router			/transactions/{id} [put]
//
// @Security Bearer
func (h *HttpTransactionHandler) Update(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.Error{
			Code:    fiber.StatusBadRequest,
			Message: "invalid transaction id",
		})
	}

	transactionReq := utilities.ExtractStructFromValidator[domain.TransactionUpdateRequest](c)

	transaction := &domain.Transaction{
		Model:              gorm.Model{ID: uint(id)},
		UserId:             transactionReq.UserId,
		CustomerId:         transactionReq.CustomerId,
		TransactionDetails: transactionReq.TransactionDetails,
	}

	if err := h.transactionSvc.Update(transaction); err != nil {
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
		Data:    transaction,
	})
}

// Delete used to delete transaction
//
//	@Summary		Delete transaction
//	@Description	Delete transaction
//	@Tags			transactions
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int				true	"transaction ID"
//	@Success		200	{object}	domain.Success	"Success delete transaction"
//	@Failure		400	{object}	domain.Error	"Bad Request"
//	@Failure		500	{object}	domain.Error	"Internal Server Error"
//	@Router			/transactions/{id} [delete]
//
// @Security Bearer
func (h *HttpTransactionHandler) Delete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.Error{
			Code:    fiber.StatusBadRequest,
			Message: "invalid transaction id",
		})
	}

	if err := h.transactionSvc.Delete(uint(id)); err != nil {
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
