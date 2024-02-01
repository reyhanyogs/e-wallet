package util

import (
	"errors"

	"github.com/reyhanyogs/e-wallet/domain"
)

func GetHttpStatus(err error) int {
	switch {
	case errors.Is(err, domain.ErrAuthFailed):
		return 401
	default:
		return 500
	}
}
