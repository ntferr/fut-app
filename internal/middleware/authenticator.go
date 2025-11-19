package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

const (
	BearerPrefix = "Bearer"
	MethodHS256  = "HS256"
)

func VerifyToken(tokenValue string, secretKey string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenValue, func(token *jwt.Token) (any, error) {
		switch token.Method.Alg() {
		case MethodHS256:
			return []byte(secretKey), nil
		default:
			return nil, NewJWTErr(nil, "token must be signed with HMAC method")
		}
	})
	if err != nil {
		return nil, NewJWTErr(err, "invalid or expired token")
	}
	if !token.Valid {
		return nil, NewJWTErr(nil, "invalid token: final check failed")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, NewJWTErr(nil, "invalid token: claims are not MapClaims")
	}

	return claims, nil
}

func RemoveBearerPrefix(authHeader string) (string, error) {
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", NewJWTErr(nil, "invalid authorization type: failed to split bearer from auth header")
	}
	tokenString := parts[1]
	return tokenString, nil
}

func JWTMiddleware(secretKey string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			authHeader := c.Request().Header.Get("Authorization")
			token, err := RemoveBearerPrefix(authHeader)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, err)
			}

			claims, err := VerifyToken(token, secretKey)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "first verify: invalid token")
			}

			user, ok := claims["username"].(string)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "final verify: invalid token")
			}

			c.Set("user", user)
			return next(c)
		}
	}
}
