package dfmt

import (
	"fmt"
	"os"

	"github.com/cjtoolkit/cfupdater/app/settings"
)

func Print(a ...interface{}) (n int, err error) {
	if !*settings.Debug {
		return
	}

	n, err = fmt.Fprint(os.Stdout, a...)
	return
}

func Printf(format string, a ...interface{}) (n int, err error) {
	if !*settings.Debug {
		return
	}

	n, err = fmt.Fprintf(os.Stdout, format, a...)
	return
}

func Println(a ...interface{}) (n int, err error) {
	if !*settings.Debug {
		return
	}

	n, err = fmt.Fprintln(os.Stdout, a...)
	return
}
