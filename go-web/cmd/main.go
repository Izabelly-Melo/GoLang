package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/izabelly/go-web/internal/handler"
	"github.com/izabelly/go-web/internal/middlewares"
	"github.com/izabelly/go-web/internal/service"
	"github.com/joho/godotenv"
)

func main() {
	rt := chi.NewRouter()
	rt.Use(middleware.Logger)
	rt.Use(middleware.Recoverer)
	
	// carrega produtos
	service := service.NewServiceProducts("products.json")
	handler := handler.NewProductHandler(service)

	err := godotenv.Load()
	if err != nil {
		log.Println("Falha ao carregar as varáveis da .env")
		return
	}

	rt.Get("/ping", handler.Ping)
	rt.Route("/products", func(rt chi.Router) {
		rt.Use(middlewares.Auth)
		rt.Get("/", handler.GetAllProducts)
		rt.Get("/{id}", handler.GetProductByID)
		rt.Get("/search", handler.SearchProduct)
		rt.Post("/", handler.CreateProduct)
		rt.Put("/{id}", handler.UpdateProduct)
		rt.Delete("/{id}", handler.DeleteProduct)
		rt.Patch("/{id}", handler.PatchProduct)
	})

	if err := http.ListenAndServe(":8080", rt); err != nil {
		panic(err)
	}
}
