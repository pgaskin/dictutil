// Package marisa provides a simplified self-contained CGO wrapper for
// marisa-trie (https://github.com/s-yata/marisa-trie).
package marisa

//go:generate go run -tags libmarisa_generate libmarisa_generate.go

//#cgo CPPFLAGS: -Wall
//#cgo LDFLAGS:
//#include <stddef.h>
//#include <stdlib.h>
//const char* marisa_read_all(int iid, char ***out_wd, size_t *out_wd_sz);
//const char* marisa_write_all(int iid, const char** wd, size_t wd_sz);
import "C"

import (
	"errors"
	"io"
	"unsafe"
)

func ReadAll(r io.Reader) ([]string, error) {
	iid := iopPut(r)
	var out_wd **C.char
	var out_wd_sz C.size_t
	err := C.marisa_read_all(
		(C.int)(iid),
		(***C.char)(unsafe.Pointer(&out_wd)),
		(*C.size_t)(unsafe.Pointer(&out_wd_sz)),
	)
	iopDel(iid)
	return gostrs(out_wd, out_wd_sz), goerr(err)
}

func WriteAll(w io.Writer, wd []string) error {
	iid := iopPut(w)
	wd_ptr, wd_sz, wd_free := cstrs(wd)
	err := C.marisa_write_all(
		(C.int)(iid),
		(**C.char)(wd_ptr),
		(C.size_t)(wd_sz),
	)
	wd_free()
	iopDel(iid)
	return goerr(err)
}

func goerr(p *C.char) (err error) {
	if p != nil {
		err = errors.New(C.GoString(p))
		C.free(unsafe.Pointer(p))
	}
	return
}

func gostrs(p **C.char, n C.size_t) (s []string) {
	if p != nil {
		s = make([]string, int(n))
		for i, v := range (*[1 << 28]*C.char)(unsafe.Pointer(p))[:int(n):int(n)] {
			s[i] = C.GoString(v)
			C.free(unsafe.Pointer(v))
		}
		C.free(unsafe.Pointer(p))
	}
	return
}

func cstrs(s []string) (p **C.char, n C.size_t, free func()) {
	n = (C.size_t)(len(s))
	if len(s) == 0 {
		free = func() {}
		return
	}
	c := make([]*C.char, len(s))
	for i, v := range s {
		c[i] = C.CString(v)
	}
	p = (**C.char)(unsafe.Pointer(&c[0]))
	free = func() {
		for _, v := range c {
			C.free(unsafe.Pointer(v))
		}
	}
	return
}
