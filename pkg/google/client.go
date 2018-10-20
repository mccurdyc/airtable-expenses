package google

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

func NewClient(ctx context.Context, cfg *oauth2.Config) (*http.Client, error) {
	tkn, err := authenticate(cfg)
	if err != nil {
		return &http.Client{}, err
	}

	return cfg.Client(ctx, tkn), nil

}

func authenticate(cfg *oauth2.Config) (*oauth2.Token, error) {
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
