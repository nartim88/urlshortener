package storage

import (
	"github.com/nartim88/urlshortener/internal/pkg/models"
)

const shortURLLen = 8

type Storage interface {
	Get(sURL models.ShortURL) (*models.FullURL, error)
	Set(fURL models.FullURL) (*models.ShortURL, error)
}
