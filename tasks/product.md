# 📦 Tarefas - Produtos e Categorias (`product`)

Gerenciamento do catálogo de produtos e agrupamento por categorias em escopo multi-tenant.

---

## 📋 Lista de Atividades

### 1. Camada de Domínio (`Domain`)
*   [ ] Definir as structs `Product` e `Category` em `model.go` com os campos:
    *   `Category`: `id`, `tenant_id`, `nome` e timestamps.
    *   `Product`: `id`, `tenant_id`, `categoria_id` (opcional/UUID), `codigo_sku`, `nome`, `preco_custo`, `preco_venda`, `unidade_medida`, `is_active` e timestamps.
*   [ ] Definir as interfaces de persistência em `model.go`:
    *   `CategoryRepository`: CRUD completo para categorias.
    *   `ProductRepository`: CRUD completo para produtos.

### 2. Camada de Infraestrutura (`Infrastructure`)
*   [ ] Implementar os repositórios em `repository.go` usando SQL Puro e `pgx/v5`.
    *   [ ] **Segurança BOLA:** Garantir filtro do `tenant_id` em toda query.
    *   [ ] **Deleção Segura:** Se a categoria for apagada do banco, os produtos vinculados a ela devem ter seu campo `categoria_id` alterado automaticamente para `NULL` (conforme especificação `ON DELETE SET NULL` do banco).
    *   [ ] **Desativação Lógica:** A exclusão de um produto realiza um update de status lógica no banco (`is_active = false`), preservando os registros nos históricos de vendas e estoque.

### 3. Camada de Aplicação (`Application`)
*   [ ] Desenvolver os DTOs em `dto.go` com validações HTTP:
    *   `ProductRequest` (nome obrigatório, SKU em maiúsculo, preço de custo >= 0, preço de venda >= preço de custo).
    *   `CategoryRequest` (nome da categoria obrigatório).
*   [ ] Implementar as lógicas de negócio em `service.go`:
    *   [ ] Validação de SKU único para produtos dentro de um mesmo tenant.
    *   [ ] Validação de nome de categoria único por tenant.

### 4. Camada de Apresentação e Roteamento (`Adapter/Router`)
*   [ ] Desenvolver o `handler.go` contendo as rotas:
    *   `POST /api/v1/categories` (Cria categoria)
    *   `GET /api/v1/categories` (Lista categorias do tenant)
    *   `POST /api/v1/products` (Cria produto)
    *   `GET /api/v1/products` (Lista produtos ativos com paginação e filtros de nome/SKU/categoria)
    *   `GET /api/v1/products/:id` (Obtém produto do tenant)
    *   `PUT /api/v1/products/:id` (Atualiza cadastro do produto)
    *   `DELETE /api/v1/products/:id` (Exclusão lógica do produto)
*   [ ] Decorar os handlers do módulo com as tags do Swagger para documentação das rotas.
*   [ ] Mapear as rotas em `routes.go` associando o roteador Echo v5.

---

## 📊 Status de Conclusão: `0%`
