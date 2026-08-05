package contact

import (
	"context"
	"errors"
	"time"
)

// Erros de domínio do módulo contact
var (
	ErrContactNotFound       = errors.New("contact not found")
	ErrDocumentAlreadyExists = errors.New("document already in use for this tenant")
)

// Contact representa a entidade de Clientes e Fornecedores no Mini ERP.
type Contact struct {
	ID               string    `json:"id"`
	TenantID         string    `json:"tenant_id"`
	IsCustomer       bool      `json:"is_customer"`
	IsSupplier       bool      `json:"is_supplier"`
	IsActive         bool      `json:"is_active"`
	TipoPessoa       string    `json:"tipo_pessoa"` // PF ou PJ
	NomeRazaoSocial  string    `json:"nome_razao_social"`
	Documento        string    `json:"documento"` // CPF ou CNPJ (apenas números)
	Email            *string   `json:"email"`
	Telefone         *string   `json:"telefone"`
	Rua              *string   `json:"rua"`
	Numero           *string   `json:"numero"`
	Bairro           *string   `json:"bairro"`
	Cidade           *string   `json:"cidade"`
	Estado           *string   `json:"estado"`
	CEP              *string   `json:"cep"`
	CriadoEm         time.Time `json:"criado_em"`
	AtualizadoEm     time.Time `json:"atualizado_em"`
}

// Repository define a interface de persistência de contatos com isolamento de tenant (prevenção BOLA).
type Repository interface {
	Create(ctx context.Context, c *Contact) error
	GetByID(ctx context.Context, id string, tenantID string) (*Contact, error)
	GetByDocument(ctx context.Context, documento string, tenantID string) (*Contact, error)
	List(ctx context.Context, tenantID string, isCustomer, isSupplier *bool) ([]*Contact, error)
	Update(ctx context.Context, c *Contact) error
	Delete(ctx context.Context, id string, tenantID string) error
}
