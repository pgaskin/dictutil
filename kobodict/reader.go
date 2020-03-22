package kobodict

import (
	"archive/zip"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
)

// Reader provides access to the contents of a dictzip file.
type Reader struct {
	Word     []string
	Dicthtml []*ReaderDicthtml
	File     []*ReaderFile
	z        *zip.Reader
	d        Decrypter
}

// ReaderDicthtml represents a dicthtml file from a Reader.
type ReaderDicthtml struct {
	Name   string
	Prefix string
	f      *zip.File
	r      *Reader
}

// ReaderDicthtml represents a raw file from a Reader (e.g. images).
type ReaderFile struct {
	Name string
	f    *zip.File
	r    *Reader
}

// Decrypter decrypts dicthtml files.
type Decrypter interface {
	// Decrypt decrypts the dicthtml bytes. It will only be called if the
	// dicthtml is not otherwise readable. An error should be returned if the
	// decryption itself encounters an error; the decryptor should not try to
	// judge if the resulting bytes are valid.
	Decrypt([]byte) ([]byte, error)
}

// NewReader returns a new dictzip reader which reads from r, with the given
// file size.
func NewReader(r io.ReaderAt, size int64) (*Reader, error) {
	zr, err := zip.NewReader(r, size)
	if err != nil {
		return nil, fmt.Errorf("open zip: %w", err)
	}

	kr := &Reader{
		z: zr,
	}

	var found bool
	for _, zf := range zr.File {
		if zf.Name == "words" {
			if fr, err := zf.Open(); err != nil {
				return nil, fmt.Errorf("open words index: %w", err)
			} else if Marisa == nil {
				return nil, fmt.Errorf("no marisa bindings found")
			} else if kr.Word, err = Marisa.ReadAll(fr); err != nil {
				return nil, fmt.Errorf("read words index: %w", err)
			}
			found = true
			break
		}
	}
	if !found {
		return nil, fmt.Errorf("not a dictzip: no words index found")
	}

	for _, f := range zr.File {
		switch {
		case !f.Mode().IsRegular():
			continue
		case f.Name == "words":
			continue
		case strings.Contains(f.Name, "/"):
			return nil, fmt.Errorf("read zip: illegal file %#v: contains slash (not in root dir)", f.Name)
		case strings.HasSuffix(f.Name, ".html"):
			kr.Dicthtml = append(kr.Dicthtml, &ReaderDicthtml{
				Name:   f.Name,
				Prefix: strings.TrimSuffix(f.Name, ".html"),
				f:      f,
				r:      kr,
			})
		default:
			kr.File = append(kr.File, &ReaderFile{
				Name: f.Name,
				f:    f,
				r:    kr,
			})
		}
	}

	return kr, nil
}

// SetDecrypter sets the Decrypter used to decrypt encrypted dicthtml files.
func (r *Reader) SetDecrypter(d Decrypter) {
	r.d = d
}

// Open returns an io.ReadCloser which reads the decoded dicthtml file. Multiple
// files can be read at once.
func (f *ReaderDicthtml) Open() (io.ReadCloser, error) {
	enc, err := func() (bool, error) {
		fr, err := f.f.Open()
		if err != nil {
			return false, fmt.Errorf("open zip entry: %v", err)
		}
		defer fr.Close()

		tmp := make([]byte, 2)
		if n, err := fr.Read(tmp); err != nil {
			return false, fmt.Errorf("read zip entry: %v", err)
		} else if n != len(tmp) {
			return false, fmt.Errorf("corrupt dicthtml: too short (%d)", n)
		}

		if tmp[0] == 0x1F && tmp[1] == 0x8B {
			return false, nil
		}

		if f.r.d == nil {
			return true, fmt.Errorf("corrupt or encrypted dicthtml: invalid header")
		}

		// maybe optimize this later?
		if buf, err := ioutil.ReadAll(io.MultiReader(bytes.NewReader(tmp), fr)); err != nil {
			return true, fmt.Errorf("read zip entry: %v", err)
		} else if dec, err := f.r.d.Decrypt(buf); err != nil {
			return true, fmt.Errorf("decrypt dicthtml: %v", err)
		} else if dec[0] != 0x1F || dec[1] != 0x8B {
			return true, fmt.Errorf("corrupt dicthtml or invalid encryption key: invalid header")
		}
		return true, nil
	}()
	if err != nil {
		return nil, err
	}

	fr, err := f.f.Open()
	if err != nil {
		return nil, fmt.Errorf("open zip entry: %v", err)
	}

	var dr io.Reader
	if enc {
		if buf, err := ioutil.ReadAll(fr); err != nil {
			return nil, fmt.Errorf("read zip entry: %v", err)
		} else if dec, err := f.r.d.Decrypt(buf); err != nil {
			return nil, fmt.Errorf("decrypt dicthtml: %v", err)
		} else if dec[0] != 0x1F || dec[1] != 0x8B {
			return nil, fmt.Errorf("corrupt dicthtml or invalid encryption key: invalid header")
		} else {
			dr = bytes.NewReader(dec)
		}
	} else {
		dr = fr
	}

	zr, err := gzip.NewReader(dr)
	if err != nil {
		return nil, fmt.Errorf("decompress dicthtml: %v", err)
	}

	return &funcReadCloser{
		Reader: zr,
		Closer: func() error {
			if err := zr.Close(); err != nil {
				fr.Close()
				return err
			}
			return fr.Close()
		},
	}, nil
}

// Open returns an io.ReadCloser which reads the contents of the file. Multiple
// files can be read at once.
func (f *ReaderFile) Open() (io.ReadCloser, error) {
	return f.f.Open()
}

type funcReadCloser struct {
	io.Reader
	Closer func() error
}

func (f *funcReadCloser) Close() error {
	if f.Closer != nil {
		return f.Closer()
	}
	return nil
}
