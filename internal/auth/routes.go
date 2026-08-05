package auth

import (
	"github.com/labstack/echo/v5"
)

// RegisterRoutes registra os endpoints do módulo auth no roteador do Echo.
// Recebe o grupo base (ex: /api/v1), o handler e os middlewares injetados para autenticação, RBAC de ADMIN e SUPER_ADMIN.
func RegisterRoutes(g *echo.Group, h *Handler, authMiddleware echo.MiddlewareFunc, adminMiddleware echo.MiddlewareFunc, superAdminMiddleware echo.MiddlewareFunc) {
	// Endpoints Públicos
	g.POST("/auth/register", h.Register)
	g.POST("/auth/login", h.Login)

	// Endpoints Administrativos Protegidos (Autenticação + RBAC de ADMIN)
	g.PATCH("/admin/users/:id/toggle", h.ToggleActive, authMiddleware, adminMiddleware)
	g.POST("/admin/users", h.CreateEmployee, authMiddleware, adminMiddleware)
	g.GET("/admin/users", h.ListEmployees, authMiddleware, adminMiddleware)

	// Endpoints de Sistema (Exclusivo Super Admin - Autenticação + RBAC de SUPER_ADMIN)
	g.PATCH("/system/users/:id/toggle", h.ToggleActiveGlobal, authMiddleware, superAdminMiddleware)
	g.GET("/system/analytics", h.GetSystemAnalytics, authMiddleware, superAdminMiddleware)
	g.GET("/system/users", h.ListAllUsers, authMiddleware, superAdminMiddleware)
}
