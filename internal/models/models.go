package models

import (
	"github.com/google/uuid"
)

type (
	// FullURL исходный url, переданный для сокращения
	FullURL string
	// ShortenID строковый идентификатор для формирования сокращенного урла
	ShortenID string
	// CorrelationID строковый идентификатор для отслеживания запроса
	CorrelationID string
	// ShortURL короткий урл
	ShortURL string
)

// FileJSONEntry структура для записи данных в файл в json формате
type FileJSONEntry struct {
	ID        *uuid.UUID `json:"id"`
	ShortenID ShortenID  `json:"shorten_id"`
	FullURL   FullURL    `json:"full_url"`
}

type User struct {
	UserID string `json:"user_id"`
}

type ShortAndFullURLs struct {
	ShortURL ShortURL `json:"short_url"`
	FullURL  FullURL  `json:"original_url"`
}

type UserIDCtxKey string

type URLsModel struct {
	ID          string    `json:"id,omitempty"`
	ShortenID   ShortenID `json:"shorten_id,omitempty"`
	FullURL     FullURL   `json:"original_url,omitempty"`
	ShortURL    ShortURL  `json:"short_url,omitempty"`
	DeletedFlag bool      `json:"deleted_flag,omitempty"`
}
