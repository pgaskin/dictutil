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

Note that currently, the only usable image method is removing them or using base64-encoding (for firmware 4.20.14601+; older versions segfault in the in-book dictionary if images are enabled), as embedded dict:/// image URLs cause the webviews to appear blank (this is a nickel bug). See https://github.com/geek1011/dictutil/issues/1 for more details.

See https://pgaskin.net/dictutil/dictgen for more information about the dictfile format.
```

## Example usage

**Building a dictzip for a dictfile:**

```
dictgen my-dictionary.df
```

If you are using Windows, you can also drag-and-drop a dictfile onto dictgen.exe. 

**Merging multiple dictfiles into a single dictzip:**

```
dictgen my-dictionary.df another.df
```

If you are using Windows, you can also drag-and-drop multiple dictfiles onto dictgen.exe. 

**Building a dictzip with images removed:**

```
dictgen -I remove my-dictionary.df
```

**Specifying a custom output filename:**

```
dictgen -o dicthtml-df.zip my-dictionary.df
```

## Dictfile format
Dictgen uses a simple, but feature-complete format for representing Kobo dictionaries.

A dictfile (with the file extension `.df`) is a plain-text file consisting of multiple entries.

Each entry represents a single definition. There can be more than one entry per word. An entry is denoted by a line starting with `@ ` followed by the headword. The headword can contain spaces, capital letters, and so on.

After the headword, zero or more header lines can be added. To add additional variants which will be matched, use `& ` followed by the word variant. The variant can be anything which could be used in a headword. This can be specified more than once, but only one variant can be specified for each `& `. Another header type is word information, denoted by a `: `. If specified, the text following it is appended after the bolded headword on the same line (see the English built-in dictionary for an example; it has things like `-verb` and the pronunciation information here). If you want to have complete control over how the entry is displayed, use `::` (without anything following it) instead of `: `. This will remove the default bolded headword at the top of the generated entry.

After the header lines, you can include the body of the entry. By default, this uses [Markdown](https://github.com/adam-p/markdown-here/wiki/Markdown-Cheatsheet) for formatting. If you want to include raw HTML, prepend the HTML with `<html>` (don't include a closing tag). This can span multiple lines, and will continue until the next entry or end of file.

In addition, you can include GIF and JPEG images in the body using the usual Markdown or HTML syntax. If the image path is relative (i.e. not a full path), it is resolved relative to the directory you run dictgen from.

You can also include custom CSS (per-entry) by including it between the `<style>` and `</style>` tags. This is supported in both HTML and Markdown mode.

## Dictfile reference

- `@ HEADWORD`: Start a new entry. The headword doesn't have to be unique, and can contain spaces.
  - Header
    - `: WORD_INFO` or `::` *(optional)*: Add extra word info after the headword, or remove it entirely.
    - `& VARIANT` *(optional)*: Add an additional word to match. Follows the same rules as the headword. Can be repeated multiple times.
  - Body
    - `MARKDOWN` or `<html> RAW_HTML`: Include a definition written in Markdown or raw HTML code.

## Examples

### Simplest

```
@ word
Definition here.
@ word 1
Definition 1 here.
@ test
Blah blah blah.
```

### Simple

```
@ no
- No means no...

@ NO
- A different definition for nitric oxide.
- Blah blah blah.

@ go
& went
& going
1. This definition is matched by three different words.
2. It's also numbered rather than bulleted.
   - With some sub-items.
   - And another.

An image:

![](image.jpg)

@ test
: this appears beside the headword
Blah blah blah.
```

### Full

```
@ word
This is the definition of a word.

@ word 2
This is the defnition of the second word.

@ water
& H2O
1. You can also use lists in Markdown.
2. And **bold text** or *italic text*.
   - Sub-items are also supported.

@ test
: -noun
Blah blah blah.

@ test
: -verb
Blah blah blah.

@ custom
::
**This is a custom word header!**

And the definition here:
- Blah blah blah.
- Blah blah blah.

@ images
Embedding an image (relative paths):

![](image.jpg)

Embedding an image (Linux/macOS style paths):

![](/path/to/image.jpg)

Embedding an image (Windows style paths):

![](C:/path/to/image.jpg)


@ raw-html
<html><p>This definition contains raw html.</p>

<p>You can split it into multiple lines for readability.</p>

<ul>
  <li>You can also use all HTML tags.</li>
  <li><span style="background: #666">This text has a dark background</span></li>
  <li><span class="test">This text is styled with CSS classes.</span></li>
</ul>

<style>
.test {
  text-decoration: underline;
}
</style>
```
