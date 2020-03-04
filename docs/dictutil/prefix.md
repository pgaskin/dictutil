---
layout: default
title: Prefix
parent: dictutil
---

# Prefix

## Usage

```
Usage: dictutil prefix [options] word...

Options:
  -f, --format string   The output format (go-slice, go-map, csv, tsv, json-array, json-object) (default "json-array")
  -h, --help            Show this help text
```

## Examples

**Get the prefix for a word:**

```sh
dictutil prefix "word"
```

**Get the prefix for multiple words:**

```sh
dictutil prefix "word1" "word2" "word3"
```

**Get the prefix for multiple words as CSV:**

```sh
dictutil prefix --format csv "word1" "word2" "word3"
```
