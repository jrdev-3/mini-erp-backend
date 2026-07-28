package auth

import (
	"github.com/labstack/echo/v5"
)

// RegisterRoutes registra os endpoints do módulo auth no roteador do Echo.
// Recebe o grupo base (ex: /api/v1), o handler e os middlewares injetados para autenticação e RBAC.
func RegisterRoutes(g *echo.Group, h *Handler, authMiddleware echo.MiddlewareFunc, adminMiddleware echo.MiddlewareFunc) {
	// Endpoints Públicos
	g.POST("/auth/register", h.Register)
	g.POST("/auth/login", h.Login)

	// Endpoints Administrativos Protegidos (Autenticação + RBAC de ADMIN)
	g.PATCH("/admin/users/:id/toggle", h.ToggleActive, authMiddleware, adminMiddleware)
}
