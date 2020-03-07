package main

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"regexp"
	"unicode"

	"github.com/geek1011/dictutil/dictgen"
	"github.com/geek1011/dictutil/kobodict"
)

// This isn't exposed as a separate package, as it's subject to change and
// highly specific to dictzip-decompile.

// The regexps used to extract data should have a similar level of strictness as
// the ones used by nickel (for simplicity, compatibility, and predictability).

// decompile decompiles a dictzip into a dictfile. External resources are not
// extracted, and must be done separately.
//
// Duplicate entries (e.g. the ones added by dictgen for fixing broken variants)
// are collapsed into one. They will be expanded again as necessary when the
// dictfile is compiled by dictgen.
func decompile(r *kobodict.Reader) (dictgen.DictFile, error) {
	var df dictgen.DictFile
	seenEntries := map[[20]byte]struct{}{}
	for _, dh := range r.Dicthtml {
		if err := func() error {
			rc, err := dh.Open()
			if err != nil {
				return fmt.Errorf("open: %w", err)
			}
			defer rc.Close()

			buf, err := ioutil.ReadAll(rc)
			if err != nil {
				return fmt.Errorf("read: %w", err)
			}

			es, err := extractEntries(buf)
			if err != nil {
				return fmt.Errorf("extract entries: %w", err)
			}

			for _, e := range es {
				ss := sha1.Sum(e)
				if _, ok := seenEntries[ss]; ok {
					continue
				}
				seenEntries[ss] = struct{}{}

				de, err := decompileEntry(e)
				if err != nil {
					return fmt.Errorf("decompile entry %#v: %w", string(e), err)
				}

				df = append(df, de)
			}

			return nil
		}(); err != nil {
			return nil, fmt.Errorf("process dicthtml %#v: %w", dh.Name, err)
		}
	}
	return df, nil
}

// The regexps/vars used by decompileEntry.
var (
	// generator matchers (match the entire entry, split into parts) (match in order) (don't include variants here)
	generator1PenelopeRe         = regexp.MustCompile(`^(?s)<a name="([^"]+)"\/><div><b>([^<]+)<\/b><br\/>(.+)<\/div>$`)                                        // also: first and second groups must be equal
	generator2KoboFrRe           = regexp.MustCompile(`^(?s)<p><a name="([^"]+)" ?(?:\/>|><\/a>)<b>\s*([^<]+)\s*<\/b>\s*(.*?)<br ?\/><br ?\/>\s*(.+)\s*<\/p>$`) // also: 2nd and 3rd (header) group must not contain "<br", "<li", "<var", "<p"; also: need to wrap returned content in a p tag
	generator3KoboEnOrDictutilRe = regexp.MustCompile(`^(?s)<p><a name="([^"]+)" ?(?:\/>|><\/a>)<b>\s*(.+?)\s*<\/b>\s*(.*?)\s*<\/p>\s*(.+)\s*$`)                // also: 2nd and 3rd (header) group must not contain "<br", "<li", "<var", "<p"
	// fallback matchers (if none of the above exist)
	headFallbackIndexWordRe = regexp.MustCompile(`<a name="([^"]+)" ?(?:\/>|><\/a>)`) // this is slightly more lenient than some of Kobo's (it makes the space before the closing optional)
	// other matchers
	variantsRe     = regexp.MustCompile(`<var>(.*?)<\/var>`)
	variantsItemRe = regexp.MustCompile(`<variant name="([^"]+)" ?(?:\/>|><\/variant>)`)
)

// decompileEntry parses an entry (it must be trimmed).
func decompileEntry(buf []byte) (*dictgen.DictFileEntry, error) {
	var entry dictgen.DictFileEntry

	// Generator-specific enhanced extraction (for making use of dictfile lines
	// starting with &, :, etc).
	var generatorMatched bool
	// -- Penelope: https://github.com/pettarin/penelope/blob/fce6dcfd899d3755ae3a5a3867d7d436105ada56/penelope/format_kobo.py#L167
	//    e.g. <w><a name="dfgdfg"/><div><b>dfgdfg</b><br/>Penelope</div>sdfsdf</div></w>
	if !generatorMatched {
		if m := generator1PenelopeRe.FindSubmatch(buf); len(m) != 0 {
			headwordIndex, headwordDisplay, contentHTML := m[1], m[2], m[3]
			if !bytes.Equal(headwordIndex, headwordDisplay) {
				// it's a false positive if those aren't identical
			} else {
				entry.Headword = string(headwordIndex)
				entry.RawHTML = true
				entry.Definition = string(contentHTML)
				generatorMatched = true
			}
		}
	}
	// -- Kobo: based on dicthtml-fr
	//    e.g. <w><p><a name="a-"/><b>a-, an-</b><br/><br/><ol> <li>Élément exprimant la négation ( pas ), ou la privation ( sans ). </li>&nbsp;&nbsp;&nbsp;⇒anormal, apolitique. </ol></p></w>
	if !generatorMatched {
		if m := generator2KoboFrRe.FindSubmatch(buf); len(m) != 0 {
			headwordIndex, headwordDisplay, headerInfo, contentHTML := m[1], m[2], m[3], m[4]
			if bytes.Contains(headwordDisplay, []byte("<br")) || bytes.Contains(headerInfo, []byte("<br")) {
				// it's a false positive if those contain line breaks
			} else if bytes.Contains(headwordDisplay, []byte("<li")) || bytes.Contains(headerInfo, []byte("<li")) {
				// it's a false positive if those contain list items
			} else if bytes.Contains(headwordDisplay, []byte("<var")) || bytes.Contains(headerInfo, []byte("<var")) {
				// it's a false positive if those contain variants
			} else if bytes.Contains(headwordDisplay, []byte("<p")) || bytes.Contains(headerInfo, []byte("<p")) {
				// it's a false positive if those contain new paragraphs
			} else {
				if bytes.EqualFold(headwordIndex, headwordDisplay) {
					entry.Headword = string(headwordDisplay)
				} else {
					entry.Headword = string(headwordIndex)
				}
				entry.RawHTML = true
				entry.HeaderInfo = string(headerInfo)
				entry.Definition = "<p>" + string(contentHTML) + "</p>"
				generatorMatched = true
			}
		}
	}
	// -- Kobo: based on dicthtml-en, a few others
	//    e.g. <w><p><a name="ab"></a><b>ab</b> [<pr>'ab</pr>] -n</p><var><variant name="variant-added-for-testing"/></var><p><ol><li>an abdominal muscle usu. used in pl.</li><li>about</li></ol></p></w>
	// -- or dictgen
	//    e.g. <w><p><a name="a" /><b>a</b> A (# emph. #).</p><var><variant name="variant-added-for-testing"/></var><ol><li>Etym: [Shortened form of an. AS. an one. See One.] An adjective, commonly called the indefinite article, and signifying one or any, but less emphatically.</li><li>&#34;At a birth&#34;; &#34;In a word&#34;; &#34;At a blow&#34;. Shak. Note: It is placed before nouns of the singular number denoting an individual object, or a quality individualized, before collective nouns, and also before plural nouns when the adjective few or the phrase great many or good many is interposed; as, a dog, a house, a man; a color; a sweetness; a hundred, a fleet, a regiment; a few persons, a great many days. It is used for an, for the sake of euphony, before words beginning with a consonant sound [for exception of certain words beginning with h, see An]; as, a table, a woman, a year, a unit, a eulogy, a ewe, a oneness, such a one, etc. Formally an was used both before vowels and consonants.</li><li>Etym: [Originally the preposition a (an, on).] In each; to or for each; as, &#34;twenty leagues a day&#34;, &#34;a hundred pounds a year&#34;, &#34;a dollar a yard&#34;, etc.</li></ol></w>
	if !generatorMatched {
		if m := generator3KoboEnOrDictutilRe.FindSubmatch(buf); len(m) != 0 {
			headwordIndex, headwordDisplay, headerInfo, contentHTML := m[1], m[2], m[3], m[4]
			if bytes.Contains(headwordDisplay, []byte("<br")) || bytes.Contains(headerInfo, []byte("<br")) {
				// it's a false positive if those contain line breaks
			} else if bytes.Contains(headwordDisplay, []byte("<li")) || bytes.Contains(headerInfo, []byte("<li")) {
				// it's a false positive if those contain list items
			} else if bytes.Contains(headwordDisplay, []byte("<var")) || bytes.Contains(headerInfo, []byte("<var")) {
				// it's a false positive if those contain variants
			} else if bytes.Contains(headwordDisplay, []byte("<p")) || bytes.Contains(headerInfo, []byte("<p")) {
				// it's a false positive if those contain new paragraphs
			} else {
				if bytes.EqualFold(headwordIndex, headwordDisplay) {
					entry.Headword = string(headwordDisplay)
				} else {
					entry.Headword = string(headwordIndex)
				}
				entry.RawHTML = true
				entry.HeaderInfo = string(headerInfo)
				entry.Definition = string(contentHTML)
				generatorMatched = true
			}
		}
	}
	// -- Fallback: extract (then remove) the first headword, rest goes in raw html definition.
	//    e.g. <w><a name="test"><p>dfkgjdlfjglkdfjg</p><var><variant name="asd"/></var></w>
	if !generatorMatched {
		entry.NoHeader = true
		entry.RawHTML = true
		entry.Definition = string(headFallbackIndexWordRe.ReplaceAllFunc(buf, func(src []byte) []byte {
			if entry.Headword != "" {
				return src // don't continue after the first headword has been found
			}
			entry.Headword = string(headFallbackIndexWordRe.FindSubmatch(src)[1])
			return nil // remove the entire a tag
		}))
		if entry.Headword == "" {
			return nil, fmt.Errorf("no headword found in %#v", string(buf))
		}
		generatorMatched = true
	}

	// Add any additional headwords (then remove) (which really shouldn't be there in the first place) as variants.
	// i.e. stray <a name="..."> tags (but not if the link has text, because then it's not a headword anymore)
	entry.Definition = string(headFallbackIndexWordRe.ReplaceAllFunc([]byte(entry.Definition), func(src []byte) []byte {
		entry.Variant = append(entry.Variant, string(headFallbackIndexWordRe.FindSubmatch(src)[1]))
		return nil // remove the entire a tag
	}))

	// Append (then remove) any variants found in the raw html.
	// i.e. <var> tags inside <variant> ones
	entry.Definition = string(variantsRe.ReplaceAllFunc([]byte(entry.Definition), func(src []byte) []byte {
		for _, m := range variantsItemRe.FindAllSubmatch(src, -1) {
			entry.Variant = append(entry.Variant, string(m[1]))
		}
		return nil // remove the entire variant tag
	}))

	return &entry, nil
}

// The regexps/vars used by extractEntries.
var (
	htmlStart = []byte("<html>")
	htmlEnd   = []byte("</html>")
	entryRe   = regexp.MustCompile(`(?s)<w>\s*(.+?)\s*<\/w>`)
)

// extractEntries gets the trimmed body of each entry in the dicthtml file.
func extractEntries(buf []byte) ([][]byte, error) {
	if idx := bytes.Index(buf, htmlStart); idx < 0 {
		return nil, fmt.Errorf("missing %s tag", string(htmlStart))
	} else {
		buf = buf[idx+len(htmlStart):]
	}

	if idx := bytes.LastIndex(buf, htmlEnd); idx < 0 {
		return nil, fmt.Errorf("missing %s tag", string(htmlStart))
	} else {
		buf = buf[:idx]
	}

	var entries [][]byte

	var cur, prev, body []int
	prev = []int{0, 0}
	for _, m := range entryRe.FindAllSubmatchIndex(buf, -1) {
		cur, body = m[0:2][:], m[2:4]
		for _, b := range buf[prev[1]:cur[0]] {
			// note: even though we might split up multi-byte utf-8 chars
			// here, it's fine, as the whitespace should be ascii if any,
			// and if there is anything else, it's an issue.
			if !unicode.IsSpace(rune(b)) {
				return nil, fmt.Errorf("non-whitespace between word entries (%#v in %#v before %#v)", string(rune(b)), string(buf[prev[1]:cur[0]]), string(buf[cur[0]:cur[1]]))
			}
		}
		prev = cur
		entries = append(entries, buf[body[0]:body[1]])
	}
	for _, b := range buf[prev[1]:] {
		if !unicode.IsSpace(rune(b)) {
			return nil, fmt.Errorf("non-whitespace after last word entry (%#v in %#v)", string(rune(b)), string(buf[prev[1]:]))
		}
	}

	return entries, nil
}
