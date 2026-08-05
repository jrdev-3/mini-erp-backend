package contact

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repository struct {
	db *pgxpool.Pool
}

// NewRepository cria uma nova instância de Repository implementada com pgxpool.
func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

// Create insere um novo contato no banco e retorna os dados atualizados com o ID gerado.
func (r *repository) Create(ctx context.Context, c *Contact) error {
	query := `
		INSERT INTO contacts (
			tenant_id, is_customer, is_supplier, is_active, tipo_pessoa, 
			nome_razao_social, documento, email, telefone, rua, 
			numero, bairro, cidade, estado, cep
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		RETURNING id, criado_em, atualizado_em
	`
	err := r.db.QueryRow(
		ctx, query,
		c.TenantID, c.IsCustomer, c.IsSupplier, c.IsActive, c.TipoPessoa,
		c.NomeRazaoSocial, c.Documento, c.Email, c.Telefone, c.Rua,
		c.Numero, c.Bairro, c.Cidade, c.Estado, c.CEP,
	).Scan(&c.ID, &c.CriadoEm, &c.AtualizadoEm)
	return err
}

// GetByID busca um contato por ID validando a posse pelo tenant_id (Prevenção BOLA).
func (r *repository) GetByID(ctx context.Context, id string, tenantID string) (*Contact, error) {
	query := `
		SELECT 
			id, tenant_id, is_customer, is_supplier, is_active, tipo_pessoa, 
			nome_razao_social, documento, email, telefone, rua, 
			numero, bairro, cidade, estado, cep, criado_em, atualizado_em
		FROM contacts
		WHERE id = $1 AND tenant_id = $2
	`
	var c Contact
	err := r.db.QueryRow(ctx, query, id, tenantID).Scan(
		&c.ID, &c.TenantID, &c.IsCustomer, &c.IsSupplier, &c.IsActive, &c.TipoPessoa,
		&c.NomeRazaoSocial, &c.Documento, &c.Email, &c.Telefone, &c.Rua,
		&c.Numero, &c.Bairro, &c.Cidade, &c.Estado, &c.CEP, &c.CriadoEm, &c.AtualizadoEm,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrContactNotFound
		}
		return nil, err
	}
	return &c, nil
}

// GetByDocument busca um contato por documento validando a posse pelo tenant_id (Prevenção BOLA).
func (r *repository) GetByDocument(ctx context.Context, documento string, tenantID string) (*Contact, error) {
	query := `
		SELECT 
			id, tenant_id, is_customer, is_supplier, is_active, tipo_pessoa, 
			nome_razao_social, documento, email, telefone, rua, 
			numero, bairro, cidade, estado, cep, criado_em, atualizado_em
		FROM contacts
		WHERE documento = $1 AND tenant_id = $2
	`
	var c Contact
	err := r.db.QueryRow(ctx, query, documento, tenantID).Scan(
		&c.ID, &c.TenantID, &c.IsCustomer, &c.IsSupplier, &c.IsActive, &c.TipoPessoa,
		&c.NomeRazaoSocial, &c.Documento, &c.Email, &c.Telefone, &c.Rua,
		&c.Numero, &c.Bairro, &c.Cidade, &c.Estado, &c.CEP, &c.CriadoEm, &c.AtualizadoEm,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrContactNotFound
		}
		return nil, err
	}
	return &c, nil
}

// List lista todos os contatos associados a um determinado inquilino aplicando filtros opcionais (Prevenção BOLA).
func (r *repository) List(ctx context.Context, tenantID string, isCustomer, isSupplier *bool) ([]*Contact, error) {
	query := `
		SELECT 
			id, tenant_id, is_customer, is_supplier, is_active, tipo_pessoa, 
			nome_razao_social, documento, email, telefone, rua, 
			numero, bairro, cidade, estado, cep, criado_em, atualizado_em
		FROM contacts
		WHERE tenant_id = $1
	`
	args := []interface{}{tenantID}
	placeholderIndex := 2

	if isCustomer != nil {
		query += fmt.Sprintf(" AND is_customer = $%d", placeholderIndex)
		args = append(args, *isCustomer)
		placeholderIndex++
	}

	if isSupplier != nil {
		query += fmt.Sprintf(" AND is_supplier = $%d", placeholderIndex)
		args = append(args, *isSupplier)
		placeholderIndex++
	}

	query += " ORDER BY nome_razao_social ASC"

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contacts []*Contact
	for rows.Next() {
		var c Contact
		err := rows.Scan(
			&c.ID, &c.TenantID, &c.IsCustomer, &c.IsSupplier, &c.IsActive, &c.TipoPessoa,
			&c.NomeRazaoSocial, &c.Documento, &c.Email, &c.Telefone, &c.Rua,
			&c.Numero, &c.Bairro, &c.Cidade, &c.Estado, &c.CEP, &c.CriadoEm, &c.AtualizadoEm,
		)
		if err != nil {
			return nil, err
		}
		contacts = append(contacts, &c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return contacts, nil
}

// Update atualiza os dados cadastrais do contato validando a posse pelo tenant_id (Prevenção BOLA).
func (r *repository) Update(ctx context.Context, c *Contact) error {
	query := `
		UPDATE contacts
		SET 
			is_customer = $1, 
			is_supplier = $2, 
			is_active = $3, 
			tipo_pessoa = $4, 
			nome_razao_social = $5, 
			documento = $6, 
			email = $7, 
			telefone = $8, 
			rua = $9, 
			numero = $10, 
			bairro = $11, 
			cidade = $12, 
			estado = $13, 
			cep = $14, 
			atualizado_em = NOW()
		WHERE id = $15 AND tenant_id = $16
	`
	result, err := r.db.Exec(
		ctx, query,
		c.IsCustomer, c.IsSupplier, c.IsActive, c.TipoPessoa,
		c.NomeRazaoSocial, c.Documento, c.Email, c.Telefone, c.Rua,
		c.Numero, c.Bairro, c.Cidade, c.Estado, c.CEP,
		c.ID, c.TenantID,
	)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrContactNotFound
	}
	return nil
}

// Delete remove fisicamente o contato do banco validando a posse pelo tenant_id (Prevenção BOLA).
func (r *repository) Delete(ctx context.Context, id string, tenantID string) error {
	query := `
		DELETE FROM contacts
		WHERE id = $1 AND tenant_id = $2
	`
	result, err := r.db.Exec(ctx, query, id, tenantID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrContactNotFound
	}
	return nil
}
