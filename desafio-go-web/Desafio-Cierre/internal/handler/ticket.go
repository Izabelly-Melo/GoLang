package handler

import (
	"app/internal/service"
	"net/http"

	"github.com/bootcamp-go/web/response"
	"github.com/go-chi/chi/v5"
)

func NewHandlerTicketDefault(sv *service.ServiceTicketDefault) *HandlerTicketDefault {
	return &HandlerTicketDefault{
		sv: sv,
	}
}

type HandlerTicketDefault struct {
	// sv is the service that will be used by the handler
	sv *service.ServiceTicketDefault
}

func (h *HandlerTicketDefault) GetTotalAmountTickets(w http.ResponseWriter, r *http.Request) {
	tickets, err := h.sv.GetTotalAmountTickets()
	if err != nil {
		response.JSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"message": "Total Countries",
		"data":    tickets,
	})
}

func (h *HandlerTicketDefault) GetTicketsAmountByDestinationCountry(w http.ResponseWriter, r *http.Request) {
	country := chi.URLParam(r, "dest")

	tickets, err := h.sv.GetTicketsAmountByDestinationCountry(country)
	if err != nil {
		response.JSON(w, http.StatusNotFound, map[string]string{
			"error": "No tickets found for the specified country: " + country,
		})
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"message": "Total tickets:",
		"data":    len(tickets),
	})
}

func (h *HandlerTicketDefault) GetAverageCountry(w http.ResponseWriter, r *http.Request) {
	country := chi.URLParam(r, "dest")
	tickets, err := h.sv.GetAverageCountry(country)
	if err != nil {
		response.JSON(w, http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"message": "Average tickets for the country " + country + ":",
		"data":    tickets,
	})

}
