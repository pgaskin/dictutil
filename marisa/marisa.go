// Package marisa provides a simplified self-contained CGO wrapper for
// marisa-trie (https://github.com/s-yata/marisa-trie).
package marisa

//go:generate go run -tags libmarisa_generate libmarisa_generate.go

//#cgo CPPFLAGS: -Wall
//#cgo LDFLAGS:
//#include <stdlib.h>
//#include "marisa.h"
import "C"

import (
	"errors"
	"io"
	"unsafe"
)

func ReadAll(r io.Reader) ([]string, error) {
	iid := iopPut(r)
	defer iopDel(iid)

	var out_wd **C.char
	var out_wd_sz C.size_t
	var out_err *C.char

	C.marisa_read_all(
		(C.int)(iid),
		(***C.char)(unsafe.Pointer(&out_wd)),
		(*C.size_t)(unsafe.Pointer(&out_wd_sz)),
		(**C.char)(unsafe.Pointer(&out_err)),
	)

	if out_wd != nil {
		defer C.marisa_wd_free(out_wd, out_wd_sz)
	}
	if out_err != nil {
		defer C.free(unsafe.Pointer(out_err))
		return nil, errors.New(C.GoString(out_err))
	}

	wd := make([]string, int(out_wd_sz))
	for i, w := range (*[1 << 28]*C.char)(unsafe.Pointer(out_wd))[:int(out_wd_sz):int(out_wd_sz)] {
		wd[i] = C.GoString(w)
	}
	return wd, nil
}

func WriteAll(w io.Writer, wd []string) error {
	iid := iopPut(w)
	defer iopDel(iid)

	in_wd := make([]*C.char, len(wd))
	for i, w := range wd {
		in_wd[i] = C.CString(w)
	}
	defer func() {
		for _, p := range in_wd {
			C.free(unsafe.Pointer(p))
		}
	}()

	var out_err *C.char

	var in_wd_ptr unsafe.Pointer
	if len(in_wd) != 0 {
		in_wd_ptr = unsafe.Pointer(&in_wd[0])
	}
	C.marisa_write_all(
		(C.int)(iid),
		(**C.char)(in_wd_ptr),
		(C.size_t)(len(in_wd)),
		(**C.char)(unsafe.Pointer(&out_err)),
	)

	if out_err != nil {
		defer C.free(unsafe.Pointer(out_err))
		return errors.New(C.GoString(out_err))
	}

	return nil
}
