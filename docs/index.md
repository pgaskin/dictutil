---
layout: default
title: Documentation
nav_order: 1
---

# Dictutil
{: .fs-9 }

A collection of documentation and tools for working with Kobo dictionaries.
{: .fs-6 .fw-300 }

[Download](https://github.com/pgaskin/dictutil/releases){: .btn .btn-primary .fs-5 .mb-4 .mb-md-0 .mr-2 } [dicthtml](./dicthtml/){: .btn .fs-5 .mb-4 .mb-md-0 } [dictgen](./dictgen/){: .btn .fs-5 .mb-4 .mb-md-0 } [dictutil](./dictutil/){: .btn .fs-5 .mb-4 .mb-md-0 }

---

**Prebuilt dictionaries:**

[GOTDict](./examples/gotdict-convert.html#download){: .btn .fs-3 .mb-1 .mb-md-0 } [Webster's 1913 Unabridged Dictionary](./examples/webster1913-convert.html#download){: .btn .fs-3 .mb-1 .mb-md-0 }

---

These tools are designed to work with v2 dictionaries (4.7.10364+).

## Getting started
If you're interested in creating dictionaries, look at the [dictgen documentation](./dictgen/). If you're interested in installing or manipulating existing dictionaries, see the [dictutil documentation](./dictutil/). Otherwise, see the [dicthtml documentation](./dicthtml/) for more information about the Kobo dictionary format.

## dicthtml
These pages are some notes I've made about the Kobo dictionary format based on reverse engineering the firmware and the official dictionaries.

- **[Format](./dicthtml/format.html):** About the Kobo dictionary format.
- **[Prefixes](./dicthtml/prefixes.html):** Details about prefix calculation.
- **[v1/v2 dictionaries](./dicthtml/v1v2.html):** Changes between v1/v2 dictionaries.
- **[Installing custom dictionaries](./dicthtml/install.html):** Notes about sideloading dictionaries.

## dictutil
dictutil is a low-level tool to unpack, pack, and perform other operations on Kobo dictzips.

- **[Dictutil](./dictutil/)**
- **[Install](./dictutil/install.html):** Install a dictzip.
- **[Uninstall](./dictutil/uninstall.html):** Uninstall a dictzip.
- **[Pack](./dictutil/pack.html):** Pack a dictzip from a dictdir.
- **[Unpack](./dictutil/unpack.html):** Unpack a dictzip into a dictdir.
- **[Prefix](./dictutil/prefix.html):** Calculate the dicthtml prefix for a word.

## dictgen
dictgen is an easy-to-use tool/library to generate Kobo dictionaries from scratch or use in conversion scripts. It deals with all the unusual bits (e.g. variant capitalization, prefix generation, etc) for you and gives warnings when it can't.

- **[Dictgen](./dictgen#usage)**
- **[Dictfile format](./dictgen#dictfile-format)**

## examples
These are some tools which make use of dictutil to convert actual dictionaries.

- **[gotdict-convert](./examples/gotdict-convert.html):** Converts [github.com/wjdp/gotdict](https://github.com/wjdp/gotdict) to a dictfile.
- **[webster1913-convert](./examples/webster1913-convert.html):** Converts [Project Gutenberg's Webster's Unabridged Dictionary](http://www.gutenberg.org/ebooks/29765.txt.utf-8) to a dictfile.
- **[dictzip-decompile](./examples/dictzip-decompile.html):** An **experimental** tool to convert a dictzip into a dictfile.
- **[bgl-convert](./examples/bgl-convert.html):** A simple tool to convert Babylon BGL dictionaries to a dictfile.

## other

- **[dictword-test](https://github.com/pgaskin/kobo-mods/tree/master/dictword-test):** Calculates word prefixes using libnickel.
- **[marisa](https://github.com/pgaskin/dictutil/tree/master/marisa):** Marisa bindings for Go.
