package service

import "fmt"

type CredentialError struct {
	inner error
	msg   string
}

func NewCredentialError(inner error, msg string) error {
	return &CredentialError{
		inner: inner,
		msg:   msg,
	}
}

func (err *CredentialError) Error() string {
	if err.inner != nil {
		return fmt.Sprintf("%s: %s", err.msg, err.inner)
	}
	return err.msg
}

func (err *CredentialError) Unwrap() error {
	return err.inner
}

func (err *CredentialError) Message() string {
	return err.msg
}
