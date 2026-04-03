package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func initDB() {
	dsn := "root:@tcp(localhost:3306)/"
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error opening connection: %v", err)
	}

	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS occ")
	if err != nil {
		log.Fatalf("Error creating database: %v", err)
	}
	db.Close()

	dsn = "root:@tcp(localhost:3306)/occ"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error connecting to 'occ': %v", err)
	}

	query := `
	CREATE TABLE IF NOT EXISTS test (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255),
		a TEXT,
		b TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`
	_, err = db.Exec(query)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}

	fmt.Println("Database and table initialized successfully!")
}
