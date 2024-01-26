package mysql

import (
	"context"
	"database/sql"
	"github.com/fredmayer/go-rest-api-template/internal/domain/models"
)

type DummyRepository struct {
	db *sql.DB
}

func NewDummy(db *sql.DB) *DummyRepository {
	return &DummyRepository{
		db,
	}
}

func (r *DummyRepository) GetById(ctx context.Context, id int) (models.Dummy, error) {
	stmt, err := r.db.Prepare("SELECT * FROM dummy WHERE dummy = ?")
	if err != nil {
		return models.Dummy{}, err
	}

	row := stmt.QueryRowContext(ctx, id)

	var model models.Dummy

	err = row.Scan(&model.Id, &model.Name, &model.Description, &model.CreatedAt)
	if err != nil {
		return models.Dummy{}, err
	}

	return model, nil
}
