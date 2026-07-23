# Diretrizes de Segurança da API (OWASP Top 10:2025)

Este documento estabelece as regras e padrões de mitigação para as principais vulnerabilidades de segurança web aplicadas ao nosso backend em Go e banco de dados Supabase.

---

## 🔒 Mitigação das Principais Vulnerabilidades

### A01: Broken Access Control (Controle de Acesso Quebrado)
*   **BOLA (Broken Object Level Authorization):** Vulnerabilidade onde um usuário tenta acessar ou modificar dados de outra empresa alterando IDs nas requisições HTTP (ex: `GET /api/v1/sales/999`).
*   **Medidas de Blindagem Multi-tenant (Baseada nas regras de controle de escopo e isolamento):**
    *   **Isolamento de Dados no Banco:** Todas as tabelas operacionais do ERP (contatos, produtos, estoque, vendas, lançamentos financeiros) devem conter a coluna `tenant_id` (UUID referenciando a empresa/conta proprietária).
    *   **Filtro Injetado via Contexto (Segurança Absoluta):** O middleware de autenticação (`middleware.Auth`) extrai o `tenant_id` seguro e validado de dentro do token JWT do usuário logado e o insere no contexto da requisição (`echo.Context` / `context.Context` do Go).
    *   **Proibição de Confiança no Cliente:** O backend **nunca** deve aceitar o `tenant_id` enviado em campos editáveis de formulários, query params ou cabeçalhos vindos do frontend. O valor do contexto é a única fonte da verdade.
    *   **Imposição de Filtro no Repositório (SQL Puro):** Toda e qualquer query SQL (queries de leitura `SELECT` ou comandos de escrita `INSERT/UPDATE/DELETE`) que atue sobre recursos do ERP deve aplicar obrigatoriamente o filtro do `tenant_id` obtido do contexto na cláusula `WHERE`. Exemplo:
        ```go
        // Exemplo seguro com SQL Puro no Repositório Go
        query := `SELECT id, name, sku, price_venda FROM products WHERE id = $1 AND tenant_id = $2`
        err := r.db.QueryRow(ctx, query, productID, tenantID).Scan(&p.ID, &p.Name, &p.SKU, &p.Price)
        ```
    *   **Ofuscação de Mensagem de Erro (Retorno 404):** Caso o usuário tente buscar ou manipular um registro com ID válido que pertença a outro `tenant_id`, a consulta SQL não retornará linhas. O backend deve responder com `404 Not Found` (e não `403 Forbidden`), evitando que atacantes descubram se aquele ID existe no banco de dados de outra empresa (prevenção contra enumeração de recursos).

### A02: Security Misconfiguration (Configuração Incorreta de Segurança)
*   **Medida:**
    *   **Swagger UI:** Desativar a documentação online em produção condicionado à flag `APP_ENV=production` no `main.go`.
    *   **Erros do Servidor:** Nunca retornar mensagens de erro internas do banco ou do sistema (ex: strings de stack trace, falha de conexão) nas respostas HTTP. Tratar o erro internamente, gerar logs detalhados no console e retornar mensagens genéricas como `"erro interno do servidor"`.
    *   **Portas de comunicação:** Não expor portas não utilizadas no Dockerfile. Expor apenas a porta HTTP do servidor (`8080`).

### A03: Software Supply Chain Failures (Falhas na Cadeia de Suprimentos)
*   **Medida:**
    *   Adotar a política "Faça você mesmo" (DIY) evitando dependências externas desnecessárias para tarefas que a biblioteca padrão (`stdlib`) do Go executa de forma estável.
    *   Manter o driver do Postgres (`pgx/v5`) sempre fixado na versão segura mais recente (`@latest`) para evitar injeções ou overflows históricos.
    *   Auditar vulnerabilidades de dependências periodicamente com `govulncheck ./...`.

### A04: Cryptographic Failures (Falhas Criptográficas)
*   **Medida:**
    *   **Tráfego de dados:** Todo o tráfego da API com o frontend e o servidor externo do Supabase deve trafegar exclusivamente em conexões criptografadas **HTTPS** (TLS 1.3 recomendado).
    *   **Dados Sensíveis:** Segredos críticos, como a string de conexão segura do PostgreSQL e a chave secreta JWT (JWT Secret) do Supabase, devem residir estritamente em variáveis de ambiente da hospedagem (Render). Eles nunca devem ser armazenados no repositório Git ou retornados em texto claro em corpos de respostas HTTP.

### A05: Injection (Injeção de SQL e Comandos)
*   **Medida:**
    *   **SQL Parameterized Queries:** Nunca concatenar strings de entrada de usuário diretamente para executar queries do banco de dados (ex: `fmt.Sprintf("SELECT * FROM users WHERE name = '%s'", name)`).
    *   **Padrão Seguro:** Sempre usar placeholders parametrizados (ex: `$1`, `$2`) fornecidos pelo driver pgx:
        ```go
        Pool.QueryRow(ctx, "SELECT id FROM users WHERE email = $1", email)
        ```

### A06: Insecure Design (Design Inseguro)
*   **Medida:**
    *   Avaliar os riscos e desenhar a arquitetura antes de iniciar o código. 
    *   Dividir as tarefas em checklists de atores e manter o arquivo [tasks/_index.md](../tasks/_index.md) atualizado para garantir que os módulos se conectem de forma ordenada e robusta.

### A07: Authentication Failures (Falhas de Autenticação)
*   **Medida:**
    *   Delegar o motor de autenticação para o Supabase Auth (que trata hashing de senhas, expiração de sessões e rate limit contra força bruta nativamente).
    *   **Validação Criptográfica Offline (Alta Performance):** Para evitar latência de rede adicional em cada chamada à API (especialmente crítico para operações rápidas no PDV), o backend em Go não fará requisições externas para o Supabase a fim de verificar sessões. Em vez disso, o middleware do Go validará a assinatura criptográfica e a expiração do JWT localmente e em memória, utilizando a chave secreta JWT (JWT Secret) do Supabase armazenada de forma segura nas variáveis de ambiente.

### A08: Software or Data Integrity Failures (Falhas de Integridade de Dados)
*   **Medida:**
    *   Garantir a integridade dos dados na camada de banco de dados do Supabase utilizando chaves estrangeiras (`REFERENCES`), restrições (`CHECK constraints`) e triggers de integridade adequadas.
    *   Executar validações manuais consistentes dos dados recebidos nos handlers (verificando limites de caracteres, formatos de e-mail e valores positivos para valores monetários) antes de persistir alterações.
