package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/reyhanyogs/e-wallet/dto"
	"github.com/reyhanyogs/e-wallet/internal/api"
	"github.com/reyhanyogs/e-wallet/internal/component"
	"github.com/reyhanyogs/e-wallet/internal/config"
	"github.com/reyhanyogs/e-wallet/internal/middleware"
	"github.com/reyhanyogs/e-wallet/internal/repository"
	"github.com/reyhanyogs/e-wallet/internal/service"
	"github.com/reyhanyogs/e-wallet/internal/sse"
)

func main() {
	config := config.Get()
	dbConnection := component.GetDatabaseConn(config)
	cacheConnection := repository.NewRedisClient(config)

	hub := &dto.Hub{
		NotificationChannel: make(map[int64]chan dto.NotificationData),
	}

	userRepository := repository.NewUser(dbConnection)
	accountRepository := repository.NewAccount(dbConnection)
	transactionRepository := repository.NewTransaction(dbConnection)
	notificationRepository := repository.NewNotification(dbConnection)
	templateRepository := repository.NewTemplate(dbConnection)
	topUpRepository := repository.NewTopUp(dbConnection)
	factorRepository := repository.NewFactor(dbConnection)

	emailService := service.NewEmail(config)
	userService := service.NewUser(userRepository, cacheConnection, emailService, factorRepository)
	notificationService := service.NewNotification(notificationRepository, templateRepository, hub)
	transactionService := service.NewTransaction(accountRepository, transactionRepository, cacheConnection, notificationService)
	midtransService := service.NewMidtrans(config)
	topUpService := service.NewTopUp(notificationService, midtransService, topUpRepository, accountRepository, transactionRepository)
	factorService := service.NewFactor(factorRepository)

	authMiddleware := middleware.Authenticate(userService)

	app := fiber.New()

	api.NewAuth(app, userService, authMiddleware)
	api.NewTransfer(app, authMiddleware, transactionService, factorService)
	api.NewNotification(app, authMiddleware, notificationService)
	api.NewTopUp(app, authMiddleware, topUpService)
	api.NewMidtrans(app, midtransService, topUpService)

	sse.NewNotification(app, authMiddleware, hub)

	_ = app.Listen(config.Server.Host + ":" + config.Server.Port)
}
