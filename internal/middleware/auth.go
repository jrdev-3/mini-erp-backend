package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
)

// Auth retorna um middleware que valida de forma offline a assinatura e expiração do token JWT.
// Extrai as claims e as injeta no contexto para uso nos handlers e isolamento multi-tenant (BOLA).
func Auth(jwtSecret []byte) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "cabeçalho de autorização ausente"})
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "formato de autorização inválido (esperado Bearer <token>)"})
			}

			tokenString := parts[1]
			token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, echo.NewHTTPError(http.StatusUnauthorized, "método de assinatura inválido")
				}
				return jwtSecret, nil
			})

			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "token de acesso inválido ou expirado"})
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "claims do token inválidas"})
			}

			tenantID, _ := claims["tenant_id"].(string)
			role, _ := claims["role"].(string)
			userID, _ := claims["sub"].(string)

			if tenantID == "" || role == "" || userID == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "dados de autorização ausentes no token"})
			}

			// Salvar claims no contexto do Echo de forma explícita
			c.Set("tenant_id", tenantID)
			c.Set("user_role", role)
			c.Set("user_id", userID)

			return next(c)
		}
	}
}
