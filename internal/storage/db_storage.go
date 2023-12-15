package storage

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"

	"github.com/nartim88/urlshortener/internal/models"
	"github.com/nartim88/urlshortener/internal/utils"
	"github.com/nartim88/urlshortener/pkg/logger"
)

type DBStorage struct {
	conn *pgx.Conn
}

func NewDBStorage(conn *pgx.Conn) StorageWithService {
	return &DBStorage{conn}
}

// Get возвращает оригинальный урл по строковому идентификатору.
// Возвращает pgx.ErrNoRows, если в базе не найдены записи, или
// ErrURLDeleted, если записи помечены как удаленные.
func (s DBStorage) Get(ctx context.Context, sID models.ShortenID) (*models.FullURL, error) {
	type queryRes struct {
		fURL      models.FullURL
		isDeleted bool
	}

	var qRes queryRes

	err := s.conn.QueryRow(ctx, `
		SELECT full_url, is_deleted
		FROM shortener 
		WHERE shorten_id=$1`,
		sID,
	).Scan(&qRes.fURL, &qRes.isDeleted)

	logger.Log.Debug().Msgf("qres: %#v", qRes)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("error while getting full_url from db: %w", err)
	}

	if qRes.isDeleted {
		return nil, ErrURLDeleted
	}

	return &qRes.fURL, nil
}

func (s DBStorage) Set(ctx context.Context, fURL models.FullURL) (*models.ShortenID, error) {
	randChars := GenerateRandChars(shortURLLen)
	newSID := models.ShortenID(randChars)
	var resSID models.ShortenID

	userID, err := utils.GetUserIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	err = s.conn.QueryRow(ctx, `
		INSERT INTO shortener (user_id, full_url, shorten_id)
		VALUES ($1, $2, $3)
		ON CONFLICT (full_url) DO UPDATE
			SET full_url = EXCLUDED.full_url
		RETURNING shorten_id;
		`,
		userID, fURL, newSID,
	).Scan(&resSID)

	if err != nil {
		return nil, fmt.Errorf("error while trying to save data in the db: %w", err)
	}

	if newSID != resSID {
		err = ErrURLExists{
			fURL,
			resSID,
		}
		return nil, err
	}

	return &newSID, nil
}

func (s DBStorage) Close(ctx context.Context) error {
	if s.conn == nil {
		return errors.New("db connection doesn't exists or already closed")
	}
	if err := s.conn.Close(ctx); err != nil {
		return err
	}
	return nil
}

func (s DBStorage) Bootstrap(ctx context.Context) (err error) {
	_, err = s.conn.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS shortener (
		    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
		    user_id uuid,
		    full_url VARCHAR(2048) NOT NULL CHECK (full_url <> ''),
		    shorten_id VARCHAR(8) NOT NULL CHECK (shorten_id <> ''),
		    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		    is_deleted BOOLEAN DEFAULT false
		);
		CREATE INDEX IF NOT EXISTS shortener_short_url_idx ON shortener (shorten_id);
		CREATE UNIQUE INDEX IF NOT EXISTS shortener_full_url_unique_idx ON shortener (full_url)
		`,
	)
	return
}

func (s DBStorage) ListURLs(ctx context.Context, u models.User) ([]SIDAndFullURL, error) {
	rows, err := s.conn.Query(ctx, `
		SELECT shorten_id, full_url FROM shortener
		WHERE user_id = $1
		`, u.UserID,
	)
	if err != nil {
		return nil, fmt.Errorf("error while trying to get URLs by user id: %w", err)
	}
	defer rows.Close()
	urls, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (SIDAndFullURL, error) {
		var (
			shortenID models.ShortenID
			fullURL   models.FullURL
		)
		if err := row.Scan(&shortenID, &fullURL); err != nil {
			return SIDAndFullURL{}, fmt.Errorf("scan error: %w", err)
		}
		return SIDAndFullURL{ShortenID: shortenID, FullURL: fullURL}, nil
	})
	if err != nil {
		return nil, err
	}
	return urls, nil
}

func (s DBStorage) MarkAsDeletedByID(ctx context.Context, IDs []models.ShortenID) error {
	var params []string

	for _, id := range IDs {
		params = append(params, fmt.Sprintf("'%s'", id))
	}

	paramsStr := strings.Join(params, ", ")
	paramsStr = "(" + paramsStr + ")"
	query := `
		UPDATE shortener SET is_deleted = true
		WHERE shorten_id IN ` + paramsStr
	tag, err := s.conn.Exec(ctx, query)

	logger.Log.Debug().Msgf("query params: %v", paramsStr)
	logger.Log.Debug().Msgf("query response tag: %v", tag)

	if err != nil {
		return err
	}

	return nil
}
