package types

import "fmt"

// Height represents a block height
type Height uint64

// HeightFromString returns a new height parsed from a string
func HeightFromString(src string) Height {
	var h Height
	fmt.Sscanf(src, "%d", &h)
	return h
}

// Valid returns true if height value is valid
func (h Height) Valid() bool {
	return h >= 0
}

// String returns height text representation
func (h Height) String() string {
	return fmt.Sprintf("%d", h)
}
