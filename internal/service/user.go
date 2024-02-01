package service

import (
	"context"
	"encoding/json"

	"github.com/reyhanyogs/e-wallet/domain"
	"github.com/reyhanyogs/e-wallet/dto"
	"github.com/reyhanyogs/e-wallet/internal/util"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	repository      domain.UserRepository
	cacheRepository domain.CacheRepository
}

func NewUser(repository domain.UserRepository, cacheRepository domain.CacheRepository) domain.UserService {
	return &userService{
		repository:      repository,
		cacheRepository: cacheRepository,
	}
}

func (s *userService) Authenticate(ctx context.Context, req dto.AuthReq) (dto.AuthRes, error) {
	user, err := s.repository.FindByUsername(ctx, req.Username)
	if err != nil {
		return dto.AuthRes{}, err
	}
	if user == (domain.User{}) {
		return dto.AuthRes{}, domain.ErrAuthFailed
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return dto.AuthRes{}, domain.ErrAuthFailed
	}

	token := util.GenerateRandomString(16)

	userJson, _ := json.Marshal(user)
	_ = s.cacheRepository.Set("user:"+token, userJson)
	return dto.AuthRes{
		Token: token,
	}, nil
}

func (s *userService) ValidateToken(ctx context.Context, token string) (dto.UserData, error) {
	data, err := s.cacheRepository.Get("user:" + token)
	if err != nil {
		return dto.UserData{}, err
	}

	var user domain.User

	err = json.Unmarshal(data, &user)
	if err != nil {
		return dto.UserData{}, err
	}

	return dto.UserData{
		ID:       user.ID,
		FullName: user.FullName,
		Phone:    user.Phone,
		Username: user.Username,
	}, nil
}
