// github.com/canonical-ledgers/bitcointax
// Copyright 2018 Canonical Ledgers, LLC. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file distributed with this source code.

package bitcointax

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	APIURL          = "https://api.bitcoin.tax/v1"
	TransactionsURI = "/transactions"
)

type Option func(*Client)

func WithHTTPClient(httpC *http.Client) Option {
	return func(c *Client) { c.client = httpC }
}

func WithURL(url string) Option {
	return func(c *Client) { c.url = url }
}

// Client stores API credentials and provides methods for querying the
// bitcoin.tax REST API. You may set any additional http.Client settings
// directly on this type.
type Client struct {
	url    string
	key    string
	secret string
	client *http.Client
}

// NewClient returns a pointer to a new Client with the given key and secret
// and with the default API endpoint (https://api.bitcoin.tax/v1) as the URL.
func NewClient(key, secret string, opts ...Option) *Client {
	c := Client{
		url:    APIURL,
		key:    key,
		secret: secret,
		client: http.DefaultClient,
	}

	for _, opt := range opts {
		opt(&c)
	}

	return &c
}

// ListTransactions makes a GET request to the /transactions endpoint and
// returns a slice of transactions for the given taxyear.Year() from the start
// transaction index and returning no more than limit transactions. If limit is
// 0 than the limit will not be specified and the API's default limit of 100
// will be used.
func (c *Client) ListTransactions(ctx context.Context, taxyear, start, limit uint) ([]Transaction, uint, error) {
	values := make(url.Values, 3)
	values.Add("taxyear", fmt.Sprintf("%d", taxyear))
	values.Add("start", fmt.Sprintf("%v", start))
	if limit != 0 {
		values.Add("limit", fmt.Sprintf("%v", limit))
	}

	response, err := c.doRequest(ctx, values, nil)
	if err != nil {
		return nil, 0, err
	}
	return response.Data.Transactions, response.Data.Total, nil
}

// AddTransactions makes a POST request to the /transactions endpoint adding
// all Transactions in txs.
func (c *Client) AddTransactions(ctx context.Context, txs ...Transaction) error {
	if len(txs) == 0 {
		return fmt.Errorf("`txs []Transaction` is empty or nil")
	}
	data, err := json.Marshal(txs)
	if err != nil {
		return fmt.Errorf("json.Marshal(%+v): %v", txs, err)
	}

	if _, err := c.doRequest(ctx, nil, bytes.NewBuffer(data)); err != nil {
		return err
	}
	return nil
}

type response struct {
	Status string `json:"status"` // "success", "fail" or "error"
	Data   struct {
		Total        uint          `json:"total"` // total number of transactions
		Transactions []Transaction `json:"transactions"`
		Message      string        `json:"message,omitempty"`
	} `json:"data"`
}

func (c *Client) doRequest(ctx context.Context, params url.Values, body io.Reader) (*response, error) {
	method := "POST"
	if body == nil {
		method = "GET"
	}
	url := c.url + TransactionsURI
	if params != nil {
		url += "?" + params.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-APIKEY", c.key)
	req.Header.Set("X-APISECRET", c.secret)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(resp.Status)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response response
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, err
	}
	const statusSuccess = "success"
	if response.Status != statusSuccess {
		return nil, fmt.Errorf("status: %q, message: %q",
			response.Status, response.Data.Message)
	}
	return &response, nil
}
