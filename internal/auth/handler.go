package auth

import (
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v5"
)

// Handler gerencia as requisições HTTP para autenticação e gestão de usuários.
type Handler struct {
	service Service
}

// NewHandler cria uma nova instância do Handler de autenticação.
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// Register godoc
// @Summary      Registrar novo Tenant e Administrador
// @Description  Cria uma nova conta de inquilino (empresa) e o primeiro usuário com cargo de ADMIN.
// @Tags         autenticacao
// @Accept       json
// @Produce      json
// @Param        request body RegisterRequest true "Dados de cadastro"
// @Success      201  {object}  User
// @Failure      400  {object}  map[string]string "Dados inválidos ou e-mail já cadastrado"
// @Failure      500  {object}  map[string]string "Erro interno no servidor"
// @Router       /api/v1/auth/register [post]
func (h *Handler) Register(c *echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "corpo da requisição inválido"})
	}

	user, err := h.service.Register(c.Request().Context(), &req)
	if err != nil {
		if errors.Is(err, ErrEmailAlreadyExists) || strings.Contains(err.Error(), "obrigatório") || strings.Contains(err.Error(), "inválido") || strings.Contains(err.Error(), "mínimo") {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "erro interno no servidor"})
	}

	return c.JSON(http.StatusCreated, user)
}

// Login godoc
// @Summary      Realizar autenticação de usuário
// @Description  Valida as credenciais do usuário e retorna um token JWT de acesso offline.
// @Tags         autenticacao
// @Accept       json
// @Produce      json
// @Param        request body LoginRequest true "Credenciais de login"
// @Success      200  {object}  LoginResponse
// @Failure      400  {object}  map[string]string "E-mail ou senha obrigatórios"
// @Failure      401  {object}  map[string]string "Credenciais inválidas"
// @Failure      403  {object}  map[string]string "Usuário inativo"
// @Failure      500  {object}  map[string]string "Erro interno no servidor"
// @Router       /api/v1/auth/login [post]
func (h *Handler) Login(c *echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "corpo da requisição inválido"})
	}

	resp, err := h.service.Login(c.Request().Context(), &req)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
		}
		if errors.Is(err, ErrUserInactive) {
			return c.JSON(http.StatusForbidden, map[string]string{"error": err.Error()})
		}
		if strings.Contains(err.Error(), "obrigatório") {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "erro interno no servidor"})
	}

	return c.JSON(http.StatusOK, resp)
}

type toggleRequest struct {
	Active bool `json:"active"`
}

// ToggleActive godoc
// @Summary      Ativar ou desativar um usuário
// @Description  Altera o status de atividade de um usuário dentro de seu próprio tenant (Requer ADMIN).
// @Tags         administracao
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "ID do Usuário"
// @Param        request body toggleRequest true "Status ativo (active)"
// @Success      200  {object}  map[string]string "Status atualizado com sucesso"
// @Failure      400  {object}  map[string]string "Parâmetros inválidos"
// @Failure      401  {object}  map[string]string "Não autorizado"
// @Failure      403  {object}  map[string]string "Acesso proibido"
// @Failure      404  {object}  map[string]string "Usuário não encontrado"
// @Failure      500  {object}  map[string]string "Erro interno no servidor"
// @Security     ApiKeyAuth
// @Router       /api/v1/admin/users/{id}/toggle [patch]
func (h *Handler) ToggleActive(c *echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "id do usuário é obrigatório"})
	}

	// Recuperar o tenant_id extraído pelo middleware JWT no contexto com segurança BOLA
	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "não autorizado"})
	}

	var req toggleRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "corpo da requisição inválido"})
	}

	err := h.service.ToggleUserActive(c.Request().Context(), id, tenantID, req.Active)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			// Prevenção a BOLA: retornar 404 caso o id não pertença ao tenant_id do executor
			return c.JSON(http.StatusNotFound, map[string]string{"error": "usuário não encontrado"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "erro interno no servidor"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "status do usuário atualizado com sucesso"})
}

// ToggleActiveGlobal godoc
// @Summary      Ativar ou desativar qualquer usuário (Super Admin)
// @Description  Altera o status de atividade de qualquer usuário de qualquer tenant na plataforma (Requer SUPER_ADMIN).
// @Tags         administracao-plataforma
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "ID do Usuário"
// @Param        request body toggleRequest true "Status ativo (active)"
// @Success      200  {object}  map[string]string "Status do usuário atualizado com sucesso"
// @Failure      400  {object}  map[string]string "Parâmetros inválidos"
// @Failure      401  {object}  map[string]string "Não autorizado"
// @Failure      403  {object}  map[string]string "Acesso proibido"
// @Failure      404  {object}  map[string]string "Usuário não encontrado"
// @Failure      500  {object}  map[string]string "Erro interno no servidor"
// @Security     ApiKeyAuth
// @Router       /api/v1/system/users/{id}/toggle [patch]
func (h *Handler) ToggleActiveGlobal(c *echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "id do usuário é obrigatório"})
	}

	var req toggleRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "corpo da requisição inválido"})
	}

	err := h.service.ToggleUserActiveGlobal(c.Request().Context(), id, req.Active)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "usuário não encontrado"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "erro interno no servidor"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "status do usuário atualizado com sucesso"})
}

// GetSystemAnalytics godoc
// @Summary      Obter métricas de saúde da plataforma (Super Admin)
// @Description  Retorna o total de usuários ativos, inativos e empresas cadastradas no ecossistema (Requer SUPER_ADMIN).
// @Tags         administracao-plataforma
// @Accept       json
// @Produce      json
// @Success      200  {object}  SystemAnalytics
// @Failure      401  {object}  map[string]string "Não autorizado"
// @Failure      403  {object}  map[string]string "Acesso proibido"
// @Failure      500  {object}  map[string]string "Erro interno no servidor"
// @Security     ApiKeyAuth
// @Router       /api/v1/system/analytics [get]
func (h *Handler) GetSystemAnalytics(c *echo.Context) error {
	analytics, err := h.service.GetSystemAnalytics(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "erro interno no servidor"})
	}

	return c.JSON(http.StatusOK, analytics)
}

// ListAllUsers godoc
// @Summary      Listar todos os usuários da plataforma (Super Admin)
// @Description  Retorna a lista completa de todos os usuários cadastrados no ecossistema (Requer SUPER_ADMIN).
// @Tags         administracao-plataforma
// @Accept       json
// @Produce      json
// @Success      200  {array}   User
// @Failure      401  {object}  map[string]string "Não autorizado"
// @Failure      403  {object}  map[string]string "Acesso proibido"
// @Failure      500  {object}  map[string]string "Erro interno no servidor"
// @Security     ApiKeyAuth
// @Router       /api/v1/system/users [get]
func (h *Handler) ListAllUsers(c *echo.Context) error {
	users, err := h.service.ListAllUsers(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "erro interno no servidor"})
	}

	return c.JSON(http.StatusOK, users)
}
