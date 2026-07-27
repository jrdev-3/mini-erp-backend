# 💵 Tarefas - Financeiro (`finance`)

Módulo responsável pelo controle de fluxo de caixa, conciliação e lançamentos de contas a pagar/receber sob isolamento multi-tenant.

---

## 📋 Lista de Atividades

### 1. Camada de Domínio (`Domain`)
*   [ ] Definir a struct `Transaction` in `model.go` com os campos:
    *   `id`, `tenant_id`, `tipo_lancamento` (RECEITA/DESPESA), `descricao`, `valor`, `data_vencimento`, `data_pagamento` (opcional/null), `status_pagamento` (PENDENTE, PAGO, ATRASADO), `origem_id` (UUID opcional para vendas/compras) e timestamps.
*   [ ] Definir a interface `Repository` in `model.go` para persistência:
    *   `Create(ctx context.Context, t *Transaction) error`
    *   `GetByID(ctx context.Context, id string, tenantID string) (*Transaction, error)`
    *   `List(ctx context.Context, tenantID string, status *string, tipo *string) ([]*Transaction, error)`
    *   `UpdateStatus(ctx context.Context, id string, tenantID string, status string, dataPagamento *time.Time) error`

### 2. Camada de Infraestrutura (`Infrastructure`)
*   [ ] Implementar o repositório em `repository.go` em SQL Puro utilizando `pgx/v5`.
    *   [ ] **Segurança contra BOLA:** Impor a cláusula `WHERE tenant_id = $x` em todas as consultas e updates de lançamentos financeiros.
    *   [ ] **Mitigação de Enumeração:** Se o registro não pertencer ao tenant logado, o repositório deve ocultar a existência do dado respondendo `404 Not Found`.

### 3. Camada de Aplicação (`Application`)
*   [ ] Criar as structs de DTO em `dto.go` com validações HTTP:
    *   `TransactionRequest` (tipo_lancamento válido, valor > 0, data_vencimento obrigatória).
*   [ ] Implementar as regras de negócios em `service.go`:
    *   [ ] Lógica para criação manual de transações.
    *   [ ] Lógica de conciliação: Ao marcar uma transação como `PAGO`, injetar automaticamente a data atual no campo `data_pagamento`.
    *   [ ] Lógica de integração: Receber chamadas do módulo `sale` para criar receitas automáticas vinculadas ao ID de origem do pedido.

### 4. Camada de Apresentação e Roteamento (`Adapter/Router`)
*   [ ] Desenvolver o `handler.go` contendo os endpoints:
    *   `POST /api/v1/finance/transactions` (Lançamento manual de receita ou despesa)
    *   `GET /api/v1/finance/transactions` (Lista fluxo de caixa com filtros de tipo e status)
    *   `PATCH /api/v1/finance/transactions/:id/pay` (Confirma recebimento/pagamento do título)
*   [ ] Decorar os handlers do módulo com as tags do Swagger para documentação das rotas.
*   [ ] Registrar as rotas do módulo em `routes.go` associando o roteador Echo v5.

---

## 📊 Status de Conclusão: `0%`
