package tickets

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Ticket struct {
	ID          string
	Name        string
	Email       string
	Destination string
	Time        time.Time
	Price       float64
}

const (
	Madrugada = "madrugada"
	Manha     = "manha"
	Tarde     = "tarde"
	Noite     = "noite"
)

var fileName string = "/Users/idmelo/Documents/GO/desafio-go-bases/tickets.csv"

// ejemplo 1
func GetTotalTickets(destination string) (int, error) {
	passagem, err := listTicket()
	if err != nil {
		fmt.Println("Erro ao preencher tickets:", err)
	}

	total := 0
	for _, ticket := range passagem {
		if ticket.Destination == destination {
			total++
		}
	}

	return total, nil
}

// ejemplo 2
func GetMornings(periodo string) (int, error) {
	passagem, err := listTicket()
	if err != nil {
		fmt.Println("Erro ao preencher tickets:", err)
	}

	var total = 0
	for _, ticket := range passagem {
		hour := ticket.Time.Hour()

		switch periodo {
		case "madrugada":
			if hour >= 0 && hour <= 6 {
				total++
			}
		case "manha":
			if hour >= 7 && hour <= 12 {
				total++
			}
		case "tarde":
			if hour >= 13 && hour <= 19 {
				total++
			}
		case "noite":
			if hour >= 20 && hour <= 23 {
				total++
			}
		}
	}

	return total, nil
}

// ejemplo 3
func AverageDestination(destination string) (int, error) {
	passagem, err := listTicket()
	if err != nil {
		fmt.Println("Erro ao preencher tickets:", err)
		return 0, nil
	}

	totalTickets := len(passagem)
	if totalTickets == 0 {
		return 0, nil
	}

	total := 0
	for _, ticket := range passagem {
		if ticket.Destination == destination {
			total++
		}
	}

	value := (total * 100) / totalTickets

	return value, nil
}

func listTicket() ([]Ticket, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	data, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return nil, err
	}

	var passagem []Ticket

	for _, dados := range data {
		if len(dados) < 6 {
			continue
		}

		price, err := strconv.ParseFloat(dados[5], 64)
		if err != nil {
			return nil, err
		}

		ticketTime, err := ParseHHMM(dados[4])
		if err != nil {
			return nil, err
		}

		ticket := Ticket{
			ID:          dados[0],
			Name:        dados[1],
			Email:       dados[2],
			Destination: dados[3],
			Time:        ticketTime,
			Price:       price,
		}

		passagem = append(passagem, ticket)
	}

	return passagem, nil
}

func ParseHHMM(timeStr string) (time.Time, error) {
	timestamp, err := time.Parse("15:04", timeStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("erro ao converter tempo: %w", err)
	}
	return timestamp, nil
}
