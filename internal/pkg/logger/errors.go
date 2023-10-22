package logger

import "fmt"

type LogLevelError struct {
	msg string
}

func (e LogLevelError) Error() string {
	return fmt.Sprint(e.msg)
}
