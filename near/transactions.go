package near

import (
	"encoding/json"
	"fmt"
)

const (
	// Description for all types:
	// https://docs.near.org/docs/concepts/transaction

	ActionCreateAccount  = "CreateAccount"  // make a new account (for a person, company, contract, car, refrigerator, etc)
	ActionDeployContract = "DeployContract" // deploy a new contract (with its own account)
	ActionFunctionCall   = "FunctionCall"   // invoke a method on a contract (with budget for compute and storage)
	ActionTransfer       = "Transfer"       // transfer tokens from one account to another
	ActionStake          = "Stake"          // express interest in becoming a proof-of-stake validator at the next available opportunity
	ActionAddKey         = "AddKey"         // add a key to an existing account (either FullAccess or FunctionCall access)
	ActionDeleteKey      = "DeleteKey"      // delete an existing key from an account
	ActionDeleteAccount  = "DeleteAccount"  // delete an account (and transfer balance to a beneficiary account)
)

type Action struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type CreateAccountAction struct {
}

type DeployContractAction struct {
	// Code []byte `json:"code"` // skipped due to large payloads
}

type FunctionCallAction struct {
	MethodName string `json:"method_name"`
	// Args       string `json:"args,omitempty"` // skipped due to large payloads
	Deposit string `json:"deposit"`
	Gas     int64  `json:"gas"`
}

type TransferAction struct {
	Deposit string `json:"deposit"`
}

type StakeAction struct {
	PublicKey string `json:"public_key"`
	Amount    string `json:"stake"`
}

type AddKeyAction struct {
	PublicKey string    `json:"public_key"`
	AccessKey AccessKey `json:"access_key"`
}

type DeleteKeyAction struct {
	PublicKey string `json:"public_key"`
}

type DeleteAccountAction struct {
	BeneficiaryID string `json:"beneficiary_id"`
}

// DecodeActions decodes all actions in the transactions
func DecodeActions(t *Transaction) []Action {
	result := make([]Action, len(t.Actions))

	for idx, act := range t.Actions {
		switch data := act.(type) {
		case string:
			switch data {
			case ActionCreateAccount:
				result[idx].Type = data
				result[idx].Data = &CreateAccountAction{}
			default:
				panic(fmt.Sprintf("unhandled action type: %v", data))
			}
		case map[string]interface{}:
			for k, v := range data {
				var dst interface{}
				var buf json.RawMessage

				b, err := json.Marshal(v)
				if err != nil {
					panic(err)
				}
				buf = b

				switch k {
				case ActionCreateAccount:
					dst = decodeAction(buf, &CreateAccountAction{})
				case ActionFunctionCall:
					dst = decodeAction(buf, &FunctionCallAction{})
				case ActionTransfer:
					dst = decodeAction(buf, &TransferAction{})
				case ActionStake:
					dst = decodeAction(buf, &StakeAction{})
				case ActionAddKey:
					dst = decodeAction(buf, &AddKeyAction{})
				case ActionDeleteKey:
					dst = decodeAction(buf, &DeleteKeyAction{})
				case ActionDeleteAccount:
					dst = decodeAction(buf, &DeleteAccountAction{})
				case ActionDeployContract:
					dst = decodeAction(buf, &DeployContractAction{})
				default:
					panic(fmt.Sprintf("unhandled action type: %v", k))
				}

				result[idx].Type = k
				result[idx].Data = dst

				break
			}
		default:
			panic(fmt.Sprintf("unhandled action type: %v", act))
		}
	}

	return result
}

func decodeAction(val json.RawMessage, base interface{}) interface{} {
	if err := json.Unmarshal(val, &base); err != nil {
		panic(err)
	}
	return base
}
