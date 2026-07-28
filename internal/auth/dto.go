package auth

import (
	"errors"
	"strings"
)

// RegisterRequest define os dados de entrada para cadastro de um novo tenant e seu primeiro usuário administrador.
type RegisterRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	CompanyName string `json:"company_name"`
}

// Validate executa validações de dados e formato na requisição de cadastro.
func (r *RegisterRequest) Validate() error {
	email := strings.TrimSpace(r.Email)
	if email == "" {
		return errors.New("o e-mail é obrigatório")
	}
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return errors.New("formato de e-mail inválido")
	}

	password := strings.TrimSpace(r.Password)
	if len(password) < 8 {
		return errors.New("a senha deve conter no mínimo 8 caracteres")
	}

	companyName := strings.TrimSpace(r.CompanyName)
	if companyName == "" {
		return errors.New("o nome da empresa é obrigatório para a criação do tenant")
	}

	return nil
}

// LoginRequest define as credenciais necessárias para autenticação do usuário.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Validate executa validações de preenchimento na requisição de login.
func (r *LoginRequest) Validate() error {
	email := strings.TrimSpace(r.Email)
	if email == "" {
		return errors.New("o e-mail é obrigatório")
	}

	password := strings.TrimSpace(r.Password)
	if password == "" {
		return errors.New("a senha é obrigatória")
	}

	return nil
}

// LoginResponse define a estrutura de resposta após autenticação bem-sucedida, contendo o token de acesso.
type LoginResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}
