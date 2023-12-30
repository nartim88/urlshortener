package storage

import (
	"errors"
	"fmt"

	"github.com/nartim88/urlshortener/internal/models"
)

// ErrURLExists урл уже существует в базе
type ErrURLExists struct {
	OriginalURL models.FullURL
	SID         models.ShortenID
}

func (u ErrURLExists) Error() string {
	return fmt.Sprintf("'%s' is already saved", u.OriginalURL)
}

var ErrURLDeleted = errors.New("url deleted")
