package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAmount(t *testing.T) {
	assert.Equal(t, "0", NewAmount("").String())
	assert.Equal(t, "0", NewAmount("0").String())
	assert.Equal(t, "10", NewAmount("10").String())
}
