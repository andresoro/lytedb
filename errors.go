package lyte

import "errors"

var (
	// ErrValueNotPtr ..
	ErrValueNotPtr = errors.New("Valued passed into Get needs to be a pointer")
)
