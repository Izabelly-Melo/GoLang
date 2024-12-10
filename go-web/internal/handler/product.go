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


func (h *HandlerProduct) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	param := chi.URLParam(r, "id")
	id, err := strconv.Atoi(param)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		resBody := model.ResBodyProduct{
			Message: "Erro ao encontrar ID",
			Data:    nil,
			Error:   true,
		}
		json.NewEncoder(w).Encode(resBody)
		return
	}

	var reqBody model.ReqBodyProduct
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resBody := model.ResBodyProduct{
			Message: "Erro ao alterar produto!",
			Data:    nil,
			Error:   true,
		}
		json.NewEncoder(w).Encode(resBody)
		return
	}

	newProduct := model.Product{
		ID:          id,
		Name:        reqBody.Name,
		Quantity:    reqBody.Quantity,
		CodeValue:   reqBody.CodeValue,
		IsPublished: reqBody.IsPublished,
		Expiration:  reqBody.Expiration,
		Price:       reqBody.Price,
	}

	products, err := h.Service.UpdateProduct(newProduct, id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resBody := model.ResBodyProduct{
			Message: "Erro ao alterar produto",
			Data:    &products,
			Error:   true,
		}
		json.NewEncoder(w).Encode(resBody)
		return
	}

	body := model.ResBodyProduct{
		Message: "Produto Alterado",
		Data:    &products,
		Error:   false,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(body)
}

func (h *HandlerProduct) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	param := chi.URLParam(r, "id")
	id, err := strconv.Atoi(param)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resBody := model.ResBodyProduct{
			Message: "ID inválido",
			Data:    nil,
			Error:   true,
		}
		json.NewEncoder(w).Encode(resBody)
		return
	}

	err = h.Service.DeleteProduct(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resBody := model.ResBodyProduct{
			Message: "Não foi possível excluir produto",
			Data:    nil,
			Error:   true,
		}
		json.NewEncoder(w).Encode(resBody)
		return
	}

	resBody := model.ResBodyProduct{
		Message: "Produto excluido",
		Data:    nil,
		Error:   false,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resBody)
}

func (h *HandlerProduct) PatchProduct(w http.ResponseWriter, r *http.Request) {
	param := chi.URLParam(r, "id")
	id, err := strconv.Atoi(param)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resBody := model.ResBodyProduct{
			Message: "ID Inválido",
			Data:    nil,
			Error:   false,
		}
		json.NewEncoder(w).Encode(resBody)
		return
	}

	var reqBody model.ReqPatchBodyProduct
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resBody := model.ResBodyProduct{
			Message: "Erro ao alterar produto!",
			Data:    nil,
			Error:   true,
		}
		json.NewEncoder(w).Encode(resBody)
		return
	}

	productSave, err := h.Service.GetProductByID(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resBody := model.ResBodyProduct{
			Message: "Erro ao tentar encontrar produto!",
			Data:    nil,
			Error:   true,
		}
		json.NewEncoder(w).Encode(resBody)
		return
	}

	name := productSave.Name
	if reqBody.Name != nil {
		name = *reqBody.Name
	}

	quantity := productSave.Quantity
	if reqBody.Quantity != nil {
		quantity = *reqBody.Quantity
	}

	codeValue := productSave.CodeValue
	if reqBody.CodeValue != nil {
		codeValue = *reqBody.CodeValue
	}

	isPublished := productSave.IsPublished
	if reqBody.IsPublished != nil {
		isPublished = *reqBody.IsPublished
	}

	expiration := productSave.Expiration
	if reqBody.Expiration != nil {
		expiration = *reqBody.Expiration
	}

	price := productSave.Price
	if reqBody.Price != nil {
		price = *reqBody.Price
	}

	resBody := model.Product{
		ID:          id,
		Name:        name,
		Quantity:    quantity,
		CodeValue:   codeValue,
		IsPublished: isPublished,
		Expiration:  expiration,
		Price:       price,
	}

	product, err := h.Service.PatchProduct(resBody, id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resBody := model.ResBodyProduct{
			Message: "Erro ao alterar",
			Data:    nil,
			Error:   true,
		}
		json.NewEncoder(w).Encode(resBody)
		return
	}

	body := model.ResBodyProduct{
		Message: "Produto Alterado",
		Data:    &product,
		Error:   false,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(body)
}

// Retornar instancia de HandlerProduct e inicializando o service com o valor
func NewProductHandler(service *service.ServiceProduct) *HandlerProduct {
	return &HandlerProduct{Service: service}
}
