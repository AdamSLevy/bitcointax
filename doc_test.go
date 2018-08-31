package bitcointax_test

import (
	"fmt"
	"time"

	"github.com/canonical-ledgers/bitcointax"
)

func ExampleClient_ListTransactions() {
	key, secret := "API KEY", "API SECRET"
	c := bitcointax.NewClient(key, secret)
	txs, total, err := c.ListTransactions(time.Now(), 0, 50)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	fmt.Printf("Total transactions %+v\n", total)
	fmt.Printf("Transactions %+v\n", txs)
}

func ExampleClient_AddTransactions() {
	key, secret := "API KEY", "API SECRET"
	c := bitcointax.NewClient(key, secret)
	txs := []bitcointax.Transaction{{
		Date:     time.Now(),
		Action:   bitcointax.IncomeTx,
		Symbol:   "FCT",
		Currency: "USD",
		Volume:   5,
		Total:    25,
	}, {
		Date:     time.Now(),
		Action:   bitcointax.IncomeTx,
		Symbol:   "FCT",
		Currency: "USD",
		Volume:   5,
		Total:    25,
	}}
	if err := c.AddTransactions(txs); err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
}
