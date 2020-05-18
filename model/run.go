package model

type Run struct {
	Model

	Height   uint64 `json:"height"`
	Success  bool   `json:"success"`
	Error    string `json:"error"`
	Duration int64  `json:"duration"`
}
