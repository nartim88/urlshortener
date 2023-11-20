package models

import "github.com/google/uuid"

type (
	// FullURL исходный url, переданный для сокращения
	FullURL string
	// ShortURL строковый идентификатор для формирования сокращенного урла
	ShortURL string
	// CorrelationID строковый идентификатор для отслеживания запроса
	CorrelationID string
)

// FileJSONEntry структура для записи данных в файл в json формате
type FileJSONEntry struct {
	ID       *uuid.UUID `json:"id"`
	ShortURL ShortURL   `json:"short_url"`
	FullURL  FullURL    `json:"full_url"`
}
