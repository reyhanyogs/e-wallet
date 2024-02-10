package repository

import (
	"context"
	"database/sql"

	"github.com/doug-martin/goqu/v9"
	"github.com/reyhanyogs/e-wallet/domain"
)

type templateRepository struct {
	db *goqu.Database
}

func NewTemplate(conn *sql.DB) domain.TemplateRepository {
	return &templateRepository{
		db: goqu.New("default", conn),
	}
}

func (r *templateRepository) FindByCode(ctx context.Context, code string) (tmp domain.Template, err error) {
	dataset := r.db.From("templates").Where(goqu.Ex{
		"code": code,
	})
	_, err = dataset.ScanStructContext(ctx, &tmp)
	return
}
