package infrastructure

import (
	"fmt"
	"book-store/internal/domain"
	"book-store/internal/utilities"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var db *gorm.DB

func dbSetup() {
	var err error
	l := gormLogger.Default.LogMode(gormLogger.Silent)
	if cfg.Database.Driver == "mysql" {
		db, err = gorm.Open(mysql.New(mysql.Config{
			DSN: cfg.Database.DSN,
		}), &gorm.Config{
			Logger: l,
		})
	}

	if err != nil {
		panic(err)
	}

	if cfg.IsDevelopment {
		fmt.Println("Development Mode")
		if err := db.AutoMigrate(
			&domain.User{},
			&domain.Book{},
			&domain.Customer{},
			&domain.Role{},
			&domain.Transaction{},
			&domain.TransactionDetail{},
		); err != nil {
			panic(err)
		}
	}

	// create initial roles
	var roleCount int64
	if err := db.Model(&domain.Role{}).Count(&roleCount).Error; err != nil {
		panic(err)
	}

	if roleCount == 0 {
		roles := []domain.Role{
			{Name: "admin", Model: gorm.Model{ID: 1}},
			{Name: "employee", Model: gorm.Model{ID: 2}},
		}

		if err := db.Create(&roles).Error; err != nil {
			panic(err)
		}
	}

	// create initial admin
	var userCount int64
	if err := db.Model(&domain.User{}).Count(&userCount).Error; err != nil {
		panic(err)
	}

	if userCount == 0 {
		userPassword, err := utilities.HashPassword("admin")
		if err != nil {
			panic(err)
		}
		user := domain.User{
			Name:     "admin",
			Email:    "admin@mail.com",
			Password: userPassword,
			RoleId:   1,
		}

		if err := db.Create(&user).Error; err != nil {
			panic(err)
		}
	}
}
