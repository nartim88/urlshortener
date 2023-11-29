package storage

import (
	"context"

	"github.com/nartim88/urlshortener/internal/pkg/models"
)

const shortURLLen = 8

// Storage базовый интерфейс для работы с данными
type Storage interface {
	// Get возвращает полный урл по строковому идентификатору
	Get(ctx context.Context, sID models.ShortenID) (*models.FullURL, error)
	// Set сохраняет в базу полный УРЛ и соответствующий ему строковой идентификатор
	Set(ctx context.Context, fURL models.FullURL) (*models.ShortenID, error)
}

// StorageWithService расширенный интерфейс для работы с данными, подходящий для работы с
// внешними сервисами, такими как бд
type StorageWithService interface {
	Storage
	// Bootstrap создание необходимых сущностей для начала работы с сервисом:
	// таблиц и индексов в бд, файлов и пр.
	Bootstrap(ctx context.Context) error
	// Close закрытие существующих соединений с внешними сервисами
	Close(ctx context.Context) error
}
