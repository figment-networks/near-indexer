package types

import "fmt"

type Height uint64

func HeightFromString(src string) Height {
	var h Height
	fmt.Sscanf(src, "%d", &h)
	return h
}

func (h Height) Valid() bool {
	return h >= 0
}
