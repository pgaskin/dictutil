// Package dictgen simplifies creating full-featured dictionaries for Kobo
// eReaders, with support for images, unicode prefixes, raw html, markdown, and
// more.
//
// A marisa implementation must be provided by
// github.com/pgaskin/kobodict/marisa or a custom one for this package to work.
package dictgen

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"strings"
	"text/template"

	"github.com/pgaskin/dictutil/kobodict"
	"github.com/russross/blackfriday/v2"
)

// WriteDictzip writes the dictfile to a kobodict.Writer, which should not have
// been used yet. The writer is not closed automatically. If the ImageHandler
// requires a file to be opened (i.e. not ImageHandlerRemove), the provided
// ImageFunc will be called.
func (df DictFile) WriteDictzip(dw *kobodict.Writer, ih ImageHandler, img ImageFunc) error {
	var prefixes []string
	prefixed := df.Prefixed()
	for pfx := range prefixed {
		prefixes = append(prefixes, pfx)
	}
	sort.Strings(prefixes)

	hbuf := bytes.NewBuffer(nil)
	for _, pfx := range prefixes {
		for _, dfe := range prefixed[pfx] {
			if err := dw.AddWord(dfe.Headword); err != nil {
				return fmt.Errorf("add word %#v: %w", dfe.Headword, err)
			}
			for _, v := range dfe.Variant {
				if err := dw.AddWord(v); err != nil {
					return fmt.Errorf("add variant %#v: %w", v, err)
				}
			}
		}
		hbuf.Reset()
		if err := prefixed[pfx].WriteKoboHTML(hbuf); err != nil {
			return fmt.Errorf("generate dicthtml for %s: %w", pfx, err)
		} else if buf, err := transformHTMLImages(ih, dw, hbuf.Bytes(), img); err != nil {
			return fmt.Errorf("generate dicthtml for %s: transform images: %w", pfx, err)
		} else if hw, err := dw.CreateDicthtml(pfx); err != nil {
			return fmt.Errorf("write dicthtml for %s: %w", pfx, err)
		} else if _, err = hw.Write(buf); err != nil {
			return fmt.Errorf("write dicthtml for %s: %w", pfx, err)
		}
	}

	return nil
}

// Prefixed shards the DictFile into the different word prefixes. The original
// DictFile is unchanged, but the entries are still pointers to the originals
// (i.e. the result will become out of date if you modify the entries).
//
// The DictFile is not validated.
//
// If a variamt has a different prefix, the entire entry is duplicated as
// necessary.
func (df DictFile) Prefixed() map[string]DictFile {
	prefixed := map[string]DictFile{}
	for _, dfe := range df {
		pfx := map[string]bool{}

		pfx[kobodict.WordPrefix(dfe.Headword)] = true
		for _, v := range dfe.Variant {
			pfx[kobodict.WordPrefix(v)] = true
		}

		for p := range pfx {
			prefixed[p] = append(prefixed[p], dfe)
		}
	}
	return prefixed
}

// WriteKoboHTML validates the DictFile and writes it to w in the dicthtml
// format.
func (df DictFile) WriteKoboHTML(w io.Writer) error {
	if err := df.Validate(); err != nil {
		return err
	}

	// must be sorted for proper matching
	dfs := df[:]
	sort.Slice(dfs, func(i int, j int) bool {
		return dfs[i].Headword < dfs[j].Headword
	})

	if _, err := w.Write([]byte("<html>")); err != nil {
		return err
	}
	for _, dfe := range dfs {
		if err := dfe.writeKoboHTML(w); err != nil {
			return err
		}
	}
	if _, err := w.Write([]byte("</html>")); err != nil {
		return err
	}

	return nil
}

// note: we don't want the html/template escaping, this isn't actually proper
// html, and also, the whitespaces in the end tags should stay EXACTLY as is
// (yes, I know there is a space before the end of the a but not the variant) to
// provide the best possible matches against the regexps Kobo uses. Also, the
// output should not have any newlines. Also, keep in mind headwords can have
// unescaped html tags in it, and they will be rendered properly by Kobo.
var koboHTMLTmpl = template.Must(template.New("").Funcs(template.FuncMap{
	"md": func(md string) string {
		return strings.TrimSpace(string(blackfriday.Run([]byte(md))))
	},
	"normhw": func(headword string) string {
		return kobodict.NormalizeWordReference(headword, false)
	},
	"normv": func(variant string) string {
		return kobodict.NormalizeWordReference(variant, true)
	},
}).Parse(`
{{- /* trim */ -}}

<w>
	{{- if .NoHeader -}}
		<a name="{{normhw .Headword}}" />
	{{- else -}}
		<p><a name="{{normhw .Headword}}" /><b>{{.Headword}}</b>{{with .HeaderInfo}} {{.}}{{end}}</p>
	{{- end -}}
	<var>
		{{- range .Variant -}}
			<variant name="{{normv .}}"/>
		{{- end -}}
	</var>
	{{- with .Definition -}}
		{{- if $.RawHTML -}}
			{{.}}
		{{- else -}}
			{{md .}}
		{{- end -}}
	{{- end -}}
	{{- with .PostRawHTML -}}
		{{.}}
	{{- end -}}
</w>

{{- /* trim */ -}}
`))

func (d DictFileEntry) writeKoboHTML(w io.Writer) error {
	return koboHTMLTmpl.Execute(w, d)
}
