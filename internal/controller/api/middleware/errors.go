package middleware

import (
	"fmt"
)

type NoCookieWithTokenErr struct {
	Err error
}

func (e *NoCookieWithTokenErr) Error() string {
	return fmt.Sprintf("cookie with the token is invalid or not provided: %v", e.Err)
}

func (e *NoCookieWithTokenErr) Unwrap() error {
	return e.Err
}

func NewNoCookieWithTokenErr(err error) error {
	return &NoCookieWithTokenErr{err}
}
