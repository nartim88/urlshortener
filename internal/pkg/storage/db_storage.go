package storage

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/nartim88/urlshortener/internal/pkg/models"
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
	return &s, nil
}

func (s *DBStorage) Get(sURL models.ShortURL) (fURL *models.FullURL, err error) {
	return nil, nil
}

func (s *DBStorage) Set(fURL models.FullURL) (sURL *models.ShortURL, err error) {
	return nil, nil
}

func (s *DBStorage) Close(ctx context.Context) error {
	if s.conn != nil {
		if err := s.conn.Close(ctx); err != nil {
			return err
		}
	}
	return errors.New("DB connection doesn't exists or already closed")
}

// createTable создание таблицы в бд
func (s *DBStorage) createTable() error {
	_, err := s.conn.Exec(
		context.Background(),
		`
		CREATE TABLE IF NOT EXISTS shortener (
		    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
		    full_url varchar(256),
		    short_url varchar(8)
		);
		CREATE INDEX IF NOT EXISTS shortener_short_url_index ON shortener (short_url);
		`,
	)
	if err != nil {
		return err
	}
	return nil
}
