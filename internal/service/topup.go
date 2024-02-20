package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/reyhanyogs/e-wallet/domain"
	"github.com/reyhanyogs/e-wallet/dto"
	"github.com/reyhanyogs/e-wallet/internal/component"
)

type topUpService struct {
	notificationService   domain.NotificationService
	midTransService       domain.MidTransService
	topUpRepository       domain.TopUpRepository
	accountRepository     domain.AccountRepository
	transactionRepository domain.TransactionRepository
}

func NewTopUp(
	notificationService domain.NotificationService,
	midtransService domain.MidTransService,
	topUpRepository domain.TopUpRepository,
	accountRepository domain.AccountRepository,
	transactionRepository domain.TransactionRepository,
) domain.TopUpService {
	return &topUpService{
		notificationService:   notificationService,
		midTransService:       midtransService,
		topUpRepository:       topUpRepository,
		accountRepository:     accountRepository,
		transactionRepository: transactionRepository,
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
		component.Log.Errorf("InitializeTopUp(GenerateSnapURL): user_id = %d: err = %s", req.UserID, err.Error())
		return dto.TopUpRes{}, err
	}

	err = s.topUpRepository.Insert(ctx, &topUp)
	if err != nil {
		component.Log.Errorf("InitializeTopUp(Insert): user_id = %d: err = %s", req.UserID, err.Error())
		return dto.TopUpRes{}, err
	}

	return dto.TopUpRes{
		SnapURL: topUp.SnapURL,
	}, nil
}

func (s *topUpService) ConfirmedTopUp(ctx context.Context, id string) error {
	topUp, err := s.topUpRepository.FindById(ctx, id)
	if err != nil {
		component.Log.Errorf("ConfirmedTopUp(FindById): user_id = %s: err = %s", id, err.Error())
		return err
	}
	if topUp == (domain.TopUp{}) {
		return domain.ErrTopUpNotFound
	}

	account, err := s.accountRepository.FindByUserID(ctx, topUp.UserID)
	if err != nil {
		component.Log.Errorf("ConfirmedTopUp(FindByUserID): user_id = %d: err = %s", topUp.UserID, err.Error())
		return err
	}
	if account == (domain.Account{}) {
		return domain.ErrAccountNotFound
	}

	err = s.transactionRepository.Insert(ctx, &domain.Transaction{
		AccountId:           account.ID,
		SofNumber:           "00",
		DofNumber:           account.AccountNumber,
		TransactionType:     "C",
		Amount:              float64(topUp.Amount),
		TransactionDatetime: time.Now(),
	})
	if err != nil {
		component.Log.Errorf("ConfirmedTopUp(Insert): user_id = %d: err = %s", topUp.UserID, err.Error())
		return err
	}

	account.Balance += float64(topUp.Amount)
	err = s.accountRepository.Update(ctx, &account)
	if err != nil {
		component.Log.Errorf("ConfirmedTopUp(Update): user_id = %d: err = %s", topUp.UserID, err.Error())
		return err
	}

	data := map[string]string{
		"amount": fmt.Sprintf("%d", topUp.Amount),
	}
	_ = s.notificationService.Insert(ctx, account.UserId, "TOPUP_SUCCESS", data)

	return nil
}
