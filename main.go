package main

import "book-store/internal/infrastructure"

// @title			Book Store API Documentation
// @version		1.0
// @host			localhost:3000
// @BasePath		/api
// @schemes		http https
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	infrastructure.Run()
}
