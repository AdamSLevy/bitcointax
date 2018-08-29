// github.com/canonical-ledgers/bitcointax v1.0.0
// Copyright 2018 Canonical Ledgers, LLC. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file distributed with this source code.

package bitcointax

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Client provides methods for querying the bitcoin.tax REST API. You may set
// any additional http.Client settings directly on this type.
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
	url := c.URL + fmt.Sprintf("%v?taxyear=%v&start=%v",
		transactionsURI, taxyear.Year(), start)
	if limit != 0 {
		url = url + fmt.Sprintf("&limit=%v", limit)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, 0, err
	}
	c.setHeaders(req.Header)
	resp, err := c.Do(req)
	if err != nil {
		return nil, 0, err
	}

	var response jsonResponse
	d := json.NewDecoder(resp.Body)
	if err := d.Decode(&response); err != nil {
		return nil, 0, err
	}
	if response.Status != successStatus {
		return nil, 0, fmt.Errorf("status: %#v, message: %#v",
			response.Status, response.Message)
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
		return err
	}
	buf := bytes.NewBuffer(jsonBytes)
	url := c.URL + transactionsURI
	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return err
	}
	c.setHeaders(req.Header)
	resp, err := c.Do(req)
	if err != nil {
		return err
	}

	var response jsonResponse
	d := json.NewDecoder(resp.Body)
	if err := d.Decode(&response); err != nil {
		return err
	}
	if response.Status != successStatus {
		return fmt.Errorf("status: %#v, message: %#v",
			response.Status, response.Message)
	}
	return nil

}

func (c Client) setHeaders(header http.Header) {
	header.Set("X-APIKEY", c.Key)
	header.Set("X-APISECRET", c.Secret)
	header.Set("Content-Type", "application/json")
}

type jsonResponse struct {
	Status  statusType `json:"status"` // "success", "fail" or "error"
	Data    dataType   `json:"data"`
	Message string     `json:"message,omitempty"`
}

type dataType struct {
	Total        uint64        `json:"total"` // total number of transactions
	Transactions []Transaction `json:"transactions"`
}

type statusType string

const (
	successStatus = statusType("success")
	failStatus    = statusType("fail")
	errorStatus   = statusType("error")
)