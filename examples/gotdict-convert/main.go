// Command gotdict-convert converts GOTDict (https://github.com/wjdp/gotdict) to
// a dictgen dictfile.
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/pflag"

	"github.com/geek1011/dictutil/dictgen"
	"github.com/geek1011/dictutil/examples/gotdict-convert/gotdict"
)

var version = "dev"

func main() {
	pflag.CommandLine.SortFlags = false
	gotdictp := pflag.StringP("gotdict", "g", "."+string(os.PathSeparator)+"gotdict", "The path to the local copy of github.com/wjdp/gotdict.")
	output := pflag.StringP("output", "o", "."+string(os.PathSeparator)+"gotdict.df", "The output filename (will be overwritten if it exists) (- is stdout)")
	images := pflag.BoolP("images", "I", false, "Include images in the generated dictfile")
	help := pflag.BoolP("help", "h", false, "Show this help text")
	pflag.Parse()

	if *help || pflag.NArg() != 0 {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\nVersion: gotdict-convert %s\n\nOptions:\n%s\nTo convert the resulting dictfile into a dictzip, use dictgen.\n", os.Args[0], version, pflag.CommandLine.FlagUsages())
		os.Exit(0)
		return
	}

	var img string
	if *images {
		fmt.Fprintf(os.Stderr, "Parsing gotdict (with images).\n")
		img = filepath.Join(*gotdictp, "images")
	} else {
		fmt.Fprintf(os.Stderr, "Parsing gotdict (no images).\n")
	}

	gd, err := gotdict.Parse(filepath.Join(*gotdictp, "_definitions"), img, true)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: parse gotdict: %v\n", err)
		os.Exit(1)
		return
	}

	fmt.Fprintf(os.Stderr, "Transforming definitions.\n")
	var df dictgen.DictFile
	for _, d := range gd {
		var hwi string
		if d.Type != "" {
			hwi = "-" + string(d.Type)
		}

		df = append(df, &dictgen.DictFileEntry{
			Headword:   d.Title,
			HeaderInfo: hwi,
			Variant:    d.Terms,
			Definition: d.Definition,
		})
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

	fmt.Fprintf(os.Stderr, "Successfully converted %d entries from gotdict %s to dictfile %s.\n", len(df), *gotdictp, *output)
	os.Exit(0)
}
