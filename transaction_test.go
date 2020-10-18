package bitcointax

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTransaction(t *testing.T) {
	require := require.New(t)
	var tx Transaction
	tx.Date = time.Now()
	data, err := tx.MarshalJSON()
	require.NoError(err)

	var tx2 Transaction
	require.NoError(tx2.UnmarshalJSON(data))
	require.Equal(tx.Date.Unix(), tx2.Date.Unix())
	tx.Date, tx2.Date = time.Time{}, time.Time{}
	require.Equal(tx, tx2)

	fmt.Println(string(data))
}
