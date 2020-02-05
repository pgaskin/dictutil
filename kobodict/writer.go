package kobodict

import (
	"archive/zip"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/geek1011/dictutil/marisa"
)

// Writer creates dictzips. It does not do any validation; it only does what it
// is told. It is up to the user to ensure the input is valid.
type Writer struct {
	z      *zip.Writer
	e      Encrypter
	words  map[string]struct{} // doesn't take up space for values
	used   map[string]struct{}
	closed bool
	last   io.WriteCloser
}

// Encrypter encrypts dicthtml files.
type Encrypter interface {
	// Encrypt encrypts the provided bytes.
	Encrypt([]byte) ([]byte, error)
}

// NewWriter creates a dictzip writer writing to w.
func NewWriter(w io.Writer) *Writer {
	return &Writer{
		z:     zip.NewWriter(w),
		words: map[string]struct{}{},
		used:  map[string]struct{}{},
	}
}

// AddWord normalizes and adds a word to the index. If the word has already been
// added, it does nothing.
func (w *Writer) AddWord(word string) error {
	if w.closed {
		return fmt.Errorf("write to closed writer")
	}
	w.words[strings.TrimSpace(word)] = struct{}{} // index words aren't normalized except for trimming spaces
	return nil
}

// CreateDicthtml adds a dicthtml file for the specified prefix and returns a
// writer which is valid until the next file is created.
func (w *Writer) CreateDicthtml(prefix string) (io.Writer, error) {
	if strings.Contains(prefix, "/") {
		return nil, fmt.Errorf("invalid prefix: must not contain slashes")
	}
	if w.closed {
		return nil, fmt.Errorf("writer already closed")
	}
	if w.last != nil {
		if err := w.last.Close(); err != nil {
			return nil, fmt.Errorf("close last file writer: %w", err)
		}
		w.last = nil
	}

	filename := prefix + ".html"
	if _, ok := w.used[filename]; ok {
		return nil, fmt.Errorf("file %#v already exists in dictzip", filename)
	}

	fw, err := w.z.Create(filename)
	if err != nil {
		return nil, fmt.Errorf("create zip entry: %w", err)
	}

	if w.e != nil {
		ew := newEncryptWriter(w.e, fw)
		zw := gzip.NewWriter(ew)

		w.last = &funcWriteCloser{
			Writer: zw,
			Closer: func() error {
				if err := zw.Close(); err != nil {
					return err
				}
				return ew.Close()
			},
		}
	} else {
		w.last = gzip.NewWriter(fw)
	}

	w.used[filename] = struct{}{}
	return w.last, nil
}

// CreateFile adds a raw file with the specified name. Note that Kobo only
// supports GIF and JPEG files starting with the "GIF" and "JFIF" magic, and the
// treatment of other files is undefined. In addition, subdirectories are not
// supported. The behaviour is undefined if a dicthtml file is added this way.
func (w *Writer) CreateFile(filename string) (io.Writer, error) {
	if strings.Contains(filename, "/") || strings.Contains(filename, "\\") {
		return nil, fmt.Errorf("invalid filename: must not contain slashes")
	} else if strings.Contains(filename, "words") {
		return nil, fmt.Errorf("invalid filename: must not be 'words'")
	} else if _, ok := w.used[filename]; ok {
		return nil, fmt.Errorf("file %#v already exists in dictzip", filename)
	}
	if w.last != nil {
		if err := w.last.Close(); err != nil {
			return nil, fmt.Errorf("close last file writer: %w", err)
		}
		w.last = nil
	}

	fw, err := w.z.Create(filename)
	if err != nil {
		return nil, fmt.Errorf("create zip entry: %w", err)
	}

	w.last = &funcWriteCloser{
		Writer: fw,
		Closer: nil,
	}
	w.used[filename] = struct{}{}
	return w.last, nil
}

// Close writes the marisa index and the zip footer. The error should not be
// ignored. It does not close the underlying writer.
func (w *Writer) Close() error {
	if w.closed {
		return fmt.Errorf("writer already closed")
	}
	if w.last != nil {
		if err := w.last.Close(); err != nil {
			return fmt.Errorf("close last file writer: %w", err)
		}
		w.last = nil
	}

	if buf, err := w.marisaBytes(); err != nil {
		return fmt.Errorf("generate index: %w", err)
	} else if fw, err := w.z.Create("words"); err != nil {
		return fmt.Errorf("create index zip entry: %w", err)
	} else if _, err := fw.Write(buf); err != nil {
		return fmt.Errorf("write index: %w", err)
	}

	if err := w.z.Close(); err != nil {
		return fmt.Errorf("close zip: %w", err)
	}
	return nil
}

// SetEncrypter sets the Encrypter used to encrypt dicthtml files. This must be
// will only apply to dicthtml files added after the encrypter is set.
func (r *Writer) SetEncrypter(e Encrypter) {
	r.e = e
}

func (w *Writer) marisaBytes() (buf []byte, err error) {
	defer func() {
		if err := recover(); err != nil {
			buf = nil
			err = fmt.Errorf("marisa: %v", err)
		}
	}()

	var words []string
	for word := range w.words {
		words = append(words, word)
	}
	sort.Strings(words) // for deterministic output

	ks := marisa.NewKeyset()
	for _, word := range words {
		ks.PushBackString(word)
	}

	td, err := ioutil.TempDir("", "marisa")
	if err != nil {
		return nil, fmt.Errorf("create temp dir: %w", err)
	}
	defer os.RemoveAll(td)

	trie := marisa.NewTrie()
	trie.Build(ks)
	trie.Save(filepath.Join(td, "words"))

	buf, err = ioutil.ReadFile(filepath.Join(td, "words"))
	if err != nil {
		return nil, fmt.Errorf("read marisa output: %w", err)
	}
	return buf, err
}

type encryptWriter struct {
	e Encrypter
	w io.Writer
	b *bytes.Buffer
	c bool
}

func newEncryptWriter(e Encrypter, w io.Writer) io.WriteCloser {
	return &encryptWriter{
		e: e,
		w: w,
		b: bytes.NewBuffer(nil),
		c: false,
	}
}

func (e encryptWriter) Write(buf []byte) (n int, err error) {
	if e.c {
		return 0, fmt.Errorf("write to closed writer")
	}
	return e.b.Write(buf)
}

// Close encrypts and writes the buffer to the underlying writer. The error
// should be checked.
func (e encryptWriter) Close() error {
	if e.c {
		return fmt.Errorf("writer already closed")
	}
	if buf, err := e.e.Encrypt(e.b.Bytes()); err != nil {
		return fmt.Errorf("encrypt bytes: %w", err)
	} else if _, err := e.w.Write(buf); err != nil {
		return fmt.Errorf("write encrypted bytes: %w", err)
	}
	return nil
}

type funcWriteCloser struct {
	io.Writer
	Closer func() error
}

func (f *funcWriteCloser) Close() error {
	if f.Closer != nil {
		return f.Closer()
	}
	return nil
}
