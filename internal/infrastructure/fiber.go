package infrastructure

import (
	"fmt"
	"gramedia-service/internal/auth"
	"gramedia-service/internal/book"
	"gramedia-service/internal/customer"
	"gramedia-service/internal/docs"
	"gramedia-service/internal/role"
	"gramedia-service/internal/transaction"
	"gramedia-service/internal/user"
	"gramedia-service/pkg/xlogger"

	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/etag"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func Run() {
	logger := xlogger.Logger

	app := fiber.New(fiber.Config{
		ProxyHeader:           cfg.ProxyHeader,
		DisableStartupMessage: true,
		ErrorHandler:          defaultErrorHandler,
		AppName:               "gramedia-service",
	})

	app.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: logger,
		Fields: cfg.LogFields,
	}))
	app.Use(recover2.New())
	app.Use(etag.New())
	app.Use(requestid.New())

	api := app.Group("api")
	docs.NewHttpHandler(api.Group("/docs"))
	customer.NewHttpHandler(api.Group("/customers"), customerService, authMiddleware)
	book.NewHttpHandler(api.Group("/books"), bookService, authMiddleware)
	role.NewHttpHandler(api.Group("/roles"), roleService)
	user.NewHttpHandler(api.Group("/users"), userService, authMiddleware)
	auth.NewHttpHandler(api.Group("/auth"), authService)
	transaction.NewHttpHandler(api.Group("/transactions"), transactionService, authMiddleware)

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	logger.Info().Msgf("Server is running on address: %s", addr)
	if err := app.Listen(addr); err != nil {
		logger.Fatal().Err(err).Msg("Server failed to start")
	}
}
