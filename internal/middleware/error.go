package middleware

import "fmt"

type JWTError struct {
	inner error
	msg   string
}

func NewJWTErr(inner error, msg string) error {
	return &JWTError{
		inner: inner,
		msg:   msg,
	}
}

func (err *JWTError) Error() string {
	if err.inner != nil {
		return fmt.Sprintf("%s: %s", err.msg, err.inner)
	}
	return err.msg
}

func (err *JWTError) Unwrap() error {
	return err.inner
}

func (err *JWTError) Message() string {
	return err.msg
}
