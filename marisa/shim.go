package marisa

import "C"

import (
	"fmt"
	"io"
	"sync"
	"unsafe"
)

// shim.go and shim.h (plus _cgo_export.h implicitly), implement a shim to
// access Go I/O interfaces efficiently, concurrently, and safely from C/C++
// code. The semantics are mostly the same as the Go ones themselves. Note that
// if any C strings are returned by the Go side, they must be freed on the C
// side.

// https://golang.org/issue/13656#issuecomment-253600758
// https://golang.org/cmd/cgo/#hdr-C_references_to_Go
// https://stackoverflow.com/a/49879469

var (
	iopMu sync.RWMutex         // for controlling access to the slice header (i.e. https://stackoverflow.com/a/49879469)
	iop   = []interface{}{nil} // the 0th element is reserved to prevent mistakes
)

// iopPut adds the io.Reader and/or io.Writer, and returns its new iid. The iid
// will be valid until iopDel is called, but will never be reused.
func iopPut(rw interface{}) int {
	switch rw.(type) {
	case io.Reader, io.Writer:
		iopMu.Lock()
		iop = append(iop, rw)
		iid := len(iop) - 1
		iopMu.Unlock()
		return iid
	default:
		panic("not a reader, writer, or both")
	}
}

// iopGet gets the interface referenced by iid. It will panic if iid has never
// been issued by iopPut, and will return nil if it has been deleted by iopDel.
func iopGet(iid int) interface{} {
	iopMu.RLock()
	if iid <= 0 || iid >= len(iop) {
		panic("invalid iid")
	}
	r := iop[iid]
	iopMu.RUnlock()
	return r
}

// iopDel sets the interface referenced by iid to nil to prevent future usage.
// It will panic if iid has never been issued by iopPut.
func iopDel(iid int) {
	iopMu.RLock()
	if iid <= 0 || iid >= len(iop) {
		panic("invalid iid")
	}
	iop[iid] = nil
	iopMu.RUnlock()
}

//export go_iop_read
func go_iop_read(iid C.int, buf *C.char, buf_n C.int, out_err **C.char) C.int {
	switch i := iopGet(int(iid)).(type) {
	case io.Reader:
		n, err := i.Read((*[1 << 28]byte)(unsafe.Pointer(buf))[:int(buf_n):int(buf_n)])
		if err == io.EOF {
			if n == 0 {
				return C.int(-1)
			}
		} else if err != nil {
			*out_err = C.CString(fmt.Sprintf("iop_read: read up to %d bytes from iid %d: %v", buf_n, int(iid), err))
		}
		return C.int(n)
	case nil:
		*out_err = C.CString(fmt.Sprintf("iop_read: iid %d has been deleted", int(iid)))
		return C.int(0)
	default:
		*out_err = C.CString(fmt.Sprintf("iop_read: iid %d is a %T, not an io.Reader", int(iid), i))
		return C.int(0)
	}
}

//export go_iop_write
func go_iop_write(iid C.int, buf *C.char, buf_n C.int, out_err **C.char) C.int {
	switch i := iopGet(int(iid)).(type) {
	case io.Writer:
		n, err := i.Write((*[1 << 28]byte)(unsafe.Pointer(buf))[:int(buf_n):int(buf_n)])
		if err == io.EOF {
			if n == 0 {
				return C.int(-1)
			}
		} else if err != nil {
			*out_err = C.CString(fmt.Sprintf("iop_write: write up to %d bytes to iid %d: %v", buf_n, int(iid), err))
		}
		return C.int(n)
	case nil:
		*out_err = C.CString(fmt.Sprintf("iop_write: iid %d has been deleted", int(iid)))
		return C.int(0)
	default:
		*out_err = C.CString(fmt.Sprintf("iop_write: iid %d is a %T, not an io.Writer", int(iid), i))
		return C.int(0)
	}
}
