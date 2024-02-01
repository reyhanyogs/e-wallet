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
