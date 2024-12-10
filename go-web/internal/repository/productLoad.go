package repository

import (
	"encoding/json"
	"log"
	"os"

	"github.com/izabelly/go-web/internal/model"
)

func LoadProducts() ([]model.Product, error) {
	file, err := os.Open("/Users/idmelo/Documents/git/GoLang/go-web/docs/products.json")
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

func AddProduct(products []model.Product) error {
	//converte os produtos para json
	file, err := json.MarshalIndent(products, "", " ")
	if err != nil {
		log.Println("Erro: ", err)
		return err
	}
	err = os.WriteFile("/Users/idmelo/Documents/git/GoLang/go-web/docs/products.json", file, 0666) //0666 permite leitura e escrita
	if err != nil {
		log.Println("Erro ao gravar no arquivo", err)
		return err
	}
	return nil
}
