# AGENTS.md

Este é o arquivo âncora de contexto para a API do Mini ERP. Ele deve ser lido integralmente no início de cada sessão para garantir o alinhamento de diretrizes técnicas e operacionais.

## 🛠️ Stack Tecnológica
*   **Backend:** Go (Golang) utilizando Clean Architecture
*   **Banco de Dados:** PostgreSQL (Docker local / gerenciado no Render)
*   **Hospedagem:** Render (deploy contínuo a partir da branch `main`)
*   **Documentação:** Swagger/OpenAPI como fonte da verdade de tipos para integrações

## 📌 Convenções e Regras Não Negociáveis
1.  **Mensagens de Commit:** Devem ser escritas sempre em **português do Brasil** (ex: `feat: adiciona conexão com banco`, `fix: corrige validação de SKU`).
2.  **Piloto Humano:** O desenvolvedor (**Joaby**) atua como o guia do projeto (planejamento, design de regras de negócio, testes manuais e direção em correções). A IA é responsável pela execução e estruturação do código técnico.
3.  **Consumo de Contexto:** Sempre leia a tarefa ativa listada no índice de tasks antes de codificar.
4.  **Imunidade a Bugs:** Erros recorrentes resolvidos devem ser documentados em `agents/skills/<nome-da-skill>/skill.md`.
5.  **Arquitetura Isolada:** A lógica de domínio e regras de negócio devem residir em `internal/`, isoladas de implementações externas de banco de dados ou roteadores HTTP.
6.  **Nomenclatura e Idioma do Código:** Todas as variáveis, pastas do projeto, constantes, structs, funções e comentários de código técnico devem ser escritos em **inglês** e possuir nomes altamente **explícitos** (sem abreviações ambíguas).
7.  **Submissão Estrita:** A IA nunca deve criar, modificar ou executar código, comandos ou configurações sem a autorização prévia ou instrução direta do usuário. A IA atua estritamente como ferramenta executora sob a tutela do piloto humano, sem impor opiniões próprias de design ou comportamento.

## 🗺️ Onde Encontrar Mais Contexto
*   **Tarefa Atual:** Veja o progresso macro em [tasks/_index.md](tasks/_index.md).
*   **Visão do Sistema:** Entenda a estrutura no mapa [architecture/overview.md](architecture/overview.md).
*   **Histórico de Decisões:** Justificativas e ADRs estão salvas em [decisions/](decisions/).
*   **Habilidades e Soluções:** Evite antipadrões consultando [agents/skills/](agents/skills/).
*   **Rotinas de Operação:** Guias passo a passo em [playbooks/](playbooks/).
