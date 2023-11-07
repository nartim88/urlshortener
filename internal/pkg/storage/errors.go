package storage

import "fmt"

// URLExistsError урл уже существует в базе
type URLExistsError struct {
	URL string
}

func (u URLExistsError) Error() string {
	return fmt.Sprintf("'%s' already saved", u.URL)
}
