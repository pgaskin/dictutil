---
layout: default
---

These tools are designed to work with v2 dictionaries (4.7.10364+).

## dicthtml
These pages are some notes I've made about the Kobo dictionary format based on reverse engineering the firmware and the official dictionaries.

- **[Format](./dicthtml/format):** About the Kobo dictionary format.
- **[Prefixes](./dicthtml/prefixes):** Details about prefix calculation.
- **[v1/v2 dictionaries](./dicthtml/v1v2):** Changes between v1/v2 dictionaries.
- **[Installing custom dictionaries](./dicthtml/install):** Notes about sideloading dictionaries.

## dictutil
dictutil is a low-level tool to unpack, pack, and perform other operations on Kobo dictzips.

- Coming soon.

## dictgen
dictgen is an easy-to-use tool/library to generate Kobo dictionaries from scratch or use in conversion scripts. It deals with all the unusual bits (e.g. variant capitalization, prefix generation, etc) for you and gives warnings when it can't.

- Coming soon.

## examples
These are some tools which make use of dictutil.

- **[gotdict-convert](./examples/gotdict-convert):** Converts [github.com/wjdp/gotdict](https://github.com/wjdp/gotdict) to a dictfile.

## other

- **[dictword-test](https://github.com/geek1011/kobo-mods/tree/master/dictword-test):** Calculates word prefixes using libnickel.
- **[marisa](https://github.com/geek1011/dictutil/tree/master/marisa):** Marisa bindings for Go.
