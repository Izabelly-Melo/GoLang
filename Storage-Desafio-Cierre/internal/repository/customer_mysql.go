package repository

import (
	"database/sql"

	"app/internal"
)

// NewCustomersMySQL creates new mysql repository for customer entity.
func NewCustomersMySQL(db *sql.DB) *CustomersMySQL {
	return &CustomersMySQL{db}
}

// CustomersMySQL is the MySQL repository implementation for customer entity.
type CustomersMySQL struct {
	// db is the database connection.
	db *sql.DB
}

// FindAll returns all customers from the database.
func (r *CustomersMySQL) FindAll() (c []internal.Customer, err error) {
	// execute the query
	rows, err := r.db.Query("SELECT `id`, `first_name`, `last_name`, `condition` FROM customers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// iterate over the rows
	for rows.Next() {
		var cs internal.Customer
		// scan the row into the customer
		err := rows.Scan(&cs.Id, &cs.FirstName, &cs.LastName, &cs.Condition)
		if err != nil {
			return nil, err
		}
		// append the customer to the slice
		c = append(c, cs)
	}
	err = rows.Err()
	if err != nil {
		return
	}

	return
}

// Save saves the customer into the database.
func (r *CustomersMySQL) Save(c *internal.Customer) (err error) {
	// execute the query
	res, err := r.db.Exec(
		"INSERT INTO customers (`first_name`, `last_name`, `condition`) VALUES (?, ?, ?)",
		(*c).FirstName, (*c).LastName, (*c).Condition,
	)
	if err != nil {
		return err
	}

	// get the last inserted id
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	// set the id
	(*c).Id = int(id)

	return
}

func (r *CustomersMySQL) GetGroupByConditions() (res []internal.ResCustomer, err error) {
	query := `SELECT CASE c.condition WHEN 0 THEN 'Inativo' WHEN 1 THEN 'Ativo' END AS Conditions,` +
		` ROUND(SUM(i.total), 2) AS Total ` +
		` FROM customers c JOIN  invoices i ON c.id = i.customer_id` +
		` GROUP BY c.condition`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// iterate over the rows
	for rows.Next() {
		var body internal.ResCustomer
		// scan the row into the customer
		err := rows.Scan(&body.Condition, &body.Total)
		if err != nil {
			return nil, err
		}
		// append the customer to the slice
		res = append(res, body)
	}
	err = rows.Err()
	if err != nil {
		return
	}

	return
}

func (r *CustomersMySQL) GetAmountCostumers() (res []internal.ResAmountCustomer, err error) {
	query := `SELECT c.first_name AS FirstName, c.last_name AS LastName, ROUND(SUM(i.total), 2) AS Amount FROM customers c` +
		` JOIN invoices i ON c.id = i.customer_id WHERE c.condition = 1 ` +
		` GROUP BY c.id, c.first_name, c.last_name ORDER BY Amount DESC LIMIT 5`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// iterate over the rows
	for rows.Next() {
		var body internal.ResAmountCustomer
		// scan the row into the customer
		err := rows.Scan(&body.FirstName, &body.LastName, &body.Amount)
		if err != nil {
			return nil, err
		}
		// append the customer to the slice
		res = append(res, body)
	}
	err = rows.Err()
	if err != nil {
		return
	}

	return
}
