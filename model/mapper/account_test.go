package mapper

import (
	"testing"
	"time"

	"github.com/figment-networks/near-indexer/model/types"
	"github.com/figment-networks/near-indexer/near"
	"github.com/stretchr/testify/assert"
)

func TestAccountFromValidator(t *testing.T) {
	block := &near.Block{
		Header: near.BlockHeader{
			Height:    10885359,
			Timestamp: 1596166782911378000,
		},
	}

	validator := &near.Validator{
		AccountID: "account",
		Stake:     "1000",
	}

	acc, err := AccountFromValidator(block, validator)

	assert.NoError(t, err)
	assert.NotNil(t, acc)
	assert.Equal(t, validator.AccountID, acc.Name)
	assert.Equal(t, types.Height(block.Header.Height), acc.StartHeight)
	assert.Equal(t, types.Height(block.Header.Height), acc.LastHeight)
	assert.Equal(t, "2020-07-31T03:39:42Z", acc.StartTime.UTC().Format(time.RFC3339))
	assert.Equal(t, "2020-07-31T03:39:42Z", acc.LastTime.UTC().Format(time.RFC3339))
	assert.Equal(t, types.NewAmount("1000"), acc.StakingBalance)
}
