package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/reyhanyogs/e-wallet/domain"
	"github.com/reyhanyogs/e-wallet/dto"
	"github.com/reyhanyogs/e-wallet/internal/util"
)

type transactionService struct {
	accountRepository     domain.AccountRepository
	transactionRepository domain.TransactionRepository
	cacheRepository       domain.CacheRepository
}

func NewTransaction(
	accountRepository domain.AccountRepository,
	transactionRepository domain.TransactionRepository,
	cacheRepository domain.CacheRepository,
) domain.TransactionService {
	return &transactionService{
		accountRepository:     accountRepository,
		transactionRepository: transactionRepository,
		cacheRepository:       cacheRepository,
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

	return nil
}
