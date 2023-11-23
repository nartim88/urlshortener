package storage

import (
	"context"
	"errors"
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
	if err = s.createTable(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s DBStorage) Get(sID models.ShortenID) (*models.FullURL, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var fURL models.FullURL
	err := s.conn.QueryRow(ctx, "SELECT full_url FROM shortener WHERE short_url=$1", sID).Scan(&fURL)
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
	sID := models.ShortenID(randChars)

	_, err := s.conn.Exec(ctx, "INSERT INTO shortener (full_url, short_url) VALUES ($1, $2);", fURL, sID)
	if err != nil {
		return nil, err
	}
	return &sID, nil
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

// createTable создание таблицы в бд
func (s DBStorage) createTable() (err error) {
	_, err = s.conn.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS shortener (
		    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
		    full_url VARCHAR(2048) NOT NULL,
		    short_url VARCHAR(8) NOT NULL,
		    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS shortener_short_url_index ON shortener (short_url);
		`,
	)
	return
}
