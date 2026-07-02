package main

type rawPrinter interface {
	Write([]byte) (int, error)
	StartRawDocument(string) error
	EndDocument() error
	Close() error
}
