package storage

import (
	"github.com/nartim88/urlshortener/internal/pkg/models"
	"github.com/nartim88/urlshortener/internal/pkg/service"
)

const shortURLLen = 8

type (
	URLs    map[models.ShortURL]models.FullURL
	Storage struct {
		URLs URLs
	}
)

// New инициализация Storage
func New() *Storage {
	return &Storage{
		URLs: make(URLs),
	}
}

// Get возвращает полный урл по сокращенному.
func (s *Storage) Get(sURL models.ShortURL) models.FullURL {
	if !s.IsExist(sURL) {
		return ""
	}
	return s.URLs[sURL]
}

// Set сохраняет в память полный УРЛ и соответствующий ему короткий УРЛ
func (s *Storage) Set(fURL models.FullURL) (models.ShortURL, error) {
	randChars := service.GenerateRandChars(shortURLLen)
	shortURL := models.ShortURL(randChars)

	if !s.IsExist(shortURL) {
		s.URLs[shortURL] = fURL
		return shortURL, nil
	}

	return "", URLExistsError{string(fURL)}
}

// IsExist проверяет сохранен ли в памяти короткий УРЛ
func (s *Storage) IsExist(sURL models.ShortURL) bool {
	if _, ok := s.URLs[sURL]; !ok {
		return false
	}
	return true
}
