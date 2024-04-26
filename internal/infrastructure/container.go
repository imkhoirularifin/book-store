package infrastructure

import (
	"gramedia-service/internal/auth"
	"gramedia-service/internal/book"
	"gramedia-service/internal/config"
	"gramedia-service/internal/customer"
	"gramedia-service/internal/domain"
	"gramedia-service/internal/middleware/jwt"
	"gramedia-service/internal/role"
	"gramedia-service/internal/transaction"
	"gramedia-service/internal/user"
	"gramedia-service/internal/utilities"
	"gramedia-service/pkg/xlogger"

	"github.com/caarlos0/env/v10"
	_ "github.com/joho/godotenv/autoload"
)

var (
	cfg config.Config

	customerRepository    domain.CustomerRepository
	bookRepository        domain.BookRepository
	roleRepository        domain.RoleRepository
	userRepository        domain.UserRepository
	transactionRepository domain.TransactionRepository

	jwtService         utilities.JwtTokenService
	customerService    domain.CustomerService
	bookService        domain.BookService
	roleService        domain.RoleService
	userService        domain.UserService
	authService        domain.AuthService
	transactionService domain.TransactionService

	authMiddleware jwt.AuthMiddleware
)

func init() {
	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}

	xlogger.Setup(cfg)
	dbSetup()

	customerRepository = customer.NewMysqlCustomerRepository(db)
	bookRepository = book.NewMysqlBookRepository(db)
	roleRepository = role.NewMysqlRoleRepository(db)
	userRepository = user.NewMysqlUserRepository(db)
	transactionRepository = transaction.NewMysqlTransactionRepository(db)

	jwtService = utilities.NewJwtTokenService(cfg)
	customerService = customer.NewCustomerService(customerRepository)
	bookService = book.NewBookService(bookRepository)
	roleService = role.NewRoleService(roleRepository)
	userService = user.NewUserService(userRepository)
	authService = auth.NewAuthService(userRepository, jwtService)
	transactionService = transaction.NewTransactionService(transactionRepository, bookRepository)

	authMiddleware = jwt.NewAuthMiddleware(jwtService)
}
