package model

import "time"

type AccessKey struct {
	Model

	StartHeight uint64    `json:"start_height"`
	StartTime   time.Time `json:"start_time"`
	PublicKey   string    `json:"public_key"`
	AccountID   string    `json:"account_id"`
	Nonce       int       `json:"nonce"`
	Permission  string    `json:"permission"`
}
