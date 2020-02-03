package main

import (
	"fmt"
	"io"
	"os"

	"github.com/geek1011/dictutil/dictgen"
	"github.com/spf13/pflag"
)

func main() {
	pflag.CommandLine.SortFlags = false
	output := pflag.StringP("output", "o", "dicthtml.zip", "The output filename (will be overwritten if it exists) (- is stdout)")
	// TODO: image-dir
	help := pflag.BoolP("help", "h", false, "Show this help text")
	pflag.Parse()

	if *help || pflag.NArg() == 0 {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] dictfile...\n\nOptions:\n%s\nIf multiple dictfiles (*.df) are provided, they will be merged (duplicate entries are fine; they will be shown in sequential order). To read from stdin, use - as the filename.\n\nThe dictfile format:\n  TODO: short doc text\n", os.Args[0], pflag.CommandLine.FlagUsages())
		os.Exit(0)
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

	fmt.Fprintf(os.Stderr, "Generating dictzip.\n")
	switch *output {
	case "-":
		if err := tdf.WriteDictzip(os.Stdout); err != nil {
			fmt.Fprintf(os.Stderr, "Error: write dictzip: %v\n", err)
			os.Exit(1)
			return
		}
	default:
		f, err := os.OpenFile(*output, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: create dictzip: %v\n", err)
			os.Exit(1)
			return
		}

		if err := tdf.WriteDictzip(f); err != nil {
			f.Close()
			fmt.Fprintf(os.Stderr, "Error: write dictzip: %v\n", err)
			os.Exit(1)
			return
		}

		if err := f.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: write dictzip: %v\n", err)
			os.Exit(1)
			return
		}
	}

	fmt.Fprintf(os.Stderr, "Successfully wrote %d entries from %d dictfile(s) to dictzip %s.\n", len(tdf), pflag.NArg(), *output)
	os.Exit(0)
}
