//+build libmarisa_generate

package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
)

func main() {
	url := "https://github.com/s-yata/marisa-trie/archive/970b20c141f11d9d7572a6bb8d0488f2e0520e22.tar.gz"
	version := "970b20c"

	if files, err := tarball(url); err != nil {
		fmt.Fprintf(os.Stderr, "Error: download tarball %#v: %v\n", url, err)
		os.Exit(1)
		return
	} else if err := func() error {
		if mr, err := libmarisa(files, version); err != nil {
			return err
		} else if mf, err := os.OpenFile("libmarisa.cc", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644); err != nil {
			return err
		} else if _, err := io.Copy(mf, mr); err != nil {
			mf.Close()
			return err
		} else {
			return mf.Close()
		}
	}(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: generate libmarisa.cc: %v\n", err)
		os.Exit(1)
		return
	} else if err := func() error {
		if mr, err := hmarisa(files, version); err != nil {
			return err
		} else if mf, err := os.OpenFile("marisa.h", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644); err != nil {
			return err
		} else if _, err := io.Copy(mf, mr); err != nil {
			mf.Close()
			return err
		} else {
			return mf.Close()
		}
	}(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: generate marisa.h: %v\n", err)
		os.Exit(1)
		return
	}
}
func hmarisa(files map[string][]byte, version string) (io.Reader, error) {
	marisaH, err := resolve(files, []string{
		"include/marisa.h",
	}, "include", "lib")
	if err != nil {
		return nil, err
	}

	fmt.Printf("Generating marisa.h\n")
	return io.MultiReader(
		// A custom header.
		strings.NewReader("// AUTOMATICALLY GENERATED, DO NOT EDIT!\n"),
		strings.NewReader("// merged from marisa-trie "+version+".\n"),
		// Include the license info.
		bytes.NewReader([]byte{'\n', '/', '/', ' '}),
		bytes.NewReader(bytes.ReplaceAll(files["COPYING.md"], []byte{'\n'}, []byte{'\n', '/', '/', ' '})),
		bytes.NewReader([]byte{'\n', '\n'}),
		// Include the header.
		bytes.NewReader(marisaH),
	), nil
}

func libmarisa(files map[string][]byte, version string) (io.Reader, error) {
	marisaGrimoireIOLib, err := resolve(files, []string{
		"lib/marisa/grimoire/io/mapper.cc",
		"lib/marisa/grimoire/io/reader.cc",
		"lib/marisa/grimoire/io/writer.cc",
	}, "include", "lib")
	if err != nil {
		return nil, err
	}

	marisaGrimoireTrieLib, err := resolve(files, []string{
		"lib/marisa/grimoire/trie/tail.cc",
		"lib/marisa/grimoire/trie/louds-trie.cc",
	}, "include", "lib")
	if err != nil {
		return nil, err
	}

	marisaGrimoireVectorLib, err := resolve(files, []string{
		"lib/marisa/grimoire/vector/bit-vector.cc",
	}, "include", "lib")
	if err != nil {
		return nil, err
	}

	marisaLib, err := resolve(files, []string{
		"lib/marisa/agent.cc",
		"lib/marisa/keyset.cc",
		"lib/marisa/trie.cc",
	}, "include", "lib")
	if err != nil {
		return nil, err
	}

	fmt.Printf("Generating libmarisa.cc\n")
	return io.MultiReader(
		// A custom header.
		strings.NewReader("// AUTOMATICALLY GENERATED, DO NOT EDIT!\n"),
		strings.NewReader("// merged from marisa-trie "+version+".\n"),
		// Include the license info.
		bytes.NewReader([]byte{'\n', '/', '/', ' '}),
		bytes.NewReader(bytes.ReplaceAll(files["COPYING.md"], []byte{'\n'}, []byte{'\n', '/', '/', ' '})),
		bytes.NewReader([]byte{'\n', '\n'}),
		// Include the warnings from the Makefile.am CXXFLAGS.
		// - Note that Clang also recognizes the GCC pragmas.
		strings.NewReader("#pragma GCC diagnostic warning \"-Wall\"\n"),
		strings.NewReader("#pragma GCC diagnostic warning \"-Weffc++\"\n"),
		strings.NewReader("#pragma GCC diagnostic warning \"-Wextra\"\n"),
		strings.NewReader("#pragma GCC diagnostic warning \"-Wconversion\"\n"),
		// Silence a warning.
		strings.NewReader("#pragma GCC diagnostic ignored \"-Wimplicit-fallthrough=\"\n"),
		// Include the libs themselves.
		bytes.NewReader(marisaGrimoireIOLib),
		bytes.NewReader(marisaGrimoireTrieLib),
		bytes.NewReader(marisaGrimoireVectorLib),
		bytes.NewReader(marisaLib),
		// Show info about the generated file.
		strings.NewReader("#line 1 \"libmarisa_generate.go\"\n"),
		strings.NewReader("#pragma GCC warning \"Using generated built-in marisa-trie "+version+".\"\n"),
	), nil
}

func tarball(url string) (map[string][]byte, error) {
	fmt.Printf("Downloading tarball from %s\n", url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	zr, err := gzip.NewReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var pfx string
	files := map[string][]byte{}

	tr := tar.NewReader(zr)
	for {
		fh, err := tr.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		if fh.Name == "pax_global_header" || fh.FileInfo().IsDir() {
			continue
		}

		if pfx == "" {
			if strings.HasPrefix(fh.Name, "./") {
				pfx = "./" + strings.Split(fh.Name, "/")[1] + "/"
			} else {
				pfx = strings.Split(fh.Name, "/")[0] + "/"
			}
		}

		if !strings.HasPrefix(fh.Name, pfx) {
			return nil, fmt.Errorf("extract file %#v: doesn't have common prefix %#v", fh.Name, pfx)
		}

		buf, err := ioutil.ReadAll(tr)
		if err != nil {
			return nil, fmt.Errorf("extract file %#v: %w", fh.Name, err)
		}

		fn := strings.TrimPrefix(fh.Name, pfx)
		files[fn] = buf

		fmt.Printf("  [D] %s\n", fn) // downloaded
	}

	return files, nil
}

func resolve(files map[string][]byte, filenames []string, includePath ...string) (resolvedFile []byte, err error) {
	fmt.Printf("Resolving C* source files %s (against:%s) (I = included, S = preserved because not found, R = skipped because already included)\n", filenames, includePath)

	var resolveFn func(indent string, files map[string][]byte, filename string, buf []byte, done []string, includePath []string) (resolvedFile []byte, err error)
	resolveFn = func(indent string, files map[string][]byte, filename string, buf []byte, done []string, includePath []string) (resolvedFile []byte, err error) {
		defer func() {
			if rerr := recover(); rerr != nil {
				resolvedFile, err = nil, rerr.(error)
			}
		}()

		resolvedFile = regexp.MustCompile(`(?m)^\s*#\s*include\s+["'<][^"'>]+["'>]$`).ReplaceAllFunc(buf, func(importBuf []byte) []byte {
			fn := string(regexp.MustCompile(`["'<]([^"'>]+)["'>]`).FindSubmatch(importBuf)[1])

			for _, ip := range includePath {
				ifn := path.Join(ip, fn)
				for _, dfn := range done {
					if m, _ := path.Match(dfn, ifn); m {
						fmt.Printf("%s[R] %s\n", indent, fn) // already included
						return nil
					}
				}

				ibuf, ok := files[ifn]
				if ok {
					fmt.Printf("%s[I] %s => %s\n", indent, fn, ifn) // include
					ibuf, err := resolveFn(indent+"    ", files, ifn, ibuf, append(done, ifn), append(includePath, path.Dir(ifn)))
					if err != nil {
						panic(fmt.Errorf("resolve %#v: %w", ifn, err))
					}
					return append(append([]byte{'\n', '\n'}, ibuf...), '\n', '\n')
				}
			}

			fmt.Printf("%s[S] %s\n", indent, fn) // preserve
			return importBuf
		})

		return
	}

	for _, fn := range filenames {
		if buf, ok := files[fn]; !ok {
			return nil, fmt.Errorf("file %#v: not found", fn)
		} else if buf, err := resolveFn("  ", files, fn, buf, []string{fn}, append(includePath, path.Dir(fn))); err != nil {
			return nil, fmt.Errorf("file %v: %w", fn, err)
		} else {
			resolvedFile = append(resolvedFile, buf...)
			resolvedFile = append(resolvedFile, '\n', '\n')
		}
	}

	return resolvedFile, nil
}
