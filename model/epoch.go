package model

import (
	"errors"
	"time"
)

type Epoch struct {
	ID              uint      `json:"-"`
	UUID            string    `json:"id"`
	StartHeight     uint64    `json:"start_height"`
	StartTime       time.Time `json:"start_time"`
	EndHeight       uint64    `json:"end_height"`
	EndTime         time.Time `json:"end_time"`
	BlocksCount     uint      `json:"blocks_count"`
	ValidatorsCount uint      `json:"validators_count"`
}

func (Epoch) TableName() string {
	return "epochs"
}

func (e Epoch) Validate() error {
	if e.UUID == "" {
		return errors.New("uuid is not provided")
	}
	if e.StartHeight == 0 {
		return errors.New("start height is not provided")
	}
	if e.StartTime.IsZero() {
		return errors.New("start time is not provided")
	}
	if e.EndHeight == 0 {
		return errors.New("end height is not provided")
	}
	if e.EndTime.IsZero() {
		return errors.New("end time is not provided")
	}
	return nil
}
