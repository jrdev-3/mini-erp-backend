package auth

import (
	"context"
	"errors"

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

// Create insere um novo usuário na tabela users e retorna os dados atualizados com o ID gerado pelo banco.
func (r *repository) Create(ctx context.Context, u *User) error {
	query := `
		INSERT INTO users (email, password_hash, role, is_active, tenant_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, criado_em, atualizado_em
	`
	err := r.db.QueryRow(ctx, query, u.Email, u.PasswordHash, u.Role, u.IsActive, u.TenantID).
		Scan(&u.ID, &u.CriadoEm, &u.AtualizadoEm)
	return err
}

// GetByEmail busca um usuário por e-mail para autenticação global (sem tenant_id prévio).
func (r *repository) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, email, password_hash, role, is_active, tenant_id, criado_em, atualizado_em
		FROM users
		WHERE email = $1
	`
	var u User
	err := r.db.QueryRow(ctx, query, email).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.Role, &u.IsActive, &u.TenantID, &u.CriadoEm, &u.AtualizadoEm,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}

// GetByID busca um usuário no banco validando a posse pelo tenant_id (Prevenção a BOLA).
func (r *repository) GetByID(ctx context.Context, id string, tenantID string) (*User, error) {
	query := `
		SELECT id, email, password_hash, role, is_active, tenant_id, criado_em, atualizado_em
		FROM users
		WHERE id = $1 AND tenant_id = $2
	`
	var u User
	err := r.db.QueryRow(ctx, query, id, tenantID).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.Role, &u.IsActive, &u.TenantID, &u.CriadoEm, &u.AtualizadoEm,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}

// UpdateActiveStatus atualiza o status ativo/inativo de um usuário dentro de seu próprio tenant (Prevenção a BOLA).
func (r *repository) UpdateActiveStatus(ctx context.Context, id string, tenantID string, active bool) error {
	query := `
		UPDATE users
		SET is_active = $1, atualizado_em = NOW()
		WHERE id = $2 AND tenant_id = $3
	`
	result, err := r.db.Exec(ctx, query, active, id, tenantID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrUserNotFound
	}
	return nil
}

// UpdateActiveStatusGlobal atualiza o status ativo/inativo de qualquer usuário no sistema (Exclusivo Super Admin).
func (r *repository) UpdateActiveStatusGlobal(ctx context.Context, id string, active bool) error {
	query := `
		UPDATE users
		SET is_active = $1, atualizado_em = NOW()
		WHERE id = $2
	`
	result, err := r.db.Exec(ctx, query, active, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrUserNotFound
	}
	return nil
}

// GetSystemAnalytics calcula as métricas agregadas da tabela users para o Super Admin.
func (r *repository) GetSystemAnalytics(ctx context.Context) (*SystemAnalytics, error) {
	query := `
		SELECT 
			COUNT(*)::bigint as total_users,
			COUNT(*) FILTER (WHERE is_active = true)::bigint as active_users,
			COUNT(*) FILTER (WHERE is_active = false)::bigint as inactive_users
		FROM users
	`
	var sa SystemAnalytics
	err := r.db.QueryRow(ctx, query).Scan(
		&sa.TotalUsers, &sa.ActiveUsers, &sa.InactiveUsers,
	)
	if err != nil {
		return nil, err
	}
	return &sa, nil
}

// ListAll retorna todos os usuários cadastrados no sistema ordenados por data de criação decrescente.
func (r *repository) ListAll(ctx context.Context) ([]*User, error) {
	query := `
		SELECT id, email, password_hash, role, is_active, tenant_id, criado_em, atualizado_em
		FROM users
		ORDER BY criado_em DESC
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		var u User
		err := rows.Scan(
			&u.ID, &u.Email, &u.PasswordHash, &u.Role, &u.IsActive, &u.TenantID, &u.CriadoEm, &u.AtualizadoEm,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
