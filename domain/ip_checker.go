package domain

import (
	"context"

	"github.com/reyhanyogs/e-wallet/dto"
)

type IpCheckerService interface {
	Query(ctx context.Context, ip string) (dto.IpChecker, error)
}
