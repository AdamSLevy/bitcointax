// github.com/canonical-ledgers/bitcointax v1.0.1
// Copyright 2018 Canonical Ledgers, LLC. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file distributed with this source code.

package bitcointax

import (
	"encoding/json"
	"time"
)

// TxType represents the valid types of transactions used by bitcoin.tax.
type TxType string

// Valid Transaction Types accepted by bitcoin.tax for "Action".
const (
	TxTypeSell       = "SELL"     // Selling crypto-currency to fiat or BTC
	TxTypeBuy        = "BUY"      // Buy crypto-currency for fiat or BTC
	TxTypeIncome     = "INCOME"   // General income
	TxTypeGiftIncome = "GIFTIN"   // Income received as a gift or tip
	TxTypeMining     = "MINING"   // Income received from mining
	TxTypeSpend      = "SPEND"    // General spending of crypto-currencies
	TxTypeGift       = "GIFT"     // Spending as a gift or tip
	TxTypeDonation   = "DONATION" // Spending to a registered charity
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

type _Transaction Transaction
type __Transaction struct {
	Date *_Time `json:"date"`
	*_Transaction
}

func (t Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(__Transaction{
		Date:         (*_Time)(&t.Date),
		_Transaction: (*_Transaction)(&t),
	})
}

func (t *Transaction) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &__Transaction{
		Date:         (*_Time)(&t.Date),
		_Transaction: (*_Transaction)(t),
	})
}

const timeFmtISO8601 = "2006-01-02T15:04:05-07:00"

type _Time time.Time

func (t *_Time) Time() *time.Time { return (*time.Time)(t) }

func (t *_Time) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	var err error
	*(*time.Time)(t), err = time.Parse(timeFmtISO8601, str)
	return err
}

func (t _Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Time().Format(timeFmtISO8601))
}
