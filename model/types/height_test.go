package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeightFromString(t *testing.T) {
	assert.Equal(t, Height(0), HeightFromString(""))
	assert.Equal(t, Height(0), HeightFromString("-100"))
	assert.Equal(t, Height(100), HeightFromString("100"))
	assert.Equal(t, Height(0), HeightFromString("foobar"))
}
func TestHeightValid(t *testing.T) {
	assert.True(t, Height(0).Valid())
}

func TestHeightString(t *testing.T) {
	assert.Equal(t, "0", HeightFromString("").String())
	assert.Equal(t, "100", Height(100).String())
}
