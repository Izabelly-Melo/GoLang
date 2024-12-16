package service

import (
	"app/internal"
	"errors"
)

// ServiceTicketDefault represents the default service of the tickets
type ServiceTicketDefault struct {
	// rp represents the repository of the tickets
	rp internal.RepositoryTicket
}

// NewServiceTicketDefault creates a new default service of the tickets
func NewServiceTicketDefault(rp internal.RepositoryTicket) *ServiceTicketDefault {
	return &ServiceTicketDefault{
		rp: rp,
	}
}

// GetTotalTickets returns the total number of tickets
func (s *ServiceTicketDefault) GetTotalAmountTickets() (total int, err error) {
	tickets, err := s.rp.Get()
	return len(tickets), err
}

func (s *ServiceTicketDefault) GetTicketsAmountByDestinationCountry(country string) (t map[int]internal.TicketAttributes, err error) {
	t, err = s.rp.GetTicketsByDestinationCountry(country)
	return
}

func (s *ServiceTicketDefault) GetAverageCountry(country string) (average float64, err error) {
	total, err := s.GetTotalAmountTickets()
	if err != nil {
		err = errors.New("failed to retrieve the total amount of tickets")
		return
	}
	dest, err := s.rp.GetTicketsByDestinationCountry(country)
	if err != nil {
		err = errors.New("the specified country was not found")
		return
	}

	if len(dest) < 1 {
		err = errors.New("no tickets available for the specified country")
		return
	}

	average = float64(len(dest)) / float64(total)

	return
}
