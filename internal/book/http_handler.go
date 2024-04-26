package book

import (
	"book-store/internal/domain"
	"book-store/internal/middleware/jwt"
	"book-store/internal/middleware/validation"
	"book-store/internal/utilities"
	"errors"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type HttpBookHandler struct {
	bookService    domain.BookService
	authMiddleware jwt.AuthMiddleware
}

func NewHttpHandler(r fiber.Router, bookService domain.BookService, authMiddleware jwt.AuthMiddleware) {
	handler := &HttpBookHandler{
		bookService:    bookService,
		authMiddleware: authMiddleware,
	}

	r.Get("/", handler.Fetch)
	r.Get("/:id", handler.GetById)
	r.Post("/", authMiddleware.RequireRole("admin"), validation.New[domain.BookStoreRequest](), handler.Store)
	r.Put("/:id", authMiddleware.RequireRole("admin"), validation.New[domain.BookUpdateRequest](), handler.Update)
	r.Delete("/:id", authMiddleware.RequireRole("admin"), handler.Delete)
}

// Fetch used to get list of book
//
//	@Summary		Get list of book
//	@Description	Get list of books
//	@Tags			books
//	@Accept			json
//	@Produce		json
//	@Param			page	query		int				false	"Page number (default 1)"
//	@Param			size	query		int				false	"Size of page (default 10)"
//	@Param			q		query		string			false	"Search query"
//	@Param			filterBy		query		string			false	"Filter by (title, author, price)"
//	@Header			200		{string}	X-Cursor		"Next page"
//	@Header			200		{string}	X-Total-Count	"Total item"
//	@Header			200		{string}	X-Max-Page		"Max page"
//	@Success		200		{array}		domain.Success	"List of books"
//	@Failure		400		{object}	domain.Error	"Bad Request"
//	@Failure		500		{object}	domain.Error	"Internal Server Error"
//	@Router			/books [get]
func (h *HttpBookHandler) Fetch(c *fiber.Ctx) error {
	page, size, query, filterBy := c.QueryInt("page", 1), c.QueryInt("size", 10), c.Query("q"), c.Query("filterBy")
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

	// read filterBy
	var filter domain.Book

	switch filterBy {
	case "title":
		filter.Title = query
	case "author":
		filter.Author = query
	case "price":
		price, err := strconv.Atoi(query)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(domain.Error{
				Code:    fiber.StatusBadRequest,
				Message: "price must be an integer",
			})
		}
		filter.Price = price
	}

	books, nextPage, err := h.bookService.Fetch(page, size, &filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.Error{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	if books == nil {
		return c.Status(fiber.StatusNotFound).JSON(domain.Error{
			Code:    fiber.StatusNotFound,
			Message: "books not found",
		})
	}

	totalItem, err := h.bookService.Count(&filter)
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
		Message: "books fetched successfully",
		Data:    books,
	})
}

// GetByID used to get book by id
//
//	@Summary		Get book by id
//	@Description	Get book by id
//	@Tags			books
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int				true	"book ID"
//	@Success		200	{object}	domain.Success	"book detail"
//	@Failure		400	{object}	domain.Error	"Bad Request"
//	@Failure		404	{object}	domain.Error	"Not Found"
//	@Failure		500	{object}	domain.Error	"Internal Server Error"
//	@Router			/books/{id} [get]
func (h *HttpBookHandler) GetById(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.Error{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	book, err := h.bookService.GetById(uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.Error{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	if book == nil {
		return c.Status(fiber.StatusNotFound).JSON(domain.Error{
			Code:    fiber.StatusNotFound,
			Message: "book not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(domain.Success{
		Code:    fiber.StatusOK,
		Message: "book fetched successfully",
		Data:    book,
	})
}

// Store used to store book
//
//	@Summary		Store book
//	@Description	Store book
//	@Tags			books
//	@Accept			json
//	@Produce		json
//	@Param			book	body		domain.BookStoreRequest	true	"book data"
//	@Success		201		{object}	domain.Success				"book detail"
//	@Failure		400		{object}	domain.Error				"Bad Request"
//	@Failure		500		{object}	domain.Error				"Internal Server Error"
//	@Router			/books [post]
//
// @Security Bearer
func (h *HttpBookHandler) Store(c *fiber.Ctx) error {
	bookReq := utilities.ExtractStructFromValidator[domain.BookStoreRequest](c)

	// Convert PublishedAt from string (dd-mm-yyyy) to time.Time
	publishedAt, err := time.Parse("02-01-2006", bookReq.PublishedAt)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.Error{
			Code:    fiber.StatusBadRequest,
			Message: "invalid published at format, should be dd-mm-yyyy",
		})
	}

	book := &domain.Book{
		Title:       bookReq.Title,
		Author:      bookReq.Author,
		Price:       bookReq.Price,
		Description: bookReq.Description,
		Pages:       bookReq.Pages,
		Isbn:        bookReq.Isbn,
		Language:    bookReq.Language,
		Stock:       bookReq.Stock,
		PublishedAt: publishedAt,
	}

	if err := h.bookService.Store(book); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.Error{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(domain.Success{
		Code:    fiber.StatusCreated,
		Message: "book created successfully",
		Data:    book,
	})
}

// Update used to update book
//
//	@Summary		Update book
//	@Description	Update book
//	@Tags			books
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int							true	"book ID"
//	@Param			book	body		domain.BookUpdateRequest	true	"book data"
//	@Success		200		{object}	domain.Success				"book detail"
//	@Failure		400		{object}	domain.Error				"Bad Request"
//	@Failure		404		{object}	domain.Error				"Not Found"
//	@Failure		500		{object}	domain.Error				"Internal Server Error"
//	@Router			/books/{id} [put]
//
// @Security Bearer
func (h *HttpBookHandler) Update(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.Error{
			Code:    fiber.StatusBadRequest,
			Message: "invalid book id",
		})
	}

	bookReq := utilities.ExtractStructFromValidator[domain.BookUpdateRequest](c)

	// Convert PublishedAt from string (dd-mm-yyyy) to time.Time
	publishedAt, err := time.Parse("02-01-2006", bookReq.PublishedAt)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.Error{
			Code:    fiber.StatusBadRequest,
			Message: "invalid published at format, should be dd-mm-yyyy",
		})
	}

	book := &domain.Book{
		Model:       gorm.Model{ID: uint(id)},
		Title:       bookReq.Title,
		Author:      bookReq.Author,
		Price:       bookReq.Price,
		Description: bookReq.Description,
		Pages:       bookReq.Pages,
		Isbn:        bookReq.Isbn,
		Language:    bookReq.Language,
		Stock:       bookReq.Stock,
		PublishedAt: publishedAt,
	}

	if err := h.bookService.Update(book); err != nil {
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
		Message: "book updated successfully",
		Data:    book,
	})
}

// Delete used to delete book
//
//	@Summary		Delete book
//	@Description	Delete book
//	@Tags			books
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int				true	"Book ID"
//	@Success		200	{object}	domain.Success	"Success delete book"
//	@Failure		400	{object}	domain.Error	"Bad Request"
//	@Failure		500	{object}	domain.Error	"Internal Server Error"
//	@Router			/books/{id} [delete]
//
// @Security Bearer
func (h *HttpBookHandler) Delete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.Error{
			Code:    fiber.StatusBadRequest,
			Message: "invalid book id",
		})
	}

	if err := h.bookService.Delete(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.Error{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(domain.Success{
		Code:    fiber.StatusOK,
		Message: "book deleted successfully",
	})
}
