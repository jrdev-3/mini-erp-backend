# ADR-0003: Uso de SQL Puro com o Driver pgx para Banco de Dados

Status: Aceito

## Contexto
Um sistema ERP lida constantemente com fluxos financeiros e de inventário cruciais que requerem total confiabilidade transacional. A camada de persistência precisa acessar o banco PostgreSQL no Supabase com velocidade máxima, além de expor as consultas SQL de forma transparente para permitir validações explícitas de banco de dados no código e em análises de portfólio.

## Decisão
Decidimos utilizar **SQL Puro** para todas as interações com o banco de dados, utilizando o driver moderno **pgx** (`github.com/jackc/pgx/v5`) em sua interface direta, gerenciando conexões por meio de um Pool de conexões (`pgxpool`).

- **Transparência e Controle:** A escrita explícita de queries SQL garante total controle de otimização, joins, transações e índices na base de dados, permitindo que o piloto humano audite e valide com facilidade a modelagem conceitual aplicada pelo assistente.
- **Máxima Performance:** A ausência de mapeamento complexo em tempo de execução via reflexão (*reflection*) evita gargalos de latência no backend.
- **Facilidade de Depuração:** Qualquer erro de banco de dados gerado retornará a query exata executada, facilitando a identificação e a correção rápida de falhas pela IA e pelo desenvolvedor.
- **Driver Otimizado:** O driver `pgx/v5` oferece suporte robusto a pools de conexões nativos (`pgxpool`), suporte automático a recursos avançados de tipo de dados do Postgres e suporte a TLS exigido pelo Supabase.

## Alternativas descartadas
- **GORM (ORM):** Descartado para evitar a opacidade de queries (consultas "mágicas" geradas implicitamente), o overhead de processamento que diminui a performance geral e a perda de controle estrito sobre como o banco é consultado em operações sensíveis do ERP.
- **sqlc (Gerador de Código):** Descartado neste estágio inicial para manter a stack simples e focada na escrita direta de SQL dentro do próprio fluxo Go, sem a necessidade de instalar ferramentas externas de compilação de código de banco de dados.

## Consequências
- Toda a lógica de persistência e repositórios em `internal/` utilizará o objeto `*pgxpool.Pool` para executar comandos SQL diretos.
- As consultas e scripts SQL (como criação de tabelas e sementes de dados) serão armazenados de forma limpa e explícita nas implementações de infraestrutura do projeto.
