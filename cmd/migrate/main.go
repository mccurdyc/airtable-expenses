package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

	spreadsheetID := os.Getenv("SHEET_ID")

	readRange := "Sheet1!A2:E"

	resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	merchantMap := make(map[string]string)
	tagMap := make(map[string]string)

	if len(resp.Values) == 0 {
		fmt.Println("No data found.")
	} else {

		for _, row := range resp.Values {
			fmt.Printf("%+v\n", row)

			gstrMerchant := strings.ToLower(row[3].(string))
			_, ok := merchantMap[gstrMerchant]
			if !ok {
				am := airtable.Merchant{
					Fields: airtable.MerchantFields{
						Name:  gstrMerchant,
						Notes: "migration from Google Sheets",
					},
				}

				b, _ := json.Marshal(am)
				br := bytes.NewReader(b)

				req, _ := http.NewRequest("POST", "https://api.airtable.com/v0/appiuHki0G4ddogkp/Merchants", br)

				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("AIRTABLE_API_KEY")))

				client := http.DefaultClient
				client.Do(req)
			}

			gstrTag := strings.ToLower(row[2].(string))
			_, ok = tagMap[gstrTag]
			if !ok {
				at := airtable.Tag{
					Fields: airtable.TagFields{
						Name: gstrTag,
					},
				}

				b, _ := json.Marshal(at)
				br := bytes.NewReader(b)

				req, _ := http.NewRequest("POST", "https://api.airtable.com/v0/appiuHki0G4ddogkp/Tags", br)

				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("AIRTABLE_API_KEY")))

				client := http.DefaultClient
				client.Do(req)
			}

			req, _ := http.NewRequest("GET", "https://api.airtable.com/v0/appiuHki0G4ddogkp/Merchants?fields[]=Name&maxRecords=1000", nil)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("AIRTABLE_API_KEY")))

			client := http.DefaultClient
			resp, _ := client.Do(req)

			body, _ := ioutil.ReadAll(resp.Body)
			resp.Body.Close()

			var merchants airtable.Merchants
			err = json.Unmarshal(body, &merchants)

			for _, m := range merchants.Records {
				merchantMap[m.Fields.Name] = m.ID
			}
			fmt.Printf("MERCHANT MAP: %+v\n", merchantMap)

			req, _ = http.NewRequest("GET", "https://api.airtable.com/v0/appiuHki0G4ddogkp/Tags?fields[]=Name&maxRecords=1000", nil)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("AIRTABLE_API_KEY")))

			client = http.DefaultClient
			resp, _ = client.Do(req)

			body, _ = ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			var tags airtable.Tags
			err = json.Unmarshal(body, &tags)

			for _, t := range tags.Records {
				tagMap[t.Fields.Name] = t.ID
			}
			fmt.Printf("TAG MAP: %+v\n", tagMap)

			date, err := time.Parse("1/2/06", row[0].(string))
			amount, _ := strconv.ParseFloat(row[1].(string), 64)

			exp := airtable.Expense{
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
			}

			b, err := json.Marshal(exp)
			fmt.Println(string(b))
			br := bytes.NewReader(b)

			req, err = http.NewRequest("POST", "https://api.airtable.com/v0/appiuHki0G4ddogkp/Purchases", br)
			if err != nil {
				fmt.Println(err)
			}

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("AIRTABLE_API_KEY")))

			client = http.DefaultClient
			client.Do(req)

			time.Sleep(1 * time.Second)
		}
	}
}
