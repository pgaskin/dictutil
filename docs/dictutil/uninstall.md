---
layout: default
title: Uninstall
parent: dictutil
---

# Uninstall

## Usage

```
Usage: dictutil uninstall [options] locale

Options:
  -k, --kobo string      KOBOeReader path (default: automatically detected)
  -b, --builtin string   How to handle built-in locales [normal = uninstall the same way as the UI] [delete = completely delete the entry (doesn't have any effect on 4.20.14601+)] [restore = download the original dictionary from Kobo again] (doesn't have any effect on 4.24.15672+) (default "normal")
  -B, --no-custom        Uninstall built-in dictionaries instead of custom ones on 4.24.15672+
  -h, --help             Show this help text
```

## Examples

**Uninstall a dictionary:**

```sh
dictutil uninstall aa
```

**Restore a overwritten built-in dictionary:**

```sh
dictutil uninstall --builtin restore fr
```

**Completely delete a built-in dictionary:**

```sh
dictutil uninstall --builtin delete fr
```

Note: You can restore the dictionary by manually downloading it and using [dictutil install](./install).

## Details
Uninstall does the following steps:

1. If the DB entry for the dictionary exists:
   - Built-in (normal): Set `Installed` to `false`.
   - Built-in (delete): Remove the row for the suffix.
   - Built-in (restore): Set `Installed` to `true`.
   - Extra: Remove the row for the suffix.
2. If the dictionary is not built-in and there is an `ExtraLocales` entry for the locale in the `.kobo/Kobo/Kobo eReader.conf`, remove it.
3. With the dictzip:
   - Built-in (normal): Delete it if it exists.
   - Built-in (delete): Delete it if it exists.
   - Built-in (restore): Delete it if it exists, then download it again from Kobo.
   - Extra: Delete it if it exists.
