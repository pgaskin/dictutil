//+build cgo

package kobodict

import (
	"io"

	"github.com/geek1011/dictutil/marisa"
)

// MarisaCGO uses the built-in CGO-based libmarisa bindings.
type MarisaCGO struct{}

func init() {
	setMarisa(new(MarisaCGO))
}

func (*MarisaCGO) ReadAll(r io.Reader) (wd []string, err error) {
	return marisa.ReadAll(r)
}

func (*MarisaCGO) WriteAll(w io.Writer, wd []string) (err error) {
	return marisa.WriteAll(w, wd)
}
