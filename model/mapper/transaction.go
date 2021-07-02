package mapper

import (
	"encoding/json"
	"fmt"
	"strings"

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
		Sender:    tx.SignerID,
		Receiver:  tx.ReceiverID,
		PublicKey: tx.PublicKey,
		Signature: tx.Signature,
		GasBurnt:  fmt.Sprintf("%v", input.TransactionOutcome.Outcome.GasBurnt),
		Fee:       util.CalculateTransactionFee(*input),
	}

	outcome, err := json.Marshal(input.TransactionOutcome)
	if err != nil {
		return nil, err
	}
	t.Outcome = outcome

	receipt, err := json.Marshal(input.ReceiptsOutcome)
	if err != nil {
		return nil, err
	}
	t.Receipt = receipt

	// Status field may be represented by different types depending on the situation
	switch val := input.Status.(type) {
	case string:
		t.Success = strings.ToLower(val) != "failure"
	case map[string]interface{}:
		t.Success = val["SuccessValue"] != nil
	default:
		return nil, fmt.Errorf("unsupported tx status attribute: %v", input.Status)
	}

	if actions := near.DecodeActions(&tx); len(actions) > 0 {
		reencoded, err := json.Marshal(actions)
		if err != nil {
			return nil, err
		}
		t.Actions = reencoded
		t.ActionsCount = len(actions)
	}

	if err := t.Validate(); err != nil {
		raw, _ := json.Marshal(tx)
		return nil, fmt.Errorf("transaction (%s) is invalid: %w", string(raw), err)
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
