package kobodict

import "io"

// Marisa is used by Reader and Writer for reading/writing Marisa tries. It is
// automatically set on supported platforms if
// github.com/geek1011/dictutil/kobodict/marisa is imported, but can be
// overridden manually.
var Marisa interface {
	MarisaReader
	MarisaWriter
}

// MarisaReader represents a simplified abstraction for reading Marisa tries.
type MarisaReader interface {
	ReadAll(io.Reader) ([]string, error)
}

// MarisaWriter represents a simplified abstraction for writing Marisa tries.
type MarisaWriter interface {
	WriteAll(io.Writer, []string) error
}
