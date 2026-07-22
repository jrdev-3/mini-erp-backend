# 🚀 Mini ERP API - Golang

> 💡 **Nota Importante para Recrutadores:** Todo o código desta aplicação é escrito e estruturado por Inteligência Artificial (**Antigravity** da Google DeepMind). O desenvolvedor (**Joaby**) atua como o guia humano, piloto e arquiteto do projeto — definindo requisitos, debatendo caminhos, conversando sobre decisões de design de software e orientando na depuração e correção de bugs. O histórico de commits e diálogos documenta essa parceria em detalhes.

> **Status do Projeto:** 🏗️ Em Desenvolvimento / Concepção

Uma API RESTful robusta, performática e escalável desenvolvida em **Go (Golang)** para servir como o núcleo de um ecossistema de **Mini ERP (Enterprise Resource Planning)**. A API será hospedada na plataforma **Render** e consumida por um front-end independente (a ser desenvolvido posteriormente).

---

## 🗺️ Visão Geral e Módulos do Sistema

O objetivo deste Mini ERP é fornecer uma base sólida para a gestão de micro e pequenas empresas. O sistema será composto pelos seguintes módulos principais:

### 👥 1. Cadastro de Clientes e Fornecedores (CRM Simplificado)
* Cadastro completo de clientes (Pessoa Física/Jurídica).
* Histórico de compras por cliente.
* Gestão de fornecedores parceiros.

### 📦 2. Controle de Estoque e Catálogo de Produtos
* Cadastro de produtos com SKU, código de barras, preço de custo, preço de venda e margem de lucro.
* Controle de estoque mínimo e alertas de reposição.
* Histórico de movimentações de estoque (entradas por compra, saídas por venda, ajustes manuais).

### 🛒 3. Gestão de Vendas (Pedidos)
* Fluxo completo de pedidos: Orçamento -> Aprovado -> Faturado -> Cancelado.
* Cálculo automático de totais, descontos e impostos simplificados.
* Associação de clientes e produtos a cada venda.

### 💰 4. Financeiro Simplificado
* Fluxo de Contas a Receber (gerado automaticamente a partir de vendas faturadas).
* Fluxo de Contas a Pagar (gerado a partir de compras ou inserções manuais).
* Fluxo de caixa básico.

---

## 🛠️ Stack Tecnológica

O projeto foi projetado utilizando tecnologias modernas e eficientes do ecossistema Go:

* **Linguagem:** [Go (Golang)](https://go.dev/) (foco em performance, concorrência nativa e baixo consumo de recursos).
* **Framework Web / Roteador:** [Go-Chi](https://github.com/go-chi/chi) ou [Gin Gonic](https://github.com/gin-gonic/gin) (a ser definido para roteamento limpo e middlewares).
* **Banco de Dados:** [PostgreSQL](https://www.postgresql.org/) (banco relacional robusto para garantir a integridade dos dados financeiros e de estoque).
* **Persistência / ORM:** [GORM](https://gorm.io/) ou SQL puro com [sqlc](https://sqlc.dev/) (garantindo consultas performáticas e type-safety).
* **Migrações:** [golang-migrate](https://github.com/golang-migrate/migrate) para controle de versão do esquema do banco de dados.
* **Documentação:** [Swagger / OpenAPI](https://swagger.io/) para que a API seja facilmente testável e integrada com o front-end.
* **Containers:** [Docker](https://www.docker.com/) e Docker Compose para padronização do ambiente de desenvolvimento.
* **Hospedagem / Infra:** [Render](https://render.com/) para deploy da API e banco de dados gerenciado.

---

## 🏛️ Arquitetura do Software

Para demonstrar qualidade técnica voltada para o mercado de trabalho, utilizaremos os princípios da **Clean Architecture** (Arquitetura Limpa) e **Domain-Driven Design (DDD)** simplificado. 

A estrutura de pastas seguirá o padrão clássico da comunidade Go:

```text
├── cmd/
│   └── api/                # Ponto de entrada da aplicação (main.go)
├── internal/
│   ├── domain/             # Entidades de negócio e regras puras
│   ├── customer/           # Módulo de Clientes (Handlers, Usecases, Repositories)
│   ├── product/            # Módulo de Produtos/Estoque
│   ├── sale/               # Módulo de Vendas
│   ├── financial/          # Módulo Financeiro
│   └── platform/
│       ├── database/       # Configuração do banco de dados e migrações
│       └── web/            # Middlewares, tratamento de erros e respostas HTTP
├── pkg/                    # Código utilitário compartilhado que pode ser exportado
├── docker-compose.yml      # Infraestrutura local (Postgres, PGAdmin, etc.)
├── Makefile                # Atalhos para comandos de desenvolvimento
└── README.md
```

### Princípios Aplicados:
1. **Independência de Frameworks:** A regra de negócio não depende de bibliotecas externas de entrega (HTTP).
2. **Testabilidade:** O código é estruturado de forma a facilitar testes unitários e de integração utilizando mocks.
3. **Inversão de Dependência:** O domínio dita as interfaces, a infraestrutura apenas as implementa.

---

## 🚀 Como Executar o Projeto Localmente

*(Instruções preliminares. Serão atualizadas conforme o código for desenvolvido).*

### Pré-requisitos
* [Go 1.22+](https://go.dev/dl/) instalado.
* [Docker](https://www.docker.com/products/docker-desktop/) & [Docker Compose](https://docs.docker.com/compose/) instalados.

### Passos
1. Clone o repositório:
   ```bash
   git clone <link-do-repositorio>
   cd mini-erp-backend
   ```
2. Configure as variáveis de ambiente criando um arquivo `.env` com base no `.env.example`.
3. Suba os serviços de infraestrutura (banco de dados):
   ```bash
   docker-compose up -d
   ```
4. Execute as migrações do banco de dados (assim que configurado):
   ```bash
   make migrate-up
   ```
5. Rode a aplicação em modo de desenvolvimento:
   ```bash
   go run cmd/api/main.go
   ```

---

## 📈 Integração Contínua & Deploy (Render)

A API será configurada para deploy automático no **Render** conectado ao branch `main` do GitHub.
Toda alteração integrada à branch principal disparará um build automático e o deploy da nova versão da aplicação, garantindo uma esteira de entrega contínua (CD).

---

## 🤝 Desenvolvimento Orientado a IA e Engenharia de Prompt

Este projeto serve também como um caso de estudo sobre como um engenheiro de software pode atuar em pareria com Inteligências Artificiais avançadas. O repositório documenta de forma transparente (através do histórico de commits e mensagens) a colaboração ativa entre o desenvolvedor (guiando as decisões de design, requisitos e regras de negócio) e o assistente **Antigravity** (responsável por propor arquitetura, gerar boilerplate, implementar lógica e realizar refatorações).

---
Criado com ❤️ por [Joaby Rodrigues da Silva](https://github.com/jrdev-3) em colaboração com Antigravity (Google DeepMind).
