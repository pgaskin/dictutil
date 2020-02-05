package dictgen

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"text/template"
)

type DictFile []*DictFileEntry

type DictFileEntry struct {
	Headword string
	Variant  []string

	NoHeader   bool
	HeaderInfo string

	RawHTML    bool
	Definition string

	line int // for internal use if parsed, zero otherwise
}

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
			return nil, fmt.Errorf("dictfile: entry at line %d: why does the definition contain a <html> tag...", dfe.line)
		}
	}

	// note: validation is done separately (and always done before generation)

	return df, nil
}

// Validate validates the entries in the DictFile.
func (df DictFile) Validate() error {
	// TODO(v0): check for empty words and variants
	// TODO(v0): fields can't have </w </html <var <variant </var </variant name="
	// TODO(v0): in addition, words and headwords can't have "
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
