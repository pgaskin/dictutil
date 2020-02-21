---
layout: default
title: gotdict-convert
parent: examples
---

# gotdict-convert
This tool converts [gotdict](https://github.com/wjdp/gotdict) to a dictfile for conversion into a Kobo dictzip.

Images are not supported yet.

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

## Pre-converted files
Ready-to-use GOTDict builds for Kobo eReaders can be found [here](https://cloud.drone.io/geek1011/dictutil). Choose the top item, click on `gotdict-convert`, then click on `upload` to find the link to the latest version.
