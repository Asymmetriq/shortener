package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Asymmetriq/shortener/internal/generated/shortener/public/model"
	"github.com/Asymmetriq/shortener/internal/generated/shortener/public/table"
	"github.com/Asymmetriq/shortener/internal/models"
	"github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
)

type dbRepository struct {
	DB *sql.DB
}

func newDBRepository(db *sql.DB) *dbRepository {
	return &dbRepository{
		DB: db,
	}
}

func (dbr *dbRepository) SetURL(ctx context.Context, entry models.StorageEntry) error {
	value := model.Urls{
		ID:          entry.ID,
		ShortURL:    entry.ShortURL,
		OriginalURL: entry.OriginalURL,
		UserID:      uuid.MustParse(entry.UserID),
	}
	insertStmnt := table.Urls.
		INSERT(
			table.Urls.ID,
			table.Urls.OriginalURL,
			table.Urls.ShortURL,
			table.Urls.UserID,
		).MODEL(value).
		ON_CONFLICT(table.Urls.ID).
		DO_UPDATE(postgres.SET(table.Urls.UserID.SET(postgres.String(entry.UserID))))

	_, err := insertStmnt.ExecContext(ctx, dbr.DB)
	return err
}

func (dbr *dbRepository) SetBatchURLs(ctx context.Context, entries []models.StorageEntry) error {
	if len(entries) == 0 {
		return nil
	}
	userID := uuid.MustParse(entries[0].UserID)

	values := make([]model.Urls, len(entries))
	for i, e := range entries {
		values[i] = model.Urls{
			ID:          e.ID,
			ShortURL:    e.ShortURL,
			OriginalURL: e.OriginalURL,
			UserID:      userID,
		}
	}
	insertStmnt := table.Urls.
		INSERT(
			table.Urls.ID,
			table.Urls.OriginalURL,
			table.Urls.ShortURL,
			table.Urls.UserID,
		).MODELS(values).ON_CONFLICT(table.Urls.ID).
		DO_UPDATE(postgres.SET(table.Urls.UserID.SET(postgres.String(userID.String()))))

	_, err := insertStmnt.ExecContext(ctx, dbr.DB)
	return err
}

func (dbr *dbRepository) GetURL(ctx context.Context, id string) (string, error) {
	selectStmnt := table.Urls.
		SELECT(table.Urls.OriginalURL).
		WHERE(table.Urls.ID.EQ(postgres.String(id))).
		LIMIT(1)

	var row model.Urls
	if err := selectStmnt.QueryContext(ctx, dbr.DB, &row); err != nil {
		return "", fmt.Errorf("no original url found with shortcut %q", id)
	}
	return row.OriginalURL, nil
}

func (dbr *dbRepository) GetAllURLs(ctx context.Context, userID string) ([]models.StorageEntry, error) {
	selectStmnt := table.Urls.
		SELECT(table.Urls.OriginalURL, table.Urls.ShortURL).FROM(table.Urls).
		WHERE(table.Urls.UserID.EQ(postgres.String(userID)))

	var rows []model.Urls
	if err := selectStmnt.QueryContext(ctx, dbr.DB, &rows); err != nil {
		return nil, fmt.Errorf("no data  found with userID %q", userID)
	}

	data := make([]models.StorageEntry, 0, len(rows))
	for _, r := range rows {
		data = append(data, models.StorageEntry{
			OriginalURL: r.OriginalURL,
			ShortURL:    r.ShortURL,
		})
	}
	return data, nil
}

func (dbr *dbRepository) Close() error {
	return dbr.DB.Close()
}

func (dbr *dbRepository) PingContext(ctx context.Context) error {
	return dbr.DB.PingContext(ctx)
}
