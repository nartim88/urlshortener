package models

import "github.com/google/uuid"

type (
	// FullURL исходный url, переданный для сокращения
	FullURL string
	// ShortenID строковый идентификатор для формирования сокращенного урла
	ShortenID string
	// CorrelationID строковый идентификатор для отслеживания запроса
	CorrelationID string
)

// FileJSONEntry структура для записи данных в файл в json формате
type FileJSONEntry struct {
	ID       *uuid.UUID `json:"id"`
	ShortURL ShortenID  `json:"short_url"`
	FullURL  FullURL    `json:"full_url"`
}
