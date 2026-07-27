# 🛒 Tarefas - Vendas e PDV (`sale`)

Orquestração de pedidos de venda, orçamentos, faturamento e vendas expressas no Ponto de Venda (PDV).

---

## 📋 Lista de Atividades

### 1. Camada de Domínio (`Domain`)
*   [ ] Definir as structs `Order` e `OrderItem` em `model.go` com os campos:
    *   `Order`: `id`, `tenant_id`, `cliente_id`, `data_pedido`, `status_pedido` (ORCAMENTO, APROVADO, CANCELADO, FATURADO), `valor_total` e timestamps.
    *   `OrderItem`: `id`, `tenant_id`, `pedido_id`, `produto_id`, `quantidade`, `preco_unitario_aplicado` e `subtotal` (calculado pelo banco).
*   [ ] Definir a interface `Repository` em `model.go` para persistência:
    *   `Create(ctx context.Context, o *Order, items []*OrderItem) error`
    *   `GetByID(ctx context.Context, id string, tenantID string) (*Order, []*OrderItem, error)`
    *   `List(ctx context.Context, tenantID string, status *string) ([]*Order, error)`
    *   `UpdateStatus(ctx context.Context, id string, tenantID string, status string) error`
*   [ ] **Inversão de Dependências (DIP):** Declarar as interfaces de estoque (`InventoryService`) e financeiro (`FinanceService`) que o módulo de vendas necessita localmente para desacoplar as integrações.

### 2. Camada de Infraestrutura (`Infrastructure`)
*   [ ] Implementar o repositório em `repository.go` em SQL Puro utilizando `pgx/v5`.
    *   [ ] **Salvamento Transacional:** Usar transações Postgres (`Tx`) ao persistir a venda e seus itens, garantindo rollback automático caso ocorra erro em qualquer linha.
    *   [ ] **Subtotal Dinâmico:** A coluna `subtotal` é gerada nativamente pelo banco de dados. O repositório Go deve ler esse valor nas consultas utilizando `.Scan()`.

### 3. Camada de Aplicação (`Application`)
*   [ ] Criar as structs de DTO em `dto.go` para JSON Binding:
    *   `OrderCreateRequest` (cliente_id obrigatório, lista de itens com produto_id, quantidade > 0 e preco_unitario >= 0).
*   [ ] Implementar o `service.go` contendo as regras de negócio:
    *   [ ] **Gatilho de Estoque:** Chamar o serviço de estoque ao transicionar a venda para `APROVADO` ou `FATURADO` para realizar a baixa física das mercadorias.
    *   [ ] **Gatilho Financeiro:** Chamar o serviço financeiro ao transicionar a venda para `FATURADO` para criar um título de Contas a Receber (Receita).
    *   [ ] **Venda Expressa PDV (ACID):** Desenvolver o fluxo do endpoint `/api/v1/sales/pdv` rodando sob uma única transação atômica que:
        1. Cria o pedido diretamente com status `FATURADO`.
        2. Efetua a baixa física no estoque.
        3. Registra a receita como lançada e paga no financeiro.
        4. Efetua rollback geral se qualquer etapa falhar (ex: falta de saldo no estoque).

### 4. Camada de Apresentação e Roteamento (`Adapter/Router`)
*   [ ] Desenvolver o `handler.go` contendo os endpoints:
    *   `POST /api/v1/sales` (Cria orçamento/pedido)
    *   `GET /api/v1/sales` (Lista pedidos do tenant)
    *   `GET /api/v1/sales/:id` (Obtém pedido e seus itens detalhados)
    *   `PUT /api/v1/sales/:id/status` (Modifica status da venda disparando integrações)
    *   `POST /api/v1/sales/pdv` (Endpoint expresso do PDV frente de caixa)
*   [ ] Decorar os handlers do módulo com as tags do Swagger para documentação das rotas.
*   [ ] Mapear as rotas em `routes.go` injetando o contexto do tenant do Echo.

---

## 📊 Status de Conclusão: `0%`
