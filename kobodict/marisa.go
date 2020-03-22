package kobodict

import "io"

// Marisa is used by Reader and Writer for reading/writing Marisa tries. It is
// automatically set on supported platforms, but can be overridden.
var Marisa interface {
	MarisaReader
	MarisaWriter
}

func setMarisa(m interface {
	MarisaReader
	MarisaWriter
}) {
	Marisa = m
}

// MarisaReader represents a simplified abstraction for reading Marisa tries.
type MarisaReader interface {
	ReadAll(io.Reader) ([]string, error)
}

// MarisaWriter represents a simplified abstraction for writing Marisa tries.
type MarisaWriter interface {
	WriteAll(io.Writer, []string) error
}
