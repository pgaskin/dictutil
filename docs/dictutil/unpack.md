---
layout: default
title: Unpack
parent: dictutil
---

# Unpack

## Usage

```
Usage: dictutil unpack [options] dictzip

Options:
  -o, --output string   The output directory (must not exist) (default: the basename of the input without the extension)
  -c, --crypt string    Decrypt the dictzip (if needed) using the specified encryption method (format: method:keyhex)
  -h, --help            Show this help text
```

## Examples

**Unpack a dictionary:**

```sh
dictutil unpack dicthtml.zip
# The output is written to ./dicthtml
```

```sh
dictutil unpack dicthtml-fr.zip
# The output is written to ./dicthtml-fr
```

**Unpack a dictionary to a custom directory:**

```
dictutil unpack --output mydictionary dicthtml.zip
```

## Details
An unpacked dictdir contains:

- `words`: The parsed marisa word list (newline-separated).
- `*.html`: The ungzipped dicthtml files.
- `*`: Any additional files as-is.
