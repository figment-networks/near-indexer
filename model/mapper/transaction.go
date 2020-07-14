package mapper

import (
	"github.com/figment-networks/near-indexer/model"
	"github.com/figment-networks/near-indexer/model/types"
	"github.com/figment-networks/near-indexer/model/util"
	"github.com/figment-networks/near-indexer/near"
)

// Transaction constructs a new transaction record from the chain input
func Transaction(block *near.Block, input *near.Transaction) (*model.Transaction, error) {
	t := &model.Transaction{
		Height:    types.Height(block.Header.Height),
		Time:      util.ParseTimeFromString(block.Header.Timestamp),
		Signer:    input.SignerID,
		SignerKey: input.PublicKey,
		Receiver:  input.ReceiverID,
		Signature: input.Signature,
	}

	for _, action := range input.Actions {
		switch action.(type) {
		case map[string]interface{}:
			// TODO
		}
	}

	return t, nil
}

// Transactions constructs a set of transactions from the chain input
func Transactions(block *near.Block, d *near.TransactionDetails) ([]model.Transaction, error) {
	result := []model.Transaction{}

	return result, nil
}
