package repository

import (
	"context"
	"database/sql"

	"github.com/doug-martin/goqu/v9"
	"github.com/reyhanyogs/e-wallet/domain"
)

type accountRepository struct {
	db *goqu.Database
}

func NewAccount(conn *sql.DB) domain.AccountRepository {
	return &accountRepository{
		db: goqu.New("default", conn),
	}
}

func (r *accountRepository) Create(ctx context.Context, account *domain.Account) error {
	executor := r.db.Insert("accounts").Rows(goqu.Record{
		"user_id":        account.UserId,
		"account_number": account.AccountNumber,
		"balance":        account.Balance,
	}).Returning("id").Executor()
	_, err := executor.ScanStructContext(ctx, account)
	return err
}

func (r *accountRepository) FindByUserID(ctx context.Context, id int64) (account domain.Account, err error) {
	dataset := r.db.From("accounts").Where(goqu.Ex{
		"user_id": id,
	})
	_, err = dataset.ScanStructContext(ctx, &account)
	return
}

func (r *accountRepository) FindByAccountNumber(ctx context.Context, accNumber string) (account domain.Account, err error) {
	dataset := r.db.From("accounts").Where(goqu.Ex{
		"account_number": accNumber,
	})
	_, err = dataset.ScanStructContext(ctx, &account)
	return
}

func (r *accountRepository) Update(ctx context.Context, account *domain.Account) error {
	executor := r.db.Update("accounts").Where(goqu.Ex{
		"id": account.ID,
	}).Set(goqu.Record{
		"balance": account.Balance,
	}).Executor()
	_, err := executor.ExecContext(ctx)
	return err
}
