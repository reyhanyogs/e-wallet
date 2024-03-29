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
	loginLogRepository := repository.NewLoginLog(dbConnection)

	queueService := service.NewQueue(config)
	emailService := service.NewEmail(queueService)
	userService := service.NewUser(userRepository, cacheConnection, emailService, factorRepository, accountRepository)
	notificationService := service.NewNotification(notificationRepository, templateRepository, hub)
	transactionService := service.NewTransaction(accountRepository, transactionRepository, cacheConnection, notificationService)
	midtransService := service.NewMidtrans(config)
	topUpService := service.NewTopUp(notificationService, midtransService, topUpRepository, accountRepository, transactionRepository)
	factorService := service.NewFactor(factorRepository)
	ipCheckerService := service.NewIpChecker()
	fdsService := service.NewFds(ipCheckerService, loginLogRepository)

	authMiddleware := middleware.Authenticate(userService)

	app := fiber.New()

	api.NewAuth(app, userService, authMiddleware, fdsService)
	api.NewTransfer(app, authMiddleware, transactionService, factorService)
	api.NewNotification(app, authMiddleware, notificationService)
	api.NewTopUp(app, authMiddleware, topUpService)
	api.NewMidtrans(app, midtransService, topUpService)

	sse.NewNotification(app, authMiddleware, hub)

	component.Log.Info("Starting application")
	_ = app.Listen(config.Server.Host + ":" + config.Server.Port)
}
