package infrastructure

import (
	"book-store/internal/auth"
	"book-store/internal/book"
	"book-store/internal/config"
	"book-store/internal/customer"
	"book-store/internal/domain"
	"book-store/internal/middleware/jwt"
	"book-store/internal/role"
	"book-store/internal/transaction"
	"book-store/internal/user"
	"book-store/internal/utilities"
	"book-store/pkg/xlogger"

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
