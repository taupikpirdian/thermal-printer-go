//go:build !windows

package main

import (
	"fmt"
	"runtime"
)

func openRawPrinter(name string) (rawPrinter, error) {
	return nil, fmt.Errorf("raw Windows printer %q is unsupported on %s", name, runtime.GOOS)
}
