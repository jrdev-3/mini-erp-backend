package auth

import (
	"context"
	"errors"
	"strings"
	"time"
)

// User representa a entidade de domínio de Usuários do Mini ERP.
type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Oculta a senha em serializações JSON por segurança
	Role         string    `json:"role"` // ADMIN ou USER
	IsActive     bool      `json:"is_active"`
	TenantID     string    `json:"tenant_id"`
	CriadoEm     time.Time `json:"criado_em"`
	AtualizadoEm time.Time `json:"atualizado_em"`
}

// Erros de domínio do módulo auth
var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already in use")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserInactive       = errors.New("user account is inactive")
)

// SystemAnalytics encapsula as métricas de saúde da plataforma para o Super Admin.
type SystemAnalytics struct {
	TotalUsers    int64 `json:"total_users"`
	ActiveUsers   int64 `json:"active_users"`
	InactiveUsers int64 `json:"inactive_users"`
}

// Repository define a interface de persistência para o módulo auth.
// O parâmetro tenantID nos métodos de consulta e escrita assegura a prevenção a BOLA no banco.
type Repository interface {
	Create(ctx context.Context, u *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id string, tenantID string) (*User, error)
	UpdateActiveStatus(ctx context.Context, id string, tenantID string, active bool) error
	UpdateActiveStatusGlobal(ctx context.Context, id string, active bool) error
	GetSystemAnalytics(ctx context.Context) (*SystemAnalytics, error)
	ListAll(ctx context.Context) ([]*User, error)
	ListByTenant(ctx context.Context, tenantID string) ([]*User, error)
}

// CreateUserRequest representa os dados necessários para o administrador cadastrar um novo funcionário.
type CreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Validate faz a validação sintática do request de criação de funcionário.
func (r *CreateUserRequest) Validate() error {
	if r.Email == "" {
		return errors.New("e-mail é obrigatório")
	}
	if !strings.Contains(r.Email, "@") || !strings.Contains(r.Email, ".") {
		return errors.New("e-mail inválido")
	}
	if r.Password == "" {
		return errors.New("senha é obrigatória")
	}
	if len(r.Password) < 6 {
		return errors.New("a senha deve ter no mínimo 6 caracteres")
	}
	return nil
}
