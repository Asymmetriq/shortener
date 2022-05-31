package repository

import (
	"context"
	"fmt"

	"github.com/Asymmetriq/shortener/internal/models"
	"github.com/jmoiron/sqlx"
)

type dbRepository struct {
	DB *sqlx.DB
}

func newDBRepository(db *sqlx.DB) *dbRepository {
	return &dbRepository{
		DB: db,
	}
}

func (dbr *dbRepository) SetURL(ctx context.Context, entry models.StorageEntry) error {
	stmnt := `
	INSERT INTO urls(id, short_url, original_url, user_id) 
	VALUES (:id, :short_url, :original_url, :user_id) 
	ON CONFLICT (original_url) DO NOTHING`

	res, err := dbr.DB.NamedExecContext(ctx, stmnt, &entry)
	if err != nil {
		return err
	}
	if n, e := res.RowsAffected(); e == nil && n == 0 {
		return models.ErrAlreadyExists
	}
	return err
}

func (dbr *dbRepository) SetBatchURLs(ctx context.Context, entries []models.StorageEntry) error {
	if len(entries) == 0 {
		return nil
	}
	stmnt := `
	INSERT INTO urls(id, short_url, original_url, user_id) 
	VALUES (:id, :short_url, :original_url, :user_id) 
	ON CONFLICT (id) DO NOTHING`

	res, err := dbr.DB.NamedExecContext(ctx, stmnt, entries)
	if err != nil {
		return err
	}
	if n, e := res.RowsAffected(); e == nil && n == 0 {
		return models.ErrAlreadyExists
	}
	return err
}

func (dbr *dbRepository) GetURL(ctx context.Context, id string) (string, error) {
	var row models.StorageEntry
	if err := dbr.DB.GetContext(ctx, &row, "SELECT original_url FROM urls WHERE id=$1", id); err != nil {
		return "", fmt.Errorf("no original url found with shortcut %q", id)
	}
	return row.OriginalURL, nil
}

func (dbr *dbRepository) GetAllURLs(ctx context.Context, userID string) ([]models.StorageEntry, error) {
	stmnt := "SELECT original_url, short_url FROM urls WHERE user_id=$1"

	var rows []models.StorageEntry
	if err := dbr.DB.SelectContext(ctx, &rows, stmnt, userID); err != nil {
		return nil, fmt.Errorf("no data  found with userID %q", userID)
	}
	return rows, nil
}

func (dbr *dbRepository) Close() error {
	return dbr.DB.Close()
}

func (dbr *dbRepository) PingContext(ctx context.Context) error {
	return dbr.DB.PingContext(ctx)
}
