package bitcointax_test

import (
	"flag"
	"fmt"
	"time"

	"github.com/canonical-ledgers/bitcointax"
)

func ExampleClient_ListTransactions() {
	key := flag.String("key", "", "bitcoin.tax key")
	secret := flag.String("secret", "", "bitcoin.tax secret")
	flag.Parse()
	c := bitcointax.NewClient(*key, *secret)
	txs, total, err := c.ListTransactions(time.Now(), 0, 50)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
	fmt.Printf("Total transactions %+v", total)
	fmt.Printf("Transactions %+v", txs)
}

func ExampleClient_AddTransactions() {
	key := flag.String("key", "", "bitcoin.tax key")
	secret := flag.String("secret", "", "bitcoin.tax secret")
	flag.Parse()
	c := bitcointax.NewClient(*key, *secret)
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
		fmt.Printf("error: %v", err)
		return
	}
}
