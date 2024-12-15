package internal

// TicketAttributes is an struct that represents a ticket
type TicketAttributes struct {
	// Name represents the name of the owner of the ticket
	Name string `json:"name"`
	// Email represents the email of the owner of the ticket
	Email string `json:"email"`
	// Country represents the destination country of the ticket
	Country string `json:"country"`
	// Hour represents the hour of the ticket
	Hour string `json:"hour"`
	// Price represents the price of the ticket
	Price float64 `json:"price"`
}

// Ticket represents a ticket
type Ticket struct {
	// Id represents the id of the ticket
	Id int `json:"id"`
	// Attributes represents the attributes of the ticket
	Attributes TicketAttributes `json:"attributes"`
}

// RepositoryTicket represents the repository interface for tickets
type RepositoryTicket interface {
	// GetAll returns all the tickets
	Get() (t map[int]TicketAttributes, err error)
	// GetTicketByDestinationCountry returns the tickets filtered by destination country
	GetTicketsByDestinationCountry(country string) (t map[int]TicketAttributes, err error)
}

type ServiceTicket interface {
	// GetTotalAmountTickets returns the total amount of tickets
	GetTotalAmountTickets() (total int, err error)

	// GetTicketsAmountByDestinationCountry returns the amount of tickets filtered by destination country
	GetTicketsAmountByDestinationCountry(country string) (t map[int]TicketAttributes, err error)
	GetAverageCountry(country string) (average float64, err error)
	// GetPercentageTicketsByDestinationCountry returns the percentage of tickets filtered by destination country
	GetPercentageTicketsByDestinationCountry(country string)
}
