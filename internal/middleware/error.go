package middleware

type JWTErr struct {
	Inner error
	Msg   string
}

func NewJWTErr(inner error, msg string) JWTErr {
	return JWTErr{
		Inner: inner,
		Msg:   msg,
	}
}

func (j JWTErr) Error() string {
	return j.Msg
}
