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
	"fmt"
	"io"
	"io/ioutil"
	"unsafe"
)

func ReadAll(r io.Reader) ([]string, error) {
	in_buf, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	var out_wd **C.char
	var out_wd_sz C.size_t
	var out_err *C.char

	var in_buf_ptr unsafe.Pointer
	if len(in_buf) != 0 {
		in_buf_ptr = unsafe.Pointer(&in_buf[0])
	}
	C.marisa_read_all(
		(*C.char)(in_buf_ptr),
		(C.size_t)(len(in_buf)),
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
	in_wd := make([]*C.char, len(wd))
	for i, w := range wd {
		in_wd[i] = C.CString(w)
	}
	defer func() {
		for _, p := range in_wd {
			C.free(unsafe.Pointer(p))
		}
	}()

	var out_buf *C.char
	var out_buf_sz C.size_t
	var out_err *C.char

	var in_wd_ptr unsafe.Pointer
	if len(in_wd) != 0 {
		in_wd_ptr = unsafe.Pointer(&in_wd[0])
	}
	C.marisa_write_all(
		(**C.char)(in_wd_ptr),
		(C.size_t)(len(in_wd)),
		(**C.char)(unsafe.Pointer(&out_buf)),
		(*C.size_t)(unsafe.Pointer(&out_buf_sz)),
		(**C.char)(unsafe.Pointer(&out_err)),
	)

	if out_buf != nil {
		defer C.free(unsafe.Pointer(out_buf))
	}
	if out_err != nil {
		defer C.free(unsafe.Pointer(out_err))
		return errors.New(C.GoString(out_err))
	}

	_, err := w.Write((*[1 << 28]byte)(unsafe.Pointer(out_buf))[:int(out_buf_sz):int(out_buf_sz)])
	return err
}

func marisa_go_test_error_helper(at int) {
	if at != -1 {
		fmt.Printf("Enabled marisa_go_test_error_helper to throw when input length is exactly %d\n", at)
	}
	C.marisa_go_test_error_helper(C.int(at), C.int(-1))
}
