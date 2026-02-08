package database

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func InitDB(connectionString string) (*sql.DB, error) {
	// Open database
	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		return nil, err
	}

	// Test connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// Set connection pool settings (optional tapi recommended)
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	log.Println("Database connected successfully")
	return db, nil
}

func RunMigrations(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS transactions (
			id SERIAL PRIMARY KEY,
			date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			total_amount INTEGER NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS transaction_details (
			id SERIAL PRIMARY KEY,
			transaction_id INTEGER REFERENCES transactions(id),
			product_id INTEGER,
			quantity INTEGER NOT NULL,
			subtotal INTEGER NOT NULL
		);`,
	}

	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			return err
		}
	}

	log.Println("Database migrations executed successfully")
	return nil
}