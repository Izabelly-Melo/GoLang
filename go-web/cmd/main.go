package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/izabelly/go-web/internal/handler"
	"github.com/izabelly/go-web/internal/service"
)

func main() {
	rt := chi.NewRouter()

	// carrega produtos
	service := service.NewServiceProducts("products.json")
	handler := handler.NewProductHandler(service)

	rt.Get("/ping", handler.Ping)
	rt.Route("/products", func(rt chi.Router) {
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
