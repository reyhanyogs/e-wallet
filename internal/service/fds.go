package service

import (
	"context"
	"log"
	"time"

	"github.com/reyhanyogs/e-wallet/domain"
	"github.com/reyhanyogs/e-wallet/dto"
	"github.com/reyhanyogs/e-wallet/internal/component"
	"github.com/reyhanyogs/e-wallet/internal/util"
)

type fdsService struct {
	ipCheckerService   domain.IpCheckerService
	loginLogRepository domain.LoginLogRepository
}

func NewFds(ipCheckerService domain.IpCheckerService, loginLogRepository domain.LoginLogRepository) domain.FdsService {
	return &fdsService{
		ipCheckerService:   ipCheckerService,
		loginLogRepository: loginLogRepository,
	}
}

func (s *fdsService) IsAuthorized(ctx context.Context, ip string, userId int64) bool {
	locationCheck, err := s.ipCheckerService.Query(ctx, ip)
	if err != nil || locationCheck == (dto.IpChecker{}) {
		component.Log.Errorf("IsAuthorized(Query): user_id = %d: err = %s", userId, err.Error())
		return false
	}

	newAccess := domain.LoginLog{
		UserID:       userId,
		IsAuthorized: false,
		IpAddress:    locationCheck.Query,
		Timezone:     locationCheck.Timezone,
		Lat:          locationCheck.Lat,
		Lon:          locationCheck.Lon,
		AccessTime:   time.Now(),
	}

	lastLogin, err := s.loginLogRepository.FindLastAuthorized(ctx, userId)
	if err != nil {
		component.Log.Errorf("IsAuthorized(FindLastAuthorized): user_id = %d: err = %s", userId, err.Error())
		_ = s.loginLogRepository.Save(ctx, &newAccess)
		return false
	}
	if lastLogin == (domain.LoginLog{}) {
		newAccess.IsAuthorized = true
		_ = s.loginLogRepository.Save(ctx, &newAccess)
		return true
	}

	distanceHour := newAccess.AccessTime.Sub(lastLogin.AccessTime)
	distanceChange := util.GetDistance(lastLogin.Lat, lastLogin.Lon, newAccess.Lat, newAccess.Lon)

	log.Printf("hour: %f, distance: %f", distanceHour.Hours(), distanceChange)

	if distanceChange/distanceHour.Hours() > 400 {
		_ = s.loginLogRepository.Save(ctx, &newAccess)
		return false
	}

	newAccess.IsAuthorized = true
	_ = s.loginLogRepository.Save(ctx, &newAccess)
	return true
}
