# 🔐 Tarefas - Autenticação e Usuários (`auth`)

Este módulo gerencia a segurança, o controle de acesso e o isolamento multi-tenant da API.

---

## 📋 Lista de Atividades

### 1. Camada de Domínio (`Domain`)
*   [ ] Definir a struct `User` em `model.go` com os campos:
    *   `id`, `email`, `password_hash`, `role` (ADMIN/USER), `is_active`, `tenant_id` (UUID), `criado_em` e `atualizado_em`.
*   [ ] Definir a interface `Repository` em `model.go` para persistência:
    *   `Create(ctx context.Context, u *User) error`
    *   `GetByEmail(ctx context.Context, email string) (*User, error)`
    *   `GetByID(ctx context.Context, id string, tenantID string) (*User, error)`
    *   `UpdateActiveStatus(ctx context.Context, id string, tenantID string, active bool) error`

### 2. Camada de Infraestrutura (`Infrastructure`)
*   [ ] Implementar o repositório em `repository.go` em SQL Puro utilizando o pool do `pgx/v5`.
    *   [ ] Garantir queries parametrizadas (placeholders `$1`, `$2`) contra SQL Injection.

### 3. Camada de Aplicação (`Application`)
*   [ ] Criar as structs de DTO em `dto.go` para JSON Binding:
    *   `RegisterRequest` (e-mail válido, senha forte, nome da empresa para gerar tenant).
    *   `LoginRequest` (e-mail e senha).
    *   `LoginResponse` (token JWT gerado com `tenant_id` e `role` embutidos).
*   [ ] Implementar o `service.go` com as regras de negócio:
    *   [ ] Hashing de senhas seguro utilizando `bcrypt` (custo padrão 10 ou mais).
    *   [ ] Validação de e-mail único.
    *   [ ] Geração de Token JWT offline assinado digitalmente com chave secreta do Supabase.

### 4. Camada de Apresentação e Roteamento (`Adapter/Router`)
*   [ ] Desenvolver o `handler.go` contendo as rotas:
    *   `POST /api/v1/auth/register` (Público - Cria tenant e primeiro usuário ADMIN)
    *   `POST /api/v1/auth/login` (Público - Retorna o JWT)
    *   `PATCH /api/v1/admin/users/:id/toggle` (Protegida - Ativa/Desativa usuário)
*   [ ] Decorar os handlers do módulo com as tags do Swagger para documentação das rotas.
*   [ ] Configurar as rotas em `routes.go` vinculando o roteador Echo v5.
*   [ ] Desenvolver os middlewares em `internal/middleware/`:
    *   [ ] `auth.go`: Middleware de validação do JWT offline decodificando `tenant_id` e `role` para o contexto da requisição.
    *   [ ] `rbac.go`: Middleware para validar cargo de `ADMIN` em rotas sensíveis.

---

## 📊 Status de Conclusão: `0%`
