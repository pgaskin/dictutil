---
layout: default
title: dictgen
has_children: false
---

# dictgen

This section contains documentation for dictgen, a high-level tool to create Kobo dictionaries.
{: .fs-6 .fw-300 }

## Usage

```
Usage: dictgen [options] dictfile...

Options:
  -o, --output string   The output filename (will be overwritten if it exists) (- is stdout) (default "dicthtml.zip")
  -c, --crypt string    Encrypt the dictzip using the specified encryption method (format: method:keyhex)
  -h, --help            Show this help text

If multiple dictfiles (*.df) are provided, they will be merged (duplicate entries are fine; they will be shown in sequential order). To read from stdin, use - as the filename.
```

## Dictfile format
TODO
