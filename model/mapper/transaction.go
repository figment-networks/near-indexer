package mapper

import (
	"encoding/json"
	"log"

	"github.com/figment-networks/near-indexer/model"
	"github.com/figment-networks/near-indexer/model/types"
	"github.com/figment-networks/near-indexer/model/util"
	"github.com/figment-networks/near-indexer/near"
)

// Transaction constructs a new transaction record from the chain input
func Transaction(block *near.Block, input *near.TransactionDetails) (*model.Transaction, error) {
	tx := input.Transaction

	t := &model.Transaction{
		Hash:      tx.Hash,
		BlockHash: block.Header.Hash,
		Height:    types.Height(block.Header.Height),
		Time:      util.ParseTime(block.Header.Timestamp),
		Signer:    tx.SignerID,
		SignerKey: tx.PublicKey,
		Receiver:  tx.ReceiverID,
		Signature: tx.Signature,
		Amount:    types.NewAmount("0"),
	}

	rawActions, err := json.Marshal(tx.Actions)
	if err != nil {
		log.Println("cant marshal actions:", err)
	} else {
		t.Actions = rawActions
	}

	return t, nil
}

// Transactions constructs a set of transactions from the chain input
func Transactions(block *near.Block, details []near.TransactionDetails) ([]model.Transaction, error) {
	result := []model.Transaction{}

	for _, t := range details {
		transaction, err := Transaction(block, &t)
		if err != nil {
			return nil, err
		}
		result = append(result, *transaction)
	}

	return result, nil
}
