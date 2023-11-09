package storage

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/nartim88/urlshortener/internal/app/shortener"
	"github.com/nartim88/urlshortener/internal/pkg/models"
)

type DBStorage struct {
	conn *pgx.Conn
}

func NewDBStorage() (Storage, error) {
	conn, err := pgx.Connect(context.Background(), shortener.App.Configs.DatabaseDSN)
	if err != nil {
		return nil, err
	}
	s := DBStorage{conn}
	return &s, nil
}

func (s *DBStorage) Get(sURL models.ShortURL) (fURL *models.FullURL, err error) {
}

func (s *DBStorage) Set(fURL models.FullURL) (sURL *models.ShortURL, err error) {
}
