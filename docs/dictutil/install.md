---
layout: default
title: Install
parent: dictutil
---

# Install

## Usage

```
Usage: dictutil install [options] dictzip

Options:
  -k, --kobo string      KOBOeReader path (default: automatically detected)
  -l, --locale string    Locale name to use (format: ALPHANUMERIC{2}; translation dictionaries are not supported) (default: detected from filename if in format dicthtml-**.zip)
  -n, --name string      Custom additional label for dictionary (ignored when replacing built-in dictionaries) (doesn't have any effect on 4.20.14601+)
  -b, --builtin string   How to handle built-in locales [replace = replace and prevent from syncing] [ignore = replace and leave syncing as-is] (default "replace")
  -h, --help             Show this help text

Note:
  If you are not replacing a built-in dictionary, the 'Enable searches on extra
  dictionaries patch' must be installed, or you will not be able to select
  your custom dictionary.
```

## Examples

**Install a dictionary with the locale in the filename (dicthtml-\*\*.zip):**

```sh
dictutil install dicthtml-aa.zip
```

**Install a dictionary with a different locale:**

```sh
dictutil install --locale aa mydictionary.zip
```

**Install a dictionary on a specific Kobo:**

```sh
dictutil install --kobo /path/to/KOBOeReader dicthtml-aa.zip
```

**Install a dictionary with a custom label (4.19.14123 and older):**

```sh
dictutil install --name "My Dictionary" dicthtml-aa.zip
```

## Details
See [installing dictionaries](../dicthtml/install) for more details on how this works.
