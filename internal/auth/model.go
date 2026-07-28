package auth

import (
	"context"
	"errors"
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

// Repository define a interface de persistência para o módulo auth.
// O parâmetro tenantID nos métodos de consulta e escrita assegura a prevenção a BOLA no banco.
type Repository interface {
	Create(ctx context.Context, u *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id string, tenantID string) (*User, error)
	UpdateActiveStatus(ctx context.Context, id string, tenantID string, active bool) error
}
