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
	client         snap.Client
	midTransConfig config.Midtrans
}

func NewMidtrans(config *config.Config) domain.MidTransService {
	var client snap.Client
	env := midtrans.Sandbox
	if config.Midtrans.IsProd {
		env = midtrans.Production
	}
	client.New(config.Midtrans.Key, env)

	return &midTransService{
		client:         client,
		midTransConfig: config.Midtrans,
	}
}

func (s *midTransService) GenerateSnapURL(ctx context.Context, t *domain.TopUp) error {
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  t.ID,
			GrossAmt: t.Amount,
		},
	}

	snapResp, err := s.client.CreateTransaction(req)
	if err != nil {
		return err
	}
	t.SnapURL = snapResp.RedirectURL
	return nil
}

func (s *midTransService) VerifyPayment(ctx context.Context, data map[string]interface{}) (bool, error) {
	var client coreapi.Client
	env := midtrans.Sandbox
	if s.midTransConfig.IsProd {
		env = midtrans.Production
	}
	client.New(s.midTransConfig.Key, env)

	// Check if order_id exist in payload
	orderId, exists := data["order_id"].(string)
	if !exists {
		return false, domain.ErrInvalidPayload
	}

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
			}
		}
	}
	return false, nil
}
