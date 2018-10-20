package airtable

import (
	"context"
	"net/http"
)

type Config struct {
	APIKey string `json:"-"`
}

func NewClient(ctx context.Context, cfg Config) (*http.Client, error) {
	return &http.Client{}, nil
}
