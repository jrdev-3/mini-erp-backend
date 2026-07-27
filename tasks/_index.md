# 🗺️ Índice de Tarefas - Mini ERP

Roadmap de implementação de infraestrutura e lógica de negócios do backend. A evolução do projeto segue a ordem lógica de dependências de domínio.

---

## ⚙️ Infraestrutura Global

*   [ ] Inicializar documentação da API com Swagger (`swag init` e rota `GET /swagger/*` no `main.go`)
*   [ ] Implementar desativação do Swagger online sob a flag de ambiente `APP_ENV=production`

---

## 📌 Progresso Macro dos Módulos

*   [ ] **Módulo 1: Autenticação e Usuários (`auth`)** ➔ [Ver Tarefas](auth.md)
*   [ ] **Módulo 2: Contatos (`contact`)** ➔ [Ver Tarefas](contact.md)
*   [ ] **Módulo 3: Produtos e Categorias (`product`)** ➔ [Ver Tarefas](product.md)
*   [ ] **Módulo 4: Estoque (`inventory`)** ➔ [Ver Tarefas](inventory.md)
*   [ ] **Módulo 5: Vendas e PDV (`sale`)** ➔ [Ver Tarefas](sale.md)
*   [ ] **Módulo 6: Financeiro (`finance`)** ➔ [Ver Tarefas](finance.md)
