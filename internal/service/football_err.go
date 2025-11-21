package service

import "fmt"

type FootballError struct {
	inner  error
	msg    string
	status int
}

func NewFootballError(inner error, msg string, status int) error {
	return &FootballError{
		inner:  inner,
		msg:    msg,
		status: status,
	}
}

func (err *FootballError) Error() string {
	if err.inner != nil {
		return fmt.Sprintf("%s: %s", err.msg, err.inner)
	}
	return err.msg
}

func (err *FootballError) Unwrap() error {
	return err.inner
}

func (err *FootballError) Message() string {
	return err.msg
}

func (err *FootballError) Status() int {
	return err.status
}
