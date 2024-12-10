package handler

import (
	"encoding/json"
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
		handleError(w, http.StatusBadRequest, "Invalid request body")
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
		handleError(w, http.StatusBadRequest, "Failed to create product")
		return
	}

	body := model.ResBodyProduct{
		Message: "Produto Criado",
		Data:    &response,
		Error:   false,
	}

	respondJSON(w, http.StatusCreated, body)
}

func (h *HandlerProduct) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	listProducts, err := h.Service.GetAllProducts()
	if err != nil {
		handleError(w, http.StatusInternalServerError, "Failed to retrieve products")
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
		handleError(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	listProducts, err := h.Service.GetProductByID(id)
	if err != nil {
		handleError(w, http.StatusInternalServerError, "Failed to retrieve product")
		return
	}

	if listProducts.ID == 0 {
		resBody := model.ResBodyProduct{
			Message: "Produto não existe",
			Data:    nil,
			Error:   false,
		}
		respondJSON(w, http.StatusOK, resBody)

	} else {
		respondJSON(w, http.StatusOK, listProducts)
	}
}

func (h *HandlerProduct) SearchProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "Application/json")
	param := r.URL.Query().Get("priceGt")
	price, err := strconv.ParseFloat(param, 64)
	if err != nil {
		handleError(w, http.StatusBadRequest, "Invalid Price format")
		return
	}

	products, err := h.Service.GetProductsPrice(price)
	if err != nil {
		handleError(w, http.StatusInternalServerError, "Failed to retrieve product")
		return
	}

	respondJSON(w, http.StatusOK, products)
}

func (h *HandlerProduct) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	param := chi.URLParam(r, "id")
	id, err := strconv.Atoi(param)
	if err != nil {
		handleError(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	var reqBody model.ReqBodyProduct
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		handleError(w, http.StatusBadRequest, "Failed to update product")
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
		handleError(w, http.StatusInternalServerError, "Failed to update product")
		return
	}

	body := model.ResBodyProduct{
		Message: "Produto Alterado",
		Data:    &products,
		Error:   false,
	}

	respondJSON(w, http.StatusOK, body)
}

func (h *HandlerProduct) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	param := chi.URLParam(r, "id")
	id, err := strconv.Atoi(param)
	if err != nil {
		handleError(w, http.StatusBadRequest, "Invalid Id format")
		return
	}

	err = h.Service.DeleteProduct(id)
	if err != nil {
		handleError(w, http.StatusInternalServerError, "Failed to delete product")
		return
	}

	resBody := model.ResBodyProduct{
		Message: "Produto excluido",
		Data:    nil,
		Error:   false,
	}

	respondJSON(w, http.StatusOK, resBody)
}

func (h *HandlerProduct) PatchProduct(w http.ResponseWriter, r *http.Request) {
	param := chi.URLParam(r, "id")
	id, err := strconv.Atoi(param)
	if err != nil {
		handleError(w, http.StatusBadRequest, "Invalid Id format")
		return
	}

	var reqBody model.ReqPatchBodyProduct
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		handleError(w, http.StatusInternalServerError, "Failed to update product")
		return
	}

	productSave, err := h.Service.GetProductByID(id)
	if err != nil {
		handleError(w, http.StatusInternalServerError, "Failed to retrieve product")
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
		handleError(w, http.StatusInternalServerError, "Failed to update product")
		return
	}

	body := model.ResBodyProduct{
		Message: "Produto Alterado",
		Data:    &product,
		Error:   false,
	}

	respondJSON(w, http.StatusOK, body)
}

// Retornar instancia de HandlerProduct e inicializando o service com o valor
func NewProductHandler(service *service.ServiceProduct) *HandlerProduct {
	return &HandlerProduct{Service: service}
}
