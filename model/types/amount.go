package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
)

var (
	reTrailingZeroes = regexp.MustCompile(`\.?0*$`)
)

// Amount represense a NEAR yocto
type Amount struct {
	raw string
}

// NewAmount returns a new amount from the given string
func NewAmount(src string) Amount {
	if src == "" {
		src = "0"
	}
	return Amount{src}
}

// NewInt64Amount returns a new amount for the given int64 value
func NewInt64Amount(val int64) Amount {
	amount := Amount{}
	amount.Scan(val)
	return amount
}

// MarshalJSON returns a JSON representation of amount
func (a Amount) MarshalJSON() ([]byte, error) {
	// Following code will render both raw and formatted values
	// return json.Marshal(map[string]string{
	// 	"raw":       a.raw,
	// 	"formatted": a.String(),
	// })
	return json.Marshal(a.raw)
}

// Valid returns true if amount is valid
func (a Amount) Valid() bool {
	return a.raw != ""
}

// Raw returns a raw amount value
func (a Amount) Raw() string {
	return a.raw
}

// String returns a formatted amount value
func (a Amount) String() string {
	return a.Format(nearFractionDigits)
}

// Value returns a serialized value
func (a Amount) Value() (driver.Value, error) {
	return a.raw, nil
}

// Scan assigns the value from interface
func (a *Amount) Scan(value interface{}) error {
	v, ok := value.(string)
	if !ok {
		return errors.New("invalid amount")
	}
	a.raw = v
	return nil
}

// Format formats the amount with a fixed digits length
func (a Amount) Format(digits int) string {
	val := bigint(a.raw)
	if val.Cmp(nearZero) == 0 {
		return "0"
	}

	roundingExp := nearNominationExp - digits - 1
	if roundingExp > 0 {
		val.Add(val, &nearRoundingOffsets[roundingExp])
	}

	strval := val.String()

	idx := len(strval) - nearNominationExp
	if idx < 0 {
		idx = 0
	}

	wholeStr := strval[0:idx]
	fractionStr := strval[idx:]
	fractionStr = fmt.Sprintf(nearFractionFormat, fractionStr)[0:digits]
	fullStr := trimTrailingZeroes(wholeStr + "." + fractionStr)

	return fullStr
}

func trimTrailingZeroes(src string) string {
	return reTrailingZeroes.ReplaceAllString(src, "")
}
