# 👥 Tarefas - Contatos (`contact`)

Módulo responsável pelo gerenciamento de Clientes e Fornecedores sob isolamento multi-tenant estrito.

---

## 📋 Lista de Atividades

### 1. Camada de Domínio (`Domain`)
*   [x] Definir a struct `Contact` em `model.go` contendo as propriedades:
    *   `id`, `tenant_id`, `is_customer`, `is_supplier`, `is_active`, `tipo_pessoa` (PF/PJ), `nome_razao_social`, `documento` (apenas números), `email`, `telefone`, `rua`, `numero`, `bairro`, `cidade`, `estado` (2 caracteres) e `cep` (8 caracteres).
*   [x] Definir a interface `Repository` em `model.go`:
    *   `Create(ctx context.Context, c *Contact) error`
    *   `GetByID(ctx context.Context, id string, tenantID string) (*Contact, error)`
    *   `List(ctx context.Context, tenantID string, isCustomer, isSupplier *bool) ([]*Contact, error)`
    *   `Update(ctx context.Context, c *Contact) error`
    *   `Delete(ctx context.Context, id string, tenantID string) error`
 
### 2. Camada de Infraestrutura (`Infrastructure`)
*   [x] Implementar o repositório em `repository.go` em SQL Puro utilizando `pgx/v5`.
    *   [x] **Segurança contra BOLA:** Todas as queries de escrita e leitura devem aplicar obrigatoriamente a cláusula `WHERE tenant_id = $x` no banco.
    *   [x] **Mitigação de Enumeração:** Buscar um contato inexistente ou que pertença a outro tenant deve retornar `pgx.ErrNoRows` e responder como `404 Not Found`.

### 3. Camada de Aplicação (`Application`)
*   [x] Criar as structs de DTO em `dto.go` contendo validações HTTP:
    *   `ContactRequest` (validação de documento CPF/CNPJ, validação de email, estado com 2 letras, CEP com 8 dígitos).
*   [x] Implementar a lógica de negócios em `service.go`:
    *   [x] Validação de documento único (CPF/CNPJ) dentro de um mesmo tenant.
    *   [x] Tratamento para aceitar o mesmo documento cadastrado em tenants diferentes (isolamento BOLA corporativo).

### 4. Camada de Apresentação e Roteamento (`Adapter/Router`)
*   [ ] Desenvolver o `handler.go` contendo os endpoints:
    *   `POST /api/v1/contacts` (Cria contato)
    *   `GET /api/v1/contacts` (Lista contatos do tenant com filtros opcionais de `is_customer` e `is_supplier`)
    *   `GET /api/v1/contacts/:id` (Obtém contato detalhado do tenant)
    *   `PUT /api/v1/contacts/:id` (Atualiza dados do contato)
    *   `DELETE /api/v1/contacts/:id` (Remove contato do tenant)
*   [ ] Decorar os handlers do módulo com as tags do Swagger para documentação das rotas.
*   [ ] Registrar as rotas no Echo em `routes.go` aplicando o middleware global de autenticação JWT e extração do `tenant_id` do contexto.

---

## 📊 Status de Conclusão: `70%`
