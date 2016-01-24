package cli

import (
	"unicode/utf8"
)

type Option struct {
	name        string
	nameLen     int
	ptrBool     *bool
	description string
}

func newOption(name, description string, ptrBool *bool) *Option {
	return &Option{name, utf8.RuneCountInString(name) + 2, ptrBool, description}
}
