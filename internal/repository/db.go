package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/Asymmetriq/shortener/internal/models"
	"github.com/jmoiron/sqlx"
)

const (
	MaxReqNum = 5
	BatchSize = 20
)

type dbRepository struct {
	DB            *sqlx.DB
	BatchChannel  chan models.DeleteRequest
	SignalChannel chan struct{}
}

func newDBRepository(db *sqlx.DB) *dbRepository {
	repo := &dbRepository{
		DB:           db,
		BatchChannel: make(chan models.DeleteRequest),
	}
	repo.backgroundDelete(context.Background())
	return repo
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
	if err := dbr.DB.GetContext(ctx, &row, "SELECT original_url, deleted FROM urls WHERE id=$1", id); err != nil {
		return "", fmt.Errorf("no original url found with shortcut %q", id)
	}
	if row.Deleted {
		return "", models.ErrDeleted
	}
	return row.OriginalURL, nil
}

func (dbr *dbRepository) GetAllURLs(ctx context.Context, userID string) ([]models.StorageEntry, error) {
	stmnt := "SELECT original_url, short_url FROM urls WHERE user_id=$1 AND deleted=false"

	var rows []models.StorageEntry
	if err := dbr.DB.SelectContext(ctx, &rows, stmnt, userID); err != nil {
		return nil, fmt.Errorf("no data  found with userID %q", userID)
	}
	return rows, nil
}

func (dbr *dbRepository) BatchDelete(ctx context.Context, req models.DeleteRequest) {
	go func(delReq models.DeleteRequest) {
		dbr.BatchChannel <- delReq
	}(req)
}

func (dbr *dbRepository) deleteBatch(ctx context.Context, req models.DeleteRequest) {
	stmnt := "UPDATE urls SET deleted=true WHERE user_id=$1 AND id=any($2);"

	for i := 0; i < len(req.IDs); i += BatchSize {
		end := i + BatchSize
		if end > len(req.IDs) {
			end = len(req.IDs)
		}
		_, err := dbr.DB.ExecContext(ctx, stmnt, req.UserID, req.IDs[i:end])
		if err != nil {
			log.Printf("async delete: %v", err)
		}
	}
}

func (dbr *dbRepository) backgroundDelete(ctx context.Context) {
	go func(ctx context.Context) {
		for {
			v := <-dbr.BatchChannel
			dbr.deleteBatch(ctx, v)
		}
	}(ctx)
}

func (dbr *dbRepository) Close() error {
	return dbr.DB.Close()
}

func (dbr *dbRepository) PingContext(ctx context.Context) error {
	return dbr.DB.PingContext(ctx)
}
