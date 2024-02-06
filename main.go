package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/reyhanyogs/e-wallet/internal/api"
	"github.com/reyhanyogs/e-wallet/internal/component"
	"github.com/reyhanyogs/e-wallet/internal/config"
	"github.com/reyhanyogs/e-wallet/internal/middleware"
	"github.com/reyhanyogs/e-wallet/internal/repository"
	"github.com/reyhanyogs/e-wallet/internal/service"
)

func main() {
	config := config.Get()
	dbConnection := component.GetDatabaseConn(config)
	cacheConnection := component.GetCacheConnection()

	userRepository := repository.NewUser(dbConnection)
	accountRepository := repository.NewAccount(dbConnection)
	transactionRepository := repository.NewTransaction(dbConnection)

	emailService := service.NewEmail(config)
	userService := service.NewUser(userRepository, cacheConnection, emailService)
	transactionService := service.NewTransaction(accountRepository, transactionRepository, cacheConnection)

	authMiddleware := middleware.Authenticate(userService)

	app := fiber.New()

	api.NewAuth(app, userService, authMiddleware)
	api.NewTransfer(app, authMiddleware, transactionService)

	_ = app.Listen(config.Server.Host + ":" + config.Server.Port)
}
