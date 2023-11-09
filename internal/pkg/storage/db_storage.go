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
