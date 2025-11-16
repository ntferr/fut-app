package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

const (
	BearerPrefix = "Bearer"
)

func VerifyToken(tokenValue string, secretKey string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(RemoveBearerPrefix(tokenValue), func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); ok {
			return []byte(secretKey), nil
		}
		// BadRequest token enviado está correto
		return nil, errors.New("invalid token")
	})
	if err != nil {
		// Unuathorized token não está válido
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func RemoveBearerPrefix(token string) string {
	if strings.HasPrefix(token, BearerPrefix+" ") {
		token = strings.Trim(BearerPrefix, token)
	}
	return token
}

func JWTMiddleware(secretKey string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			authHeader := c.Request().Header.Get("Authorization")
			claims, err := VerifyToken(authHeader, secretKey)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Token inválido")
			}

			user, ok := claims["user"].(string)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "Token inválido")
			}

			// Adicionar as claims ao contexto para uso posterior
			c.Set("user", user)
			c.Set("claims", claims)
			return next(c)
		}

		// return echo.NewHTTPError(http.StatusUnauthorized, "Token inválido")
	}
}
