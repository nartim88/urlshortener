package storage

import (
	"github.com/nartim88/urlshortener/internal/pkg/models"
)

const shortURLLen = 8

type Storage interface {
	// Get возвращает полный урл по сокращенному.
	Get(sURL models.ShortenID) (*models.FullURL, error)
	// Set сохраняет в базу полный УРЛ и соответствующий ему короткий УРЛ
	Set(fURL models.FullURL) (*models.ShortenID, error)
}
