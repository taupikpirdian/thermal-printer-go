//go:build windows

package main

import "github.com/alexbrainman/printer"

func openRawPrinter(name string) (rawPrinter, error) {
	return printer.Open(name)
}
