package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/reyhanyogs/e-wallet/domain"
)

type transactionRepository struct {
	db *goqu.Database
}

func NewTransaction(conn *sql.DB) domain.TransactionRepository {
	return &transactionRepository{
		db: goqu.New("default", conn),
	}
}

func (r *transactionRepository) Insert(ctx context.Context, transaction *domain.Transaction) error {
	executor := r.db.Insert("transactions").Rows(goqu.Record{
		"account_id":           transaction.AccountId,
		"sof_number":           transaction.SofNumber,
		"dof_number":           transaction.DofNumber,
		"transaction_type":     transaction.TransactionType,
		"amount":               transaction.Amount,
		"transaction_datetime": time.Now(),
	}).Returning("id").Executor()
	_, err := executor.ScanStructContext(ctx, transaction)
	return err
}
