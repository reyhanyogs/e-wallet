package repository

import (
	"context"
	"database/sql"

	"github.com/doug-martin/goqu/v9"
	"github.com/reyhanyogs/e-wallet/domain"
)

type factorRepository struct {
	db *goqu.Database
}

func NewFactor(conn *sql.DB) domain.FactorRepository {
	return &factorRepository{
		db: goqu.New("default", conn),
	}
}

func (r *factorRepository) FindByUser(ctx context.Context, id int64) (factor domain.Factor, err error) {
	dataset := r.db.From("factors").Where(goqu.Ex{
		"user_id": id,
	})
	_, err = dataset.ScanStructContext(ctx, &factor)
	return
}

func (r *factorRepository) Insert(ctx context.Context, factor *domain.Factor) error {
	executor := r.db.Insert("factors").Rows(goqu.Record{
		"user_id": factor.UserID,
		"pin":     factor.PIN,
	}).Returning("id").Executor()
	_, err := executor.ScanStructContext(ctx, factor)
	return err
}
