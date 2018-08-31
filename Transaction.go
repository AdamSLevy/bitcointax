// github.com/canonical-ledgers/bitcointax v1.0.1
// Copyright 2018 Canonical Ledgers, LLC. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file distributed with this source code.

package bitcointax

import (
	"encoding/json"
	"fmt"
	"time"
)

// Transaction represents the transaction JSON object sent and received by the
// bitcoin.tax REST API.
type Transaction struct {
	Date        time.Time `json:"date"`                  // Date of transaction in ISO 8601 or unix timestamp [required]
	Action      TxType    `json:"action"`                // Type of transaction, e.g. "BUY" or "SELL"         [required]
	Symbol      string    `json:"symbol"`                // Crypto-currency symbol                            [required]
	Currency    string    `json:"currency"`              // ISO 4217 currency symbol or "BTC", "LTC" or "XRP" [required]
	Volume      float64   `json:"volume"`                // Number of units of symbol                         [required]
	Exchange    string    `json:"exchange,omitempty"`    // Exchange or wallet name, e.g. "Coinbase"
	ExchangeID  string    `json:"exchangeid,omitempty"`  // Exchange or wallet's unique transaction id
	Price       float64   `json:"price,omitempty"`       // Price of symbol in units of currency (if total is unknown)   default: total/volume or average daily rate
	Total       float64   `json:"total,omitempty"`       // Total value of transaction in currency (if price is unknown) default: price * volume
	Fee         float64   `json:"fee,omitempty"`         // Fee for transaction in units of currency                     default: 0
	FeeCurrency string    `json:"feecurrency,omitempty"` // Currency for transaction fee                                 default: currency
	Memo        string    `json:"memo,omitempty"`        // Note for transaction
	TxHash      string    `json:"txhash,omitempty"`      // Hash value from symbol's blockchain
	Sender      string    `json:"sender,omitempty"`      // Coin address of sender
	Recipient   string    `json:"recipient,omitempty"`   // Coin address of recipient
	ID          string    `json:"id,omitempty"`          // Unique ID assigned internally by bitcoin.tax
}

type transaction struct {
	Date        timeT   `json:"date"`
	Action      TxType  `json:"action"`
	Symbol      string  `json:"symbol"`
	Currency    string  `json:"currency"`
	Volume      float64 `json:"volume"`
	Exchange    string  `json:"exchange,omitempty"`
	ExchangeID  string  `json:"exchangeid,omitempty"`
	Price       float64 `json:"price,omitempty"`
	Total       float64 `json:"total,omitempty"`
	Fee         float64 `json:"fee,omitempty"`
	FeeCurrency string  `json:"feecurrency,omitempty"`
	Memo        string  `json:"memo,omitempty"`
	TxHash      string  `json:"txhash,omitempty"`
	Sender      string  `json:"sender,omitempty"`
	Recipient   string  `json:"recipient,omitempty"`
	ID          string  `json:"id,omitempty"`
}

// timeT is a wrapper around time.Time that implements the json.Marshaler
// interface such that it is compatible with what the bitcoin.tax API expects
// in the Transaction.Date field.
type timeT time.Time

// TxType represents the valid types of transactions used by bitcoin.tax.
type TxType string

// Valid TxTypes
const (
	SellTx       = TxType("SELL")     // Selling crypto-currency to fiat or BTC
	BuyTx        = TxType("BUY")      // Buy crypto-currency for fiat or BTC
	IncomeTx     = TxType("INCOME")   // General income
	GiftIncomeTx = TxType("GIFTIN")   // Income received as a gift or tip
	MiningTx     = TxType("MINING")   // Income received from mining
	SpendTx      = TxType("SPEND")    // General spending of crypto-currencies
	GiftTx       = TxType("GIFT")     // Spending as a gift or tip
	DonationTx   = TxType("DONATION") // Spending to a registered charity
)

// MarshalJSON implements the json.Marshaler interface such that it is
// compatible with what the bitcoin.tax API expects in the Transaction.Date
// field.
func (t Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(transaction{
		Date:        timeT(t.Date),
		Action:      t.Action,
		Symbol:      t.Symbol,
		Currency:    t.Currency,
		Volume:      t.Volume,
		Exchange:    t.Exchange,
		ExchangeID:  t.ExchangeID,
		Price:       t.Price,
		Total:       t.Total,
		Fee:         t.Fee,
		FeeCurrency: t.FeeCurrency,
		Memo:        t.Memo,
		TxHash:      t.TxHash,
		Sender:      t.Sender,
		Recipient:   t.Recipient,
		ID:          t.ID,
	})
}

// MarshalJSON marshals t into a byte slice with a string representation of the
// unix timestamp returned by time.Time(t).Unix().
func (t timeT) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%v", time.Time(t).Unix())), nil
}
