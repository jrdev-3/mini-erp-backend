# Stage 1: Build (Compilação do binário Go)
FROM golang:1.26-alpine AS builder

WORKDIR /build

# Copiar arquivos de dependências primeiro para aproveitar o cache do Docker
COPY go.mod go.sum ./
RUN go mod download

# Copiar o restante do código-fonte do projeto
COPY . .

# Compilar o binário estático otimizado para produção (desativa CGO para portabilidade)
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/mini-erp cmd/api/main.go

# Stage 2: Runtime (Imagem final mínima e segura)
FROM alpine:3.19

# Instalar ca-certificates obrigatórios para permitir conexões seguras (SSL/TLS) com o Supabase
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copiar apenas o binário compilado do estágio anterior (reduz tamanho final)
COPY --from=builder /app/mini-erp .

# Porta padrão de escuta (Render injetará a variável PORT dinamicamente)
EXPOSE 8080

ENTRYPOINT ["./mini-erp"]
