package middleware

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v5"
)

// RBAC valida se o papel (Role) extraído do contexto do usuário condiz com as permissões exigidas.
func RBAC(allowedRoles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			role, ok := c.Get("user_role").(string)
			if !ok || role == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "não autorizado"})
			}

			allowed := false
			for _, r := range allowedRoles {
				if role == r {
					// Dupla validação de segurança para SUPER_ADMIN (verifica se o tenant_id é o tenant mestre da plataforma)
					if r == "SUPER_ADMIN" {
						tenantID, _ := c.Get("tenant_id").(string)
						adminTenantID := os.Getenv("ADMIN_TENANT_ID")
						if adminTenantID == "" || tenantID != adminTenantID {
							// Ignora a autorização caso o tenant não seja o administrativo
							continue
						}
					}
					allowed = true
					break
				}
			}

			if !allowed {
				return c.JSON(http.StatusForbidden, map[string]string{"error": "permissão insuficiente para este recurso"})
			}

			return next(c)
		}
	}
}
