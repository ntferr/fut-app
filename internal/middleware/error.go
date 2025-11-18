package middleware

type JWTErr struct {
	inner       error
	description string
}

func NewJWTErr(inner error, description string) JWTErr {
	return JWTErr{
		inner:       inner,
		description: description,
	}
}

func (j JWTErr) Error() string {
	return j.description
}
