package main

import (
	"fmt"

	"github.com/bootcamp-go/desafio-go-bases/internal/tickets"
)

func main() {
	total, err := tickets.GetTotalTickets("Brazil")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("=== Exercício 1 ===")
	fmt.Println(total)

	fmt.Println("=== Exercício 2 ===")
	totalTime, err := tickets.GetMornings("tarde")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(totalTime)

	fmt.Println("=== Exercício 3 ===")
	totalDes, err := tickets.AverageDestination("Brazil")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(totalDes)
}
