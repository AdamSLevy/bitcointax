package bitcointax

import (
	"fmt"
	"time"
)

// Transaction represents the transaction JSON object sent and received by
// bitcoin.tax.
type Transaction struct {
	Date        Time    `json:"date"`                  // Date of transaction in ISO 8601 or unix timestamp
	Action      TxType  `json:"action"`                // Type of transaction, e.g. "BUY" or "SELL"
	Symbol      string  `json:"symbol"`                // Crypto-currency symbol
	Currency    string  `json:"currency"`              // ISO 4217 currency symbol or "BTC", "LTC" or "XRP"
	Volume      float64 `json:"volume"`                // Number of units of symbol
	Exchange    string  `json:"exchange,omitempty"`    // Exchange or wallet name, e.g. "Coinbase"
	Exchangeid  string  `json:"exchangeid,omitempty"`  // Exchange or wallet's unique transaction id
	Price       float64 `json:"price,omitempty"`       // Price of symbol in units of currency (if total is unknown)	default: total/volume or average daily rate
	Total       float64 `json:"total,omitempty"`       // Total value of transaction in currency (if price is unknown)	default: price * volume
	Fee         float64 `json:"fee,omitempty"`         // Fee for transaction in units of currency	                        default: 0
	Feecurrency string  `json:"feecurrency,omitempty"` // Currency for transaction fee	                                default: currency
	Memo        string  `json:"memo,omitempty"`        // Note for transaction
	Txhash      string  `json:"txhash,omitempty"`      // Hash value from symbol's blockchain
	Sender      string  `json:"sender,omitempty"`      // Coin address of sender
	Recipient   string  `json:"recipient,omitempty"`   // Coin address of recipient
	ID          string  `json:"id,omitempty"`
	Message     string  `json:"message,omitempty"`
}

// TxType is a string type that represents the valid types of
// transactions used by bitcoin.tax.
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

type Time time.Time

func (t Time) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%#v", time.Time(t).Format("2006-01-02T15:04:05-0700"))),
		nil

}

func (t *Time) UnmarshalJSON(b []byte) error {
	time := time.Time{}
	if err := time.UnmarshalJSON(b); err != nil {
		return err
	}
	*t = Time(time)
	return nil
}
