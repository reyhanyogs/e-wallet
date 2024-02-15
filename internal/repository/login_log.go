package repository

import (
	"context"
	"database/sql"

	"github.com/doug-martin/goqu/v9"
	"github.com/reyhanyogs/e-wallet/domain"
)

type loginLogRepository struct {
	db *goqu.Database
}

func NewLoginLog(conn *sql.DB) domain.LoginLogRepository {
	return &loginLogRepository{
		db: goqu.New("default", conn),
	}
}

func (r *loginLogRepository) FindLastAuthorized(ctx context.Context, userId int64) (loginLog domain.LoginLog, err error) {
	dataset := r.db.From("login_log").Where(goqu.Ex{
		"user_id":       userId,
		"is_authorized": true,
	}).Order(goqu.I("id").Desc()).Limit(1)
	_, err = dataset.ScanStructContext(ctx, &loginLog)
	return
}

func (r *loginLogRepository) Save(ctx context.Context, loginLog *domain.LoginLog) error {
	executor := r.db.Insert("login_log").Rows(goqu.Record{
		"user_id":       loginLog.UserID,
		"is_authorized": loginLog.IsAuthorized,
		"ip_address":    loginLog.IpAddress,
		"timezone":      loginLog.Timezone,
		"lat":           loginLog.Lat,
		"lon":           loginLog.Lon,
		"access_time":   loginLog.AccessTime,
	}).Returning("id").Executor()
	_, err := executor.ScanStructContext(ctx, loginLog)
	return err
}
