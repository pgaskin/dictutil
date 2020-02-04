package main

import (
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/geek1011/dictutil/dictgen"
	"github.com/geek1011/dictutil/kobodict"
	"github.com/spf13/pflag"
)

func main() {
	pflag.CommandLine.SortFlags = false
	output := pflag.StringP("output", "o", "dicthtml.zip", "The output filename (will be overwritten if it exists) (- is stdout)")
	encrypt := pflag.StringP("encrypt", "e", "", "Encrypt the dictzip using the specified encryption method (format: method:keyhex)")
	// TODO: image-dir
	help := pflag.BoolP("help", "h", false, "Show this help text")
	pflag.Parse()

	if *help || pflag.NArg() == 0 {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] dictfile...\n\nOptions:\n%s\nIf multiple dictfiles (*.df) are provided, they will be merged (duplicate entries are fine; they will be shown in sequential order). To read from stdin, use - as the filename.\n\nThe dictfile format:\n  TODO: short doc text\n", os.Args[0], pflag.CommandLine.FlagUsages())
		os.Exit(0)
		return
	}

	var e kobodict.Crypter
	if *encrypt != "" {
		if spl := strings.SplitN(*encrypt, ":", 2); len(spl) < 2 {
			fmt.Fprintf(os.Stderr, "Error: invalid format for --encrypt: no ':' found.\n")
			os.Exit(2)
			return
		} else if key, err := hex.DecodeString(spl[1]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: invalid format for --encrypt: decode hex: %v.\n", err)
			os.Exit(2)
			return
		} else if enc, err := kobodict.NewCrypter(spl[0], key); err != nil {
			fmt.Fprintf(os.Stderr, "Error: invalid format for --encrypt: initialize encrypter: %v.\n", err)
			os.Exit(2)
			return
		} else {
			e = enc
		}
	}

	var tdf dictgen.DictFile

	fmt.Fprintf(os.Stderr, "Parsing dictfiles.\n")
	var seenStdin bool
	for _, fn := range pflag.Args() {
		if fn == "-" {
			if seenStdin {
				fmt.Fprintf(os.Stderr, "Error: stdin can only be specified once.\n")
				os.Exit(1)
				return
			}
			seenStdin = true
		}

		if err := func() error {
			var fr io.Reader
			if fn == "-" {
				fr = os.Stdin
			} else {
				f, err := os.OpenFile(fn, os.O_RDONLY, 0)
				if err != nil {
					return err
				}
				defer f.Close()
				fr = f
			}

			if df, err := dictgen.ParseDictFile(fr); err != nil {
				return err
			} else if err := df.Validate(); err != nil {
				return err
			} else {
				tdf = append(tdf, df...)
			}

			return nil
		}(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: input %#v: %v.\n", fn, err)
			os.Exit(1)
			return
		}
	}

	fmt.Fprintf(os.Stderr, "Opening output.\n")
	var f io.WriteCloser
	switch *output {
	case "-":
		f = os.Stdout
	default:
		if ff, err := os.OpenFile(*output, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error: create dictzip: %v\n", err)
			os.Exit(1)
			return
		} else {
			f = ff
		}
	}

	fmt.Fprintf(os.Stderr, "Generating dictzip.\n")
	dw := kobodict.NewWriter(f)
	dw.SetEncrypter(e)
	if err := tdf.WriteDictzip(dw); err != nil {
		f.Close()
		fmt.Fprintf(os.Stderr, "Error: write dictzip: %v\n", err)
		os.Exit(1)
		return
	} else if err := dw.Close(); err != nil {
		f.Close()
		fmt.Fprintf(os.Stderr, "Error: write dictzip: %v\n", err)
		os.Exit(1)
		return
	} else if err := f.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: write dictzip: %v\n", err)
		os.Exit(1)
		return
	}

	fmt.Fprintf(os.Stderr, "Successfully wrote %d entries from %d dictfile(s) to dictzip %s.\n", len(tdf), pflag.NArg(), *output)
	os.Exit(0)
}
