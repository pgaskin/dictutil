package kobodict

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

// Unpack is a helper function to unpack the contents of a Reader to a folder
// on-disk. The provided dir must be non-existent. Unpack will not close the
// reader.
func Unpack(r *Reader, dir string) error {
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		return fmt.Errorf("dir %#v already exists", dir)
	}
	if err := os.Mkdir(dir, 0755); err != nil {
		return fmt.Errorf("create dir %#v: %w", dir, err)
	}
	for _, f := range r.File {
		if err := unpackFile(dir, f.Open, f.Name); err != nil {
			return fmt.Errorf("unpack file %#v: %w", f.Name, err)
		}
	}
	for _, f := range r.Dicthtml {
		if err := unpackFile(dir, f.Open, f.Name); err != nil {
			return fmt.Errorf("unpack dicthtml %#v (prefix: %s): %w", f.Name, f.Prefix, err)
		}
	}
	if err := ioutil.WriteFile(filepath.Join(dir, "words"), []byte(strings.Join(r.Word, "\n")), 0644); err != nil {
		return fmt.Errorf("write words file: %w", err)
	}
	return nil
}

func unpackFile(dir string, open func() (io.ReadCloser, error), name string) error {
	fr, err := open()
	if err != nil {
		return fmt.Errorf("read contents: %w", err)
	}
	defer fr.Close()

	fw, err := os.OpenFile(filepath.Join(dir, name), os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0644)
	if err != nil {
		return fmt.Errorf("create output file: %w", err)
	}
	defer fw.Close()

	if _, err := io.Copy(fw, fr); err != nil {
		return fmt.Errorf("write output file: %w", err)
	}

	if err := fw.Close(); err != nil {
		return fmt.Errorf("write output file: %w", err)
	}

	return nil
}

// Pack is a helper function to pack the contents a folder unpacked using Unpack
// into a Writer. It is assumed that the writer has not been used. The provided
// file will be overwritten if it exists and is a regular file, or created if it
// doesn't exist. Pack will not close the writer.
func Pack(w *Writer, dir string) error {
	if fi, err := os.Stat(filepath.Join(dir, "words")); os.IsNotExist(err) || (err == nil && fi.IsDir()) {
		return fmt.Errorf("dir %#v is not an unpacked dictzip (no words file)", dir)
	}

	fis, err := ioutil.ReadDir(dir) // note: this is sorted
	if err != nil {
		return fmt.Errorf("read dir %#v: %w", dir, err)
	}

	for _, fi := range fis {
		switch {
		case fi.IsDir():
			return fmt.Errorf("invalid dir %#v: dirs are not supported", fi.Name())
		case fi.Name() == "words":
			continue
		case strings.HasSuffix(fi.Name(), ".html"):
			if err := func() error {
				fr, err := os.OpenFile(filepath.Join(dir, fi.Name()), os.O_RDONLY, 0)
				if err != nil {
					return fmt.Errorf("open file: %w", err)
				}
				defer fr.Close()

				tmp := make([]byte, 2)
				if _, err := fr.Read(tmp); err != nil {
					return fmt.Errorf("read file: %w", err)
				} else if tmp[0] == 0x1F && tmp[1] == 0x8B {
					return fmt.Errorf("invalid unpacked dicthtml file: already compressed")
				} else if _, err := fr.Seek(0, os.SEEK_SET); err != nil {
					return fmt.Errorf("read file: %w", err)
				}

				fw, err := w.CreateDicthtml(strings.TrimSuffix(fi.Name(), ".html"))
				if err != nil {
					return fmt.Errorf("create dictzip entry: %w", err)
				}

				if _, err := io.Copy(fw, fr); err != nil {
					return fmt.Errorf("write file: %w", err)
				}

				return nil
			}(); err != nil {
				return fmt.Errorf("add dicthtml %#v: %w", fi.Name(), err)
			}
		default:
			if err := func() error {
				fr, err := os.OpenFile(filepath.Join(dir, fi.Name()), os.O_RDONLY, 0)
				if err != nil {
					return fmt.Errorf("open file: %w", err)
				}
				defer fr.Close()

				fw, err := w.CreateFile(strings.TrimSuffix(fi.Name(), ".html"))
				if err != nil {
					return fmt.Errorf("create dictzip entry: %w", err)
				}

				if _, err := io.Copy(fw, fr); err != nil {
					return fmt.Errorf("write file: %w", err)
				}

				return nil
			}(); err != nil {
				return fmt.Errorf("add file %#v: %w", fi.Name(), err)
			}
		}
	}

	if err := func() error {
		fr, err := os.OpenFile(filepath.Join(dir, "words"), os.O_RDONLY, 0)
		if err != nil {
			return fmt.Errorf("open words file: %w", err)
		}
		defer fr.Close()

		sc := bufio.NewScanner(fr)
		for sc.Scan() {
			if !utf8.Valid(sc.Bytes()) {
				return fmt.Errorf("invalid word: %#v", sc.Text())
			}
			if word := strings.TrimSpace(sc.Text()); len(word) != 0 {
				if err := w.AddWord(word); err != nil {
					return fmt.Errorf("add word %#v: %s", word, err)
				}
			}
		}
		if sc.Err() != nil {
			return fmt.Errorf("read words file: %w", err)
		}

		return nil
	}(); err != nil {
		return fmt.Errorf("add words index: %w", err)
	}

	return nil
}
