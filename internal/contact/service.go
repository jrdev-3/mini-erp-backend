package contact

import (
	"context"
	"errors"
)

// Service define os casos de uso para a camada de aplicação de contatos.
type Service interface {
	CreateContact(ctx context.Context, tenantID string, req *ContactRequest) (*Contact, error)
	GetContactByID(ctx context.Context, id string, tenantID string) (*Contact, error)
	ListContacts(ctx context.Context, tenantID string, isCustomer, isSupplier *bool) ([]*Contact, error)
	UpdateContact(ctx context.Context, id string, tenantID string, req *ContactRequest) (*Contact, error)
	DeleteContact(ctx context.Context, id string, tenantID string) error
}

type service struct {
	repo Repository
}

// NewService cria uma nova instância da camada de aplicação do módulo contact.
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// CreateContact valida as regras de negócio e persiste o novo contato vinculando-o ao tenant.
func (s *service) CreateContact(ctx context.Context, tenantID string, req *ContactRequest) (*Contact, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// 1. Validar unicidade do documento (CPF/CNPJ) dentro do mesmo tenant (Prevenção BOLA e duplicidade)
	existing, err := s.repo.GetByDocument(ctx, req.Documento, tenantID)
	if err == nil && existing != nil {
		return nil, ErrDocumentAlreadyExists
	}
	if err != nil && !errors.Is(err, ErrContactNotFound) {
		return nil, err
	}

	// Definir status ativo por padrão caso não tenha sido enviado
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	// 2. Instanciar a entidade Contact
	c := &Contact{
		TenantID:        tenantID,
		IsCustomer:      req.IsCustomer,
		IsSupplier:      req.IsSupplier,
		IsActive:        isActive,
		TipoPessoa:      req.TipoPessoa,
		NomeRazaoSocial: req.NomeRazaoSocial,
		Documento:       req.Documento,
		Email:           req.Email,
		Telefone:        req.Telefone,
		Rua:             req.Rua,
		Numero:          req.Numero,
		Bairro:          req.Bairro,
		Cidade:          req.Cidade,
		Estado:          req.Estado,
		CEP:             req.CEP,
	}

	// 3. Salvar no repositório
	if err := s.repo.Create(ctx, c); err != nil {
		return nil, err
	}

	return c, nil
}

// GetContactByID recupera o contato detalhado validando a posse pelo tenant (BOLA).
func (s *service) GetContactByID(ctx context.Context, id string, tenantID string) (*Contact, error) {
	if id == "" || tenantID == "" {
		return nil, errors.New("id do contato e tenant_id são obrigatórios")
	}
	return s.repo.GetByID(ctx, id, tenantID)
}

// ListContacts lista os contatos pertencentes ao tenant aplicando filtros opcionais.
func (s *service) ListContacts(ctx context.Context, tenantID string, isCustomer, isSupplier *bool) ([]*Contact, error) {
	if tenantID == "" {
		return nil, errors.New("tenant_id é obrigatório")
	}
	return s.repo.List(ctx, tenantID, isCustomer, isSupplier)
}

// UpdateContact valida as regras de negócio de alteração e atualiza o contato.
func (s *service) UpdateContact(ctx context.Context, id string, tenantID string, req *ContactRequest) (*Contact, error) {
	if id == "" || tenantID == "" {
		return nil, errors.New("id do contato e tenant_id são obrigatórios")
	}

	if err := req.Validate(); err != nil {
		return nil, err
	}

	// 1. Verificar se o contato existe e pertence ao tenant
	c, err := s.repo.GetByID(ctx, id, tenantID)
	if err != nil {
		return nil, err
	}

	// 2. Se o documento mudou, validar a unicidade do novo documento no mesmo tenant
	if c.Documento != req.Documento {
		existing, err := s.repo.GetByDocument(ctx, req.Documento, tenantID)
		if err == nil && existing != nil {
			return nil, ErrDocumentAlreadyExists
		}
		if err != nil && !errors.Is(err, ErrContactNotFound) {
			return nil, err
		}
	}

	// Definir status ativo caso tenha sido enviado
	isActive := c.IsActive
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	// 3. Atualizar os dados da entidade
	c.IsCustomer = req.IsCustomer
	c.IsSupplier = req.IsSupplier
	c.IsActive = isActive
	c.TipoPessoa = req.TipoPessoa
	c.NomeRazaoSocial = req.NomeRazaoSocial
	c.Documento = req.Documento
	c.Email = req.Email
	c.Telefone = req.Telefone
	c.Rua = req.Rua
	c.Numero = req.Numero
	c.Bairro = req.Bairro
	c.Cidade = req.Cidade
	c.Estado = req.Estado
	c.CEP = req.CEP

	// 4. Salvar no repositório
	if err := s.repo.Update(ctx, c); err != nil {
		return nil, err
	}

	return c, nil
}

// DeleteContact remove o contato do banco validando a posse pelo tenant (BOLA).
func (s *service) DeleteContact(ctx context.Context, id string, tenantID string) error {
	if id == "" || tenantID == "" {
		return errors.New("id do contato e tenant_id são obrigatórios")
	}
	return s.repo.Delete(ctx, id, tenantID)
}
