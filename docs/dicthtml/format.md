---
layout: default
title: Dictionary format
parent: dicthtml
---

# Dictionary format

This document will refer to a packed `dicthtml*.zip` files as a *dictzip*, an unpacked and deindexed one as a *dictdir*, and an individual dictionary html file as a *dicthtml* file.

This document describes [version 2](./v1v2.html) of Kobo dictionaries. Japanese dictionaries are not the focus of this document, and information about such may be incomplete.

Kobo dictzips are standard ZIP files. In general, they are named `dicthtml-LOCALE.zip` where **LOCALE** is the ISO-639-1 language code (although custom locale codes can be used if added to ExtraLocales in the Kobo config file). In addition, a dictzip may be named `dicthtml-LOCALE-LOCALE.zip` for translation dictionaries. Some official dictionaries have slightly different locale codes (e.g. `jaxxdjs`), but these are exceptions.

As of firmware 4.22.15190, the official dictionaries are hosted on `https://kbdownload1-a.akamaihd.net/ereader/dictionaries/v2/dicthtml*.zip`, and support the following locales:

| Locale     | Language |
| ---        | --- |
| de         | Deutsch |
| de-en      | Deutsch - English |
| en-de      | English - Deutsch |
| en         | English |
| en-es      | English - Español |
| en-fr      | English - Français |
| en-it      | English - Italiano |
| en-ja      | English - 日本語（ジーニアス） |
| en-ja-pgs  | English - 日本語（プログレッシブ） |
| en-nl      | English - Nederlands |
| en-pt      | English - Português |
| en-tr      | English - Türkçe |
| es-en      | Español - English |
| es         | Español |
| fr-en      | Français - English |
| fr         | Français |
| it-en      | Italiano - English |
| it         | Italiano |
| jaxxdjs    | 日本語 |
| nl         | Nederlands |
| pt-en      | Português - English |
| pt         | Português |

Note that due to licensing reasons, some locales, like `nl`, aren't directly searchable. In addition, the official dictionaries are encrypted (there won't be any details on this here).

## Dictzip files
A dictzip is a ZIP archive containing multiple files. These files should only be in the top-level of the archive. The filenames must be encoded in UTF-8, especially when the prefixes include accents and other special characters.

### words (index)
The words file is a binary dump of a [marisa trie](https://github.com/s-yata/marisa-trie). The trie contains each headword/variant in the dictionary with the whitespace trimmed, but otherwise left as-is (i.e. no case normalization). This trie is used when autocompleting words (the word will be displayed as it is in the index) and checking if they are in the dictionary. Each word has an equal weighting.

In a dictdir (an unpacked dictzip) unpacked by dictutil, the words file is a plaintext list of words. The order does not matter.

To build the words file, you can use something like the following C++ code:

```cpp
#include <algorithm>
#include <string>
#include <vector>
#include <marisa.h>

void save_words(std::string filename, std::vector<std::string> words) {
    std::sort(words.begin(), words.end());
    std::unique(words.begin(), words.end());

    auto ks = new marisa::Keyset();
    for (auto const& w: words)
        ks->push_back(w.c_str());
    
    auto tr = new marisa::Trie();
    tr->build(*ks);
    tr->save(filename.c_str());
}

int main() {
    save_words("words", std::vector<std::string>{
        "word1",
        "word2",
        "word3",
    });
}
```

For more details on word matching, see [here](./matching.html).

### *.html (dicthtml)
Each dictzip contains one or more UTF-8 encoded dicthtml files named `PREFIX.html`, where **PREFIX** is the prefix calculated from the words within it.

This file contains an `<html>` tag containing one or more `<w>` tags. Each `<w>` tag contains a single definition. A definition consists of almost-XHTML code. The `<w>` tag must contain a tag like `<a name="HEADWORD" />`, where **HEADWORD** is the headword normalized by trimming spaces. The **HEADWORD** can contain a limited set of unescaped HTML tags (even though this wouldn't be strictly valid XHTML), including `<sup></sup>` and `<sub></sub>`. Somewhere after the `<a>` tag, an optional list of variants can be included by enclosing repeated `<variant name="VARIANT"/>` tags inside a `<var></var>` tag. The **VARIANT** must be normalized by trimming whitespace and lowercasing it (following unicode normalization rules, i.e. normalizing accented characters as well). The rest of the `<w>` tag can contain valid HTML which is displayed as-is in the definition.

Note that these tags must be included EXACTLY as specified above (including whitespace), as dicthtml files are parsed using regexps by Kobo. As of firmware 4.19.14123, the following regexps are used (they are unlikely to change for v2 dictionaries):

| Regexp | Purpose |
| --- | --- |
| `<w>` | Locating the start of a word entry (DictionaryParser::searchWordMultipleCases). |
| `<a name=".*" />` | Locating headwords (DictionaryParser::findAllDefinitionsForWord). |
| `<variant name="%1"` | Locating a specific variant (DictionaryParser::searchWordMultipleCases). |
| `</w>` | Splitting at the end of an entry (DictionaryParser::searchWordMultipleCases). |
| `(<a name="` | Kanji stuff (not discussed in detail here). |
| `([%1-%2]?)` | Kanji stuff (not discussed in detail here). |
| `(<a name="%1" />.*</w>)` | Kanji stuff (not discussed in detail here). |
| `(<a name="%1".*</w>)` | Extracting a definition (DictionaryParser23searchWordMultipleCases) (note: this is why it's generally a good practice to keep the headword and variants as close to the start of the entry's `<w>` tag as possible) |

In addition, if there are multiple `<w>` tags with duplicate headwords or variants, they will all be displayed in the order they appear in the dicthtml when looking up the word.

If you want to match the official dictionaries, you can use the following template (note that newlines have been added for readability). If there are not any variants, leave the `<var>` tag in, but empty.

```xml
<w>
    <p><a name="HEADWORD_TRIMMED" /><b>HEADWORD</b> -noun (or pronunciations, nothing, or whatever)</p>
    <var>
        <variant name="VARIANT_TRIMMED_LOWERCASED"/>
        <variant name="VARIANT_TRIMMED_LOWERCASED"/>
    </var>
    definition html
</w>
```

For best results, words with multiple variants with different prefixes should have their entire `<w>` entry copied into the dicthtml files for each prefix. If a headword/variant does not match the expected prefix for the dicthtml file it is in, it will be ignored silently.

It is not necessary to explicitly specify variants for things like plurals with a `s` suffix, as if a word cannot be found, the first prefix match will be used (i.e. `tests` will match a headword/variant named `test`). For this reason, words should be sorted from shortest to longest in alphabetical order in general (to prevent accidental matches of the wrong word).

A dicthtml file can optionally be encrypted using AES-128-ECB encryption and PKCS#7 padding. 

### *.gif, *.jpg, etc
Dictzips can also contain a few specific types of resources. As of firmware 4.19.14123, only GIF (with the `.gif` extension and the `GIF` magic) and JPEG (with the `.jpg` extension and the `JFIF` magic) images are supported.

To reference the images, you must use a URL like `dict:///example.gif`; just `example.gif` won't work. But beware: if anything after `dict:` doesn't exist, nickel will segfault and reboot when you try to view the entry.

Another way to add images is to reference them with a data URL (e.g. `<img src="data:image/gif;base64,...">`). The advantage is that you don't need to include images separately, and any filetype will work, but the disadvantage is that large images will significantly increase the loading time for the dictionary.

If you have control over the target device, you can also use `file:///...` URLs to reference a local file (which can also be of any filetype).

As of firmware 4.19.14123, all of these methods are too buggy (due to bugs in libnickel) to be usable. The only one which works is base64-encoded images in the full-screen dictionary view (i.e. not the in-book dictionary). The `dict:///` URLs cause the webview to appear blank, and the base64-encoded and file URLs cause nickel to segfault in the in-book dictionary view. See [#1](https://github.com/geek1011/dictutil/issues/1) for more details.

Starting in firmware 4.20.14601, the base64 method works perfectly in both views (yay!). Other URLs don't segfault the in-book dictionary view anymore, and `dict:///` URLs still blank the webview.

## Example

This is an example dictdir using most of the things mentioned above.

**words**

```
H<sub>2</sub>O
h2o
h2o1
Test Word
dihydrogen monoxide
example
test word 1
test-image
test-image-base64
testing
testing 1
```

**11.html**

```xml
<html>
    <w>
        <p><a name="h<sub>2</sub>o" /><b>H<sub>2</sub>O</b> -example</p>
        <var>
            <variant name="h2o"/>
            <variant name="h2o1"/>
            <variant name="dihydrogen monoxide"/>
        </var>
        <p>Water</p>
    </w>
</html>
```

**di.html**

```xml
<html>
    <w>
        <p><a name="h<sub>2</sub>o" /><b>H<sub>2</sub>O</b> -example</p>
        <var>
            <variant name="h2o"/>
            <variant name="h2o1"/>
            <variant name="dihydrogen monoxide"/>
        </var>
        <p>Water</p>
    </w>
</html>
```

**ex.html**

```xml
<html>
    <w>
        <p><a name="Test Word" /><b>Test Word</b> -example</p>
        <var>
            <variant name="test word 1"/>
            <variant name="example"/>
        </var>
        <p>
            <ul>
                <li>Lorem ipsum dolor.</li>
                <li>Blah blah blah.</li>
            </ul>
        </p>
    </w>
</html>
```

**te.html**

```xml
<html>
    <w>
        <p><a name="Test Word" /><b>Test Word</b> -example</p>
        <var>
            <variant name="test word 1"/>
            <variant name="example"/>
        </var>
        <p>
            <ul>
                <li>Lorem ipsum dolor.</li>
                <li>Blah blah blah.</li>
            </ul>
        </p>
    </w>
    <w>
        <a name="testing" />
        <p>The simplest possible entry.</p>
    </w>
    <w>
        <a name="testing" />
        <var></var>
        <p>This will also appear another definition.</p>
    </w>
    <w>
        <a name="test-image" />
        <var></var>
        <p>Image (a black square) (this won't work on current firmware versions): <img src="dict:///example.gif" /></p>
    </w>
    <w>
        <a name="test-image-base64" />
        <var></var>
        <p>Image (a black square): <img src="data:image/gif;base64,R0lGODlhPAA8AIABAAAAAP///yH5BAEKAAEALAAAAAA8ADwAAAJAhI+py+0Po5y02ouz3rz7D4biSJbmiabqyrbuC8fyTNf2jef6zvf+DwwKh8Si8YhMKpfMpvMJjUqn1Kr1ihUVAAA7" /></p>
    </w>
    <w>
        <a name="testing 1" />
        <p><span style="background: black; color: white;">Test</span></p>
    </w>
</html>
```

**example.gif**

<img src="data:image/gif;base64,R0lGODlhPAA8AIABAAAAAP///yH5BAEKAAEALAAAAAA8ADwAAAJAhI+py+0Po5y02ouz3rz7D4biSJbmiabqyrbuC8fyTNf2jef6zvf+DwwKh8Si8YhMKpfMpvMJjUqn1Kr1ihUVAAA7" />
