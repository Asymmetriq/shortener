package repository

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Asymmetriq/shortener/internal/models"
	"github.com/jmoiron/sqlx"
)

const (
	MaxReqNum = 5
	BatchSize = 20
)

type dbRepository struct {
	DB *sqlx.DB

	batchChannel    chan models.DeleteRequest
	signalTimer     *time.Timer
	groupedRequests map[string][]string
	isTimerRunning  bool
	once            sync.Once
}

func newDBRepository(db *sqlx.DB) *dbRepository {
	repo := &dbRepository{
		DB:              db,
		batchChannel:    make(chan models.DeleteRequest, MaxReqNum),
		groupedRequests: make(map[string][]string),
	}
	return repo
}
func (dbr *dbRepository) Signal() {
	if !dbr.isTimerRunning {
		dbr.signalTimer = time.NewTimer(3 * time.Second)
		dbr.isTimerRunning = true
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

func (dbr *dbRepository) BatchDelete(req models.DeleteRequest) {
	go func(delReq models.DeleteRequest) {
		dbr.batchChannel <- delReq
	}(req)
	dbr.Signal()
	dbr.once.Do(func() {
		dbr.backgroundDelete()
	})
}

func (dbr *dbRepository) PingContext(ctx context.Context) error {
	return dbr.DB.PingContext(ctx)
}

func (dbr *dbRepository) deleteBatch(userID string, IDs []string) {
	stmnt := "UPDATE urls SET deleted=true WHERE user_id=$1 AND id=any($2);"

	for i := 0; i < len(IDs); i += BatchSize {
		end := i + BatchSize
		if end > len(IDs) {
			end = len(IDs)
		}
		_, err := dbr.DB.Exec(stmnt, userID, IDs[i:end])
		if err != nil {
			log.Printf("async delete: %v", err)
		}
	}
}

func (dbr *dbRepository) backgroundDelete() {
	go func() {
		defer func() {
			if p := recover(); p != nil {
				log.Printf("recovered from %v", p)
			}
		}()

		for {
			select {
			case <-dbr.signalTimer.C:
				for userID, IDs := range dbr.groupedRequests {
					dbr.deleteBatch(userID, IDs)
				}
				dbr.groupedRequests = make(map[string][]string)
				dbr.isTimerRunning = false

			case req, ok := <-dbr.batchChannel:
				if !ok {
					return
				}
				dbr.groupRequests(req)
			}
		}
	}()
}
func (dbr *dbRepository) groupRequests(req models.DeleteRequest) {
	if _, ok := dbr.groupedRequests[req.UserID]; ok {
		dbr.groupedRequests[req.UserID] = append(dbr.groupedRequests[req.UserID], req.IDs...)
	}
	dbr.groupedRequests[req.UserID] = req.IDs
}

func (dbr *dbRepository) Close() error {
	close(dbr.batchChannel)
	return dbr.DB.Close()
}
