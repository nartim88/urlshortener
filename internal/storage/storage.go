package storage

import (
	"context"
	"math/rand"
	"time"

	"github.com/nartim88/urlshortener/config"
	"github.com/nartim88/urlshortener/internal/models"
)

const shortURLLen = 8

type Row struct {
	ShortURL string
	FullURL  models.FullURL
}

type SIDAndFullURL struct {
	models.ShortenID
	models.FullURL
}

// Storage базовый интерфейс для работы с данными
type Storage interface {
	// Get возвращает полный урл по строковому идентификатору
	Get(ctx context.Context, sID models.ShortenID) (*models.FullURL, error)
	// Set сохраняет в базу полный УРЛ и возвращает соответствующий ему строковой идентификатор
	Set(ctx context.Context, fURL models.FullURL) (*models.ShortenID, error)
	// ListURLs возвращает все записи переданного пользователя
	ListURLs(ctx context.Context, u models.User) ([]SIDAndFullURL, error)
	// MarkAsDeletedByID помечает урлы как удаленные по строковым идентификаторам
	MarkAsDeletedByID(ctx context.Context, IDs []models.ShortenID) error
	// SetBatch сохраняет в базу батч урлов по аналогии с Set
	SetBatch(ctx context.Context, fURLs []models.FullURL) (map[models.FullURL]models.ShortenID, error)
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

// GenerateRandChars возвращает строку из случайных символов.
func GenerateRandChars(n int) []byte {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	shortKey := make([]byte, n)
	for i := 0; i < n; i++ {
		shortKey[i] = config.Charset[rnd.Intn(len(config.Charset))]
	}
	return shortKey
}
