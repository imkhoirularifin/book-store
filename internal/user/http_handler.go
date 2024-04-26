package user

import (
	"errors"
	"gramedia-service/internal/domain"
	"gramedia-service/internal/middleware/jwt"
	"gramedia-service/internal/middleware/validation"
	"gramedia-service/internal/utilities"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type HttpUserHandler struct {
	userSvc        domain.UserService
	authMiddleware jwt.AuthMiddleware
}

func NewHttpHandler(r fiber.Router, userSvc domain.UserService, authMiddleware jwt.AuthMiddleware) {
	handler := &HttpUserHandler{
		userSvc:        userSvc,
		authMiddleware: authMiddleware,
	}

	r.Get("/", handler.Fetch)
	r.Get("/:id", handler.GetById)
	r.Post("/", authMiddleware.RequireRole("admin"), validation.New[domain.UserStoreRequest](), handler.Store)
	r.Put("/:id", authMiddleware.RequireRole("admin"), validation.New[domain.UserUpdateRequest](), handler.Update)
	r.Delete("/:id", authMiddleware.RequireRole("admin"), handler.Delete)
}

// Fetch used to get list of user
//
//	@Summary		Get list of user
//	@Description	Get list of users
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			page	query		int				false	"Page number (default 1)"
//	@Param			size	query		int				false	"Size of page (default 10)"
//	@Param			q		query		string			false	"Search query"
//	@Header			200		{string}	X-Cursor		"Next page"
//	@Header			200		{string}	X-Total-Count	"Total item"
//	@Header			200		{string}	X-Max-Page		"Max page"
//	@Success		200		{array}		domain.Success	"List of users"
//	@Failure		400		{object}	domain.Error	"Bad Request"
//	@Failure		500		{object}	domain.Error	"Internal Server Error"
//	@Router			/users [get]
func (h *HttpUserHandler) Fetch(c *fiber.Ctx) error {
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

	filter := &domain.User{Name: query}
	users, nextPage, err := h.userSvc.Fetch(page, size, filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.Error{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	if users == nil {
		return c.Status(fiber.StatusNotFound).JSON(domain.Error{
			Code:    fiber.StatusNotFound,
			Message: "users not found",
		})
	}

	totalItem, err := h.userSvc.Count(filter)
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
		Data:    users,
	})
}

// GetByID used to get user by id
//
//	@Summary		Get user by id
//	@Description	Get user by id
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int				true	"user ID"
//	@Success		200	{object}	domain.Success	"user detail"
//	@Failure		400	{object}	domain.Error	"Bad Request"
//	@Failure		404	{object}	domain.Error	"Not Found"
//	@Failure		500	{object}	domain.Error	"Internal Server Error"
//	@Router			/users/{id} [get]
func (h *HttpUserHandler) GetById(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.Error{
			Code:    fiber.StatusBadRequest,
			Message: "invalid user id",
		})
	}

	user, err := h.userSvc.GetById(uint(id))
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
		Data:    user,
	})
}

// Store used to store user
//
//	@Summary		Store user
//	@Description	Store user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			user	body		domain.UserStoreRequest	true	"user data"
//	@Success		201		{object}	domain.Success				"user detail"
//	@Failure		400		{object}	domain.Error				"Bad Request"
//	@Failure		500		{object}	domain.Error				"Internal Server Error"
//	@Router			/users [post]
//
// @Security Bearer
func (h *HttpUserHandler) Store(c *fiber.Ctx) error {
	userReq := utilities.ExtractStructFromValidator[domain.UserStoreRequest](c)

	user := &domain.User{
		Name:     userReq.Name,
		Email:    userReq.Email,
		Password: userReq.Password,
		RoleId:   userReq.RoleId,
	}

	if err := h.userSvc.Store(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.Error{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(domain.Success{
		Code:    fiber.StatusCreated,
		Message: "success",
		Data:    user,
	})
}

// Update used to update user
//
//	@Summary		Update user
//	@Description	Update user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int							true	"user ID"
//	@Param			user	body		domain.UserUpdateRequest	true	"user data"
//	@Success		200		{object}	domain.Success				"user detail"
//	@Failure		400		{object}	domain.Error				"Bad Request"
//	@Failure		404		{object}	domain.Error				"Not Found"
//	@Failure		500		{object}	domain.Error				"Internal Server Error"
//	@Router			/users/{id} [put]
//
// @Security Bearer
func (h *HttpUserHandler) Update(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.Error{
			Code:    fiber.StatusBadRequest,
			Message: "invalid user id",
		})
	}

	userReq := utilities.ExtractStructFromValidator[domain.UserUpdateRequest](c)

	user := &domain.User{
		Model:    gorm.Model{ID: uint(id)},
		Name:     userReq.Name,
		Email:    userReq.Email,
		Password: userReq.Password,
		RoleId:   userReq.RoleId,
	}

	if err := h.userSvc.Update(user); err != nil {
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
		Data:    user,
	})
}

// Delete used to delete user
//
//	@Summary		Delete User
//	@Description	Delete User
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int				true	"User ID"
//	@Success		200	{object}	domain.Success	"Success delete user"
//	@Failure		400	{object}	domain.Error	"Bad Request"
//	@Failure		500	{object}	domain.Error	"Internal Server Error"
//	@Router			/users/{id} [delete]
//
// @Security Bearer
func (h *HttpUserHandler) Delete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.Error{
			Code:    fiber.StatusBadRequest,
			Message: "invalid user id",
		})
	}

	if err := h.userSvc.Delete(uint(id)); err != nil {
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
