package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapValue(t *testing.T) {
	m := NewMap()
	m["a"] = "hello"
	m["b"] = []string{"hello"}
	m["c"] = nil

	val, err := m.Value()

	assert.NoError(t, err)
	assert.NotNil(t, val)
}

func TestMapScan(t *testing.T) {
	data := `{"a":"hello","b":["hello"],"c":null}`
	m := NewMap()

	assert.Error(t, errMapInvalidSource, m.Scan(data))
	assert.NoError(t, m.Scan([]byte(data)))
	assert.Equal(t, "hello", m["a"])
	assert.Equal(t, []interface{}{"hello"}, m["b"])
	assert.Equal(t, nil, m["c"])
}
