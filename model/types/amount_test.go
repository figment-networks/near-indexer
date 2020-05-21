package types

import (
	"testing"

	"github.com/bmizerany/assert"
)

func TestNewAmount(t *testing.T) {
	assert.Equal(t, "0", NewAmount("").Raw())
	assert.Equal(t, "0", NewAmount("0").Raw())
	assert.Equal(t, "10", NewAmount("10").Raw())
}

func TestValid(t *testing.T) {
	assert.Equal(t, true, NewAmount("0").Valid())
	assert.Equal(t, true, NewAmount("1").Valid())
	assert.Equal(t, true, NewAmount("").Valid())
}

func TestFormat(t *testing.T) {
	assert.Equal(t, "104000.00001", NewAmount("104000000000000000000000000000").Format(5))
	assert.Equal(t, "480.00001", NewAmount("480000000972146165000000000").Format(5))
}
