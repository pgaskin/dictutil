//+build cgo

package kobodict

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/geek1011/dictutil/marisa"
)

// MarisaCGO uses the built-in CGO-based libmarisa bindings.
type MarisaCGO struct{}

func init() {
	setMarisa(new(MarisaCGO))
}

func (*MarisaCGO) ReadAll(r io.Reader) (wd []string, err error) {
	defer func() {
		if err := recover(); err != nil {
			wd = nil
			err = fmt.Errorf("marisa (CGO): %v", err)
		}
	}()

	tf, err := ioutil.TempFile("", "marisa")
	if err != nil {
		return nil, fmt.Errorf("marisa (CGO): create temp file: %w", err)
	}
	defer os.Remove(tf.Name())
	defer tf.Close()

	if _, err := io.Copy(tf, r); err != nil {
		return nil, fmt.Errorf("marisa (CGO): write trie to temp file: %w", err)
	}

	// based on marisa-dump
	trie := marisa.NewTrie()
	defer marisa.DeleteTrie(trie)

	trie.Load(tf.Name())

	agent := marisa.NewAgent()
	defer marisa.DeleteAgent(agent)

	agent.SetQueryString("")

	for trie.PredictiveSearch(agent) {
		wd = append(wd, agent.Key().Str())
	}

	return wd, nil
}

func (*MarisaCGO) WriteAll(w io.Writer, wd []string) (err error) {
	defer func() {
		if err := recover(); err != nil {
			err = fmt.Errorf("marisa (CGO): %v", err)
		}
	}()

	ks := marisa.NewKeyset()
	for _, v := range wd {
		ks.PushBackString(v)
	}
	defer marisa.DeleteKeyset(ks)

	td, err := ioutil.TempDir("", "marisa")
	if err != nil {
		return fmt.Errorf("marisa (CGO): create temp dir: %w", err)
	}
	defer os.RemoveAll(td)

	trie := marisa.NewTrie()
	defer marisa.DeleteTrie(trie)

	trie.Build(ks)
	trie.Save(filepath.Join(td, "trie"))

	f, err := os.Open(filepath.Join(td, "trie"))
	if err != nil {
		return fmt.Errorf("marisa (CGO): read marisa output: %w", err)
	}
	defer f.Close()

	if _, err := io.Copy(w, f); err != nil {
		return fmt.Errorf("marisa (CGO): copy marisa output: %w", err)
	}

	return err
}
