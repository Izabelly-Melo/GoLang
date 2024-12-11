package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/izabelly/go-web/internal/repository"
	"github.com/izabelly/go-web/internal/service"
	"github.com/stretchr/testify/require"
)

func TestHandlerProduct_GetAllProducts(t *testing.T) {
	t.Run("success to get all product", func(t *testing.T) {
		//Arrange/Given
		repo := repository.NewRepositoryProduct("../../docs/products_test.json")
		service := service.NewServiceProducts(repo)
		handlers := NewProductHandler(service)

		//Act/When
		req := httptest.NewRequest("GET", "/products", nil)
		req.Header.Set("API_TOKEN", "02101998")
		res := httptest.NewRecorder()
		handlers.GetAllProducts(res, req)
		//Assert/Then
		expectedCode := http.StatusOK
		expectedBody := `[
				{
					"id": 1,
					"name": "Oil - Margarine",
					"quantity": 439,
					"code_value": "S82254D",
					"is_published": true,
					"expiration": "15/12/2021",
					"price": 71.42
					},
					{
					"id": 2,
					"name": "Pineapple - Canned, Rings",
					"quantity": 345,
					"code_value": "M4637",
					"is_published": true,
					"expiration": "09/08/2021",
					"price": 352.79
   				}
			]`

		expectedHeader := http.Header{"Content-Type": []string{"application/json"}}
		require.Equal(t, expectedCode, res.Code)           //verifica se o código de status está correto
		require.JSONEq(t, expectedBody, res.Body.String()) //verifica se o corpo da resposta corresponde ao esperado
		require.Equal(t, expectedHeader, res.Header())     //verifica se o cabeçalho da resposta está correto
	})

}

func TestHandlerProduct_GetProductByID(t *testing.T) {

	t.Run("success to get product ID", func(t *testing.T) {
		// Arrange/Given
		repo := repository.NewRepositoryProduct("../../docs/products_test.json")
		service := service.NewServiceProducts(repo)
		handlers := NewProductHandler(service)

		rt := chi.NewRouteContext()
		rt.URLParams.Add("id", "1")

		// Act/When
		req := httptest.NewRequest("GET", "/products/{id}", nil)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rt))
		req.Header.Set("API_TOKEN", "02101998")
		res := httptest.NewRecorder()
		handlers.GetProductByID(res, req)

		// Assert/Then
		expectedCode := http.StatusOK
		expectedBody := `
		{
			"id": 1,
			"name": "Oil - Margarine",
			"quantity": 439,
			"code_value": "S82254D",
			"is_published": true,
			"expiration": "15/12/2021",
			"price": 71.42
		}`

		expectedHeader := http.Header{"Content-Type": []string{"application/json"}}
		require.Equal(t, expectedCode, res.Code)           // verifica se o código de status está correto
		require.JSONEq(t, expectedBody, res.Body.String()) // verifica se o corpo da resposta corresponde ao esperado
		require.Equal(t, expectedHeader, res.Header())     // verifica se o cabeçalho da resposta está correto
	})

	t.Run("missing param ID product", func(t *testing.T) {
		// Arrange/Given
		repo := repository.NewRepositoryProduct("../../docs/products_test.json")
		service := service.NewServiceProducts(repo)
		handlers := NewProductHandler(service)

		r := chi.NewRouter()
		r.Get("/products/{id}", handlers.GetProductByID)

		// Act/When
		req := httptest.NewRequest("GET", "/products/", nil)
		req.Header.Set("API_TOKEN", "02101998")
		res := httptest.NewRecorder()
		handlers.GetProductByID(res, req)

		// Assert/Then
		expectedCode := http.StatusBadRequest
		expectedBody := `
		{
			"message": "Invalid ID format",
			"error": true
		}`

		expectedHeader := http.Header{"Content-Type": []string{"application/json"}}
		require.Equal(t, expectedCode, res.Code)           // verifica se o código de status está correto
		require.JSONEq(t, expectedBody, res.Body.String()) // verifica se o corpo da resposta corresponde ao esperado
		require.Equal(t, expectedHeader, res.Header())     // verifica se o cabeçalho da resposta está correto
	})

	t.Run("not found ID product", func(t *testing.T) {
		// Arrange/Given
		repo := repository.NewRepositoryProduct("../../docs/products_test.json")
		service := service.NewServiceProducts(repo)
		handlers := NewProductHandler(service)

		r := chi.NewRouter()
		r.Get("/products/{id}", handlers.GetProductByID)

		// Act/When
		req := httptest.NewRequest("GET", "/products/10", nil)
		req.Header.Set("API_TOKEN", "02101998")
		res := httptest.NewRecorder()
		handlers.GetProductByID(res, req)

		// Assert/Then
		expectedCode := http.StatusBadRequest
		expectedBody := `
		{
			"message": "Invalid ID format",
			"error": true
		}`

		expectedHeader := http.Header{"Content-Type": []string{"application/json"}}
		require.Equal(t, expectedCode, res.Code)           // verifica se o código de status está correto
		require.JSONEq(t, expectedBody, res.Body.String()) // verifica se o corpo da resposta corresponde ao esperado
		require.Equal(t, expectedHeader, res.Header())     // verifica se o cabeçalho da resposta está correto
	})

	t.Run("StatusUnauthorized", func(t *testing.T) {
		// Arrange/Given
		repo := repository.NewRepositoryProduct("../../docs/products_test.json")
		service := service.NewServiceProducts(repo)
		handlers := NewProductHandler(service)

		rt := chi.NewRouter()
		rt.Get("/products/{id}", handlers.GetProductByID)

		// Act/When
		req := httptest.NewRequest("GET", "/products/10", nil)
		req.Header.Set("API_TOKEN", "1")
		res := httptest.NewRecorder()

		// Assert/Then
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

func TestHandlerProduct_SearchProduct(t *testing.T) {

	t.Run("success to get products with price > 100", func(t *testing.T) {
		// Arrange/Given
		repo := repository.NewRepositoryProduct("../../docs/products_test.json")
		service := service.NewServiceProducts(repo)
		handlers := NewProductHandler(service)

		rt := chi.NewRouter()
		rt.Get("/products/search", handlers.SearchProduct)

		// Act/When
		req := httptest.NewRequest("GET", "/products/search?priceGt=100", nil)
		req.Header.Set("API_TOKEN", "02101998")
		res := httptest.NewRecorder()

		// Assert/Then
		rt.ServeHTTP(res, req)
		expectedCode := http.StatusOK
		expectedBody := `[{
			"id": 2,
			"name": "Pineapple - Canned, Rings",
			"quantity": 345,
			"code_value": "M4637",
			"is_published": true,
			"expiration": "09/08/2021",
			"price": 352.79
		}]`

		expectedHeader := http.Header{"Content-Type": []string{"application/json"}}
		require.Equal(t, expectedCode, res.Code)           // verifica se o código de status está correto
		require.JSONEq(t, expectedBody, res.Body.String()) // verifica se o corpo da resposta corresponde ao esperado
		require.Equal(t, expectedHeader, res.Header())     // verifica se o cabeçalho da resposta está correto
	})

	t.Run("missing param Price", func(t *testing.T) {
		// Arrange/Given
		repo := repository.NewRepositoryProduct("../../docs/products_test.json")
		service := service.NewServiceProducts(repo)
		handlers := NewProductHandler(service)

		rt := chi.NewRouter()
		rt.Get("/products/search", handlers.SearchProduct)

		// Act/When
		req := httptest.NewRequest("GET", "/products/search", nil)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rt))
		req.Header.Set("API_TOKEN", "02101998")
		res := httptest.NewRecorder()

		// Assert/Then
		rt.ServeHTTP(res, req)
		expectedCode := http.StatusBadRequest
		expectedBody := `
		{
			"message": "Invalid Price format",
			"error": true
		}`

		expectedHeader := http.Header{"Content-Type": []string{"application/json"}}
		require.Equal(t, expectedCode, res.Code)           // verifica se o código de status está correto
		require.JSONEq(t, expectedBody, res.Body.String()) // verifica se o corpo da resposta corresponde ao esperado
		require.Equal(t, expectedHeader, res.Header())     // verifica se o cabeçalho da resposta está correto
	})
}
