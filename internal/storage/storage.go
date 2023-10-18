package storage

import (
	"github.com/nartim88/urlshortener/internal/service"
)

const shortURLLen = 8

type FullURL string
type ShortURL string

type Storage struct {
	URLs map[ShortURL]FullURL
}

func New() *Storage {
	store := &Storage{
		URLs: make(map[ShortURL]FullURL),
	}
	return store
}

// Get возвращает полный урл по сокращенному.
func (s *Storage) Get(shortURL ShortURL) FullURL {
	if !s.IsExist(shortURL) {
		return ""
	}
	return s.URLs[shortURL]
}

// Set сохраняет в память полный УРЛ и соответствующий ему короткий УРЛ
func (s *Storage) Set(fullURL FullURL) error {
	randChars := service.GenerateRandChars(shortURLLen)
	shortURL := ShortURL(randChars)

	if !s.IsExist(shortURL) {
		s.URLs[shortURL] = fullURL
		return nil
	}

	return URLExistsError{string(fullURL)}
}

// IsExist проверяет сохранен ли в памяти короткий УРЛ
func (s *Storage) IsExist(shortURL ShortURL) bool {
	if _, ok := s.URLs[shortURL]; !ok {
		return false
	}
	return true
}
