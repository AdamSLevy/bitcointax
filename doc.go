// github.com/canonical-ledgers/bitcointax
// Copyright 2018 Canonical Ledgers, LLC. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file distributed with this source code.

// Package bitcointax implements a Transaction type which properly Marshals and
// Unmarshals into the valid JSON format for the bitcoin.tax REST API, as well
// as a Client with methods for listing and adding Transactions.
//
// The bitcoin.tax REST API is documented here: https://www.bitcoin.tax/api
package bitcointax
