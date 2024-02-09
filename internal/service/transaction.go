package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/reyhanyogs/e-wallet/domain"
	"github.com/reyhanyogs/e-wallet/dto"
	"github.com/reyhanyogs/e-wallet/internal/util"
)

type transactionService struct {
	accountRepository      domain.AccountRepository
	transactionRepository  domain.TransactionRepository
	cacheRepository        domain.CacheRepository
	notificationRepository domain.NotificationRepository
	hub                    *dto.Hub
}

func NewTransaction(
	accountRepository domain.AccountRepository,
	transactionRepository domain.TransactionRepository,
	cacheRepository domain.CacheRepository,
	notificationRepository domain.NotificationRepository,
	hub *dto.Hub,
) domain.TransactionService {
	return &transactionService{
		accountRepository:      accountRepository,
		transactionRepository:  transactionRepository,
		cacheRepository:        cacheRepository,
		notificationRepository: notificationRepository,
		hub:                    hub,
	}
}

func (s *transactionService) TransferInquiry(ctx context.Context, req dto.TransferInquiryReq) (dto.TransferInquiryRes, error) {
	user := ctx.Value("x-user").(dto.UserData)
	myAccount, err := s.accountRepository.FindByUserID(ctx, user.ID)
	if err != nil {
		return dto.TransferInquiryRes{}, err
	}

	if myAccount == (domain.Account{}) {
		return dto.TransferInquiryRes{}, domain.ErrAccountNotFound
	}

	dofAccount, err := s.accountRepository.FindByAccountNumber(ctx, req.AccountNumber)
	if err != nil {
		return dto.TransferInquiryRes{}, err
	}

	if dofAccount == (domain.Account{}) {
		return dto.TransferInquiryRes{}, domain.ErrAccountNotFound
	}

	if myAccount.Balance < req.Amount {
		return dto.TransferInquiryRes{}, domain.ErrInsufficientBalance
	}

	inquiryKey := util.GenerateRandomString(32)
	jsonData, _ := json.Marshal(req)
	_ = s.cacheRepository.Set(inquiryKey, jsonData)

	return dto.TransferInquiryRes{
		InquiryKey: inquiryKey,
	}, nil
}

func (s *transactionService) TransferExecute(ctx context.Context, req dto.TransferExecuteReq) error {
	data, err := s.cacheRepository.Get(req.InquiryKey)
	if err != nil {
		return domain.ErrInquiryNotFound
	}

	var reqInq dto.TransferInquiryReq
	_ = json.Unmarshal(data, &reqInq)
	if reqInq == (dto.TransferInquiryReq{}) {
		return domain.ErrInquiryNotFound
	}

	user := ctx.Value("x-user").(dto.UserData)
	myAccount, err := s.accountRepository.FindByUserID(ctx, user.ID)
	if err != nil {
		return domain.ErrAccountNotFound
	}

	dofAccount, err := s.accountRepository.FindByAccountNumber(ctx, reqInq.AccountNumber)
	if err != nil {
		return err
	}

	debitTransaction := domain.Transaction{
		Amount:              reqInq.Amount,
		AccountId:           myAccount.ID,
		SofNumber:           myAccount.AccountNumber,
		DofNumber:           dofAccount.AccountNumber,
		TransactionType:     "D",
		TransactionDatetime: time.Now(),
	}
	err = s.transactionRepository.Insert(ctx, &debitTransaction)
	if err != nil {
		return err
	}

	creditTransaction := domain.Transaction{
		Amount:              reqInq.Amount,
		AccountId:           dofAccount.ID,
		SofNumber:           myAccount.AccountNumber,
		DofNumber:           dofAccount.AccountNumber,
		TransactionType:     "C",
		TransactionDatetime: time.Now(),
	}
	err = s.transactionRepository.Insert(ctx, &creditTransaction)
	if err != nil {
		return err
	}

	myAccount.Balance -= reqInq.Amount
	err = s.accountRepository.Update(ctx, &myAccount)
	if err != nil {
		return err
	}

	dofAccount.Balance += reqInq.Amount
	err = s.accountRepository.Update(ctx, &dofAccount)
	if err != nil {
		return err
	}

	go s.notificationAfterTransfer(myAccount, dofAccount, reqInq.Amount)
	return nil
}

func (s *transactionService) notificationAfterTransfer(sofAccount domain.Account, dofAccount domain.Account, amount float64) {
	notificationSender := domain.Notification{
		UserID:    sofAccount.UserId,
		Title:     "Transfer berhasil",
		Body:      fmt.Sprintf("Transfer senilai %.2f berhasil", amount),
		IsRead:    0,
		Status:    1,
		CreatedAt: time.Now(),
	}
	notificationReceiver := domain.Notification{
		UserID:    dofAccount.UserId,
		Title:     "Dana diterima",
		Body:      fmt.Sprintf("Dana diterima senilai %.2f", amount),
		IsRead:    0,
		Status:    1,
		CreatedAt: time.Now(),
	}

	_ = s.notificationRepository.Insert(context.Background(), &notificationSender)
	if channel, ok := s.hub.NotificationChannel[sofAccount.UserId]; ok {
		channel <- dto.NotificationData{
			ID:        notificationSender.ID,
			Title:     notificationSender.Title,
			Body:      notificationSender.Body,
			Status:    notificationSender.Status,
			IsRead:    notificationSender.IsRead,
			CreatedAt: notificationSender.CreatedAt,
		}
	}

	_ = s.notificationRepository.Insert(context.Background(), &notificationReceiver)
	if channel, ok := s.hub.NotificationChannel[dofAccount.UserId]; ok {
		channel <- dto.NotificationData{
			ID:        notificationReceiver.ID,
			Title:     notificationReceiver.Title,
			Body:      notificationReceiver.Body,
			Status:    notificationReceiver.Status,
			IsRead:    notificationReceiver.IsRead,
			CreatedAt: notificationReceiver.CreatedAt,
		}
	}
}
