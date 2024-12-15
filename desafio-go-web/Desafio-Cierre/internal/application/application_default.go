package application

import (
	"app/internal/handler"
	"app/internal/loader"
	"app/internal/repository"
	"app/internal/service"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// ConfigServerChi is a struct that represents the configuration for ServerChi
type ConfigAppDefault struct {
	// serverAddr represents the address of the server
	ServerAddr string
	// dbFile represents the path to the database file
	DbFile string
}

// NewApplicationDefault creates a new default application
func NewApplicationDefault(cfg *ConfigAppDefault) *ApplicationDefault {
	// default values
	defaultRouter := chi.NewRouter()
	defaultConfig := &ConfigAppDefault{
		ServerAddr: ":8080",
	}
	if cfg != nil {
		if cfg.ServerAddr != "" {
			defaultConfig.ServerAddr = cfg.ServerAddr
		}
		if cfg.DbFile != "" {
			defaultConfig.DbFile = cfg.DbFile
		}
	}

	return &ApplicationDefault{
		rt:         defaultRouter,
		serverAddr: defaultConfig.ServerAddr,
		dbFile:     defaultConfig.DbFile,
	}
}

// ApplicationDefault represents the default application
type ApplicationDefault struct {
	// router represents the router of the application
	rt *chi.Mux
	// serverAddr represents the address of the server
	serverAddr string
	// dbFile represents the path to the database file
	dbFile string
}

// Run is a method that runs the application
func (a *ApplicationDefault) SetUp() (err error) {
	// dependencies
	db := loader.NewLoaderTicketCSV(a.dbFile)
	tickets, err := db.Load()
	if err != nil {
		log.Println("failed to load")
		return
	}
	rp := repository.NewRepositoryTicketMap(0, tickets)
	// service ...
	service := service.NewServiceTicketDefault(rp)
	// handler ...
	handler := handler.NewHandlerTicketDefault(service)

	// routes
	(*a).rt.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("OK"))
	})

	(*a).rt.Route("/ticket", func(rt chi.Router) {
		// - GET /ticket
		rt.Get("/", handler.GetTotalAmountTickets)
		rt.Get("/getByCountry/{dest}", handler.GetTicketsAmountByDestinationCountry)
		rt.Get("/getAverage/{dest}", handler.GetAverageCountry)
	})
	return
}

// Run runs the application
func (a *ApplicationDefault) Run() (err error) {
	err = http.ListenAndServe(a.serverAddr, a.rt)
	return
}
