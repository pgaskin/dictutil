package main

import (
	"fmt"
	"os"

	"github.com/pgaskin/dictutil/kobodict"
	"github.com/spf13/pflag"
)

func init() {
	commands = append(commands, &command{Name: "prefix", Short: "x", Description: "Calculate the prefix for a word", Main: prefixMain})
}

func prefixMain(args []string, fs *pflag.FlagSet) int {
	fs.SortFlags = false
	format := fs.StringP("format", "f", "json-array", "The output format (go-slice, go-map, csv, tsv, json-array, json-object)")
	help := fs.BoolP("help", "h", false, "Show this help text")
	fs.Parse(args[1:])

	if *help || fs.NArg() == 0 {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] word...\n\nOptions:\n%s", args[0], fs.FlagUsages())
		return 0
	}

	if *format != "go-slice" && *format != "go-map" && *format != "csv" && *format != "tsv" && *format != "json-array" && *format != "json-object" {
		fmt.Fprintf(os.Stderr, "Error: invalid format %#v, see --help for more details.\n", *format)
		return 2
	}

	switch *format {
	case "go-slice":
		fmt.Printf("[][]string{\n")
	case "go-map":
		fmt.Printf("map[string]string{\n")
	case "csv", "tsv":
		break
	case "json-array":
		fmt.Printf("[\n")
	case "json-object":
		fmt.Printf("{\n")
	default:
		panic("invalid output format")
	}

	for i, word := range fs.Args() {
		prefix := kobodict.WordPrefix(word)
		last := i == fs.NArg()-1

		switch *format {
		case "go-slice":
			fmt.Printf("\t{%#v, %#v},\n", word, prefix)
		case "go-map":
			fmt.Printf("\t%#v: %#v,\n", word, prefix)
		case "csv":
			fmt.Printf("%s,%s\n", word, prefix)
		case "tsv":
			fmt.Printf("%s\t%s\n", word, prefix)
		case "json-array":
			fmt.Printf("    [%#v, %#v]", word, prefix)
			if last {
				fmt.Printf("\n")
			} else {
				fmt.Printf(",\n")
			}
		case "json-object":
			fmt.Printf("    %#v: %#v", word, prefix)
			if last {
				fmt.Printf("\n")
			} else {
				fmt.Printf(",\n")
			}
		default:
			panic("invalid output format")
		}
	}

	switch *format {
	case "csv", "tsv":
		break
	case "json-array":
		fmt.Printf("]\n")
	case "json-object", "go-slice", "go-map":
		fmt.Printf("}\n")
	default:
		panic("invalid output format")
	}

	return 0
}
