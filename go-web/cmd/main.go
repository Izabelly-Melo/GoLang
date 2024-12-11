package main

import (
	"log"
	"net/http"

	"github.com/izabelly/go-web/internal/handler"
	"github.com/izabelly/go-web/internal/repository"
	"github.com/izabelly/go-web/internal/routes"
	"github.com/izabelly/go-web/internal/service"
	"github.com/joho/godotenv"
)

func main() {
	// carrega produtos
	repo := repository.NewRepositoryProduct("./docs/products.json")
	service := service.NewServiceProducts(repo)
	handler := handler.NewProductHandler(service)

	err := godotenv.Load()
	if err != nil {
		log.Println("Falha ao carregar as varáveis da .env")
		return
	}
	rt := routes.Routes(handler)

	if err := http.ListenAndServe(":8080", rt); err != nil {
		panic(err)
	}
}
