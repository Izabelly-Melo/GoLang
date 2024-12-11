package repository

import (
	"encoding/json"
	"log"
	"os"

	"github.com/izabelly/go-web/internal/model"
)

type RepositoryProduct struct {
	FilePath string
}

func NewRepositoryProduct(filePath string) *RepositoryProduct {
	return &RepositoryProduct{FilePath: filePath}
}

func (r *RepositoryProduct) LoadProducts() ([]model.Product, error) {
	file, err := os.Open(r.FilePath)
	if err != nil {
		log.Println("Erro ao abrir arquivo", err)
		return nil, err
	}

	defer file.Close()

	var products []model.Product
	jsonParser := json.NewDecoder(file)

	err = jsonParser.Decode(&products)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (r *RepositoryProduct) AddProduct(products []model.Product) error {
	//converte os produtos para json
	file, err := json.MarshalIndent(products, "", " ")
	if err != nil {
		log.Println("Erro: ", err)
		return err
	}
	err = os.WriteFile(r.FilePath, file, 0666) //0666 permite leitura e escrita
	if err != nil {
		log.Println("Erro ao gravar no arquivo", err)
		return err
	}
	return nil
}
