package marisa

import (
	"sort"
	"testing"
)

func TestCommonPrefixSearch(t *testing.T) {
	keyset := NewKeyset()
	keyset.PushBackStringWithWeight("a", 2)
	keyset.PushBackString("app")
	keyset.PushBackString("apple")

	trie := NewTrie()
	trie.Build(keyset)

	agent := NewAgent()
	agent.SetQueryString("apple")

	tcs := []struct {
		prefix string
		id     int64
	}{
		{"a", 0},
		{"app", 1},
		{"apple", 2},
	}

	var i int
	for i = 0; i < len(tcs) && trie.CommonPrefixSearch(agent); i++ {
		key, tc := agent.Key(), tcs[i]

		if prefix := key.Str(); prefix != tc.prefix {
			t.Errorf("expected prefix %s, got %s", tc.prefix, prefix)
		} else if id := key.Id(); id != tc.id {
			t.Errorf("expected id %d, got %d", tc.id, id)
		}
	}
	if i != len(tcs) {
		t.Errorf("got %d prefixes, expected %d", i, len(tcs))
	}
}

func TestPredictiveSearch(t *testing.T) {
	ks := NewKeyset()
	for _, v := range []string{"foo", "foobar", "foobaz", "abcdef"} {
		ks.PushBackString(v)
	}

	tr := NewTrie()
	tr.Build(ks)

	for _, tc := range []struct {
		in  string
		out []string
	}{
		{"fo", []string{"foo", "foobar", "foobaz"}},
		{"abcd", []string{"abcdef"}},
		{"123", []string{}},
		{"", []string{"abcdef", "foo", "foobar", "foobaz"}},
	} {
		agent := NewAgent()
		agent.SetQueryString(tc.in)

		var found []string
		for tr.PredictiveSearch(agent) {
			found = append(found, agent.Key().Str())
		}

		sort.Strings(found)
		sort.Strings(tc.out)

		var diff bool
		if len(found) != len(tc.out) {
			diff = true
		} else {
			for i := range found {
				if found[i] != tc.out[i] {
					diff = true
					break
				}
			}
		}

		if diff {
			t.Errorf("%#v: expected %#v, got %#v", tc.in, tc.out, found)
		}
	}
}
