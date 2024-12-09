package service

import (
	"fmt"

	"github.com/izabelly/go-web/internal/model"
	"github.com/izabelly/go-web/internal/repository"
	"github.com/izabelly/go-web/pkg/validations"
)

type ServiceProduct struct {
	FilePath string
}

func (s *ServiceProduct) AddProduct(product model.Product) (model.Product, error) {
	listProduct, err := repository.LoadProducts()
	if err != nil {
		return model.Product{}, err
	}

	validCodeValue := validations.ValidCodeValue(product.CodeValue)
	validDate := validations.ValidDate(product.Expiration)
	validName := validations.ValidName(product.Name)
	validPrice := validations.ValidPrice(product.Price)
	validQuantity := validations.ValidQuantity(product.Quantity)

	if !validCodeValue || !validName || !validDate || !validPrice || !validQuantity {
		return model.Product{}, fmt.Errorf("Produto inválido: %+v", product)
	}

	product.ID = len(listProduct) + 1

	err = repository.AddProduct(append(listProduct, product))
	if err != nil {
		return model.Product{}, err
	}

	return product, nil
}

func (s *ServiceProduct) GetAllProducts() ([]model.Product, error) {
	listProduct, err := repository.LoadProducts()
	if err != nil {
		return nil, err
	}

	return listProduct, nil
}

func (s *ServiceProduct) GetProductByID(id int) (model.Product, error) {
	listProduct, err := repository.LoadProducts()
	if err != nil {
		return model.Product{}, err
	}

	var getProduct model.Product
	for _, prod := range listProduct {
		if prod.ID == id {
			getProduct = prod
		}
	}
	return getProduct, nil
}

func (s *ServiceProduct) GetProductsPrice(price float64) ([]model.Product, error) {
	listProduct, err := repository.LoadProducts()
	if err != nil {
		return nil, err
	}
	var getProduct []model.Product
	for _, prod := range listProduct {
		if prod.Price > price {
			getProduct = append(getProduct, prod)
		}
	}
	return getProduct, nil
}

// Retornar instancia de ServiceProduct e inicializando o FilePath com o valor
func NewServiceProducts(filePath string) *ServiceProduct {
	return &ServiceProduct{FilePath: filePath}
}
