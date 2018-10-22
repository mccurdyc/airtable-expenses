package airtable

import (
	"encoding/json"
)

type Tags struct {
	Records []Tag `json:"records"`
}

type Tag struct {
	ID     string    `json:"id,omitempty"`
	Fields TagFields `json:"fields"`
}

type TagFields struct {
	Name      string   `json:"Name"`
	Purchases []string `json:"Purchases,omitempty"`

	NumPurchases  int     `json:"NumPurchases,omitempty"`
	SumPurchases  float64 `json:"SumPurchases,omitempty"`
	AverageAmount float64 `json:"AverageAmount,omitempty"`
	MinAmount     float64 `json:"MinAmount,omitempty"`
	MaxAmount     float64 `json:"MaxAmount,omitempty"`
}

func (c *Client) CreateUniqueTag(n string) error {
	c.URL.Path = "/Tags"

	nt := Tag{
		Fields: TagFields{
			Name: n,
		},
	}

	b, _ := json.Marshal(nt)

	var t Tags
	return c.createUnique(n, &t, b)
}
