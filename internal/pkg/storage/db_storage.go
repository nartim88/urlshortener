package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/nartim88/urlshortener/internal/pkg/models"
	"github.com/nartim88/urlshortener/internal/pkg/service"
)

type DBStorage struct {
	conn *pgx.Conn
}

func NewDBStorage(dsn string) (Storage, error) {
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return nil, err
	}
	s := DBStorage{conn}
	if err = s.bootstrap(); err != nil {
		return nil, err
	}
	return &s, nil
}

func (s DBStorage) Get(sID models.ShortenID) (*models.FullURL, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var fURL models.FullURL
	err := s.conn.QueryRow(ctx, `
		SELECT full_url 
		FROM shortener 
		WHERE short_url=$1`,
		sID,
	).Scan(&fURL)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &fURL, nil
}

func (s DBStorage) Set(fURL models.FullURL) (*models.ShortenID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	randChars := service.GenerateRandChars(shortURLLen)
	newSID := models.ShortenID(randChars)
	var resSID models.ShortenID

	err := s.conn.QueryRow(ctx, `
		INSERT INTO shortener (full_url, short_url)
		VALUES ($1, $2)
		ON CONFLICT (full_url) DO UPDATE
			SET full_url = EXCLUDED.full_url
		RETURNING short_url;
		`,
		fURL, newSID,
	).Scan(&resSID)
	if err != nil {
		return nil, fmt.Errorf("error while trying to save data in the db: %w", err)
	}
	if newSID != resSID {
		err = URLExistsError{
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

// bootstrap создание необходимых таблиц и индексов в бд
func (s DBStorage) bootstrap() (err error) {
	_, err = s.conn.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS shortener (
		    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
		    full_url VARCHAR(2048) NOT NULL CHECK (full_url <> ''),
		    short_url VARCHAR(8) NOT NULL CHECK (short_url <> ''),
		    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS shortener_short_url_idx ON shortener (short_url);
		CREATE UNIQUE INDEX IF NOT EXISTS shortener_full_url_unique_idx ON shortener (full_url)
		`,
	)
	return
}
