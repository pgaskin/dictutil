package kobodict

import (
	"archive/zip"
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
	words  map[string]struct{} // doesn't take up space for values
	closed bool
	last   io.WriteCloser
}

// NewWriter creates a dictzip writer writing to w.
func NewWriter(w io.Writer) *Writer {
	return &Writer{
		z:     zip.NewWriter(w),
		words: map[string]struct{}{},
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
	fw, err := w.z.Create(prefix + ".html")
	if err != nil {
		return nil, fmt.Errorf("create zip entry: %w", err)
	}
	w.last = gzip.NewWriter(fw)
	return w.last, nil
}

// CreateFile adds a raw file with the specified name. Note that Kobo only
// supports GIF and JPEG files starting with the "GIF" and "JFIF" magic, and the
// treatment of other files is undefined. In addition, subdirectories are not
// supported.
func (w *Writer) CreateFile(filename string) (io.Writer, error) {
	if strings.Contains(filename, "/") {
		return nil, fmt.Errorf("invalid filename: must not contain slashes")
	}
	if w.last != nil {
		if err := w.last.Close(); err != nil {
			return nil, fmt.Errorf("close last file writer: %w", err)
		}
		w.last = nil
	}
	return w.z.Create(filename)
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
