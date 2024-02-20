package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/reyhanyogs/e-wallet/domain"
	"github.com/reyhanyogs/e-wallet/dto"
	"github.com/reyhanyogs/e-wallet/internal/util"
)

type topUpApi struct {
	topUpService domain.TopUpService
}

func NewTopUp(app *fiber.App, authMid fiber.Handler, topUpService domain.TopUpService) {
	h := &topUpApi{
		topUpService: topUpService,
	}

	app.Post("/topup/initialize", authMid, h.InitializeTopUp)
}

func (handler *topUpApi) InitializeTopUp(ctx *fiber.Ctx) error {
	var req dto.TopUpReq
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(400).JSON(dto.Response{
			Message: err.Error(),
		})
	}

	user := ctx.Locals("x-user").(dto.UserData)
	req.UserID = user.ID

	res, err := handler.topUpService.InitializeTopUp(ctx.Context(), req)
	if err != nil {
		return ctx.Status(util.GetHttpStatus(err)).JSON(dto.Response{
			Message: err.Error(),
		})
	}

	return ctx.Status(200).JSON(res)
}
