package dictgen

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

type testcase struct {
	What string

	In  string
	Err error

	Out DictFile

	OutDictFile string
	OutKoboHTML string
}

// TODO(v1): more specific tests
var testcases = []testcase{{
	What: "some of everything",
	In: `@ blank

@ headword
: info
& variant1
&variant2
test
test

@ custom
& NORMALIZEME
::
<html>
<b>custom word:</b>
<p>test</p>
@ markdown
:-test
1. Definition point 1.
  - Blah
  - Blah
2. Blah blah blah.
3. Blah *blah* **blah**!

Blah blah blah.`,
	Out: DictFile{
		{Headword: "blank", Variant: []string(nil), NoHeader: false, HeaderInfo: "", RawHTML: false, Definition: "", line: 1},
		{Headword: "headword", Variant: []string{"variant1", "variant2"}, NoHeader: false, HeaderInfo: "info", RawHTML: false, Definition: "test\ntest", line: 3},
		{Headword: "custom", Variant: []string{"NORMALIZEME"}, NoHeader: true, HeaderInfo: "", RawHTML: true, Definition: "<b>custom word:</b>\n<p>test</p>", line: 10},
		{Headword: "markdown", Variant: []string(nil), NoHeader: false, HeaderInfo: "-test", RawHTML: false, Definition: "1. Definition point 1.\n  - Blah\n  - Blah\n2. Blah blah blah.\n3. Blah *blah* **blah**!\n\nBlah blah blah.", line: 16},
	},
	OutDictFile: `@ blank

@ headword
: info
& variant1
& variant2
test
test

@ custom
::
& NORMALIZEME
<html>
<b>custom word:</b>
<p>test</p>

@ markdown
: -test
1. Definition point 1.
  - Blah
  - Blah
2. Blah blah blah.
3. Blah *blah* **blah**!

Blah blah blah.

`,
	OutKoboHTML: `<html><w><p><a name="blank" /><b>blank</b></p><var></var></w><w><a name="custom" /><var><variant name="normalizeme"/></var><b>custom word:</b>
<p>test</p></w><w><p><a name="headword" /><b>headword</b> info</p><var><variant name="variant1"/><variant name="variant2"/></var><p>test
test</p></w><w><p><a name="markdown" /><b>markdown</b> -test</p><var></var><ol>
<li>Definition point 1.

<ul>
<li>Blah</li>
<li>Blah</li>
</ul></li>
<li>Blah blah blah.</li>
<li>Blah <em>blah</em> <strong>blah</strong>!</li>
</ol>

<p>Blah blah blah.</p></w></html>`,
}}

func TestDictFile(t *testing.T) {
	for _, tc := range testcases {
		t.Logf("case %#v", tc.What)

		df, err := ParseDictFile(strings.NewReader(tc.In))
		if tc.Err == nil && err != nil {
			t.Fatalf("case %#v: parse dictfile: unexpected error: %v", tc.What, err)
		} else if tc.Err != nil && err == nil {
			t.Fatalf("case %#v: parse dictfile: expected error (%v)", tc.What, tc.Err)
		} else if tc.Err != nil && tc.Err.Error() != err.Error() {
			t.Fatalf("case %#v: parse dictfile: expected error (%v), got: %v", tc.What, tc.Err, err)
		}

		exp, err := json.MarshalIndent(tc.Out, "| ", "    ")
		if err != nil {
			panic(err)
		}

		act, err := json.MarshalIndent(df, "| ", "    ")
		if err != nil {
			panic(err)
		}

		if !reflect.DeepEqual(exp, act) {
			for _, dfe := range df {
				fmt.Printf("%#v,\n", dfe)
			}
			t.Fatalf("case %#v: expected:\n%s\n\ngot:\n%s", tc.What, exp, act)
		}

		buf := bytes.NewBuffer(nil)
		if err := df.WriteDictFile(buf); err != nil {
			t.Fatalf("case %#v: write dictfile: unexpected error: %v", tc.What, err)
		} else if tc.OutDictFile != buf.String() {
			fmt.Printf("expected:\n`%s`\n\ngot:\n`%s`", tc.OutDictFile, buf.String())
			t.Fatalf("case %#v: unexpected dictfile output", tc.What)
		}

		buf.Reset()
		if err := df.WriteKoboHTML(buf); err != nil {
			t.Fatalf("case %#v: write kobo html: unexpected error: %v", tc.What, err)
		} else if tc.OutKoboHTML != buf.String() {
			fmt.Printf("expected:\n`%s`\n\ngot:\n`%s`", tc.OutKoboHTML, buf.String())
			t.Fatalf("case %#v: unexpected kobo html output", tc.What)
		}
	}
}
