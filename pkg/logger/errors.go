package logger

import "fmt"

type LogLevelError struct {
	msg  string
	data any
}

func (e LogLevelError) Error() string {
	return fmt.Sprintf("%s: %v", e.msg, e.data)
}
