package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/oauth2"
	sheets "google.golang.org/api/sheets/v4"

	"github.com/mccurdyc/airtable-expenses/pkg/airtable"
	"github.com/mccurdyc/airtable-expenses/pkg/expenses"
	"github.com/mccurdyc/airtable-expenses/pkg/google"
)

func main() {
	ctx := expenses.Ctx{
		Context: context.Background(),
	}

	googleClient, err := google.NewClient(ctx.Context, &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/spreadsheets.readonly"},
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URI"),
		Endpoint: oauth2.Endpoint{
			AuthURL:  os.Getenv("GOOGLE_AUTH_URI"),
			TokenURL: os.Getenv("GOOGLE_TOKEN_URI"),
		},
	})

	srv, err := sheets.New(googleClient)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	resp, err := srv.Spreadsheets.Values.Get(os.Getenv("SHEET_ID"), os.Getenv("SHEET_READ_RANGE")).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	if len(resp.Values) == 0 {
		fmt.Println("No data found.")
		os.Exit(0)
	}

	airtableClient := airtable.NewClient(context.TODO,
		Config{
			APIKey: os.Getenv("AIRTABLE_API_KEY"),
			Host:   os.Getenv("AIRTABLE_HOST"),
		})

	for _, row := range resp.Values {
		fmt.Printf("%+v\n", row)

		for i, e := range row {
			row[i] = strings.Trim(strings.ToLower(row[i].(string)), " ")
		}

		merchant := airtableClient.CreateUniqueMerchant(gstrMerchant)
		tag := airtableClient.CreateUniqueTag(gstrTag)

		date, err := time.Parse("1/2/06", row[0].(string))
		amount, _ := strconv.ParseFloat(row[1].(string), 64)

		airtableClient.CreateExpense(airtable.Expense{
			Fields: airtable.Fields{
				Date:   airtable.JSONTime(date),
				Amount: amount,
				Merchant: []string{
					merchantMap[gstrMerchant],
				},
				Tag: []string{
					tagMap[gstrTag],
				},
				Purchaser: airtable.Purchaser{
					ID:    "usr5jWb7gSI9FRp2D",
					Email: "mccurdyc22@gmail.com",
					Name:  "Colton McCurdy",
				},
				Notes: "migration from Google Sheets",
			},
		})

		time.Sleep(1 * time.Second)
	}
}
