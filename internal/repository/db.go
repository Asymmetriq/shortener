package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Asymmetriq/shortener/internal/generated/shortener/public/model"
	"github.com/Asymmetriq/shortener/internal/generated/shortener/public/table"
	"github.com/Asymmetriq/shortener/internal/shorten"
	"github.com/go-jet/jet/v2/postgres"
)

type dbRepository struct {
	DB *sql.DB
}

func newDBRepository(db *sql.DB) *dbRepository {
	return &dbRepository{
		DB: db,
	}
}

func (dbr *dbRepository) SetURL(ctx context.Context, url, userID, host string) (string, error) {
	id := shorten.Shorten(url)
	shortURL := buildURL(host, id)

	insertStmnt := table.Urls.
		INSERT(
			table.Urls.ID,
			table.Urls.OriginalURL,
			table.Urls.ShortURL,
			table.Urls.UserID,
		).VALUES(
		id,
		url,
		shortURL,
		userID,
	).ON_CONFLICT(table.Urls.ID).DO_UPDATE(postgres.SET(table.Urls.UserID.SET(postgres.String(userID))))

	_, err := insertStmnt.ExecContext(ctx, dbr.DB)
	if err != nil {
		return "", err
	}

	return shortURL, nil
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

func (dbr *dbRepository) GetAllURLs(ctx context.Context, userID string) ([]Data, error) {
	selectStmnt := table.Urls.
		SELECT(table.Urls.OriginalURL, table.Urls.ShortURL).FROM(table.Urls).
		WHERE(table.Urls.UserID.EQ(postgres.String(userID)))

	var rows []model.Urls
	if err := selectStmnt.QueryContext(ctx, dbr.DB, &rows); err != nil {
		return []Data{}, fmt.Errorf("no data  found with userID %q", userID)
	}
	data := make([]Data, 0, len(rows))
	for _, r := range rows {
		data = append(data, Data{
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
