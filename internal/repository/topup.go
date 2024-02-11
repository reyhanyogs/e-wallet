package repository

import (
	"context"
	"database/sql"

	"github.com/doug-martin/goqu/v9"
	"github.com/reyhanyogs/e-wallet/domain"
)

type topUpRepository struct {
	db *goqu.Database
}

func NewTopUp(conn *sql.DB) domain.TopUpRepository {
	return &topUpRepository{
		db: goqu.New("default", conn),
	}
}

func (r *topUpRepository) FindById(ctx context.Context, id string) (topUp domain.TopUp, err error) {
	dataset := r.db.From("topup").Where(goqu.Ex{
		"id": id,
	})
	_, err = dataset.ScanStructContext(ctx, &topUp)
	return
}

func (r *topUpRepository) Insert(ctx context.Context, t *domain.TopUp) error {
	executor := r.db.Insert("topup").Rows(goqu.Record{
		"id":       t.ID,
		"user_id":  t.UserID,
		"amount":   t.Amount,
		"status":   t.Status,
		"snap_url": t.SnapURL,
	}).Executor()
	_, err := executor.ExecContext(ctx)
	return err
}

func (r *topUpRepository) Update(ctx context.Context, t *domain.TopUp) error {
	executor := r.db.Update("topup").Where(goqu.Ex{
		"id": t.ID,
	}).Set(goqu.Record{
		"amount":   t.Amount,
		"status":   t.Status,
		"snap_url": t.SnapURL,
	}).Executor()
	_, err := executor.ExecContext(ctx)
	return err
}
