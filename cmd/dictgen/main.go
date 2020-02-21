package main

import (
	"encoding/hex"
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"strings"

	"github.com/geek1011/dictutil/dictgen"
	"github.com/geek1011/dictutil/kobodict"
	"github.com/spf13/pflag"
)

var version = "dev"

func main() {
	pflag.CommandLine.SortFlags = false
	output := pflag.StringP("output", "o", "dicthtml.zip", "The output filename (will be overwritten if it exists) (- is stdout)")
	crypt := pflag.StringP("crypt", "c", "", "Encrypt the dictzip using the specified encryption method (format: method:keyhex)")
	imageMethod := pflag.StringP("image-method", "I", "base64", "How to handle images (if an image path is relative, it is loaded from the current dir) (base64 - optimize and encode as base64, embed - add to dictzip, remove)")
	help := pflag.BoolP("help", "h", false, "Show this help text")
	pflag.Parse()

	if *help || pflag.NArg() == 0 {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] dictfile...\n\nVersion: dictgen %s\n\nOptions:\n%s\nIf multiple dictfiles (*.df) are provided, they will be merged (duplicate entries are fine; they will be shown in sequential order). To read from stdin, use - as the filename.\n\nSee https://pgaskin.net/dictutil/dictgen for more information about the dictfile format.\n", os.Args[0], version, pflag.CommandLine.FlagUsages())
		os.Exit(0)
		return
	}

	var e kobodict.Crypter
	if *crypt != "" {
		if spl := strings.SplitN(*crypt, ":", 2); len(spl) < 2 {
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

	var ih dictgen.ImageHandler
	switch *imageMethod {
	case "base64":
		ih = new(dictgen.ImageHandlerBase64)
	case "embed":
		ih = new(dictgen.ImageHandlerEmbed)
	case "remove":
		ih = new(dictgen.ImageHandlerRemove)
	default:
		fmt.Fprintf(os.Stderr, "Error: invalid value for --image-method, see --help for details.")
		os.Exit(2)
		return
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
		ff, err := os.OpenFile(*output, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: create dictzip: %v\n", err)
			os.Exit(1)
			return
		}
		f = ff
	}

	fmt.Fprintf(os.Stderr, "Generating dictzip.\n")
	dw := kobodict.NewWriter(f)
	dw.SetEncrypter(e)
	if e != nil {
		fmt.Fprintf(os.Stderr, "  Using encryption.\n")
	}
	switch v := ih.(type) {
	case *dictgen.ImageHandlerBase64:
		fmt.Fprintf(os.Stderr, "  Using image method: optimize and encode as base64 data URL (max_width=%d, max_height=%d, grayscale=%t, jpeg_quality=%d).\n", v.MaxSize.X, v.MaxSize.Y, !v.NoGrayscale, v.JPEGQuality)
	case *dictgen.ImageHandlerEmbed:
		fmt.Fprintf(os.Stderr, "  Using image method: add to dictzip as-is (warning: nickel is buggy with this as of firmware 4.19.14123).\n")
	case *dictgen.ImageHandlerRemove:
		fmt.Fprintf(os.Stderr, "  Using image method: remove images.\n")
	default:
		fmt.Fprintf(os.Stderr, "  Using image method: %#v.\n", v)
	}
	if err := tdf.WriteDictzip(dw, ih, dictgen.ImageFuncFilesystem); err != nil {
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
