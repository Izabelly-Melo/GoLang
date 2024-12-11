package handler

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/izabelly/go-web/internal/repository"
	"github.com/izabelly/go-web/internal/service"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func TestHandlerProduct_CreateProduct(t *testing.T) {

	t.Run("Success to create new product", func(t *testing.T) {
		repo := repository.NewRepositoryProduct("../../docs/products_test.json")
		service := service.NewServiceProducts(repo)
		handlers := NewProductHandler(service)

		rt := chi.NewRouter()
		rt.Post("/products", handlers.CreateProduct)

		json := `
		{
			"name": "album",
			"quantity": 43,
			"code_value": "AAAA",
			"is_published": true,
			"expiration": "02/10/2021",
			"price": 800
		}`

		req := httptest.NewRequest("POST", "/products", bytes.NewReader([]byte(json)))
		req.Header.Set("API_TOKEN", "02101998")
		res := httptest.NewRecorder()
		rt.ServeHTTP(res, req)

		log.Println("Response Status Code:", res.Code)
		log.Println("Response Body:", res.Body.String())
		expectedCode := http.StatusCreated
		expectedBody := `{
				"message": "Produto Criado",
				"data": {
					"id": 3,
					"name": "album",
					"quantity": 43,
					"code_value": "AAAA",
					"is_published": true,
					"expiration": "02/10/2021",
					"price": 800
				},
				"error": false
			}`

		expectedHeader := http.Header{"Content-Type": []string{"application/json"}}
		require.Equal(t, expectedCode, res.Code)           // verifica se o código de status está correto
		require.JSONEq(t, expectedBody, res.Body.String()) // verifica se o corpo da resposta corresponde ao esperado
		require.Equal(t, expectedHeader, res.Header())     // verifica se o cabeçalho da resposta está correto
	})

	t.Run("failed to create new product", func(t *testing.T) {
		repo := repository.NewRepositoryProduct("../../docs/products_test.json")
		service := service.NewServiceProducts(repo)
		handlers := NewProductHandler(service)

		rt := chi.NewRouter()
		rt.Post("/products", handlers.CreateProduct)

		json := `
		{
				"name": "album",
				"quantity": 43,
				"code_value": "AAAA",
				"is_published": true,
				"expiration": "02/10/2021"
		}`

		req := httptest.NewRequest("POST", "/products", bytes.NewReader([]byte(json)))
		req.Header.Set("API_TOKEN", "02101998")
		res := httptest.NewRecorder()
		rt.ServeHTTP(res, req)

		expectedCode := http.StatusBadRequest
		expectedBody := `{
				"message": "Failed to update product",
				"error": true
			}`

		expectedHeader := http.Header{"Content-Type": []string{"application/json"}}
		require.Equal(t, expectedCode, res.Code)           // verifica se o código de status está correto
		require.JSONEq(t, expectedBody, res.Body.String()) // verifica se o corpo da resposta corresponde ao esperado
		require.Equal(t, expectedHeader, res.Header())     // verifica se o cabeçalho da resposta está correto
	})

	t.Run("StatusUnauthorized", func(t *testing.T) {
		repo := repository.NewRepositoryProduct("../../docs/products_test.json")
		service := service.NewServiceProducts(repo)
		handlers := NewProductHandler(service)
		godotenv.Load("API_TOKEN")
		os.Setenv("API_TOKEN", "11")
		rt := chi.NewRouter()
		rt.Post("/products", handlers.CreateProduct)

		json := `
		{
				"name": "album",
				"quantity": 43,
				"code_value": "AAAA",
				"is_published": true,
				"expiration": "02/10/2021"
		}`

		req := httptest.NewRequest("POST", "/products", bytes.NewReader([]byte(json)))
		res := httptest.NewRecorder()

		rt.ServeHTTP(res, req)
		expectedCode := http.StatusUnauthorized
		expectedBody := `
		{
			"message": "unauthorized",
			"error": true
		}`

		expectedHeader := http.Header{"Content-Type": []string{"application/json"}}
		require.Equal(t, expectedCode, res.Code)           // verifica se o código de status está correto
		require.JSONEq(t, expectedBody, res.Body.String()) // verifica se o corpo da resposta corresponde ao esperado
		require.Equal(t, expectedHeader, res.Header())     // verifica se o cabeçalho da resposta está correto
	})
}
