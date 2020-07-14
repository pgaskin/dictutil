---
layout: default
title: webster1913-convert
parent: examples
---

# webster1913-convert
This tool converts [Project Gutenberg's Webster's Unabridged Dictionary](http://www.gutenberg.org/ebooks/29765.txt.utf-8) into a dictfile for conversion into a Kobo dictzip.

## Download
Pre-built dictionaries can be downloaded from the following links:
- Webster's 1913 Dictionary: [dictzip (dicthtml-wb.zip)](https://ci.appveyor.com/api/projects/pgaskin/dictutil/artifacts/webster1913/dicthtml-wb.zip?branch=master&all=false&pr=false), [source dictfile (webster1913.df)](https://ci.appveyor.com/api/projects/pgaskin/dictutil/artifacts/webster1913/webster1913.df?branch=master&all=false&pr=false)

You can use [dictutil](../dictutil/install.html) to install the dictionaries, or see [here](../dicthtml/install.html) for manual installation instructions.

## Usage

```
Usage: webster1913-convert [options] gutenberg_webster1913_path

Options:
  -o, --output string   The output filename (will be overwritten if it exists) (- is stdout) (default "./webster1913.df")
      --dump            Instead of converting, dump the parsed dictionary to stdout as JSON (for debugging)
  -h, --help            Show this help text

Arguments:
  gutenberg_webster1913_path is the path to Project Gutenberg's Webster's 1913 dictionary. Use - to read from stdin.

To convert the resulting dictfile into a dictzip, use dictgen.
```

The source dictionary can be downloaded [here](http://www.gutenberg.org/ebooks/29765.txt.utf-8) or [here](https://github.com/pgaskin/dictserver/raw/master/data/dictionary.txt).

You can also use the parser as a [Go library](https://pkg.go.dev/github.com/pgaskin/dictutil/examples/webster1913-convert/webster1913).
