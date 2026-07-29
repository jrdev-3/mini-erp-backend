package auth

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Service define os casos de uso para a camada de aplicação de autenticação.
type Service interface {
	Register(ctx context.Context, req *RegisterRequest) (*User, error)
	Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error)
	ToggleUserActive(ctx context.Context, id string, tenantID string, active bool) error
	ToggleUserActiveGlobal(ctx context.Context, id string, active bool) error
	GetSystemAnalytics(ctx context.Context) (*SystemAnalytics, error)
	ListAllUsers(ctx context.Context) ([]*User, error)
}

type service struct {
	repo      Repository
	jwtSecret []byte
}

// NewService cria uma nova instância de Service.
func NewService(repo Repository, jwtSecret string) Service {
	return &service{
		repo:      repo,
		jwtSecret: []byte(jwtSecret),
	}
}

// Register cria um novo tenant único e registra o primeiro usuário ADMIN correspondente.
func (s *service) Register(ctx context.Context, req *RegisterRequest) (*User, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// 1. Verificar se o e-mail já está cadastrado no banco globalmente
	_, err := s.repo.GetByEmail(ctx, req.Email)
	if err == nil {
		return nil, ErrEmailAlreadyExists
	}
	if err != nil && !errors.Is(err, ErrUserNotFound) {
		return nil, err
	}

	// 2. Criptografar a senha com bcrypt (custo padrão recomendado de 10)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 3. Gerar um novo ID de inquilino (Tenant ID) seguro e único via UUID v4
	tenantID := uuid.NewString()

	// 4. Instanciar o usuário administrador
	newUser := &User{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Role:         "ADMIN",
		IsActive:     true,
		TenantID:     tenantID,
	}

	// 5. Salvar na persistência
	if err := s.repo.Create(ctx, newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}

// Login autentica as credenciais, valida se a conta está ativa e gera um token JWT de acesso offline.
func (s *service) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// 1. Obter usuário pelo e-mail
	u, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			// Mitiga enumeração de e-mails respondendo com credenciais inválidas genéricas
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	// 2. Comparar as senhas
	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// 3. Validar se a conta do usuário está ativa no ERP
	if !u.IsActive {
		return nil, ErrUserInactive
	}

	// 4. Gerar o token JWT offline embutindo sub (ID do usuário), tenant_id e role nas claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":       u.ID,
		"tenant_id": u.TenantID,
		"role":      u.Role,
		"exp":       time.Now().Add(time.Hour * 24).Unix(), // Expira em 24 horas
	})

	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token: tokenString,
		User:  u,
	}, nil
}

// ToggleUserActive ativa ou desativa um usuário do tenant (validação BOLA imposta no nível do repositório).
func (s *service) ToggleUserActive(ctx context.Context, id string, tenantID string, active bool) error {
	if id == "" || tenantID == "" {
		return errors.New("id do usuário e tenant_id são obrigatórios")
	}
	return s.repo.UpdateActiveStatus(ctx, id, tenantID, active)
}

// ToggleUserActiveGlobal ativa ou desativa qualquer usuário do sistema (Exclusivo Super Admin).
func (s *service) ToggleUserActiveGlobal(ctx context.Context, id string, active bool) error {
	if id == "" {
		return errors.New("id do usuário é obrigatório")
	}
	return s.repo.UpdateActiveStatusGlobal(ctx, id, active)
}

// GetSystemAnalytics obtém as métricas de saúde da plataforma para o Super Admin.
func (s *service) GetSystemAnalytics(ctx context.Context) (*SystemAnalytics, error) {
	return s.repo.GetSystemAnalytics(ctx)
}

// ListAllUsers obtém a lista de todos os usuários cadastrados (Exclusivo Super Admin).
func (s *service) ListAllUsers(ctx context.Context) ([]*User, error) {
	return s.repo.ListAll(ctx)
}
