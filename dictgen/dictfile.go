package dictgen

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"text/template"
)

// A DictFile is a high-level representation of a Kobo dictionary.
type DictFile []*DictFileEntry

// DictFileEntry represents a single entry in the DictFile.
type DictFileEntry struct {
	Headword string
	Variant  []string

	NoHeader   bool
	HeaderInfo string

	RawHTML    bool
	Definition string

	line int // for internal use if parsed, zero otherwise
}

// ParseDictFile parses a DictFile from it's textual representation (usually
// stored in a file with the extension .df).
func ParseDictFile(r io.Reader) (DictFile, error) {
	var df DictFile
	var dfe *DictFileEntry

	br := bufio.NewScanner(r)
	var line int

	for br.Scan() {
		buf := br.Bytes()
		line++

		if len(buf) == 0 {
			// if in a block and after the metadata (in the definition),
			// preserve the blank line
			if dfe != nil && len(dfe.Definition) != 0 {
				dfe.Definition += "\n"
			}
			continue
		}

		switch buf[0] {
		case '@':
			// start another one
			dfe = new(DictFileEntry)

			// add the headword and line info
			dfe.Headword = strings.TrimSpace(string(buf[1:]))
			dfe.line = line

			// but error if the headword is blank (note that duplicates are
			// acceptable, and encouraged in some cases; Kobo will merge it;
			// try looking up 'be' in the English dictionary)
			if len(dfe.Headword) == 0 {
				return nil, fmt.Errorf("dictfile: line %d: empty headword after @", line)
			}

			// otherwise, add it to the dictfile (remember it's a pointer, it'll
			// still get updated)
			df = append(df, dfe)
		case ':':
			// if not in a block (before the first @), return an error
			if dfe == nil {
				return nil, fmt.Errorf("dictfile: line %d: header info (: or ::) specified before word (@)", line)
			}

			// if already after the metadata (in the definition), return an error
			if len(dfe.Definition) != 0 {
				return nil, fmt.Errorf("dictfile: line %d: header info (: or ::) specified within definition content (prepend a space if this was intended to be part of the definition itself)", line)
			}

			// if already seen the header info (a line starting with :)
			if dfe.NoHeader || len(dfe.HeaderInfo) != 0 {
				return nil, fmt.Errorf("dictfile: line %d: multiple header infos (: or ::) specified in definition block", line)
			}

			// put the trimmed text in the header info, or disable the header if
			// it is ::
			if len(buf) >= 2 {
				if buf[1] == ':' {
					if len(strings.TrimSpace(string(buf[2:]))) != 0 {
						return nil, fmt.Errorf("dictfile: line %d: extra data after no header specified (::)", line)
					}
					dfe.NoHeader = true
				} else {
					dfe.HeaderInfo = strings.TrimSpace(string(buf[1:]))
				}
			} else {
				dfe.HeaderInfo = ""
			}
		case '&':
			// if not in a block, error
			if dfe == nil {
				return nil, fmt.Errorf("dictfile: line %d: variant (&) specified before word (@)", line)
			}

			// if already after the metadata (in the definition), error
			if len(dfe.Definition) != 0 {
				return nil, fmt.Errorf("dictfile: line %d: variant (&) specified within definition content (prepend a space if this was intended to be part of the definition itself)", line)
			}

			// trim the rest of the line (error if nothing left)
			v := strings.TrimSpace(string(buf[1:]))
			if len(v) == 0 {
				return nil, fmt.Errorf("dictfile: line %d: no word after variant specifier (&)", line)
			}

			// and add it to the variant list
			dfe.Variant = append(dfe.Variant, v)
		default:
			// if not in a block, error
			if dfe == nil {
				return nil, fmt.Errorf("dictfile: line %d: definition specified before word (@)", line)
			}

			// append the line to the definition
			dfe.Definition += string(buf) + "\n"
		}
	}

	// check for read errors
	if err := br.Err(); err != nil {
		return nil, err
	}

	// and finally, update the raw html flag and cleanup whitespace
	for _, dfe := range df {
		dfe.Definition = strings.TrimSpace(dfe.Definition)

		if v := strings.TrimSpace(strings.TrimPrefix(dfe.Definition, "<html>")); v != dfe.Definition {
			if strings.HasSuffix(v, "</html>") {
				return nil, fmt.Errorf("dictfile: entry at line %d: raw HTML definitions are specified with <html>, but SHOULD NOT be a full HTML document ending with </html>", dfe.line)
			}
			dfe.RawHTML = true
			dfe.Definition = v
		} else if strings.Contains(dfe.Definition, "<html>") {
			return nil, fmt.Errorf("dictfile: entry at line %d: why does the definition contain a <html> tag ... to make it raw HTML, it should be at the very beginning", dfe.line)
		}
	}

	// note: validation is done separately (and always done before generation)

	return df, nil
}

// Validate validates the entries in the DictFile. Note that duplicate entries
// are fine, and are encouraged if necessary (Kobo will merge them).
func (df DictFile) Validate() error {
	illegal := func(s string, word bool) error {
		if word && strings.Contains(s, "\"") {
			return fmt.Errorf("must not contain %#v", "\"")
		}
		for _, c := range []string{
			"<w", "</w",
			"<html", "</html",
			"<var", "</var",
			"<a name=",
		} {
			// TODO: optimize
			if strings.Contains(s, c) {
				return fmt.Errorf("must not contain %#v", c)
			}
		}
		return nil
	}
	for i, dfe := range df {
		if strings.TrimSpace(dfe.Headword) == "" {
			return fmt.Errorf("word %#v (i:%d, dfe:%#v): headword must not be blank", dfe.Headword, i, dfe)
		} else if err := illegal(dfe.Headword, true); err != nil {
			return fmt.Errorf("word %#v (i:%d): headword contains illegal string: %w", dfe.Headword, i, err)
		}
		for _, v := range dfe.Variant {
			if strings.TrimSpace(v) == "" {
				return fmt.Errorf("word %#v (i:%d): variant %#v must not be blank", dfe.Headword, i, v)
			} else if err := illegal(v, true); err != nil {
				return fmt.Errorf("word %#v (i:%d): variant %#v contains illegal string : %w", dfe.Headword, i, v, err)
			}
		}
		if err := illegal(dfe.HeaderInfo, false); err != nil {
			return fmt.Errorf("word %#v (i:%d): header info %#v contains illegal string : %w", dfe.Headword, i, dfe.HeaderInfo, err)
		}
		if err := illegal(dfe.Definition, false); err != nil {
			return fmt.Errorf("word %#v (i:%d): definition %#v contains illegal string : %w", dfe.Headword, i, dfe.Definition, err)
		}
	}
	return nil
}

// WriteDictFile validates the DictFile and writes it to w in the dictfile
// format.
func (df DictFile) WriteDictFile(w io.Writer) error {
	if err := df.Validate(); err != nil {
		return err
	}

	for _, dfe := range df {
		if err := dfe.writeDictFileEntry(w); err != nil {
			return err
		}
		// for consistency with template if git converted newlines
		if _, err := w.Write([]byte(`
`)); err != nil {
			return err
		}
	}
	return nil
}

// note: this assumes the entry is valid
var dictFileEntryTmpl = template.Must(template.New("").Funcs(template.FuncMap{
	"dfesc": func(str string) string {
		return strings.NewReplacer(
			"\n@", "\n @",
			"\n:", "\n :",
			"\n&", "\n &",
		).Replace(str)
	},
}).Parse(`
{{- /* trim leading whitespace from template */ -}}

{{with .Headword}}@ {{.}}{{end -}}

{{with .NoHeader}}
::{{else}}{{with .HeaderInfo}}
: {{.}}{{end}}{{end -}}

{{range .Variant}}
& {{.}}{{end -}}

{{with .RawHTML}}
<html>{{end -}}

{{with .Definition}}
{{dfesc .}}{{end -}}

{{- /* keep trailing newline at end of template */}}
`))

func (d DictFileEntry) writeDictFileEntry(w io.Writer) error {
	return dictFileEntryTmpl.Execute(w, d)
}
