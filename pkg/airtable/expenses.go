package airtable

import (
	"fmt"
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
