# 🏪 Tarefas - Estoque (`inventory`)

Módulo responsável pelo controle físico de saldo de estoque dos produtos e logs históricos de movimentação.

---

## 📋 Lista de Atividades

### 1. Camada de Domínio (`Domain`)
*   [ ] Definir as structs `Stock` e `Movement` em `model.go` com os campos:
    *   `Stock`: `id`, `tenant_id`, `produto_id`, `quantidade_atual`, `quantidade_minima` e timestamps.
    *   `Movement`: `id`, `tenant_id`, `produto_id`, `tipo_movimentacao` (ENTRADA_COMPRA, SAIDA_VENDA, AJUSTE_INVENTARIO, PERDA), `quantidade`, `origem_id` (UUID opcional para venda/compra) e data.
*   [ ] Definir a interface `Repository` em `model.go` para persistência:
    *   `GetStockByProduct(ctx context.Context, productID string, tenantID string) (*Stock, error)`
    *   `UpdateStock(ctx context.Context, stock *Stock) error`
    *   `AddMovement(ctx context.Context, m *Movement) error`
    *   `ListMovementsByProduct(ctx context.Context, productID string, tenantID string) ([]*Movement, error)`

### 2. Camada de Infraestrutura (`Infrastructure`)
*   [ ] Implementar o repositório em `repository.go` em SQL Puro utilizando `pgx/v5`.
    *   [ ] **Controle Concorrente:** Utilizar controle transacional seguro no Postgres (como `SELECT ... FOR UPDATE` se necessário) ao atualizar quantidades de estoque concorrentemente.
    *   [ ] **Segurança BOLA:** Filtrar todas as consultas de estoque e logs com a cláusula obrigatória `tenant_id`.

### 3. Camada de Aplicação (`Application`)
*   [ ] Criar as structs de DTO em `dto.go` para requisições de movimentação:
    *   `ManualMovementRequest` (produto_id obrigatório, tipo de movimentação válida, quantidade > 0).
*   [ ] Implementar a lógica de negócios em `service.go`:
    *   [ ] **Impedir Estoque Negativo:** O serviço deve barrar e retornar erro se qualquer saída resultar em `quantidade_atual < 0` (segurança garantida pela `CHECK constraint` do banco).
    *   [ ] **Gatilho de Atualização Automática:** Lógica para receber chamadas de outros módulos (como `sale` ao faturar venda) e atualizar o estoque, registrando o log correspondente com a chave `origem_id` vinculada.

### 4. Camada de Apresentação e Roteamento (`Adapter/Router`)
*   [ ] Desenvolver o `handler.go` contendo as rotas:
    *   `GET /api/v1/inventory/stock` (Lista saldos de estoque do tenant)
    *   `GET /api/v1/inventory/stock/:product_id` (Obtém saldo de um produto específico do tenant)
    *   `POST /api/v1/inventory/movements` (Registra movimentação manual - ex: ajuste ou perda)
    *   `GET /api/v1/inventory/movements/:product_id` (Exibe extrato de movimentações de um produto do tenant)
*   [ ] Decorar os handlers do módulo com as tags do Swagger para documentação das rotas.
*   [ ] Configurar as rotas no roteador Echo v5 em `routes.go` aplicando o middleware JWT.

---

## 📊 Status de Conclusão: `0%`
