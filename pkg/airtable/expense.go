package airtable

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Expense struct {
	Fields Fields `json:"fields"`
}

type Fields struct {
	ID          int       `json:"ID,omitempty"`
	Date        JSONTime  `json:"Date"`
	Amount      float64   `json:"Amount"`
	Merchant    []string  `json:"Merchant"`
	Tag         []string  `json:"Tag"`
	Purchaser   Purchaser `json:"Purchaser"`
	Notes       string    `json:"Notes,omitempty"`
	Attachments []string  `json:"Attachments,omitempty"`
}

type JSONTime time.Time

func (t JSONTime) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format("01/02/2006"))
	return []byte(stamp), nil
}

func (c *Client) CreateExpense(exp Expense) {
	c.URL.Path = "/Purchases"

	b, _ := json.Marshal(exp)
	buf := bytes.NewBuffer(b)

	req := &http.Request{
		Method: "POST",
		URL:    c.URL,
		Header: c.header,
		Body:   ioutil.NopCloser(buf),
	}

	c.client.Do(req)
}
