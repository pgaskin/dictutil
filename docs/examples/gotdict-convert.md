---
layout: default
title: gotdict-convert
parent: examples
---

# gotdict-convert
This tool converts [gotdict](https://github.com/wjdp/gotdict) to a dictfile for conversion into a Kobo dictzip.

Images are supported on firmware 4.20.14601+.

## Download
Pre-built dictionaries can be downloaded from the following links:
- GOTDict *(with images, firmware 4.20.14601+)*: [dictzip (dicthtml-gt.zip)](https://ci.appveyor.com/api/projects/geek1011/dictutil/artifacts/gotdict/dicthtml-gt.zip?branch=master&all=false&pr=false), [source dictfile (gotdict.df)](https://ci.appveyor.com/api/projects/geek1011/dictutil/artifacts/gotdict/gotdict.df?branch=master&all=false&pr=false)
- GOTDict *(without images)*: [dictzip (dicthtml-gt.noimg.zip)](https://ci.appveyor.com/api/projects/geek1011/dictutil/artifacts/gotdict/dicthtml-gt.noimg.zip?branch=master&all=false&pr=false), [source dictfile (gotdict.noimg.df)](https://ci.appveyor.com/api/projects/geek1011/dictutil/artifacts/gotdict/gotdict.noimg.df?branch=master&all=false&pr=false)

You can use [dictutil](../dictutil/install.html) to install the dictionaries, or see [here](../dicthtml/install.html) for manual installation instructions.

## Usage

```
Usage: gotdict-convert [options]

Version: dev

Options:
  -g, --gotdict string   The path to the local copy of github.com/wjdp/gotdict. (default "./gotdict")
  -o, --output string    The output filename (will be overwritten if it exists) (- is stdout) (default "./gotdict.df")
  -I, --images           Include images in dictfile
  -h, --help             Show this help text

To convert the resulting dictfile into a dictzip, use dictgen.
```

You can also use the parser as a [Go library](https://pkg.go.dev/github.com/geek1011/dictutil/examples/gotdict-convert/gotdict).
