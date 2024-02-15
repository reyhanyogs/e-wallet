package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/reyhanyogs/e-wallet/domain"
	"github.com/reyhanyogs/e-wallet/dto"
)

type ipCheckerService struct {
}

func NewIpChecker() domain.IpCheckerService {
	return &ipCheckerService{}
}

func (s *ipCheckerService) Query(ctx context.Context, ip string) (checker dto.IpChecker, err error) {
	url := fmt.Sprintf("http://ip-api.com/json/%s?fields=status,lat,lon,timezone,query", ip)
	resp, err := http.Get(url)
	if err != nil {
		return dto.IpChecker{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return dto.IpChecker{}, err
	}

	err = json.Unmarshal(body, &checker)
	return
}
