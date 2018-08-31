// github.com/canonical-ledgers/bitcointax
// Copyright 2018 Canonical Ledgers, LLC. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file distributed with this source code.

package bitcointax

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Client stores API credentials and provides methods for querying the
// bitcoin.tax REST API. You may set any additional http.Client settings
// directly on this type.
type Client struct {
	URL    string
	Key    string
	Secret string
	http.Client
}

const bitcoinTaxAPIURL = "https://api.bitcoin.tax/v1"
const transactionsURI = "/transactions"

// NewClient returns a pointer to a new Client with the given key and secret
// and with the default API endpoint (https://api.bitcoin.tax/v1) as the URL.
func NewClient(key, secret string) *Client {
	return &Client{URL: bitcoinTaxAPIURL, Key: key, Secret: secret}
}

// ListTransactions makes a GET request to the /transactions endpoint and
// returns a slice of transactions for the given taxyear.Year() from the start
// transaction index and returning no more than limit transactions. If limit is
// 0 than the limit will not be specified and the API's default limit of 100
// will be used.
func (c Client) ListTransactions(taxyear time.Time,
	start, limit uint64) ([]Transaction, uint64, error) {
	values := make(url.Values)
	values.Add("taxyear", taxyear.Format("2006"))
	values.Add("start", fmt.Sprintf("%v", start))
	if limit != 0 {
		values.Add("limit", fmt.Sprintf("%v", limit))
	}

	response, err := c.doRequest(c.URL+transactionsURI+"?"+values.Encode(), nil)
	if err != nil {
		return nil, 0, err
	}
	return response.Data.Transactions, response.Data.Total, nil
}

// AddTransactions makes a POST request to the /transactions endpoint adding
// all Transactions in txs.
func (c Client) AddTransactions(txs []Transaction) error {
	if len(txs) == 0 {
		return fmt.Errorf("`txs []Transaction` is empty or nil")
	}
	jsonBytes, err := json.Marshal(txs)
	if err != nil {
		return fmt.Errorf("json.Marshal(%+v): %v", txs, err)
	}
	buf := bytes.NewBuffer(jsonBytes)

	if _, err := c.doRequest(c.URL+transactionsURI, buf); err != nil {
		return err
	}
	return nil
}

func (c Client) doRequest(url string, body io.Reader) (*responseT, error) {
	method := "POST"
	if body == nil {
		method = "GET"
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest(%#v, %#v, %+v): %v",
			method, url, body, err)
	}
	c.setHeaders(req.Header)
	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http.Client.Do(http.Request%+v): %v", req, err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(resp.Status)
	}

	response := &responseT{}
	d := json.NewDecoder(resp.Body)
	if err := d.Decode(response); err != nil {
		return nil, fmt.Errorf("json.Decoder.Decode(): %v", err)
	}
	if response.Status != successStatus {
		return nil, fmt.Errorf("status: %#v, message: %#v",
			response.Status, response.Data.Message)
	}
	return response, nil
}

func (c Client) setHeaders(header http.Header) {
	header.Set("X-APIKEY", c.Key)
	header.Set("X-APISECRET", c.Secret)
	header.Set("Content-Type", "application/json")
}

type responseT struct {
	Status statusT `json:"status"` // "success", "fail" or "error"
	Data   dataT   `json:"data"`
}

type dataT struct {
	Total        uint64        `json:"total"` // total number of transactions
	Transactions []Transaction `json:"transactions"`
	Message      string        `json:"message,omitempty"`
}

type statusT string

const (
	successStatus = statusT("success")
	failStatus    = statusT("fail")
	errorStatus   = statusT("error")
)
