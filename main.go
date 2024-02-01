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
	userService := service.NewUser(userRepository, cacheConnection)

	authMiddleware := middleware.Authenticate(userService)

	app := fiber.New()

	api.NewAuth(app, userService, authMiddleware)

	_ = app.Listen(config.Server.Host + ":" + config.Server.Port)
}
