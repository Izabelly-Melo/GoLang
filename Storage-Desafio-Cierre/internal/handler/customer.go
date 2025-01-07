package handler

import (
	"net/http"

	"app/internal"

	"github.com/bootcamp-go/web/request"
	"github.com/bootcamp-go/web/response"
)

// NewCustomersDefault returns a new CustomersDefault
func NewCustomersDefault(sv internal.ServiceCustomer) *CustomersDefault {
	return &CustomersDefault{sv: sv}
}

// CustomersDefault is a struct that returns the customer handlers
type CustomersDefault struct {
	// sv is the customer's service
	sv internal.ServiceCustomer
}

// CustomerJSON is a struct that represents a customer in JSON format
type CustomerJSON struct {
	Id        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Condition int    `json:"condition"`
}

type CustomerConditionJSON struct {
	Condition string  `json:"condition"`
	Total     float64 `json:"total"`
}

// GetAll returns all customers
func (h *CustomersDefault) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		// ...

		// process
		c, err := h.sv.FindAll()
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "error getting customers")
			return
		}

		// response
		// - serialize
		csJSON := make([]CustomerJSON, len(c))
		for ix, v := range c {
			csJSON[ix] = CustomerJSON{
				Id:        v.Id,
				FirstName: v.FirstName,
				LastName:  v.LastName,
				Condition: v.Condition,
			}
		}
		response.JSON(w, http.StatusOK, map[string]any{
			"message": "customers found",
			"data":    csJSON,
		})
	}
}

// RequestBodyCustomer is a struct that represents the request body for a customer
type RequestBodyCustomer struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Condition int    `json:"condition"`
}

// Create creates a new customer
func (h *CustomersDefault) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		// - body
		var reqBody RequestBodyCustomer
		err := request.JSON(r, &reqBody)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "error deserializing request body")
			return
		}

		// process
		// - deserialize
		c := internal.Customer{
			CustomerAttributes: internal.CustomerAttributes{
				FirstName: reqBody.FirstName,
				LastName:  reqBody.LastName,
				Condition: reqBody.Condition,
			},
		}
		// - save
		err = h.sv.Save(&c)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "error saving customer")
			return
		}

		// response
		// - serialize
		cs := CustomerJSON{
			Id:        c.Id,
			FirstName: c.FirstName,
			LastName:  c.LastName,
			Condition: c.Condition,
		}
		response.JSON(w, http.StatusCreated, map[string]any{
			"message": "customer created",
			"data":    cs,
		})
	}
}

func (h *CustomersDefault) GetConditionsCustomer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		c, err := h.sv.GetConditionsCustomer()
		if err != nil {
			response.Error(w, http.StatusInternalServerError, err.Error())
			return
		}

		csJSON := make([]CustomerConditionJSON, len(c))
		for ix, v := range c {
			csJSON[ix] = CustomerConditionJSON{
				Condition: v.Condition,
				Total:     v.Total,
			}
		}
		response.JSON(w, http.StatusOK, map[string]any{
			"message": "ativos e inativos",
			"data":    csJSON,
		})
	}
}

type CustomerMoreActivesJSON struct {
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Amount    float64 `json:"amount"`
}

func (h *CustomersDefault) GetCustomersMoreActives() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := h.sv.GetCustomersMoreActives()
		if err != nil {
			response.Error(w, http.StatusInternalServerError, err.Error())
			return
		}

		csJSON := make([]CustomerMoreActivesJSON, len(c))
		for ix, v := range c {
			csJSON[ix] = CustomerMoreActivesJSON{
				FirstName: v.FirstName,
				LastName:  v.LastName,
				Amount:    v.Amount,
			}
		}
		response.JSON(w, http.StatusOK, map[string]any{
			"message": "customer",
			"data":    csJSON,
		})
	}
}
