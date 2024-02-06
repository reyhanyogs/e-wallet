package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/reyhanyogs/e-wallet/domain"
	"github.com/reyhanyogs/e-wallet/dto"
	"github.com/reyhanyogs/e-wallet/internal/util"
)

type transferApi struct {
	transactionService domain.TransactionService
}

func NewTransfer(app *fiber.App, authMid fiber.Handler, transactionService domain.TransactionService) {
	h := &transferApi{
		transactionService: transactionService,
	}

	app.Post("transfer/inquiry", authMid, h.TransferInquiry)
	app.Post("transfer/execute", authMid, h.TransferExecute)
}

func (h *transferApi) TransferInquiry(ctx *fiber.Ctx) error {
	var req dto.TransferInquiryReq
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(400).JSON(dto.Response{
			Message: "invalid body",
		})
	}

	inquiry, err := h.transactionService.TransferInquiry(ctx.Context(), req)
	if err != nil {
		return ctx.Status(400).JSON(dto.Response{
			Message: err.Error(),
		})
	}

	return ctx.Status(200).JSON(inquiry)
}

func (h *transferApi) TransferExecute(ctx *fiber.Ctx) error {
	var req dto.TransferExecuteReq
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(util.GetHttpStatus(err)).JSON(dto.Response{
			Message: err.Error(),
		})
	}

	err := h.transactionService.TransferExecute(ctx.Context(), req)
	if err != nil {
		return ctx.Status(400).JSON(dto.Response{
			Message: "invalid body",
		})
	}

	return ctx.SendStatus(200)
}
