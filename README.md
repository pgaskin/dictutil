<h1 align="center">dictutil</h1>

[![](https://img.shields.io/github/v/release/geek1011/dictutil?include_prereleases)](https://github.com/geek1011/dictutil/releases) [![](https://img.shields.io/drone/build/geek1011/dictutil/master)](https://cloud.drone.io/geek1011/dictutil) [![](https://img.shields.io/badge/godoc-reference-blue.svg)](https://pkg.go.dev/mod/github.com/geek1011/dictutil?tab=versions) [![](https://goreportcard.com/badge/github.com/geek1011/dictutil)](https://goreportcard.com/report/github.com/geek1011/dictutil)

This repository contains a collection of tools and libraries to work with Kobo dictionaries, plus comprehensive documentation of Kobo's dictionary format.

Unlike previous attempts at working with Kobo dictionaries, dictutil has full support for all features supported by nickel (word prefixes, unicode, variants, images, etc), with a focus on simplicity, correctness (prefix generation and other features are directly tested against libnickel's code and regexps, v1/v2 dictionaries are differentiated), and completeness (most of the research was done by reverse-engineering libnickel).

Dictutil consists of multiple tools and libraries:
- [**dictutil**](https://pgaskin.net/dictutil/dictutil) provides commands for installing, removing, unpacking, packing, and performing low-level modifications and tests on Kobo dictionaries. All operations are intended to be correct, lossless, and deterministic.
- [**dictgen**](https://pgaskin.net/dictutil/dictgen) simplifies creating full-featured dictionaries for Kobo eReaders, with support for images, unicode prefixes, raw html, markdown, and more.
- [**examples/gotdict-convert**](https://pgaskin.net/dictutil/examples/gotdict-convert) is a working example of using dictutil to convert [GOTDict](https://github.com/wjdp/gotdict) into a Kobo dictionary.
- *Library:* [**kobodict**](https://pkg.go.dev/github.com/geek1011/dictutil/kobodict) provides support for reading, writing, encrypting, and decrypting Kobo dictionaries.
- *Library:* [**dictgen**](https://pkg.go.dev/github.com/geek1011/dictutil/dictgen) provides the functionality of dictgen as a library.
- *Library:* [**marisa**](./marisa) provides self-contained CGO bindings for [marisa-trie](https://github.com/s-yata/marisa-trie).

Dictutil implements [version 2](https://pgaskin.net/dictutil/dicthtml/v1v2) of the Kobo dictionary format, which supports firmware versions 4.7.10364+.

For more information, see the [documentation](https://pgaskin.net/dictutil). If you just want a quick overview of the utilities provided, continue reading below.

## Usage
See the [documentation](https://pgaskin.net/dictutil) for more detailed information and examples.

### dictutil

```
Usage: dictutil command [options] [arguments]

Dictutil provides low-level utilities to manipulate Kobo dictionaries (v2).

Commands:
  install (I)          Install a dictzip file
  pack (p)             Pack a dictzip file
  prefix (x)           Calculate the prefix for a word
  uninstall (U)        Uninstall a dictzip file
  unpack (u)           Unpack a dictzip file
  help                 Show help for all commands

Options:
  -h, --help   Show this help text
```

```
Usage: dictutil install [options] dictzip

Options:
  -k, --kobo string      KOBOeReader path (default: automatically detected)
  -l, --locale string    Locale name to use (format: ALPHANUMERIC{2}; translation dictionaries are not supported) (default: detected from filename if in format dicthtml-**.zip)
  -n, --name string      Custom additional label for dictionary (ignored when replacing built-in dictionaries)
  -b, --builtin string   How to handle built-in locales [replace = replace and prevent from syncing] [ignore = replace and leave syncing as-is] (default "replace")
  -h, --help             Show this help text

Note:
  If you are not replacing a built-in dictionary, the 'Enable searches on extra
  dictionaries patch' must be installed, or you will not be able to select
  your custom dictionary.
```

```
Usage: dictutil uninstall [options] locale

Options:
  -k, --kobo string      KOBOeReader path (default: automatically detected)
  -b, --builtin string   How to handle built-in locales [normal = uninstall the same way as the UI] [delete = completely delete the entry] [restore = download the original dictionary from Kobo again] (default "normal")
  -h, --help             Show this help text
```

```
Usage: dictutil pack [options] dictdir

Options:
  -o, --output string   The output dictzip filename (will be overwritten if it exists) (default "dicthtml.zip")
  -c, --crypt string    Encrypt the dictzip using the specified encryption method (format: method:keyhex)
  -h, --help            Show this help text
```

```
Usage: dictutil unpack [options] dictzip

Options:
  -o, --output string   The output directory (must not exist) (default: the basename of the input without the extension)
  -c, --crypt string    Decrypt the dictzip (if needed) using the specified encryption method (format: method:keyhex)
  -h, --help            Show this help text
```

```
Usage: dictutil prefix [options] word...

Options:
  -f, --format string   The output format (go-slice, go-map, csv, tsv, json-array, json-object) (default "json-array")
  -h, --help            Show this help text
```

### dictgen

```
Usage: dictgen [options] dictfile...

Options:
  -o, --output string         The output filename (will be overwritten if it exists) (- is stdout) (default "dicthtml.zip")
  -c, --crypt string          Encrypt the dictzip using the specified encryption method (format: method:keyhex)
  -I, --image-method string   How to handle images (if an image path is relative, it is loaded from the current dir) (base64 - optimize and encode as base64, embed - add to dictzip, remove) (default "base64")
  -h, --help                  Show this help text

If multiple dictfiles (*.df) are provided, they will be merged (duplicate entries are fine; they will be shown in sequential order). To read from stdin, use - as the filename.

Note that the only usable image method is currently removing them or using base64-encoding (for firmware 4.20.14601+; older versions segfault in the in-book dictionary), as embedded dict:/// image URLs cause the webviews to appear blank (this is a nickel bug). See https://github.com/geek1011/dictutil/issues/1 for more details.

See https://pgaskin.net/dictutil/dictgen for more information about the dictfile format.
```

**See [here](https://pgaskin.net/dictutil/dictgen) for information and examples of the dictfile format.**

### gotdict-convert

```
Usage: gotdict-convert [options]

Options:
  -g, --gotdict string   The path to the local copy of github.com/wjdp/gotdict. (default "./gotdict")
  -o, --output string    The output filename (will be overwritten if it exists) (- is stdout) (default "./gotdict.df")
  -I, --images           Include images in dictfile
  -h, --help             Show this help text

To convert the resulting dictfile into a dictzip, use dictgen.
```