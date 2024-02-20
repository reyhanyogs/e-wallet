package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/reyhanyogs/e-wallet/domain"
	"github.com/reyhanyogs/e-wallet/dto"
	"github.com/reyhanyogs/e-wallet/internal/util"
)

type authApi struct {
	userService domain.UserService
	fdsService  domain.FdsService
}

func NewAuth(app *fiber.App, userService domain.UserService, authMid fiber.Handler, fdsService domain.FdsService) {
	h := authApi{
		userService: userService,
		fdsService:  fdsService,
	}

	app.Post("token/generate", h.GenerateToken)
	app.Post("token/validate", authMid, h.ValidateToken)
	app.Post("user/register", h.RegisterUser)
	app.Post("user/validate-otp", h.ValidateOTP)
}

func (handler *authApi) GenerateToken(ctx *fiber.Ctx) error {
	var req dto.AuthReq
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(400).JSON(dto.Response{
			Message: err.Error(),
		})
	}

	token, err := handler.userService.Authenticate(ctx.Context(), req)
	if err != nil {
		return ctx.Status(util.GetHttpStatus(err)).JSON(dto.Response{
			Message: err.Error(),
		})
	}

	if !handler.fdsService.IsAuthorized(ctx.Context(), ctx.Get("X-Forwarded-For"), token.UserId) {
		return ctx.SendStatus(401)
	}

	return ctx.Status(200).JSON(token)
}

func (handler *authApi) ValidateToken(ctx *fiber.Ctx) error {
	user := ctx.Locals("x-user")
	return ctx.Status(200).JSON(user)
}

func (handler *authApi) RegisterUser(ctx *fiber.Ctx) error {
	var req dto.UserRegisterReq
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(400).JSON(dto.Response{
			Message: err.Error(),
		})
	}

	res, err := handler.userService.Register(ctx.Context(), req)
	if err != nil {
		return ctx.Status(util.GetHttpStatus(err)).JSON(dto.Response{
			Message: err.Error(),
		})
	}
	return ctx.Status(200).JSON(res)
}

func (handler *authApi) ValidateOTP(ctx *fiber.Ctx) error {
	var req dto.ValidateOtpReq
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(400).JSON(dto.Response{
			Message: err.Error(),
		})
	}

	err := handler.userService.ValidateOTP(ctx.Context(), req)
	if err != nil {
		return ctx.Status(util.GetHttpStatus(err)).JSON(dto.Response{
			Message: err.Error(),
		})
	}
	return ctx.SendStatus(200)
}
