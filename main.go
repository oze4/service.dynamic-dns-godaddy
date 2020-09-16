package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/oze4/godaddygo"
	"github.com/oze4/godaddygo/pkg/endpoints"
)

func main() {
	godotenv.Load()

	api := newIcanhazip()

	apires, err := api.get()
	if err != nil {
		explainStrExit("Error getting public IP from https://icanhazip.com "+err.Error(), 1)
	}

	k := os.Getenv("GODADDY_APIKEY")
	s := os.Getenv("GODADDY_APISECRET")
	d := os.Getenv("GODADDY_DOMAIN")

	gd := newGoDaddy(k, s, d)

	gdres, err := gd.get()
	if err != nil {
		explainStrExit("Error getting IP from GoDaddy: "+err.Error(), 1)
	}

	if apires == gdres {
		explainStrExit("Public IP has not changed: "+apires, 0)
	}

	explainStr("Public IP has changed, updating GoDaddy now. Old: " + gdres + " New: " + apires)

	if err := gd.update(apires); err != nil {
		explainStrExit("Error updating GoDaddy DNS: "+err.Error(), 1)
	}

	explainStrExit("\nDone\n", 0)
}

/**
 * helper functions
 */

func explainStr(message string) {
	fmt.Printf("%s\n", message)
}

func explainStrExit(message string, exitCode int) {
	fmt.Printf("%s\n", message)
	os.Exit(exitCode)
}

/**
 * icanhazip
 */

func newIcanhazip() icanhazip {
	return icanhazip{}
}

type icanhazip struct{}

func (i icanhazip) get() (string, error) {
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

/**
 * godaddy
 */

func newGoDaddy(k, s, d string) goDaddy {
	baseline := os.Getenv("BASELINE_RECORD")
	api := godaddygo.ConnectProduction(k, s).V1().Domain(d).Records()
	return goDaddy{baseline, api}
}

type goDaddy struct {
	baseline string // BASELINE_RECORD
	recs     endpoints.Records
}

// get returns the *value* of your `BASELINE_RECORD`
func (g *goDaddy) get() (string, error) {
	r, e := g.recs.GetByTypeName("A", g.baseline)
	if e != nil {
		return "", e
	}
	return strings.TrimSpace((*r)[0].Data), nil
}

// update sets all A records to your new public IP
func (g goDaddy) update(newIP string) error {
	r, e := g.recs.GetByType("A")
	if e != nil {
		return e
	}
	for _, d := range *r {
		if e := g.recs.SetValue(d.Type, d.Name, newIP); e != nil {
			return errors.New("Error setting record: " + d.Name)
		}
	}
	return nil
}
