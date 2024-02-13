package service

import (
	"context"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
	"github.com/reyhanyogs/e-wallet/domain"
	"github.com/reyhanyogs/e-wallet/internal/config"
)

type midTransService struct {
	config config.Midtrans
	env    midtrans.EnvironmentType
}

func NewMidtrans(config *config.Config) domain.MidTransService {
	env := midtrans.Sandbox
	if config.Midtrans.IsProd {
		env = midtrans.Production
	}

	return &midTransService{
		config: config.Midtrans,
		env:    env,
	}
}

func (s *midTransService) GenerateSnapURL(ctx context.Context, t *domain.TopUp) error {
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  t.ID,
			GrossAmt: t.Amount,
		},
	}

	var client snap.Client
	client.New(s.config.Key, s.env)

	snapResp, err := client.CreateTransaction(req)
	if err != nil {
		return err
	}
	t.SnapURL = snapResp.RedirectURL
	return nil
}

func (s *midTransService) VerifyPayment(ctx context.Context, orderId string) (bool, error) {
	var client coreapi.Client
	client.New(s.config.Key, s.env)

	// Check transaction to Midtrans with param order_id
	transactionStatusResp, err := client.CheckTransaction(orderId)
	if err != nil {
		return false, err
	} else {
		if transactionStatusResp != nil {
			if transactionStatusResp.TransactionStatus == "settlement" {
				return true, nil
			} else if transactionStatusResp.TransactionStatus == "capture" {
				if transactionStatusResp.FraudStatus == "accept" {
					return true, nil
				}
			} else if transactionStatusResp.TransactionStatus == "pending" {
				// set db status to pending
			}
		}
	}
	return false, nil
}
