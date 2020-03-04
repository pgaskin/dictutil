---
layout: default
title: Pack
parent: dictutil
---

# Pack

## Usage

```
Usage: dictutil pack [options] dictdir

Options:
  -o, --output string   The output dictzip filename (will be overwritten if it exists) (default "dicthtml.zip")
  -c, --crypt string    Encrypt the dictzip using the specified encryption method (format: method:keyhex)
  -h, --help            Show this help text
```

## Examples

**Pack a dictdir:**

```sh
dictutil pack /path/to/dictdir
# the output is written to dicthtml.zip
```

**Pack a dictdir to a specific filename:**

```sh
dictutil pack --output "dicthtml-aa.zip" /path/to/dictdir
```

## Input format
The input dictdir is the same as the output of [dictutil unpack](./unpack).
