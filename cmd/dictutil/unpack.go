package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/geek1011/dictutil/kobodict"
	"github.com/spf13/pflag"
)

func init() {
	commands = append(commands, &command{Name: "unpack", Short: "u", Description: "Unpack a dictzip file", Main: unpackMain})
}

func unpackMain(args []string, fs *pflag.FlagSet) int {
	fs.SortFlags = false
	output := fs.StringP("output", "o", "", "The output directory (must not exist) (default: the basename of the input without the extension)")
	crypt := fs.StringP("crypt", "c", "", "Decrypt the dictzip (if needed) using the specified encryption method (format: method:keyhex)")
	help := fs.BoolP("help", "h", false, "Show this help text")
	fs.Parse(args[1:])

	if *help || fs.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] dictzip\n\nOptions:\n%s", args[0], fs.FlagUsages())
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

	ofn := *output
	if ofn == "" {
		ofn = strings.TrimSuffix(filepath.Base(fn), filepath.Ext(fn))
	}

	fmt.Printf("Opening input dictzip.\n")
	f, err := os.Open(fn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: open input file %#v: %v.\n", fn, err)
		return 1
	}
	defer f.Close()

	s, err := f.Stat()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: stat input file %#v: %v.\n", fn, err)
		return 1
	}

	fmt.Printf("Parsing dictzip.\n")
	dr, err := kobodict.NewReader(f, s.Size())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: parse input file %#v: %v.\n", fn, err)
		return 1
	}
	dr.SetDecrypter(c)

	fmt.Printf("Unpacking dictzip.\n")
	if err := kobodict.Unpack(dr, ofn); err != nil {
		fmt.Fprintf(os.Stderr, "Error: unpack input file %#v to %#v: %v.\n", fn, ofn, err)
		return 1
	}

	fmt.Printf("Successfully unpacked dictzip %#v to dictdir %#v.\n", fn, ofn)
	return 0
}
