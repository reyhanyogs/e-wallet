package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/reyhanyogs/e-wallet/domain"
	"github.com/reyhanyogs/e-wallet/dto"
	"github.com/reyhanyogs/e-wallet/internal/component"
	"github.com/reyhanyogs/e-wallet/internal/util"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	repository        domain.UserRepository
	cacheRepository   domain.CacheRepository
	emailService      domain.EmailService
	factorRepository  domain.FactorRepository
	accountRepository domain.AccountRepository
}

func NewUser(
	repository domain.UserRepository,
	cacheRepository domain.CacheRepository,
	emailService domain.EmailService,
	factorRepository domain.FactorRepository,
	accountRepository domain.AccountRepository,
) domain.UserService {
	return &userService{
		repository:        repository,
		cacheRepository:   cacheRepository,
		emailService:      emailService,
		factorRepository:  factorRepository,
		accountRepository: accountRepository,
	}
}

func (s *userService) Authenticate(ctx context.Context, req dto.AuthReq) (dto.AuthRes, error) {
	user, err := s.repository.FindByUsername(ctx, req.Username)
	if err != nil {
		component.Log.Errorf("Authenticate(FindByUsername): username = %s: err = %s", req.Username, err.Error())
		return dto.AuthRes{}, err
	}
	if user == (domain.User{}) {
		return dto.AuthRes{}, domain.ErrAuthFailed
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return dto.AuthRes{}, domain.ErrAuthFailed
	}

	if !user.EmailVerifiedAtDB.Valid {
		return dto.AuthRes{}, domain.ErrAuthFailed
	}

	token := util.GenerateRandomString(16)

	userJson, _ := json.Marshal(user)
	_ = s.cacheRepository.Set("user:"+token, userJson)
	return dto.AuthRes{
		UserId: user.ID,
		Token:  token,
	}, nil
}

func (s *userService) ValidateToken(ctx context.Context, token string) (dto.UserData, error) {
	data, err := s.cacheRepository.Get("user:" + token)
	if err != nil {
		component.Log.Errorf("ValidateToken(Get): token = %s: err = %s", token, err.Error())
		return dto.UserData{}, err
	}

	var user domain.User

	err = json.Unmarshal(data, &user)
	if err != nil {
		component.Log.Errorf("ValidateToken(Unmarshal): token = %s: err = %s", token, err.Error())
		return dto.UserData{}, err
	}

	return dto.UserData{
		ID:       user.ID,
		FullName: user.FullName,
		Phone:    user.Phone,
		Username: user.Username,
	}, nil
}

func (s *userService) Register(ctx context.Context, req dto.UserRegisterReq) (dto.UserRegisterRes, error) {
	exist, err := s.repository.FindByUsername(ctx, req.Username)
	if err != nil {
		component.Log.Errorf("Register(FindByUsername): username = %s: err = %s", req.Username, err.Error())
		return dto.UserRegisterRes{}, err
	}

	if exist != (domain.User{}) {
		return dto.UserRegisterRes{}, domain.ErrUsernameTaken
	}

	hashedPass, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	hashedPIN, _ := bcrypt.GenerateFromPassword([]byte(req.PIN), 12)

	user := domain.User{
		FullName: req.FullName,
		Phone:    req.Phone,
		Email:    req.Email,
		Username: req.Username,
		Password: string(hashedPass),
	}

	otpCode := util.GenerateRandomNumber(4)
	referenceId := util.GenerateRandomString(16)

	err = s.emailService.Send(req.Email, "OTP Code", "Your OTP Code are: "+otpCode)
	if err != nil {
		component.Log.Errorf("Register(Send): email = %s: err = %s", req.Email, err.Error())
		return dto.UserRegisterRes{}, err
	}

	err = s.repository.Insert(ctx, &user)
	if err != nil {
		component.Log.Errorf("Register(Insert): username = %s: err = %s", req.Username, err.Error())
		return dto.UserRegisterRes{}, err
	}

	factor := domain.Factor{
		UserID: user.ID,
		PIN:    string(hashedPIN),
	}

	err = s.factorRepository.Insert(ctx, &factor)
	if err != nil {
		component.Log.Errorf("Register(Insert): factor_user_id = %d: err = %s", factor.UserID, err.Error())
		return dto.UserRegisterRes{}, err
	}

	account := &domain.Account{
		UserId:        user.ID,
		AccountNumber: util.GenerateRandomNumber(12),
		Balance:       0,
	}

	err = s.accountRepository.Create(ctx, account)
	if err != nil {
		component.Log.Errorf("Register(Create): user_id = %d: err = %s", user.ID, err.Error())
		return dto.UserRegisterRes{}, err
	}

	_ = s.cacheRepository.Set("otp:"+referenceId, []byte(otpCode))
	_ = s.cacheRepository.Set("user-ref:"+referenceId, []byte(user.Username))
	return dto.UserRegisterRes{
		ReferenceID: referenceId,
	}, nil
}

func (s *userService) ValidateOTP(ctx context.Context, req dto.ValidateOtpReq) error {
	data, err := s.cacheRepository.Get("otp:" + req.ReferenceID)
	if err != nil {
		component.Log.Error(fmt.Sprintf("ValidateOTP(Get): otp_reference_id = %s: err = %s", req.ReferenceID, err.Error()))
		return domain.ErrOtpNotFound
	}

	otp := string(data)
	if otp != req.OTP {
		return domain.ErrOtpInvalid
	}

	data, err = s.cacheRepository.Get("user-ref:" + req.ReferenceID)
	if err != nil {
		component.Log.Error(fmt.Sprintf("ValidateOTP(Get): user_ref_reference_id = %s: err = %s", req.ReferenceID, err.Error()))
		return domain.ErrUsernameNotFound
	}
	user, err := s.repository.FindByUsername(ctx, string(data))
	if err != nil {
		component.Log.Error(fmt.Sprintf("ValidateOTP(FindByUsername): username = %s: err = %s", string(data), err.Error()))
		return err
	}

	user.EmailVerifiedAt = time.Now()
	_ = s.repository.Update(ctx, &user)

	return nil
}
