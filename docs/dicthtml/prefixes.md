---
layout: default
title: Prefixes
parent: dicthtml
---

# Prefixes
Kobo dictionaries are sharded by a prefix derived from the headword.

The information in this document is based on reverse engineering DictionaryParser::htmlForWord.

**Note:** Kobo will only look in the file matching the word's prefix, so if a variant has a different prefix, it must be duplicated into each matching file (note that duplicate words aren't an issue).

**Note:** This document only covers the algorithm used for non-Japanese (Kanji) dictionaries.

## Prefix algorithm
Prefixes are calculated using the following steps. Note that "character" refers to a single Unicode code point, not a byte.

1. Trim the word at the first null byte, if any (i.e. treat it as a C string).
2. Discard everything but the first two characters.
3. Convert the characters to lowercase using the Unicode case mapping rules.
4. Trim all whitespace characters on the left and right sides.
5. If the string is empty, return "11".
6. If the first of the remaining characters is in the Unicode Cyrillic character class, return them as-is.
7. Right-pad the remaining characters to 2 characters long using "`a`"s.
8. If either of the first two characters are not in the Unicode Letter character class, return "11".
9. Return the characters as-is.

## Examples

<!-- dictutil x -fjson-array word | jq -r '.[] | "| \"`" + .[0] + "`\" | \"`" + .[1] + "`\" | |"' -->

| Word | Prefix | Notes |
| --- | --- | --- |
| "`test`" | "`te`" | |
| "`a`" | "`aa`" | |
| "`Èe`" | "`èe`" | The word is made lowercase using unicode rules (i.e. accented characters are included). |
| "`multiple words`" | "`mu`" | |
| "`àççèñts`" | "`àç`" | |
| "`à`" | "`àa`" | |
| "`ç`" | "`ça`" | |
| "" | "`11`" | |
| "`  `" | "`11`" | Space trimming is done after taking the first 2 characters. |
| "` x`" | "`xa`" | |
| "`   123`" | "`11`" | |
| "`x   23`" | "`xa`" | |
| "`д `" | "`д`" | "д" is a Cyrillic character, and it's the first character of the word (after trimming spaces), so it isn't padded with "a"s. |
| "`дaд`" | "`дa`" | |
| "`未未`" | "`未未`" | |
| "`未`" | "`未a`" | Even though "未" is a two-byte character, it is a single unicode rune (and the characters are counted, not bytes). |
| "`  未`" | "`11`" | Space trimming is done after taking the first 2 characters. |
| "` 未`" | "`未a`" | The two-byte "未" character isn't split up when taking the first 2 characters. |

## Testing
You can test Kobo's prefix algorithm directly using [dictword-test](https://github.com/geek1011/kobo-mods/tree/master/dictword-test/).

If you just want an easy way to generate prefixes for words, use the [dictutil prefix](../dictutil/prefix.html) command

## Sample implementation
Here is the Go implementation used in dictutil:

```go
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
			return "11" // if either of the first 2 chars are letters, return "11"
		}
	}

	return string(pfx)
}
```
