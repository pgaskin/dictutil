package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/geek1011/dictutil/dictgen"
	"github.com/spf13/pflag"
)

func main() {
	pflag.CommandLine.SortFlags = false
	gotdict := pflag.StringP("gotdict", "g", "."+string(os.PathSeparator)+"gotdict", "The path to the local copy of github.com/wjdp/gotdict.")
	output := pflag.StringP("output", "o", "."+string(os.PathSeparator)+"gotdict.df", "The output file path (will be overwritten if it exists) (- is stdout)")
	help := pflag.BoolP("help", "h", false, "Show this help text")
	pflag.Parse()

	if *help || pflag.NArg() != 0 {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\nOptions:\n%s\nTo convert the resulting dictfile into a dictzip, use dictgen.\n", os.Args[0], pflag.CommandLine.FlagUsages())
		os.Exit(0)
		return
	}

	gd, err := ParseGOTDict(filepath.Join(*gotdict, "_definitions"), "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: parse gotdict: %v\n", err)
		os.Exit(1)
		return
	}

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

	if *output == "-" {
		if err := df.WriteDictFile(os.Stdout); err != nil {
			fmt.Fprintf(os.Stderr, "Error: write dictfile: %v\n", err)
			os.Exit(1)
			return
		}
	} else {
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

	fmt.Fprintf(os.Stderr, "Successfully converted %d entries from %s to %s.\n", len(df), *gotdict, *output)
	os.Exit(0)
}
