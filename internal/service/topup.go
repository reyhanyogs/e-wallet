package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/reyhanyogs/e-wallet/domain"
	"github.com/reyhanyogs/e-wallet/dto"
)

type topUpService struct {
	notificationService domain.NotificationService
	midTransService     domain.MidTransService
	topUpRepository     domain.TopUpRepository
	accountRepository   domain.AccountRepository
}

func NewTopUp(
	notificationService domain.NotificationService,
	midtransService domain.MidTransService,
	topUpRepository domain.TopUpRepository,
	accountRepository domain.AccountRepository,
) domain.TopUpService {
	return &topUpService{
		notificationService: notificationService,
		midTransService:     midtransService,
		topUpRepository:     topUpRepository,
		accountRepository:   accountRepository,
	}
}

func (s *topUpService) InitializeTopUp(ctx context.Context, req dto.TopUpReq) (dto.TopUpRes, error) {
	topUp := domain.TopUp{
		ID:     uuid.NewString(),
		UserID: req.UserID,
		Status: 0,
		Amount: req.Amount,
	}
	err := s.midTransService.GenerateSnapURL(ctx, &topUp)
	if err != nil {
		return dto.TopUpRes{}, err
	}

	err = s.topUpRepository.Insert(ctx, &topUp)
	if err != nil {
		return dto.TopUpRes{}, err
	}

	return dto.TopUpRes{
		SnapURL: topUp.SnapURL,
	}, nil
}

func (s *topUpService) ConfirmedTopUp(ctx context.Context, id string) error {
	topUp, err := s.topUpRepository.FindById(ctx, id)
	if err != nil {
		return err
	}
	if topUp == (domain.TopUp{}) {
		return domain.ErrTopUpNotFound
	}

	account, err := s.accountRepository.FindByUserID(ctx, topUp.UserID)
	if err != nil {
		return err
	}
	if account == (domain.Account{}) {
		return domain.ErrAccountNotFound
	}

	account.Balance += float64(topUp.Amount)
	err = s.accountRepository.Update(ctx, &account)
	if err != nil {
		return err
	}

	data := map[string]string{
		"amount": fmt.Sprintf("%.2f", topUp.Amount),
	}
	_ = s.notificationService.Insert(ctx, account.UserId, "TOPUP_SUCCESS", data)

	return nil
}
