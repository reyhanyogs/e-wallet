package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/reyhanyogs/e-wallet/domain"
	"github.com/reyhanyogs/e-wallet/dto"
	"github.com/reyhanyogs/e-wallet/internal/util"
)

type authApi struct {
	userService domain.UserService
}

func NewAuth(app *fiber.App, userService domain.UserService, authMid fiber.Handler) {
	h := authApi{
		userService: userService,
	}

	app.Post("token/generate", h.GenerateToken)
	app.Post("token/validate", authMid, h.ValidateToken)
}

func (handler *authApi) GenerateToken(ctx *fiber.Ctx) error {
	var req dto.AuthReq
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.SendStatus(400)
	}

	token, err := handler.userService.Authenticate(ctx.Context(), req)
	if err != nil {
		return ctx.SendStatus(util.GetHttpStatus(err))
	}

	return ctx.Status(200).JSON(token)
}

func (handler *authApi) ValidateToken(ctx *fiber.Ctx) error {
	user := ctx.Locals("x-user")
	return ctx.Status(200).JSON(user)
}
