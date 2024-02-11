package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/reyhanyogs/e-wallet/domain"
	"github.com/reyhanyogs/e-wallet/dto"
	"github.com/reyhanyogs/e-wallet/internal/util"
)

type midtransApi struct {
	midtransService domain.MidTransService
	topUpService    domain.TopUpService
}

func NewMidtrans(app *fiber.App, midtransService domain.MidTransService, topUpService domain.TopUpService) {
	h := &midtransApi{
		midtransService: midtransService,
		topUpService:    topUpService,
	}
	app.Post("/midtrans/payment-callback", h.paymentHandlerNotification)
}

func (h *midtransApi) paymentHandlerNotification(ctx *fiber.Ctx) error {
	var payload map[string]interface{}
	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(400).JSON(dto.Response{
			Message: err.Error(),
		})
	}

	orderId, exists := payload["order_id"].(string)
	if !exists {
		return ctx.Status(400).JSON(dto.Response{
			Message: domain.ErrInvalidPayload.Error(),
		})
	}

	success, err := h.midtransService.VerifyPayment(ctx.Context(), payload)
	if err != nil {
		return ctx.Status(util.GetHttpStatus(err)).JSON(dto.Response{
			Message: err.Error(),
		})
	}
	if !success {
		return ctx.SendStatus(400)
	}

	err = h.topUpService.ConfirmedTopUp(ctx.Context(), orderId)
	if err != nil {
		return ctx.Status(util.GetHttpStatus(err)).JSON(dto.Response{
			Message: err.Error(),
		})
	}

	return ctx.SendStatus(200)
}
