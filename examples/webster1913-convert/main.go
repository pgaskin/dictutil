// Command webster1913-convert converts Project Gutenberg's Webster's 1913
// Unabridged Dictionary to a dictgen dictfile.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"os"

	"github.com/spf13/pflag"

	"github.com/pgaskin/dictutil/dictgen"
	"github.com/pgaskin/dictutil/examples/webster1913-convert/webster1913"
)

var version = "dev"

var deftmpl = template.Must(template.New("").Funcs(template.FuncMap{
	"spldc": func(s string) []string {
		for i, c := range s {
			if c == '.' || c == ',' || c == '(' {
				return []string{s[:i], s[i:]}
			}
		}
		return []string{"", s}
	},
}).Parse(`
	{{- with .Etymology}}<p><i>{{.}}</i></p>{{end -}}
	{{- with .Meanings}}<ol>{{range .}}<li>{{.Text}}{{with .Example}}<br/><br/>{{.}}{{end}}</li>{{end}}</ol>{{end -}}
	{{- with .PhraseDefns}}<p>{{range $n, $v := .}}{{if $n}} {{end}}{{range $x, $y := (spldc $v)}}{{if $x}}<span>{{$y}}</span>{{else}}<b>{{$y}}</b>{{end}}{{end}}{{end}}</p>{{end -}}
	{{- with .Synonyms}}<p>{{range $n, $v := .}}{{if $n}} {{end}}{{$v}}{{end}}</p>{{end -}}
	{{- with .Extra}}<p>{{.}}</p>{{end -}}
`))

func main() {
	pflag.CommandLine.SortFlags = false
	output := pflag.StringP("output", "o", "."+string(os.PathSeparator)+"webster1913.df", "The output filename (will be overwritten if it exists) (- is stdout)")
	dump := pflag.Bool("dump", false, "Instead of converting, dump the parsed dictionary to stdout as JSON (for debugging)")
	help := pflag.BoolP("help", "h", false, "Show this help text")
	pflag.Parse()

	if *help || pflag.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] gutenberg_webster1913_path\n\nVersion: webster1913-convert %s\n\nOptions:\n%s\nArguments:\n  gutenberg_webster1913_path is the path to Project Gutenberg's Webster's 1913 dictionary. Use - to read from stdin.\n\nTo convert the resulting dictfile into a dictzip, use dictgen.\n", os.Args[0], version, pflag.CommandLine.FlagUsages())
		os.Exit(0)
		return
	}

	fmt.Fprintf(os.Stderr, "Opening input file.\n")
	var r io.Reader
	switch v := pflag.Args()[0]; v {
	case "-":
		r = os.Stdin
	default:
		f, err := os.Open(v)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: open input %#v: %v\n", v, err)
			os.Exit(1)
			return
		}
		defer f.Close()
		r = f
	}

	fmt.Fprintf(os.Stderr, "Parsing dictionary.\n")
	wd, err := webster1913.Parse(r, func(i int, word string) {
		if i%1000 == 0 {
			fmt.Fprintf(os.Stderr, "[% 5d] %s\n", i, word)
		}
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: parse webster1913: %v\n", err)
		os.Exit(1)
		return
	}

	if *dump {
		fmt.Fprintf(os.Stderr, "Dumping JSON to stdout.\n")
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "    ")
		enc.Encode(wd)
		os.Exit(0)
		return
	}

	fmt.Fprintf(os.Stderr, "Transforming definitions.\n")
	var df dictgen.DictFile
	dbuf := bytes.NewBuffer(nil)
	for _, d := range wd {
		dbuf.Reset()
		if err := deftmpl.Execute(dbuf, d); err != nil {
			fmt.Fprintf(os.Stderr, "Error: render definition %#v: %v\n", d, err)
			os.Exit(1)
			return
		}
		df = append(df, &dictgen.DictFileEntry{
			Headword:   d.Headword,
			Variant:    d.Variant,
			RawHTML:    true,
			HeaderInfo: d.Info,
			Definition: dbuf.String(),
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

	fmt.Fprintf(os.Stderr, "Successfully converted %d entries from Webster's 1913 dictionary %#v to dictfile %s.\n", len(df), pflag.Args()[0], *output)
	os.Exit(0)
}
