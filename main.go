package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/oze4/godaddygo"
)

// GoDaddy holds godaddy api creds
type GoDaddy struct {
	key            string
	secret         string
	domain         string
	baselineRecord string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Unable to load environmental variables")
		os.Exit(1)
	}

	fromapi, err := getFromAPI()
	if err != nil {
		log.Fatalf("Error getting public IP from https://icanhazip.com %s\n", err.Error())
		os.Exit(1)
	}

	gd := GoDaddy{
		key:            os.Getenv("GODADDY_APIKEY"),
		secret:         os.Getenv("GODADDY_APISECRET"),
		domain:         os.Getenv("GODADDY_DOMAIN"),
		baselineRecord: os.Getenv("BASELINE_RECORD"),
	}

	fromgodaddy, err := getFromGoDaddy(gd)
	if err != nil {
		log.Fatalf("Error getting IP from GoDaddy %s\n", err.Error())
		os.Exit(1)
	}

	tfromapi := strings.TrimSpace(fromapi)
	tfromgodaddy := strings.TrimSpace(fromgodaddy)

	if tfromapi != tfromgodaddy {
		if err := updateGoDaddy(gd, tfromapi); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
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

func getFromGoDaddy(g GoDaddy) (string, error) {
	r, e := godaddygo.ConnectProduction(g.key, g.secret).V1().Domain(g.domain).Records().GetByTypeName("A", g.baselineRecord)
	if e != nil {
		return "", e
	}
	return (*r)[0].Data, nil
}

func updateGoDaddy(g GoDaddy, newIP string) error {
	recs := godaddygo.ConnectProduction(g.key, g.secret).V1().Domain(g.domain).Records()
	r, e := recs.GetByType("A")
	if e != nil {
		return e
	}
	for _, d := range *r {
		if e := recs.SetValue(d.Type, d.Name, newIP); e != nil {
			return errors.New("Error setting record: " + d.Name)
		}
	}
	return nil
}
