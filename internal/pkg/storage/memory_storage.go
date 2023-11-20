package storage

import (
	"errors"

	"github.com/nartim88/urlshortener/internal/pkg/models"
	"github.com/nartim88/urlshortener/internal/pkg/service"
)

type MemStorage struct {
	Memory map[models.ShortenID]models.FullURL
}

// NewMemStorage инициализация Storage в памяти
func NewMemStorage() Storage {
	s := MemStorage{
		make(map[models.ShortenID]models.FullURL),
	}
	return &s
}

func (s *MemStorage) Get(sURL models.ShortenID) (*models.FullURL, error) {
	if !s.isExist(sURL) {
		return nil, nil
	}
	fURL := s.Memory[sURL]
	return &fURL, nil
}

func (s *MemStorage) Set(fURL models.FullURL) (*models.ShortenID, error) {
	randChars := service.GenerateRandChars(shortURLLen)
	shortURL := models.ShortenID(randChars)

	if s.isExist(shortURL) {
		return nil, errors.New("URL already exists")
	}
	s.Memory[shortURL] = fURL
	return &shortURL, nil
}

// isExist проверяет сохранен ли в памяти короткий УРЛ
func (s *MemStorage) isExist(sURL models.ShortenID) bool {
	_, ok := s.Memory[sURL]
	return ok
}
