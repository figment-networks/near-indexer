package model

import "github.com/figment-networks/near-indexer/model/types"

type Run struct {
	Model

	Height   types.Height `json:"height"`
	Success  bool         `json:"success"`
	Error    *string      `json:"error"`
	Duration int64        `json:"duration"`
}
