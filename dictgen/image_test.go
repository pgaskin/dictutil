package dictgen

import (
	"reflect"
	"testing"
)

func TestImgTagRe(t *testing.T) {
	inHTML := `
		<img src="test">
		<img src="test" />
		<img src="test" alt="asd" />
		<img height="10" width="10" src="test" alt="asd" />
		<img height="10" width="10"
src = "test"
alt="asd" />
	`
	exImg := [][]string{
		{`<img src="`, `test`, `">`},
		{`<img src="`, `test`, `" />`},
		{`<img src="`, `test`, `" alt="asd" />`},
		{`<img height="10" width="10" src="`, `test`, `" alt="asd" />`},
		{`<img height="10" width="10"
src = "`, `test`, `"
alt="asd" />`},
	}

	acMatch := imgTagRe.FindAllStringSubmatch(inHTML, -1)
	acImg := make([][]string, len(acMatch))
	for i, m := range acMatch {
		acImg[i] = m[1:]
	}

	if !reflect.DeepEqual(exImg, acImg) {
		t.Errorf("Expected %#v, got %#v.", exImg, acImg)
	}
}

// TODO(v1): test the image handlers, especially the one which does the replacements
