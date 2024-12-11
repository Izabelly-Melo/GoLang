package service

import (
	"fmt"

	"github.com/izabelly/go-web/internal/model"
	"github.com/izabelly/go-web/internal/repository"
	"github.com/izabelly/go-web/pkg/validations"
)

type ServiceProduct struct {
	Repository *repository.RepositoryProduct
}

func (s *ServiceProduct) AddProduct(product model.Product) (model.Product, error) {
	listProduct, err := s.Repository.LoadProducts()

	validCodeValue := validations.ValidCodeValue(product.CodeValue, s.Repository)
	validDate := validations.ValidDate(product.Expiration)
	validName := validations.ValidName(product.Name)
	validPrice := validations.ValidPrice(product.Price)
	validQuantity := validations.ValidQuantity(product.Quantity)

	if !validCodeValue || !validName || !validDate || !validPrice || !validQuantity {
		return model.Product{}, fmt.Errorf("Produto inválido: %+v", product)
	}

	product.ID = len(listProduct) + 1

	err = s.Repository.AddProduct(append(listProduct, product))
	if err != nil {
		return model.Product{}, err
	}

	return product, nil
}

func (s *ServiceProduct) GetAllProducts() ([]model.Product, error) {
	listProduct, err := s.Repository.LoadProducts()
	if err != nil {
		return nil, err
	}

	return listProduct, nil
}

func (s *ServiceProduct) GetProductByID(id int) (model.Product, error) {
	listProduct, err := s.Repository.LoadProducts()
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
	listProduct, err := s.Repository.LoadProducts()
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

func (s *ServiceProduct) UpdateProduct(newProduct model.Product, id int) (model.Product, error) {
	listProduct, err := s.Repository.LoadProducts()
	if err != nil {
		return newProduct, err
	}

	if id == 0 {
		return newProduct, fmt.Errorf("Id inválido")
	}

	var updatedList []model.Product
	var updatedProd model.Product

	for _, prod := range listProduct {
		if prod.ID == id {
			validDate := validations.ValidDate(newProduct.Expiration)
			validName := validations.ValidName(newProduct.Name)
			validPrice := validations.ValidPrice(newProduct.Price)
			validQuantity := validations.ValidQuantity(newProduct.Quantity)

			if !validName || !validDate || !validPrice || !validQuantity {
				return model.Product{}, fmt.Errorf("Produto inválido: %+v", newProduct)
			}

			updatedProd = model.Product{
				ID:          id,
				Name:        newProduct.Name,
				Quantity:    newProduct.Quantity,
				CodeValue:   newProduct.CodeValue,
				IsPublished: newProduct.IsPublished,
				Expiration:  newProduct.Expiration,
				Price:       newProduct.Price,
			}
			updatedList = append(updatedList, updatedProd)
		} else {
			updatedList = append(updatedList, prod)
		}
	}

	err = s.Repository.AddProduct(updatedList)
	if err != nil {
		return model.Product{}, err
	}

	return updatedProd, nil
}

func (s *ServiceProduct) DeleteProduct(id int) error {
	listProduct, err := s.Repository.LoadProducts()
	if err != nil {
		return err
	}

	if id == 0 {
		return fmt.Errorf("Id inválido")
	}

	var updatedList []model.Product
	for i, prod := range listProduct {
		if prod.ID == id {
			updatedList = append(updatedList[:i], listProduct[i+1:]...)
		} else {
			updatedList = append(updatedList, prod)
		}
	}

	err = s.Repository.AddProduct(updatedList)
	if err != nil {
		return err
	}

	return nil
}

func (s *ServiceProduct) PatchProduct(product model.Product, id int) (model.Product, error) {
	if id == 0 {
		return model.Product{}, fmt.Errorf("Id inválido")
	}

	prod, err := s.UpdateProduct(product, id)
	if err != nil {
		return model.Product{}, err
	}

	return prod, nil
}

// Retornar instancia de ServiceProduct e inicializando o FilePath com o valor
func NewServiceProducts(repo *repository.RepositoryProduct) *ServiceProduct {
	return &ServiceProduct{Repository: repo}
}
