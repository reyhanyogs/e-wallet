package repository

import (
	"context"
	"database/sql"

	"github.com/doug-martin/goqu/v9"
	"github.com/reyhanyogs/e-wallet/domain"
)

type userRepository struct {
	db *goqu.Database
}

func NewUser(conn *sql.DB) domain.UserRepository {
	return &userRepository{
		db: goqu.New("default", conn),
	}
}

func (r *userRepository) FindByID(ctx context.Context, id int64) (user domain.User, err error) {
	data := r.db.From("users").Where(goqu.Ex{
		"id": id,
	})

	_, err = data.ScanStructContext(ctx, &user)
	return
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (user domain.User, err error) {
	data := r.db.From("users").Where(goqu.Ex{
		"username": username,
	})

	_, err = data.ScanStructContext(ctx, &user)
	return
}

func (r *userRepository) Insert(ctx context.Context, user *domain.User) error {
	executor := r.db.Insert("users").Rows(goqu.Record{
		"full_name": user.FullName,
		"username":  user.Username,
		"password":  user.Password,
		"phone":     user.Phone,
		"email":     user.Email,
	}).Returning("id").Executor()

	_, err := executor.ScanStructContext(ctx, user)
	return err
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	user.EmailVerifiedAtDB = sql.NullTime{
		Time:  user.EmailVerifiedAt,
		Valid: true,
	}

	executor := r.db.Update("users").Where(goqu.Ex{
		"id": user.ID,
	}).Set(goqu.Record{
		"full_name":         user.FullName,
		"username":          user.Username,
		"password":          user.Password,
		"phone":             user.Phone,
		"email":             user.Email,
		"email_verified_at": user.EmailVerifiedAtDB,
	}).Executor()

	_, err := executor.ExecContext(ctx)
	return err
}
