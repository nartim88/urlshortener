package storage

import (
	"github.com/nartim88/urlshortener/internal/pkg/service"
)

const shortURLLen = 8

type FullURL string
type ShortURL string
type URLs map[ShortURL]FullURL

type Storage struct {
	URLs URLs
}

// New инициализация Storage
func New() *Storage {
	return &Storage{
		URLs: make(URLs),
	}
}

// Get возвращает полный урл по сокращенному.
func (s *Storage) Get(sURL ShortURL) FullURL {
	if !s.IsExist(sURL) {
		return ""
	}
	return s.URLs[sURL]
}

// Set сохраняет в память полный УРЛ и соответствующий ему короткий УРЛ
func (s *Storage) Set(fURL FullURL) (ShortURL, error) {
	randChars := service.GenerateRandChars(shortURLLen)
	shortURL := ShortURL(randChars)

	if !s.IsExist(shortURL) {
		s.URLs[shortURL] = fURL
		return shortURL, nil
	}

	return "", URLExistsError{string(fURL)}
}

// IsExist проверяет сохранен ли в памяти короткий УРЛ
func (s *Storage) IsExist(sURL ShortURL) bool {
	if _, ok := s.URLs[sURL]; !ok {
		return false
	}
	return true
}
