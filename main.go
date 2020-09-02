package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Unable to load environmental variables")
		os.Exit(1)
	}
    
    publicIPFromDB, err := getPublicIPAddressFromDB()
    if err != nil {
        panic(err.Error())
    }
    
    fmt.Printf("Successfully connected to database!\r\nCurrent public IP from database: %s\n", publicIPFromDB)
}

func getPublicIPAddressFromDB() (string, error) {
	var (
		host     = os.Getenv("PG_HOST")
		port     = os.Getenv("PG_PORT")
		user     = os.Getenv("PG_USER")
		password = os.Getenv("PG_PASSWORD")
		dbname   = os.Getenv("PG_DATABASE")
	)

	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s "+"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return "", err
    }
    
	defer db.Close()

    // Force connection to PG
	err = db.Ping()
	if err != nil {
		return "", err
    }

    var public string
    err = db.QueryRow(`SELECT public FROM ip_addresses WHERE id = $1;`, 1).Scan(&public)
    if err != nil {
      return "", err
    }

    return public, nil
}
