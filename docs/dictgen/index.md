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
  -o, --output string         The output filename (will be overwritten if it exists) (- is stdout) (default "dicthtml.zip")
  -c, --crypt string          Encrypt the dictzip using the specified encryption method (format: method:keyhex)
  -I, --image-method string   How to handle images (if an image path is relative, it is loaded from the current dir) (base64 - optimize and encode as base64, embed - add to dictzip, remove) (default "base64")
  -h, --help                  Show this help text

If multiple dictfiles (*.df) are provided, they will be merged (duplicate entries are fine; they will be shown in sequential order). To read from stdin, use - as the filename.

Note that the only usable image method is currently removing them or using base64-encoding (for firmware 4.20.14601+; older versions segfault in the in-book dictionary if images are enabled), as embedded dict:/// image URLs cause the webviews to appear blank (this is a nickel bug). See https://github.com/geek1011/dictutil/issues/1 for more details.

See https://pgaskin.net/dictutil/dictgen for more information about the dictfile format.
```

## Dictfile format
TODO
