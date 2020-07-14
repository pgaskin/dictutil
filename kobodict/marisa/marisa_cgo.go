//+build cgo

package marisa

import (
	"io"

	"github.com/pgaskin/dictutil/marisa"
)

func (*platform) ReadAll(r io.Reader) (wd []string, err error) {
	return marisa.ReadAll(r)
}

func (*platform) WriteAll(w io.Writer, wd []string) (err error) {
	return marisa.WriteAll(w, wd)
}
