package cli

import (
	"unicode/utf8"
)

type Argument struct {
	mandatory   bool
	name        string
	nameLen     int
	ptr         interface{}
	description string
}

func newArgument(mandatory bool, name, description string, ptr interface{}) *Argument {
	return &Argument{mandatory, name, utf8.RuneCountInString(name), ptr, description}
}
