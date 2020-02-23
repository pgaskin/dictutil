package kobodict

import (
	"strings"
	"unicode"
)

// NormalizeWordReference normalizes a word for use in an dicthtml headword
// (<a name="...") or variant (<variant name="..."). It matches the way Kobo
// finds words in a file.
//
// The logic is reversed from DictionaryParser::htmlForWord in libnickel.
//
// Note: Headwords are prefix-matched against the query, the uppercased query,
// the lowercased query, and the lowercased query with the first letter
// uppercased. Variants are only prefix-matched against the lowercased query.
//
// Note: The matching is only done in the file matching the prefix for the query.
func NormalizeWordReference(w string, variant bool) string {
	if variant {
		// variants must always be lowercase (the match is only checked against
		// the lowercased query)
		w = strings.ToLower(w)
	}

	// trim leading and trailing whitespace
	return strings.TrimSpace(w)
}

// WordPrefix gets the prefix of a word for sharding dicthtml files.
//
// This is not to be used with Kanji, as those are handled by a separate
// function for Japanese dictionaries.
//
// WordPrefix is a simplification of the logic reversed from
// DictionaryParser::htmlForWord (see wordPrefix), but with performance and
// cleaner code. It is should have the exact same results.
func WordPrefix(word string) string {
	pfx := []rune(word)

	for i, c := range pfx {
		if i >= 2 || c == '\x00' { // limit to 2 chars, also cut at null
			pfx = pfx[:i] // trim up to current char
			break
		}
		pfx[i] = unicode.ToLower(c) // this includes accented chars
	}

	for len(pfx) != 0 {
		if unicode.IsSpace(pfx[0]) {
			pfx = pfx[1:] // trim left space
		} else {
			break
		}
	}

	for len(pfx) != 0 {
		if unicode.IsSpace(pfx[len(pfx)-1]) {
			pfx = pfx[:len(pfx)-1] // trim right space
		} else {
			break
		}
	}

	if len(pfx) == 0 {
		return "11" // if empty, return "11"
	}

	if !unicode.Is(unicode.Cyrillic, pfx[0]) {
		for len(pfx) < 2 {
			pfx = append(pfx, 'a') // pad right with 'a's to 2 chars
		}
		if !unicode.IsLetter(pfx[0]) || !unicode.IsLetter(pfx[1]) {
			return "11" // if neither of the first 2 chars are letters, return "11"
		}
	}

	return string(pfx)
}

// wordPrefix gets the prefix of a word for sharding dicthtml files.
//
// This is not to be used with Kanji, as those are handled by a separate
// function for Japanese dictionaries.
//
// The logic is reversed from DictionaryParser::htmlForWord in libnickel. It
// matches it as closely as possible.
func wordPrefix(w string) string {
	// w
	// QString::toLower()
	w = strings.ToLower(w)

	// QString::leftRef(2)
	if len(w) > 2 {
		w = string([]rune(w)[:2])
	}

	// QString::trimmed()
	w = strings.TrimSpace(w)

	// simplify the following code by converting to rune slice
	r := []rune(w)

	// A null byte is a valid Unicode character, but in C, it's treated as
	// the end of a string. To keep compatibility with libnickel, we need to
	// end a string there if necessary.
	for i, c := range r {
		if c == '\x00' {
			r = r[:i]
			break
		}
	}

	// DictionaryParser::isCyrillic(w[0])
	// skip if true
	if !(len(r) != 0 && unicode.Is(unicode.Cyrillic, r[0])) {
		// add an 'a' for right padding if not 2 chars
		if len(r) != 2 {
			r = append(r, 'a')
		}
	}

	// DictionaryParser::isCyrillic(w[0])
	// skip if != false
	switch {
	case !(len(r) != 0 && unicode.Is(unicode.Cyrillic, r[0])):
		// inlined QChar::isLetter(w[0]), QChar::isLetter(w[1]), unnecessary length check
		// skip if both true
		if (len(r) >= 1 && unicode.IsLetter(r[0])) && (len(r) >= 2 && unicode.IsLetter(r[1])) {
			break
		}
		fallthrough
	case len(r) == 0:
		// w = QString::fromLatin1_helper("11"..., 2)
		return "11"
	}

	return string(r)
}
