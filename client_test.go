// github.com/canonical-ledgers/bitcointax
// Copyright 2018 Canonical Ledgers, LLC. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file distributed with this source code.

package bitcointax

import (
	"context"
	"flag"
	"fmt"
	"math"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	key, secret string

	c *Client
)

func TestMain(m *testing.M) {
	flag.StringVar(&key, "key", "", "bitcoin.tax key")
	flag.StringVar(&secret, "secret", "", "bitcoin.tax secret")
	flag.Parse()
	c = NewClient(key, secret, WithHTTPClient(&http.Client{Timeout: 5 * time.Second}))
	m.Run()
}

func TestListTransactions(t *testing.T) {
	require.NotEmpty(t, key, "key")
	require.NotEmpty(t, secret, "secret")

	testListTransactions(t, 0, 0)
	testListTransactions(t, 0, 50)
	testListTransactions(t, 25, 125)
	testListTransactions(t, 0, 2)
	testListTransactions(t, 20, math.MaxUint64)
}

func testListTransactions(t *testing.T, start, length uint) {
	now := time.Now()
	msg := fmt.Sprintf("ListTransactions(%v, %v, %v)", now.Year(), start, length)
	t.Run(msg, func(t *testing.T) {
		t.Parallel()
		txs, total, err := c.ListTransactions(context.Background(), uint(now.Year()), start, length)
		if length == 0 {
			length = 100
		}
		require.NoError(t, err)
		assert.Lenf(t, txs, int(min(total-min(start, total), length)),
			"len(txs) == min(total(%v)-min(start(%v), total), length(%v))",
			total, start, length)
	})
}

func min(l, r uint) uint {
	if l < r {
		return l
	}
	return r
}

func TestAddTransactions(t *testing.T) {
	require := require.New(t)
	require.NotEmpty(key, "key")
	require.NotEmpty(secret, "secret")
	assert := assert.New(t)

	err := c.AddTransactions(context.Background())
	assert.Error(err, "AddTransactions(nil)")

	txs := []Transaction{}
	err = c.AddTransactions(context.Background(), txs...)
	assert.Error(err, "AddTransactions(nil)")

	txs = []Transaction{{
		Date:     time.Now(),
		Action:   TxTypeIncome,
		Symbol:   "FCT",
		Currency: "USD",
		Volume:   5,
		Total:    25,
	}, {
		Date:     time.Now(),
		Action:   TxTypeIncome,
		Symbol:   "FCT",
		Currency: "USD",
		Volume:   5,
		Total:    25,
	}}
	err = c.AddTransactions(context.Background(), txs...)
	assert.NoError(err, "AddTransactions")
}
