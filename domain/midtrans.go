package domain

import "context"

type MidTransService interface {
	GenerateSnapURL(ctx context.Context, t *TopUp) error
	VerifyPayment(ctx context.Context, orderId string) (bool, error)
}
