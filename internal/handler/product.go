package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/izabelly/go-web/internal/model"
	"github.com/izabelly/go-web/internal/service"
)

type HandlerProduct struct {
	Service *service.ServiceProduct
}

func (h *HandlerProduct) Ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong\n"))
}

func (h *HandlerProduct) CreateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	var reqBody model.ReqBodyProduct
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resBody := model.ResBodyProduct{
			Message: "Erro ao criar produto!",
			Data:    nil,
			Error:   true,
		}
		json.NewEncoder(w).Encode(resBody)
		return
	}

	productBody := model.Product{
		Name:        reqBody.Name,
		Quantity:    reqBody.Quantity,
		CodeValue:   reqBody.CodeValue,
		IsPublished: reqBody.IsPublished,
		Expiration:  reqBody.Expiration,
		Price:       reqBody.Price,
	}

	response, err := h.Service.AddProduct(productBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		resBody := model.ResBodyProduct{
			Message: "Erro ao criar produto(s)!",
			Data:    nil,
			Error:   true,
		}
		json.NewEncoder(w).Encode(resBody)
		return
	}

	body := model.ResBodyProduct{
		Message: "Produto Criado",
		Data:    &response,
		Error:   false,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(body)
}

func (h *HandlerProduct) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	listProducts, err := h.Service.GetAllProducts()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Panicln(err)
		resBody := model.ResBodyProduct{
			Message: "Erro ao listar produtos!",
			Data:    nil,
			Error:   true,
		}
		json.NewEncoder(w).Encode(resBody)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(listProducts)
}

func (h *HandlerProduct) GetProductByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	param := chi.URLParam(r, "id")
	id, err := strconv.Atoi(param)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Panicln(err)
		resBody := model.ResBodyProduct{
			Message: "Erro ao encontrar ID",
			Data:    nil,
			Error:   true,
		}
		json.NewEncoder(w).Encode(resBody)
		return
	}

	listProducts, err := h.Service.GetProductByID(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		resBody := model.ResBodyProduct{
			Message: "Erro",
			Data:    nil,
			Error:   true,
		}
		json.NewEncoder(w).Encode(resBody)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(listProducts)

}

func (h *HandlerProduct) SearchProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "Application/json")
	param := r.URL.Query().Get("priceGt")
	price, err := strconv.ParseFloat(param, 64)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		resBody := model.ResBodyProduct{
			Message: "Formato de preço inválido",
			Data:    nil,
			Error:   true,
		}
		json.NewEncoder(w).Encode(resBody)
		return
	}

	products, err := h.Service.GetProductsPrice(price)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		resBody := model.ResBodyProduct{
			Message: "Falha ao listar produtos",
			Data:    nil,
			Error:   true,
		}
		json.NewEncoder(w).Encode(resBody)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}

// Retornar instancia de HandlerProduct e inicializando o service com o valor
func NewProductHandler(service *service.ServiceProduct) *HandlerProduct {
	return &HandlerProduct{Service: service}
}
