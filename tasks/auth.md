# 🔐 Tarefas - Autenticação e Usuários (`auth`)

Este módulo gerencia a segurança, o controle de acesso e o isolamento multi-tenant da API.

---

## 📋 Lista de Atividades

### 1. Camada de Domínio (`Domain`)
*   [x] Definir a struct `User` em `model.go` com os campos:
    *   `id`, `email`, `password_hash`, `role` (ADMIN/USER), `is_active`, `tenant_id` (UUID), `criado_em` e `atualizado_em`.
*   [x] Definir a interface `Repository` em `model.go` para persistência:
    *   `Create(ctx context.Context, u *User) error`
    *   `GetByEmail(ctx context.Context, email string) (*User, error)`
    *   `GetByID(ctx context.Context, id string, tenantID string) (*User, error)`
    *   `UpdateActiveStatus(ctx context.Context, id string, tenantID string, active bool) error`

### 2. Camada de Infraestrutura (`Infrastructure`)
*   [x] Implementar o repositório em `repository.go` em SQL Puro utilizando o pool do `pgx/v5`.
    *   [x] Garantir queries parametrizadas (placeholders `$1`, `$2`) contra SQL Injection.

### 3. Camada de Aplicação (`Application`)
*   [x] Criar as structs de DTO em `dto.go` para JSON Binding:
    *   `RegisterRequest` (e-mail válido, senha forte, nome da empresa para gerar tenant).
    *   `LoginRequest` (e-mail e senha).
    *   `LoginResponse` (token JWT gerado com `tenant_id` e `role` embutidos).
*   [x] Implementar o `service.go` com as regras de negócio:
    *   [x] Hashing de senhas seguro utilizando `bcrypt` (custo padrão 10 ou mais).
    *   [x] Validação de e-mail único.
    *   [x] Geração de Token JWT offline assinado digitalmente com chave secreta do Supabase.

### 4. Camada de Apresentação e Roteamento (`Adapter/Router`)
*   [x] Desenvolver o `handler.go` contendo as rotas:
    *   `POST /api/v1/auth/register` (Público - Cria tenant e primeiro usuário ADMIN)
    *   `POST /api/v1/auth/login` (Público - Retorna o JWT)
    *   `PATCH /api/v1/admin/users/:id/toggle` (Protegida - Ativa/Desativa colaborador)
    *   `POST /api/v1/admin/users` (Protegida - Cadastra funcionário no mesmo tenant)
    *   `GET /api/v1/admin/users` (Protegida - Lista funcionários do comércio)
    *   `PATCH /api/v1/system/users/:id/toggle` (Exclusivo Super Admin - Ativa/Desativa qualquer usuário do sistema)
    *   `GET /api/v1/system/analytics` (Exclusivo Super Admin - Métricas agregadas da plataforma)
    *   `GET /api/v1/system/users` (Exclusivo Super Admin - Lista todos os usuários do ecossistema)
*   [x] Decorar os handlers do módulo com as tags do Swagger para documentação das rotas.
*   [x] Configurar as rotas em `routes.go` vinculando o roteador Echo v5.
*   [x] Desenvolver os middlewares em `internal/middleware/`:
    *   [x] `auth.go`: Middleware de validação do JWT offline decodificando `tenant_id` e `role` para o contexto da requisição.
    *   [x] `rbac.go`: Middleware para validar cargo de `ADMIN` e realizar a dupla validação de `SUPER_ADMIN` contra o `ADMIN_TENANT_ID` da plataforma.

---

## 📊 Status de Conclusão: `100%`
