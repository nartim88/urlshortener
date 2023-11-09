package storage

import "github.com/nartim88/urlshortener/internal/pkg/models"

type DBStorage struct {
}

func NewDBStorage() (Storage, error) {
	s := DBStorage{}
	return &s, nil
}

func (s *DBStorage) Get(sURL models.ShortURL) (fURL *models.FullURL, err error) {
}

func (s *DBStorage) Set(fURL models.FullURL) (sURL *models.ShortURL, err error) {
}
