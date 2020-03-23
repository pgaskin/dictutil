package marisa

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"io"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

func TestTrieIO(t *testing.T) {
	emptyBuf := bytes.NewBuffer(nil)
	emptyS := "1aa6c451104c2c1b24ecb66ecb84bde2403c49b1" // marisa-build </dev/null | sha1sum -

	normalWd := []string{"asd", "bnm", "cvb", "dfg"} // for n in asd bnm cvb dfg; do echo $n; done | marisa-build | sha1sum -
	normalBuf := bytes.NewBuffer(nil)
	normalS := "bdf9be48216379734fa0256263467ba6ab2e0931"

	t.Run("WriteAll", func(t *testing.T) {
		t.Run("Error", func(t *testing.T) {
			err := WriteAll(new(errIO), normalWd)
			t.Logf("err=%v", err)
			if v := "MARISA_IO_ERROR"; err == nil || !strings.Contains(err.Error(), v) {
				t.Errorf("expected err to contain `%v`, got `%v`", v, err)
			}
		})
		t.Run("Empty", func(t *testing.T) {
			ss := sha1.New()
			if err := WriteAll(io.MultiWriter(emptyBuf, ss), nil); err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			t.Logf("sum=%x", ss.Sum(nil))
			if runtime.GOARCH == "amd64" {
				if v := hex.EncodeToString(ss.Sum(nil)); v != emptyS {
					t.Errorf("output sha1 mismatch: expected %s, got %s", emptyS, v)
				}
			} else {
				t.Logf("skipping sha1 check on non-amd64 architecture, as the correct file differs slightly on each one (usually by ~4 bytes)")
			}
		})
		t.Run("Normal", func(t *testing.T) {
			ss := sha1.New()
			if err := WriteAll(io.MultiWriter(normalBuf, ss), normalWd); err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			t.Logf("sum=%x", ss.Sum(nil))
			if runtime.GOARCH == "amd64" {
				if v := hex.EncodeToString(ss.Sum(nil)); v != normalS {
					t.Errorf("output sha1 mismatch: expected %s, got %s", normalS, v)
				}
			} else {
				t.Logf("skipping sha1 check on non-amd64 architecture, as the correct file differs slightly on each one (usually by ~4 bytes)")
			}
		})
	})
	t.Run("ReadAll", func(t *testing.T) {
		t.Run("Error", func(t *testing.T) {
			wd, err := ReadAll(new(errIO))
			if v := "MARISA_IO_ERROR"; err == nil || !strings.Contains(err.Error(), v) {
				t.Errorf("expected err to contain `%v`, got `%v`", v, err)
			}
			t.Logf("err=%v", err)
			if wd != nil {
				t.Errorf("expected returned slice to be nil, got %#v", wd)
			}
		})
		t.Run("Empty", func(t *testing.T) {
			wd, err := ReadAll(emptyBuf)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			t.Logf("wd=%+s", wd)
			if len(wd) != 0 {
				t.Errorf("expected no words to be returned")
			}
		})
		t.Run("Normal", func(t *testing.T) {
			wd, err := ReadAll(normalBuf)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			t.Logf("wd=%+s", wd)
			if !reflect.DeepEqual(wd, normalWd) {
				t.Errorf("expected %#v, got %#v", normalWd, wd)
			}
		})
	})
}

type errIO struct{}

func (*errIO) Write([]byte) (int, error) { return 0, errors.New("go_test_error") }
func (*errIO) Read([]byte) (int, error)  { return 0, errors.New("go_test_error") }
