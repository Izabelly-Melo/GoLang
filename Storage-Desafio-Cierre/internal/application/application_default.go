package application

import (
	"app/internal/handler"
	"app/internal/repository"
	"app/internal/service"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-sql-driver/mysql"
)

// ConfigApplicationDefault is the configuration for NewApplicationDefault.
type ConfigApplicationDefault struct {
	// Db is the database configuration.
	Db *mysql.Config
	// Addr is the server address.
	Addr string
}

// NewApplicationDefault creates a new ApplicationDefault.
func NewApplicationDefault(config *ConfigApplicationDefault) *ApplicationDefault {
	// default values
	defaultCfg := &ConfigApplicationDefault{
		Db:   nil,
		Addr: ":8080",
	}
	if config != nil {
		if config.Db != nil {
			defaultCfg.Db = config.Db
		}
		if config.Addr != "" {
			defaultCfg.Addr = config.Addr
		}
	}

	return &ApplicationDefault{
		cfgDb:   defaultCfg.Db,
		cfgAddr: defaultCfg.Addr,
	}
}

// ApplicationDefault is an implementation of the Application interface.
type ApplicationDefault struct {
	// cfgDb is the database configuration.
	cfgDb *mysql.Config
	// cfgAddr is the server address.
	cfgAddr string
	// db is the database connection.
	db *sql.DB
	// router is the chi router.
	router *chi.Mux
}

// SetUp sets up the application.
func (a *ApplicationDefault) SetUp() (err error) {
	// dependencies
	// - db: init
	a.db, err = sql.Open("mysql", a.cfgDb.FormatDSN())
	if err != nil {
		return
	}
	// - db: ping
	err = a.db.Ping()
	if err != nil {
		return
	}

	customers, err := a.CheckTableContainsData("fantasy_products.customers")
	if err != nil {
		return fmt.Errorf("erro ao verificar customers: %w", err)
	}
	if !customers {
		err = a.ImportCustomersFromJSON("./docs/db/json/customers.json")
		if err != nil {
			return fmt.Errorf("erro ao importar customers: %w", err)
		}
	}

	// Verifique se as invoices já existem
	invoices, err := a.CheckTableContainsData("fantasy_products.invoices")
	if err != nil {
		return fmt.Errorf("erro ao verificar invoices: %w", err)
	}
	if !invoices {
		err = a.ImportInvoicesFromJSON("./docs/db/json/invoices.json")
		if err != nil {
			return fmt.Errorf("erro ao importar invoices: %w", err)
		}
	}

	// Verifique se os products já existem
	products, err := a.CheckTableContainsData("fantasy_products.products")
	if err != nil {
		return fmt.Errorf("erro ao verificar products: %w", err)
	}
	if !products {
		err = a.ImportProductsFromJSON("./docs/db/json/products.json")
		if err != nil {
			return fmt.Errorf("erro ao importar products: %w", err)
		}
	}

	// Verifique se as sales já existem
	sales, err := a.CheckTableContainsData("fantasy_products.sales")
	if err != nil {
		return fmt.Errorf("erro ao verificar sales: %w", err)
	}
	if !sales {
		err = a.ImportSalesFromJSON("./docs/db/json/sales.json")
		if err != nil {
			return fmt.Errorf("erro ao importar sales: %w", err)
		}
	}

	// - repository
	rpCustomer := repository.NewCustomersMySQL(a.db)
	rpProduct := repository.NewProductsMySQL(a.db)
	rpInvoice := repository.NewInvoicesMySQL(a.db)
	rpSale := repository.NewSalesMySQL(a.db)
	// - service
	svCustomer := service.NewCustomersDefault(rpCustomer)
	svProduct := service.NewProductsDefault(rpProduct)
	svInvoice := service.NewInvoicesDefault(rpInvoice)
	svSale := service.NewSalesDefault(rpSale)
	// - handler
	hdCustomer := handler.NewCustomersDefault(svCustomer)
	hdProduct := handler.NewProductsDefault(svProduct)
	hdInvoice := handler.NewInvoicesDefault(svInvoice)
	hdSale := handler.NewSalesDefault(svSale)

	// routes
	// - router
	a.router = chi.NewRouter()
	// - middlewares
	a.router.Use(middleware.Logger)
	a.router.Use(middleware.Recoverer)
	// - endpoints
	a.router.Route("/customers", func(r chi.Router) {
		// - GET /customers
		r.Get("/", hdCustomer.GetAll())
		r.Get("/conditions", hdCustomer.GetConditionsCustomer())
		r.Get("/actives", hdCustomer.GetCustomersMoreActives())
		// - POST /customers
		r.Post("/", hdCustomer.Create())
	})
	a.router.Route("/products", func(r chi.Router) {
		// - GET /products
		r.Get("/", hdProduct.GetAll())
		r.Get("/sold", hdProduct.GetProductsMoreSold())
		// - POST /products
		r.Post("/", hdProduct.Create())
	})
	a.router.Route("/invoices", func(r chi.Router) {
		// - GET /invoices
		r.Get("/", hdInvoice.GetAll())
		// - POST /invoices
		r.Post("/", hdInvoice.Create())
	})
	a.router.Route("/sales", func(r chi.Router) {
		// - GET /sales
		r.Get("/", hdSale.GetAll())
		// - POST /sales
		r.Post("/", hdSale.Create())
	})

	return
}

// Run runs the application.
func (a *ApplicationDefault) Run() (err error) {
	defer a.db.Close()

	err = http.ListenAndServe(a.cfgAddr, a.router)
	return
}

func (a *ApplicationDefault) ImportCustomersFromJSON(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("erro ao abrir o arquivo: %w", err)
	}
	defer file.Close()

	var customers []handler.CustomerJSON
	if err := json.NewDecoder(file).Decode(&customers); err != nil {
		return fmt.Errorf("erro ao decodificar JSON: %w", err)
	}

	for _, customer := range customers {
		_, err := a.db.Exec("INSERT INTO fantasy_products.customers (`id`, `first_name`, `last_name`, `condition`) VALUES (?, ?, ?, ?)", customer.Id, customer.FirstName, customer.LastName, customer.Condition)
		if err != nil {
			return fmt.Errorf("erro ao inserir customer %s: %w", customer.FirstName, err)
		}
	}

	return nil
}

func (a *ApplicationDefault) ImportInvoicesFromJSON(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("erro ao abrir o arquivo: %w", err)
	}
	defer file.Close()

	var invoices []handler.InvoiceJSON
	if err := json.NewDecoder(file).Decode(&invoices); err != nil {
		return fmt.Errorf("erro ao decodificar JSON: %w", err)
	}

	for _, invoice := range invoices {
		_, err := a.db.Exec("INSERT INTO fantasy_products.invoices (`id`, `datetime`, `customer_id`, `total`) VALUES (?, ?, ?, ?)", invoice.Id, invoice.Datetime, invoice.CustomerId, invoice.Total)
		if err != nil {
			return fmt.Errorf("erro ao inserir invoice com o ID %d: %w", invoice.Id, err)
		}
	}

	return nil
}

func (a *ApplicationDefault) ImportProductsFromJSON(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("erro ao abrir o arquivo: %w", err)
	}
	defer file.Close()

	var products []handler.ProductJSON
	if err := json.NewDecoder(file).Decode(&products); err != nil {
		return fmt.Errorf("erro ao decodificar JSON: %w", err)
	}

	for _, product := range products {
		_, err := a.db.Exec("INSERT INTO fantasy_products.products (`id`, `description`, `price`) VALUES (?, ?, ?)", product.Id, product.Description, product.Price)
		if err != nil {
			return fmt.Errorf("erro ao inserir product com o ID %d: %w", product.Id, err)
		}
	}

	return nil
}

func (a *ApplicationDefault) ImportSalesFromJSON(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("erro ao abrir o arquivo: %w", err)
	}
	defer file.Close()

	var sales []handler.SaleJSON
	if err := json.NewDecoder(file).Decode(&sales); err != nil {
		return fmt.Errorf("erro ao decodificar JSON: %w", err)
	}

	for _, sale := range sales {
		_, err := a.db.Exec("INSERT INTO fantasy_products.sales (`id`,`quantity`,`invoice_id`,`product_id`) VALUES (?, ?, ?, ?)", sale.Id, sale.Quantity, sale.InvoiceId, sale.ProductId)
		if err != nil {
			return fmt.Errorf("erro ao inserir sale com o ID %d: %w", sale.Id, err)
		}
	}

	return nil
}

func (a *ApplicationDefault) CheckTableContainsData(tableName string) (bool, error) {
	var count int
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)
	err := a.db.QueryRow(query).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("erro ao verificar a tabela %s: %w", tableName, err)
	}
	return count > 0, nil
}
