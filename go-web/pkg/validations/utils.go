package validations

import (
	"log"
	"time"

	"github.com/izabelly/go-web/internal/repository"
)

func ValidCodeValue(codeValue string, repository *repository.RepositoryProduct) bool {
	listProd, err := repository.LoadProducts()
	if err != nil {
		log.Println(err)
		return false
	}

	for _, product := range listProd {
		if codeValue == product.CodeValue {
			return false
		}
	}
	return true
}

func ValidDate(expiration string) bool {
	_, err := time.Parse("02/01/2006", expiration) // formatar a data sempre ser essa data
	if err != nil {
		return false
	}

	if len(expiration) > 10 {
		return false
	}

	return true
}

func ValidName(name string) bool {
	return name != ""
}

func ValidQuantity(quantity int) bool {
	return quantity != 0
}

func ValidPrice(price float64) bool {
	return price != 0.0
}
