// github.com/canonical-ledgers/bitcointax v1.0.0
// Copyright 2018 Canonical Ledgers, LLC. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file distributed with this source code.

package bitcointax

import (
	"flag"
	"fmt"
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
	c = NewClient(key, secret)
	m.Run()
}

func TestListTransactions(t *testing.T) {
	t.Parallel()
	require.NotEmpty(t, key, "key")
	require.NotEmpty(t, secret, "secret")

	t.Run("start: 0 length: 0", func(t *testing.T) { testListTransactions(t, 0, 0) })
	testListTransactions(t, 0, 50)
	testListTransactions(t, 25, 125)
	testListTransactions(t, 0, 2)
}

func testListTransactions(t *testing.T, start, length uint64) {
	now := time.Now()
	msg := fmt.Sprintf("ListTransactions(%v, %v, %v)", now.UTC(), start, length)
	t.Run(msg, func(t *testing.T) {
		t.Parallel()
		txs, total, err := c.ListTransactions(now, start, length)
		if length == 0 {
			length = 100
		}
		require.NoError(t, err)
		assert.True(t, len(txs) == int(min(total-min(start, total), length)),
			"len(txs) == min(total-min(start, total), length)")
	})
}

func min(l, r uint64) uint64 {
	if l < r {
		return l
	}
	return r
}

func TestAddTransactions(t *testing.T) {
	t.Parallel()
	require := require.New(t)
	require.NotEmpty(key, "key")
	require.NotEmpty(secret, "secret")
	assert := assert.New(t)

	err := c.AddTransactions(nil)
	assert.Error(err, "AddTransactions(nil)")

	txs := []Transaction{}
	err = c.AddTransactions(txs)
	assert.Error(err, "AddTransactions(nil)")

	txs = []Transaction{{
		Date:     Time(time.Now()),
		Action:   IncomeTx,
		Symbol:   "FCT",
		Currency: "USD",
		Volume:   5,
		Total:    25,
	}, {
		Date:     Time(time.Now()),
		Action:   IncomeTx,
		Symbol:   "FCT",
		Currency: "USD",
		Volume:   5,
		Total:    25,
	}}
	err = c.AddTransactions(txs)
	assert.NoError(err, "AddTransactions")
}
