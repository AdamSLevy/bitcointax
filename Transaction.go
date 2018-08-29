package bitcointax

type Transaction struct {
	Date        string          `json:"date"`                  // Date of transaction in ISO 8601 or unix timestamp
	Action      TransactionType `json:"action"`                // Type of transaction, e.g. "BUY" or "SELL"
	Symbol      string          `json:"symbol"`                // Crypto-currency symbol
	Currency    string          `json:"currency"`              // ISO 4217 currency symbol or "BTC", "LTC" or "XRP"
	Volume      float64         `json:"volume"`                // Number of units of symbol
	Exchange    string          `json:"exchange,omitempty"`    // Exchange or wallet name, e.g. "Coinbase"
	Exchangeid  string          `json:"exchangeid,omitempty"`  // Exchange or wallet's unique transaction id
	Price       float64         `json:"price,omitempty"`       // Price of symbol in units of currency (if total is unknown)	default: total/volume or average daily rate
	Total       float64         `json:"total,omitempty"`       // Total value of transaction in currency (if price is unknown)	default: price * volume
	Fee         float64         `json:"fee,omitempty"`         // Fee for transaction in units of currency	                        default: 0
	Feecurrency string          `json:"feecurrency,omitempty"` // Currency for transaction fee	                                default: currency
	Memo        string          `json:"memo,omitempty"`        // Note for transaction
	Txhash      string          `json:"txhash,omitempty"`      // Hash value from symbol's blockchain
	Sender      string          `json:"sender,omitempty"`      // Coin address of sender
	Recipient   string          `json:"recipient,omitempty"`   // Coin address of recipient

}

type TransactionType string

const (
	Sell       = TransactionType("SELL")     // Selling crypto-currency to fiat or BTC
	Buy        = TransactionType("BUY")      // Buy crypto-currency for fiat or BTC
	Income     = TransactionType("INCOME")   // General income
	GiftIncome = TransactionType("GIFTIN")   // Income received as a gift or tip
	Mining     = TransactionType("MINING")   // Income received from mining
	Spend      = TransactionType("SPEND")    // General spending of crypto-currencies
	Gift       = TransactionType("GIFT")     // Spending as a gift or tip
	Donation   = TransactionType("DONATION") // Spending to a registered charity
)
