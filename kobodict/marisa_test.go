package kobodict

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"reflect"
	"runtime"
	"testing"
)

func TestMarisa(t *testing.T) {
	if Marisa == nil {
		t.Skipf("warning: Marisa not supported on platform GOOS=%s GOARCH=%s and must be provided externally", runtime.GOOS, runtime.GOARCH)
	}

	w := []string{
		"asd",
		"dfg",
		"sdf",
	}

	buf := bytes.NewBuffer(nil)
	if err := Marisa.WriteAll(buf, w); err != nil {
		t.Fatalf("unexpected error when writing trie: %v", err)
	} else if buf.Len() == 0 {
		t.Errorf("written trie is empty")
	}

	ss := sha1.New()

	nw, err := Marisa.ReadAll(io.TeeReader(buf, ss))
	if err != nil {
		t.Fatalf("unexpected error when reading written trie: %v", err)
	} else if len(nw) == 0 {
		t.Errorf("read trie is empty")
	} else if !reflect.DeepEqual(nw, w) {
		t.Errorf("read tree: expected %+s, got %+s", w, nw)
	}

	if runtime.GOARCH == "amd64" {
		if x, y := hex.EncodeToString(ss.Sum(nil)), "ea7252fc4e86585dea884e4bcb5ce7be90676474"; x != y {
			t.Errorf("trie output is incorrect or non-determinstic, expected sha1 %s, got %s", y, x)
		}
	} else {
		t.Logf("skipping sha1 check on non-amd64 architecture, as the correct file differs slightly on each one (usually by ~4 bytes)")
	}
}
