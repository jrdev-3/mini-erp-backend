package main

import (
	"log"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func main() {
	e := echo.New()

	// Middleware de recuperação de panic do Echo v5
	e.Use(middleware.Recover())

	// Grupo versionado v1 da API
	v1 := e.Group("/api/v1")
	_ = v1

	// Inicialização do servidor com tratamento de erro nativo do Go
	if err := e.Start(":8080"); err != nil {
		log.Fatalf("shutting down the server: %v", err)
	}
}
