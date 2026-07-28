package main

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	echoSwagger "github.com/swaggo/echo-swagger/v2"

	_ "github.com/jrdev-3/mini-erp-backend/docs"
	"github.com/jrdev-3/mini-erp-backend/internal/auth"
	customMiddleware "github.com/jrdev-3/mini-erp-backend/internal/middleware"
)

// @title           Mini ERP API
// @version         1.0
// @description     API de backend para o sistema de Mini ERP, desenvolvida em Go seguindo Clean Architecture.
// @termsOfService  http://swagger.io/terms/

// @contact.name   Suporte JR DEV
// @contact.url    https://github.com/jrdev-3

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        Authorization
// @description                 Insira o token JWT no formato: Bearer <TOKEN>

func main() {
	e := echo.New()

	// Middlewares globais nativos do Echo v5
	e.Use(middleware.Recover())

	// 1. Carregar variáveis de ambiente (injetadas dinamicamente em produção no Render)
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("A variável de ambiente DATABASE_URL é obrigatória")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("A variável de ambiente JWT_SECRET é obrigatória")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 2. Inicializar o pool de conexões com o PostgreSQL (Supabase) via pgxpool
	dbPool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("Erro ao inicializar o pool de conexões do pgx: %v", err)
	}
	defer dbPool.Close()

	// Testar o ping de inicialização da conexão
	if err := dbPool.Ping(context.Background()); err != nil {
		log.Fatalf("Falha de conexão física com o banco de dados (ping): %v", err)
	}

	// Grupo versionado v1 da API
	v1 := e.Group("/api/v1")

	// 3. Montar a fiação da Clean Architecture do módulo de Autenticação (auth)
	authRepo := auth.NewRepository(dbPool)
	authService := auth.NewService(authRepo, jwtSecret)
	authHandler := auth.NewHandler(authService)

	// 4. Instanciar middlewares customizados para injeção de dependência nas rotas
	authMiddleware := customMiddleware.Auth([]byte(jwtSecret))
	adminMiddleware := customMiddleware.RBAC("ADMIN")

	// 5. Registrar as rotas no Echo
	auth.RegisterRoutes(v1, authHandler, authMiddleware, adminMiddleware)

	// 6. Expor documentação do Swagger UI em ambientes locais/staging (oculta em produção)
	if os.Getenv("APP_ENV") != "production" {
		e.GET("/swagger/*", echoSwagger.WrapHandler)
	}

	// Inicialização do servidor HTTP com suporte à porta dinâmica do Render
	log.Printf("[SERVER] Servidor do Mini ERP iniciado com sucesso na porta %s", port)
	if err := e.Start(":" + port); err != nil {
		log.Fatalf("shutting down the server: %v", err)
	}
}
