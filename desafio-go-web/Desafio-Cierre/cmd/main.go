package main

import (
	"app/internal/application"
	"fmt"

	"github.com/joho/godotenv"
)

func main() {
	// env
	godotenv.Load()
	// application
	// - config
	cfg := &application.ConfigAppDefault{
		ServerAddr: ":8080",
		DbFile:     "docs/db/tickets.csv",
	}
	app := application.NewApplicationDefault(cfg)

	// - setup
	err := app.SetUp()
	if err != nil {
		fmt.Println(err)
		return
	}

	// - run
	err = app.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
}
