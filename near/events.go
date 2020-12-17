package near

import (
	"encoding/json"
	"fmt"
)

type ReasonName string

const (
	// Validator didn't produce enough blocks
	NotEnoughBlocks ReasonName = "NotEnoughBlocks"

	// Validator didn't produce enough chunks
	NotEnoughChunks ReasonName = "NotEnoughChunks"

	// Validator stake is now below threshold
	NotEnoughStake ReasonName = "NotEnoughStake"

	// Validator unstaked themselves
	Unstaked ReasonName = "Unstaked"

	// Enough stake but is not chosen because of seat limits.
	DidNotGetASeat ReasonName = "DidNotGetASeat"

	// Slashed validators are kicked out.
	Slashed ReasonName = "Slashed"
)

type KickoutReason struct {
	Name ReasonName
	Data map[string]interface{}
}

func (r *KickoutReason) UnmarshalJSON(data []byte) error {
	var dst interface{}
	if err := json.Unmarshal(data, &dst); err != nil {
		return err
	}

	switch t := dst.(type) {
	case string:
		r.Name = ReasonName(t)
	case map[string]interface{}:
		for k, v := range t {
			r.Name = ReasonName(k)
			r.Data = v.(map[string]interface{})
			break
		}
	default:
		return fmt.Errorf("unexpected type: %v", t)
	}

	return nil
}
