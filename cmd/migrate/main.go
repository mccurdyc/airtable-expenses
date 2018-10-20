package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	sheets "google.golang.org/api/sheets/v4"
)

func main() {
	cfg := getConfig()
	tkn, err := getToken(cfg)
	if err != nil {
		fmt.Println(err)
	}

	client := getClient(cfg, tkn)

	srv, err := sheets.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	spreadsheetID := os.Getenv("SHEET_ID")

	readRange := "Sheet1!A2:E"

	resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	if len(resp.Values) == 0 {
		fmt.Println("No data found.")
	} else {

		for _, row := range resp.Values {
			// Print columns A and E, which correspond to indices 0 and 4.
			fmt.Printf("%s, %s\n", row[0], row[4])
		}
	}
}

func getConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/spreadsheets.readonly"},
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URI"),
		Endpoint: oauth2.Endpoint{
			AuthURL:  os.Getenv("GOOGLE_AUTH_URI"),
			TokenURL: os.Getenv("GOOGLE_TOKEN_URI"),
		},
	}
}

func getToken(cfg *oauth2.Config) (*oauth2.Token, error) {
	authURL := cfg.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	fmt.Printf("visit authURL: %s\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		return &oauth2.Token{}, errors.Wrap(err, "unable to read auth code")
	}

	tok, err := cfg.Exchange(context.TODO(), authCode)
	if err != nil {
		return &oauth2.Token{}, errors.Wrap(err, "unable to retreive token from the web")
	}

	return tok, nil
}

func getClient(cfg *oauth2.Config, tkn *oauth2.Token) *http.Client {
	return cfg.Client(context.Background(), tkn)
}
