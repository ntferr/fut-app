package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

const (
	BearerPrefix = "Bearer"
)

func VerifyToken(tokenValue string, secretKey string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenValue, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de assinatura inesperado: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
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

func RemoveBearerPrefix(authHeader string) (string, error) {
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("failed to split bearer from auth header")
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
				return echo.NewHTTPError(http.StatusUnauthorized, "Token inválido")
			}

			user, ok := claims["username"].(string)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "Token inválido")
			}

			// Adicionar as claims ao contexto para uso posterior
			c.Set("user", user)
			c.Set("claims", claims)
			return next(c)
		}
	}
}
