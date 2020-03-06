// Package marisa provides a self-contained SWIG wrapper for marisa-trie
// (https://github.com/s-yata/marisa-trie).
package marisa

//go:generate go run -tags libmarisa_generate libmarisa_generate.go

//#cgo CPPFLAGS:
//#cgo LDFLAGS:

import "C"
