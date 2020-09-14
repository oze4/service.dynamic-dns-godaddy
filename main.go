package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/oze4/godaddygo"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Unable to load environmental variables")
		os.Exit(1)
	}

	fromdb, err := getFromDB()
	if err != nil {
		log.Fatalf("Error getting public IP from postgres: %s\n", err.Error())
		os.Exit(1)
	}

	fromapi, err := getFromAPI()
	if err != nil {
		log.Fatalf("Error getting public IP from https://icanhazip.com %s\n", err.Error())
	 	os.Exit(1)
    }

	if fromdb != fromapi {
        if err := updatePublicIP(fromapi); err != nil {
            fmt.Println(err.Error())
            os.Exit(1)
        }

		api := godaddygo.ConnectProduction(
			os.Getenv("GODADDY_APIKEY"),
			os.Getenv("GODADDY_APISECRET"),
		)

		zone := api.V1().Domain(os.Getenv("GODADDY_DOMAIN")).Records()

		records, err := zone.GetByType("A")
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		for _, rec := range *records {
			fmt.Printf("Updating record: %s %s %s\n", rec.Type, rec.Name, fromapi)
			if err := zone.SetValue(rec.Type, rec.Name, fromapi); err != nil {
				fmt.Println(err.Error())
				fmt.Println("Despite this error we will continue.")
			}
		}
	} else {
		fmt.Printf("Public IP has not changed! %s\n", fromdb)
	}

	os.Exit(0)
}

func getFromAPI() (string, error) {
	resp, err := http.Get("https://icanhazip.com")
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(body)), nil
}

func getFromDB() (string, error) {
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

	return strings.TrimSpace(public), nil
}

func updatePublicIP(ip string) error {
    // Yes, I know there is duplicate code, this is a tiny "microservice"
    // so we should be fine...
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
		return err
	}

	defer db.Close()

	// Force connection to PG
	err = db.Ping()
	if err != nil {
		return err
    }

	if _, err := db.Exec(`UPDATE ip_addresses SET public = $2 WHERE id = $1;`, 1, ip); err != nil {
		return err
	}

	return nil
}
