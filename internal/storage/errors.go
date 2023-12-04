package storage

import (
	"fmt"

	"github.com/nartim88/urlshortener/internal/models"
)

// URLExistsError урл уже существует в базе
type URLExistsError struct {
	OriginalURL models.FullURL
	SID         models.ShortenID
}

func (u URLExistsError) Error() string {
	return fmt.Sprintf("'%s' is already saved", u.OriginalURL)
}
