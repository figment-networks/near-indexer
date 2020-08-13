package near

import (
	"encoding/json"
	"fmt"
)

// DecodeActions decodes all actions in the transactions
func DecodeActions(t *Transaction) []Action {
	result := make([]Action, len(t.Actions))

	for idx, act := range t.Actions {
		switch act.(type) {
		case string:
			name := act.(string)
			switch name {
			case ActionCreateAccount:
				result[idx].Type = name
				result[idx].Data = &CreateAccountAction{}
			default:
				panic(fmt.Sprintf("unhandled action type: %v", name))
			}
		case map[string]interface{}:
			actmap := act.(map[string]interface{})
			for k, v := range actmap {
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
