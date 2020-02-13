package bitcointax

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTransaction(t *testing.T) {
	var tx Transaction
	tx.Date = time.Now()
	data, err := tx.MarshalJSON()
	require.NoError(t, err)

	fmt.Println(string(data))
}
