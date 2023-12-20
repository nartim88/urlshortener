package storage

import (
	"context"

	"github.com/nartim88/urlshortener/internal/models"
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

func (s *MemStorage) Get(ctx context.Context, sID models.ShortenID) (*models.FullURL, error) {
	if !s.isExist(sID) {
		return nil, nil
	}
	fURL := s.Memory[sID]
	return &fURL, nil
}

func (s *MemStorage) Set(ctx context.Context, fURL models.FullURL) (*models.ShortenID, error) {
	randChars := GenerateRandChars(shortURLLen)
	shortURL := models.ShortenID(randChars)

	if s.isExist(shortURL) {
		return nil, ErrURLExists{
			OriginalURL: fURL,
			SID:         shortURL,
		}
	}

	s.Memory[shortURL] = fURL

	return &shortURL, nil
}

func (s *MemStorage) ListURLs(ctx context.Context, u models.User) ([]SIDAndFullURL, error) {
	return nil, nil
}

// isExist проверяет сохранен ли в памяти короткий УРЛ
func (s *MemStorage) isExist(sID models.ShortenID) bool {
	_, ok := s.Memory[sID]
	return ok
}

func (s *MemStorage) MarkAsDeletedByID(ctx context.Context, IDs []models.ShortenID) error {
	return nil
}

func (s *MemStorage) SetBatch(ctx context.Context, fURLs []models.FullURL) (map[models.FullURL]models.ShortenID, error) {
	return nil, nil
}
