package main

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/geek1011/dictutil/kobodict"
	"github.com/spf13/pflag"
)

func init() {
	commands = append(commands, &command{Name: "pack", Short: "p", Description: "Pack a dictzip file", Main: packMain})
}

func packMain(args []string, fs *pflag.FlagSet) int {
	fs.SortFlags = false
	output := fs.StringP("output", "o", "dicthtml.zip", "The output dictzip filename (will be overwritten if it exists)")
	crypt := fs.StringP("crypt", "c", "", "Encrypt the dictzip using the specified encryption method (format: method:keyhex)")
	help := fs.BoolP("help", "h", false, "Show this help text")
	fs.Parse(args[1:])

	if *help || fs.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] dictdir\n\nOptions:\n%s", args[0], fs.FlagUsages())
		return 0
	}

	var c kobodict.Crypter
	if *crypt != "" {
		if spl := strings.SplitN(*crypt, ":", 2); len(spl) < 2 {
			fmt.Fprintf(os.Stderr, "Error: invalid format for --encrypt: no ':' found.\n")
			return 2
		} else if key, err := hex.DecodeString(spl[1]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: invalid format for --encrypt: decode hex: %v.\n", err)
			return 2
		} else if enc, err := kobodict.NewCrypter(spl[0], key); err != nil {
			fmt.Fprintf(os.Stderr, "Error: invalid format for --encrypt: initialize encrypter: %v.\n", err)
			return 2
		} else {
			c = enc
		}
	}

	fn, err := filepath.Abs(fs.Args()[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: resolve input path %#v: %v.\n", fs.Args()[0], err)
		return 2
	}

	ofn, err := filepath.Abs(*output)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: resolve output path %#v: %v.\n", *output, err)
		return 2
	}

	if fi, err := os.Stat(fn); err != nil {
		fmt.Fprintf(os.Stderr, "Error: inaccessible input dir %#v: %v.\n", fn, err)
		return 2
	} else if !fi.IsDir() {
		fmt.Fprintf(os.Stderr, "Error: input %#v is not a dir.\n", fn)
		return 2
	}

	fmt.Printf("Creating output temp file\n")
	f, err := ioutil.TempFile(filepath.Dir(ofn), "tmp_dicthtml.*.zip")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: create output temp file: %v.\n", err)
		return 2
	}
	defer os.Remove(f.Name())
	defer f.Close()

	fmt.Printf("Packing dictzip.\n")
	dw := kobodict.NewWriter(f)
	defer dw.Close()

	dw.SetEncrypter(c)

	if err := kobodict.Pack(dw, fn); err != nil {
		fmt.Fprintf(os.Stderr, "Error: pack input dir %#v to %#v: %v.\n", fn, ofn, err)
		return 1
	}

	if err := dw.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: pack input dir %#v to %#v: %v.\n", fn, ofn, err)
		return 1
	}

	fmt.Printf("Renaming output file.\n")
	if err := f.Chmod(0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error: rename output file: %v.\n", err)
		return 2
	}
	if err := f.Sync(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: rename output file: %v.\n", err)
		return 2
	}
	if err := f.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: rename output file: %v.\n", err)
		return 2
	}
	if err := os.Rename(f.Name(), ofn); err != nil { // this will replace existing files properly on Go1.5+
		fmt.Fprintf(os.Stderr, "Error: rename output file: %v.\n", err)
		return 2
	}

	fmt.Printf("Successfully packed dictdir %#v to dictzip %#v.\n", fn, ofn)
	return 0
}
