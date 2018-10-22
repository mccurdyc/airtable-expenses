package airtable

import (
	"encoding/json"
)

type Merchants struct {
	Records []Merchant `json:"records"`
}

type Merchant struct {
	ID     string         `json:"id,omitempty"`
	Fields MerchantFields `json:"fields"`
}

type MerchantFields struct {
	Name        string   `json:"Name"`
	Notes       string   `json:"Notes,omitempty"`
	Attachments []string `json:"Attachments,omitempty"`
	Purchases   []string `json:"Purchases,omitempty"`

	NumPurchases  int     `json:"NumPurchases,omitempty"`
	SumPurchases  float64 `json:"SumPurchases,omitempty"`
	AverageAmount float64 `json:"AverageAmount,omitempty"`
	MinAmount     float64 `json:"MinAmount,omitempty"`
	MaxAmount     float64 `json:"MaxAmount,omitempty"`
}

type merchantMap map[string]string

func (c *Client) CreateUniqueMerchant(n string) error {
	c.URL.Path = "/Merchants"

	nm := Merchant{
		Fields: MerchantFields{
			Name: n,
		},
	}

	b, _ := json.Marshal(nm)

	var ms Merchants
	return c.createUnique(n, &ms, b)
}
