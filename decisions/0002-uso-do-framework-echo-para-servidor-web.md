# ADR-0002: Uso do Framework Echo para Servidor Web e Roteamento HTTP

Status: Aceito

## Contexto
Toda API RESTful necessita de um mecanismo para gerenciar rotas, middlewares de segurança, desserialização de dados de requisição (JSON) e formatação de respostas. Embora a biblioteca padrão de Go (`net/http`) seja excelente, o desenvolvimento de um ERP exige uma estrutura de middlewares robusta (CORS, logs, recuperação de panic, segurança) e facilidade de manipulação de requisições sem redundância de código.

## Decisão
Decidimos utilizar o framework **Echo** (`github.com/labstack/echo/v4`) para gerenciar as rotas, middlewares e o ciclo de requisição/resposta HTTP da nossa API.

As justificativas para a escolha do Echo são:
- **Alta Performance:** O roteamento baseado em árvore Radix minimiza o consumo de memória e a latência de despacho de rotas.
- **Contexto Centralizado:** O uso do `echo.Context` agrupa operações de parsing de JSON, respostas HTTP e tratamento de query params de forma consistente. Isso fornece alta previsibilidade de código para o trabalho cooperativo com IA.
- **Middlewares Prontos:** O ecossistema do Echo provê middlewares estáveis e testados para Logger, Recover, CORS e segurança, poupando a necessidade de escrever soluções proprietárias para estes fins básicos.

## Alternativas descartadas
- **Go-Chi:** Descartado pois, embora leve, exige a escrita manual ou integração de bibliotecas de terceiros para parsing avançado de dados e formatação simplificada de respostas.
- **Gin Gonic:** Descartado pois, apesar de prático, acopla muito a aplicação a padrões não-padrão de Go e tem um runtime ligeiramente mais pesado que o Echo.
- **net/http (Biblioteca Padrão do Go):** Descartada devido à necessidade de escrever muito código boilerplate para lidar com tarefas básicas como middlewares, bindings de JSON e padronização de respostas de erro.

## Consequências
- Os handlers da API seguirão a assinatura padrão `func(c echo.Context) error`.
- Toda a configuração de middlewares e roteamento centralizados ficará localizada na inicialização do servidor HTTP, facilitando a manutenção e a legibilidade do código.
