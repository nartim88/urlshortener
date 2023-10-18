package storage

import "fmt"

type URLExistsError struct {
	URL string
}

func (u URLExistsError) Error() string {
	return fmt.Sprintf("'%s' already saved", u.URL)
}
