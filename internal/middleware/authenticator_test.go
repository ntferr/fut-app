//go:build unit

package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/fut-app/internal/model"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJWTMiddleware(t *testing.T) {
	secretKey := "test-secret-key"
	middleware := JWTMiddleware(secretKey)

	tests := []struct {
		name           string
		setupRequest   func() *http.Request
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "should return unauthorized when no authorization header",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("GET", "/", nil)
				return req
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Token de autorização necessário",
		},
		{
			name: "should return unauthorized when invalid bearer format",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("GET", "/", nil)
				req.Header.Set("Authorization", "InvalidToken")
				return req
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Formato de token inválido. Use: Bearer <token>",
		},
		{
			name: "should return unauthorized when token is expired",
			setupRequest: func() *http.Request {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, &model.CustomClaims{
					AuthRequest: model.AuthRequest{
						User: "test-user",
					},
					RegisteredClaims: jwt.RegisteredClaims{
						ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)), // Token expirado
					},
				})
				tokenString, _ := token.SignedString([]byte(secretKey))

				req := httptest.NewRequest("GET", "/", nil)
				req.Header.Set("Authorization", "Bearer "+tokenString)
				return req
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Token expirado",
		},
		{
			name: "should return unauthorized when invalid token signature",
			setupRequest: func() *http.Request {
				// Token assinado com chave diferente
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, &model.CustomClaims{
					AuthRequest: model.AuthRequest{
						User: "test-user",
					},
				})
				tokenString, _ := token.SignedString([]byte("wrong-secret-key"))

				req := httptest.NewRequest("GET", "/", nil)
				req.Header.Set("Authorization", "Bearer "+tokenString)
				return req
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Token inválido",
		},
		{
			name: "should successfully authenticate with valid token",
			setupRequest: func() *http.Request {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, &model.CustomClaims{
					AuthRequest: model.AuthRequest{
						User:     "test-user",
						Password: "test-pass",
					},
					RegisteredClaims: jwt.RegisteredClaims{
						ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
					},
				})
				tokenString, _ := token.SignedString([]byte(secretKey))

				req := httptest.NewRequest("GET", "/", nil)
				req.Header.Set("Authorization", "Bearer "+tokenString)
				return req
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "success",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := tt.setupRequest()
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Handler mock que será chamado se o middleware passar
			mockHandler := func(c echo.Context) error {
				return c.String(http.StatusOK, "success")
			}

			// Executar o middleware
			handler := middleware(mockHandler)
			err := handler(c)

			if tt.expectedStatus != http.StatusOK {
				require.Error(t, err)
				httpErr, ok := err.(*echo.HTTPError)
				require.True(t, ok, "Error should be of type *echo.HTTPError")
				assert.Equal(t, tt.expectedStatus, httpErr.Code)
				assert.Contains(t, httpErr.Message, tt.expectedBody)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, rec.Code)
				assert.Equal(t, tt.expectedBody, rec.Body.String())

				// Verificar se os dados do usuário foram setados no contexto
				user := c.Get("user")
				assert.NotNil(t, user)

				claims := c.Get("claims")
				assert.NotNil(t, claims)
			}
		})
	}
}

func TestJWTMiddleware_InvalidSigningMethod(t *testing.T) {
	secretKey := "test-secret-key"
	middleware := JWTMiddleware(secretKey)

	// Token com método de assinatura inválido (RS256 em vez de HS256)
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, &model.CustomClaims{
		AuthRequest: model.AuthRequest{
			User: "test-user",
		},
	})

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+token.Raw) // Raw token não assinado

	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)

	mockHandler := func(c echo.Context) error {
		return c.String(http.StatusOK, "success")
	}

	handler := middleware(mockHandler)
	err := handler(c)

	require.Error(t, err)
	httpErr := err.(*echo.HTTPError)
	assert.Equal(t, http.StatusUnauthorized, httpErr.Code)
	assert.Contains(t, httpErr.Message, "Método de assinatura inválido")
}

func TestJWTMiddleware_InvalidTokenStructure(t *testing.T) {
	secretKey := "test-secret-key"
	middleware := JWTMiddleware(secretKey)

	// Token inválido
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer invalid-token-string")

	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)

	mockHandler := func(c echo.Context) error {
		return c.String(http.StatusOK, "success")
	}

	handler := middleware(mockHandler)
	err := handler(c)

	require.Error(t, err)
	httpErr := err.(*echo.HTTPError)
	assert.Equal(t, http.StatusUnauthorized, httpErr.Code)
	assert.Contains(t, httpErr.Message, "Token inválido")
}
