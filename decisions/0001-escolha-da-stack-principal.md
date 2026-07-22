# ADR-0001: Escolha do Go (Golang) e Supabase (PostgreSQL) como Stack Principal

Status: Aceito

## Contexto
Para o desenvolvimento do Mini ERP, é necessário escolher uma stack de backend que seja altamente performática, segura e de baixo consumo de recursos, viabilizando hospedagem econômica em produção. Além disso, a tecnologia precisa oferecer uma alta previsibilidade e rigidez de compilação para otimizar o fluxo de codificação cooperativo com Inteligência Artificial, e um banco de dados relacional que garanta a integridade referencial dos dados financeiros e de estoque.

## Decisão
Decidimos utilizar a linguagem **Go (Golang)** para a construção da API do backend e o **Supabase** (que provê instâncias gerenciadas de PostgreSQL) como o nosso banco de dados relacional.

As justificativas para a escolha de **Go** são:
- **Previsibilidade e Rigidez:** A simplicidade sintática do Go e a rigidez do seu compilador reduzem drasticamente bugs em produção e tornam o código gerado em cooperação com a IA previsível e robusto.
- **Eficiência e Custo:** Baixíssimo consumo de memória e processamento, facilitando o deploy em instâncias gratuitas ou econômicas (como no Render).
- **Segurança:** O ecossistema de dependências enxuto minimiza o risco de ataques na cadeia de suprimentos (*supply chain*).

As justificativas para a escolha do **Supabase** são:
- **PostgreSQL Nativo:** Garante a robustez relacional exigida por um ERP (transações ACID, chaves estrangeiras, consistência rígida).
- **Custo Inicial Zero:** O plano gratuito permite colocar a base de dados online em ambiente de produção sem qualquer custo imediato.
- **Pronto para Escalar:** Facilidade de escalabilidade futura caso o volume de dados do ERP cresça.

## Alternativas descartadas
- **Node.js (TypeScript) + MongoDB:** Descartado devido ao maior risco de instabilidade decorrente do dinamismo da linguagem, maior consumo de recursos de servidor, dependências excessivas no ecossistema npm e falta de integridade referencial estrita nativa do MongoDB (incompatível com as necessidades de um ERP).
- **PostgreSQL Gerenciado nativo no Render ou AWS RDS:** Descartados devido ao custo financeiro imediato exigido para manter instâncias ativas na nuvem, inviabilizando o deploy inicial gratuito do projeto.

## Consequências
- A API backend será desenvolvida em Go, compilada nativamente.
- A persistência utilizará a instância PostgreSQL do Supabase, conectada via string de conexão segura.
- A configuração da aplicação local e remota precisará suportar a integração com credenciais de banco de dados externas via variáveis de ambiente.
