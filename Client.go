package bitcointax

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Client embeds http.Client and provides methods for querying bitcoin.tax.
type Client struct {
	URL    string
	Key    string
	Secret string
	http.Client
}

const bitcoinTaxAPIURL = "https://api.bitcoin.tax/v1"
const transactionsURI = "/transactions"

// NewClient returns a new Client with the given key and secret and with the
// default API endpoint (https://api.bitcoin.tax/v1) as the URL.
func NewClient(key, secret string) *Client {
	return &Client{URL: bitcoinTaxAPIURL, Key: key, Secret: secret}
}

// ListTransactions makes a GET request to the /transactions endpoint and
// returns a slice of transactions for the given taxyear.Year() from the start
// transaction index and returning no more than limit transactions. If limit is
// 0 than the default limit of 100 will be used.
func (c Client) ListTransactions(taxyear time.Time,
	start, limit uint64) ([]Transaction, uint64, error) {
	if limit == 0 {
		limit = 100
	}
	url := c.URL + fmt.Sprintf("%v?taxyear=%v&start=%v&limit=%v",
		transactionsURI, taxyear.Year(), start, limit)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, 0, err
	}
	c.setHeaders(req.Header)
	resp, err := c.Do(req)
	if err != nil {
		return nil, 0, err
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}
	var response jsonResponse
	if err := json.Unmarshal(respBytes, &response); err != nil {
		return nil, 0, err
	}
	if response.Status != successStatus {
		var response map[string]interface{}
		json.Unmarshal(respBytes, &response)
		return nil, 0, fmt.Errorf("%+v", response)
	}
	return response.Data.Transactions, response.Data.Total, nil
}

// AddTransactions makes a POST request to the /transactions endpoint adding the
// Transactions t.
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

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var response jsonResponse
	if err := json.Unmarshal(respBytes, &response); err != nil {
		return err
	}
	if response.Status != successStatus {
		var response map[string]interface{}
		json.Unmarshal(respBytes, &response)
		return fmt.Errorf("%+v", response)
	}
	return nil

}

func (c Client) setHeaders(header http.Header) {
	header.Set("X-APIKEY", c.Key)
	header.Set("X-APISECRET", c.Secret)
	header.Set("Content-Type", "application/json")
}

type jsonResponse struct {
	Status statusType `json:"status"` // "success", "fail" or "error"
	Data   dataType   `json:"data"`
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
