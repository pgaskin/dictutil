// Command dictzip-decompile converts a dictzip into a dictfile. The regenerated
// dictzip from the dictfile may not match exactly, but it will look the same,
// and certain bugs with prefixes and variants will be implicitly fixed by the
// conversion process (i.e. variant in wrong file, incorrect prefix, missing
// words in index file). All output is in raw HTML, not Markdown.
//
// This is an experimental tool, and the output may not be perfect on complex
// dictionaries.
package main

import (
	"fmt"
	"io"
	"os"

	"github.com/pgaskin/dictutil/kobodict"
	"github.com/spf13/pflag"

	_ "github.com/pgaskin/dictutil/kobodict/marisa"
)

var version = "dev"

func main() {
	pflag.CommandLine.SortFlags = false
	output := pflag.StringP("output", "o", "."+string(os.PathSeparator)+"decompiled.df", "The output filename (will be overwritten if it exists) (- is stdout)")
	resources := pflag.BoolP("resources", "r", false, "Also extract referenced resources to the current directory (warning: any existing files will be overwritten, so it is recommended to run in an empty directory if enabled)")
	help := pflag.BoolP("help", "h", false, "Show this help text")
	pflag.Parse()

	if *help || pflag.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] dictzip\n\nVersion: dictzip-decompile %s\n\nOptions:\n%s\nArguments:\n  dictzip is the path to the dictzip to decompile.\n\nTo convert the resulting dictfile into a dictzip, use dictgen.\n\nNote: The regenerated dictzip from the dictfile may not match exactly, but it will look the same, and certain bugs with prefixes and variants will be implicitly fixed by the conversion process (i.e. variant in wrong file, incorrect prefix, missing words in index file). All output is in raw HTML, not Markdown.\n\nThis is an experimental tool, and the output may not be perfect on complex dictionaries.\n", os.Args[0], version, pflag.CommandLine.FlagUsages())
		if pflag.NArg() != 0 {
			os.Exit(2)
		} else {
			os.Exit(0)
		}
		return
	}

	fn := pflag.Args()[0]

	fmt.Fprintf(os.Stderr, "Opening input dictzip.\n")
	f, err := os.Open(fn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: open input file %#v: %v.\n", fn, err)
		os.Exit(1)
		return
	}
	defer f.Close()

	s, err := f.Stat()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: stat input file %#v: %v.\n", fn, err)
		os.Exit(1)
		return
	}

	fmt.Fprintf(os.Stderr, "Parsing dictzip.\n")
	dr, err := kobodict.NewReader(f, s.Size())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: parse input file %#v: %v.\n", fn, err)
		os.Exit(1)
		return
	}

	fmt.Fprintf(os.Stderr, "Decompiling dictzip.\n")
	df, err := decompile(dr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: decompile dictzip %#v: %v.\n", fn, err)
		os.Exit(1)
		return
	}

	if *resources {
		fmt.Fprintf(os.Stderr, "Extracting resources.\n")
		for _, f := range dr.File {
			fmt.Fprintf(os.Stderr, "  ./%s\n", f.Name)
			if err := func() error {
				rc, err := f.Open()
				if err != nil {
					return fmt.Errorf("open: %w", err)
				}
				defer rc.Close()

				f, err := os.OpenFile(f.Name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
				if err != nil {
					return fmt.Errorf("create output: %w", err)
				}
				defer f.Close()

				if _, err := io.Copy(f, rc); err != nil {
					return fmt.Errorf("copy: %w", err)
				}

				if err := f.Close(); err != nil {
					return fmt.Errorf("write output: %w", err)
				}

				return nil
			}(); err != nil {
				fmt.Fprintf(os.Stderr, "Error: extract resource %#v: %v.\n", f.Name, err)
				os.Exit(1)
				return
			}
		}
	} else {
		if len(dr.File) != 0 {
			fmt.Fprintf(os.Stderr, "Warning: dictfile contains %d resources, but skipping because resource extraction is not enabled (see --help for more details).\n", len(dr.File))
		}
	}

	fmt.Fprintf(os.Stderr, "Writing dictfile.\n")
	switch *output {
	case "-":
		if err := df.WriteDictFile(os.Stdout); err != nil {
			fmt.Fprintf(os.Stderr, "Error: write dictfile: %v\n", err)
			os.Exit(1)
			return
		}
	default:
		f, err := os.OpenFile(*output, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: create dictfile: %v\n", err)
			os.Exit(1)
			return
		}

		if err := df.WriteDictFile(f); err != nil {
			f.Close()
			fmt.Fprintf(os.Stderr, "Error: write dictfile: %v\n", err)
			os.Exit(1)
			return
		}

		if err := f.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: write dictfile: %v\n", err)
			os.Exit(1)
			return
		}
	}

	fmt.Fprintf(os.Stderr, "Successfully converted %d entries from dictzip %#v to dictfile %s.\n", len(df), fn, *output)
	os.Exit(0)
}
