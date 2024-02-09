package sse

import (
	"bufio"
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/reyhanyogs/e-wallet/dto"
)

type notificationSse struct {
	hub *dto.Hub
}

func NewNotification(app *fiber.App, authMid fiber.Handler, hub *dto.Hub) {
	h := notificationSse{
		hub: hub,
	}
	app.Get("sse/notification-stream", authMid, h.StreamNotification)
}

func (h *notificationSse) StreamNotification(ctx *fiber.Ctx) error {
	ctx.Set("Content-Type", "text/event-stream")

	user := ctx.Locals("x-user").(dto.UserData)
	h.hub.NotificationChannel[user.ID] = make(chan dto.NotificationData)

	ctx.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		event := fmt.Sprintf("event: %s\n"+"data: \n\n", "initial")
		_, _ = fmt.Fprintf(w, event)
		_ = w.Flush()
		for notification := range h.hub.NotificationChannel[user.ID] {
			data, _ := json.Marshal(notification)
			event = fmt.Sprintf("event: %s\n"+"data: %s\n\n", "notification-updated", data)
			_, _ = fmt.Fprintf(w, event)
			_ = w.Flush()
		}
	})
	return nil
}
