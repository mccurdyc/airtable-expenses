package airtable

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Config struct {
	Host   string
	APIKey string
}

type Client struct {
	URL    *url.URL
	client *http.Client
	header http.Header
}

func NewClient(ctx context.Context, cfg Config) *Client {
	return &Client{
		URL: &url.URL{
			Scheme: "https",
			Host:   cfg.Host,
		},
		client: http.DefaultClient,
		header: http.Header{
			"Content-Type":  []string{"application/json"},
			"Authorization": []string{fmt.Sprintf("Bearer %s", cfg.APIKey)},
		},
	}
}

func (c *Client) createUnique(n string, t interface{}, b []byte) error {
	if len(c.URL.Path) == 0 {
		return errors.New("URL path must be set")
	}

	c.URL.RawQuery = "?fields[]=Name&maxRecords=1000"

	req := &http.Request{
		Method: "GET",
		URL:    c.URL,
		Header: c.header,
	}

	resp, _ := c.client.Do(req)

	// TODO: (@mccurdyc) handle pagination of results
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = json.Unmarshal(body, t)
	if err != nil {
		return err
	}

	switch e := t.(type) {
	case Merchants:
		t = t.(Merchants)
	case Tags:
		t = t.(Tags)
	default:
		return fmt.Errorf("type %v not found", e)
	}

	m := make(map[string]string)
	for _, j := range t.Records {
		m[j.Fields.Name] = j.ID
	}

	_, ok := m[n]
	if !ok {
		c.URL.RawQuery = ""

		if len(b) == 0 {
			return errors.New("new object cannot be empty")
		}

		buf := bytes.NewBuffer(b)

		req = &http.Request{
			Method: "POST",
			URL:    c.URL,
			Header: c.header,
			Body:   ioutil.NopCloser(buf),
		}

		c.client.Do(req)
	}

	return nil
}
