package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/fut-app/internal/model"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func JWTMiddleware(secretKey string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Token de autorização necessário")
			}

			if !strings.HasPrefix(authHeader, "Bearer ") {
				return echo.NewHTTPError(http.StatusUnauthorized, "Formato de token inválido. Use: Bearer <token>")
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			// Validar o token
			token, err := jwt.ParseWithClaims(tokenString, &model.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
				// Verificar o método de assinatura
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, echo.NewHTTPError(http.StatusUnauthorized, "Método de assinatura inválido")
				}
				return []byte(secretKey), nil
			})

			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Token inválido: "+err.Error())
			}

			// Verificar se o token é válido e extrair claims
			if claims, ok := token.Claims.(*model.CustomClaims); ok && token.Valid {
				if claims.ExpiresAt != nil {
					if time.Now().After(claims.ExpiresAt.Time) {
						return echo.NewHTTPError(http.StatusUnauthorized, "Token expirado")
					}
				}
				// Adicionar as claims ao contexto para uso posterior
				c.Set("user", claims.AuthRequest)
				c.Set("claims", claims)
				return next(c)
			}

			return echo.NewHTTPError(http.StatusUnauthorized, "Token inválido")
		}
	}
}
