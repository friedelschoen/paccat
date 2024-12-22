package recipe

import "errors"

var (
	NoAttributeError = errors.New("unknown attribute")
	NoOutputError    = errors.New("did not produce output")
)
