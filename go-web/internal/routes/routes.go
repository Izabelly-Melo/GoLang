package routes

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/izabelly/go-web/internal/handler"
	"github.com/izabelly/go-web/internal/middlewares"
)

func Routes(h *handler.HandlerProduct) http.Handler {
	rt := chi.NewRouter()
	rt.Use(middleware.Logger)
	rt.Use(middleware.Recoverer)

	rt.Route("/products", func(rt chi.Router) {
		rt.Use(middlewares.Auth)
		rt.Get("/", h.GetAllProducts)
		rt.Get("/{id}", h.GetProductByID)
		rt.Get("/search", h.SearchProduct)
		rt.Post("/", h.CreateProduct)
		rt.Put("/{id}", h.UpdateProduct)
		rt.Delete("/{id}", h.DeleteProduct)
		rt.Patch("/{id}", h.PatchProduct)
	})

	return rt
}
