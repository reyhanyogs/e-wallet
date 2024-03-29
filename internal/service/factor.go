package service

import (
	"context"

	"github.com/reyhanyogs/e-wallet/domain"
	"github.com/reyhanyogs/e-wallet/dto"
	"github.com/reyhanyogs/e-wallet/internal/component"
	"golang.org/x/crypto/bcrypt"
)

type factorService struct {
	factorRepository domain.FactorRepository
}

func NewFactor(factorRepository domain.FactorRepository) domain.FactorService {
	return &factorService{
		factorRepository: factorRepository,
	}
}

func (s *factorService) ValidatePIN(ctx context.Context, req dto.ValidatePinReq) error {
	factor, err := s.factorRepository.FindByUser(ctx, req.UserID)
	if err != nil {
		component.Log.Errorf("ValidatePIN(FindByUser): user_id = %d: err = %s", req.UserID, err.Error())
		return err
	}

	if factor == (domain.Factor{}) {
		return domain.ErrPinInvalid
	}

	err = bcrypt.CompareHashAndPassword([]byte(factor.PIN), []byte(req.PIN))
	if err != nil {
		component.Log.Errorf("ValidatePIN(CompareHashAndPassword): user_id = %d: err = %s", req.UserID, err.Error())
		return domain.ErrPinInvalid
	}

	return nil
}
