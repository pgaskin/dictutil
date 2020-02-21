package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"unicode"

	"gopkg.in/yaml.v2"
)

// GOTDict represents the GOTDict.
type GOTDict []*GOTDef

// GOTDef represents a definition.
type GOTDef struct {
	// Title is the main title of the definition (it may contain spaces) (i.e. Tyrion Lannister).
	Title string
	// Terms are other forms of the title which should be recognized.
	Terms []string
	// Type is the record type.
	Type GOTType
	// Images contains referenced image files.
	Images map[string][]byte
	// Definition contains the Markdown definition.
	Definition string
}

// GOTType is a GOTDict record type.
type GOTType string

const (
	// GOTTypeUnknown is used for definitions without a type set (i.e. before types were used).
	GOTTypeUnknown GOTType = ""
	// GOTTypeCharacter is a character (e.g. Jon, Tyrion).
	GOTTypeCharacter GOTType = "character"
	// GOTTypeHouse is a house (e.g. Lannister, Stark).
	GOTTypeHouse GOTType = "house"
	// GOTTypeEvent is an event in time.
	GOTTypeEvent GOTType = "event"
	// GOTTypeCity is a city.
	GOTTypeCity GOTType = "city"
	// GOTTypeLocation is a location (e.g. King's Landing).
	GOTTypeLocation GOTType = "location"
	// GOTTypeRiver is a river.
	GOTTypeRiver GOTType = "river"
	// GOTTypeShip is a ship.
	GOTTypeShip GOTType = "ship"
	// GOTTypeWord is an uncommon or ASOIAF-specific word.
	GOTTypeWord GOTType = "word"
)

// ParseGOTDict parses the GOTDict. If imgdir is an empty string, images are
// removed. If imgref is true, image paths are set to the full filepath rather
// than reading the images to memory.
func ParseGOTDict(defdir, imgdir string, imgref bool) (GOTDict, error) {
	var dict GOTDict

	fis, err := ioutil.ReadDir(defdir)
	if err != nil {
		return nil, err
	}

	seen := map[string]*GOTDef{}
	for _, fi := range fis {
		if filepath.Ext(fi.Name()) != ".mdd" {
			continue
		}

		if fi.Name() == "reach.mdd" {
			continue // FIX: older duplicate of the-reach.mdd
		}

		buf, err := ioutil.ReadFile(filepath.Join(defdir, fi.Name()))
		if err != nil {
			return nil, err
		}

		var obj struct {
			Title string   `yaml:"title"`
			Terms []string `yaml:"terms"`
			Type  GOTType  `yaml:"type"`

			// FIX: some entries have a messed up format
			Term  []string `yaml:"term"`
			House string   `yaml:"house"`
		}

		md, err := unmarshalStrictFrontMatter(buf, &obj)
		if err != nil {
			return nil, fmt.Errorf("parse %s frontmatter: %w", fi.Name(), err)
		} else if obj.Title == "" {
			return nil, fmt.Errorf("parse %s frontmatter: title not set", fi.Name())
		}

		def := &GOTDef{}

		obj.Title = strings.TrimSpace(obj.Title)
		if odef, ok := seen[obj.Title]; ok {
			return nil, fmt.Errorf("parse %s: already seen %#v in other def %#v", fi.Name(), def.Title, odef)
		}
		seen[obj.Title] = def
		def.Title = obj.Title

		for _, terms := range [][]string{obj.Terms, obj.Term} {
			for _, term := range terms {
				term = strings.TrimSpace(term)
				if term == obj.Title {
					continue // FIX: some entries duplicate the title as a term
				} else if term == "Jon Umber" {
					continue // FIX: duplicated terms
				} else if odef, ok := seen[term]; ok {
					return nil, fmt.Errorf("parse %s: already seen term %#v in other def %#v", fi.Name(), term, odef)
				}
				seen[term] = def
				def.Terms = append(def.Terms, term)
			}
		}

		def.Type = GOTType(strings.TrimSpace(string(obj.Type)))
		def.Images = map[string][]byte{}
		def.Definition = string(md)

		if imgdir == "" {
			def.Definition = regexp.MustCompile(`(\s*Map on [Nn]ext [Pp]age\.?)|(\s*\(Map on [Nn]ext [Pp]age\.?\))|(!\[[^]]*\]\([^)]+\))`).ReplaceAllLiteralString(def.Definition, "")
		} else {
			var repl []string
			for _, img := range regexp.MustCompile(`!\[[^]]*\]\((images/)?([^)]+)\)`).FindAllStringSubmatch(def.Definition, -1) {
				if img[1] == "" {
					return nil, fmt.Errorf("parse %s: unknown image path %#v", fi.Name(), img[1])
				}
				fn, err := filepath.Abs(filepath.Join(imgdir, img[2]))
				if err != nil {
					return nil, fmt.Errorf("parse %s: resolve image %#v: %w", fi.Name(), img[1], err)
				}
				if imgref {
					if _, err := os.Stat(fn); err != nil {
						return nil, fmt.Errorf("parse %s: stat image %#v: %w", fi.Name(), img[1], err)
					}
					repl = append(repl, "("+img[1]+img[2]+")", "("+fn+")")
				} else {
					imgbuf, err := ioutil.ReadFile(fn)
					if err != nil {
						return nil, fmt.Errorf("parse %s: read image %#v: %w", fi.Name(), img[1], err)
					}
					def.Images[img[2]] = imgbuf
					repl = append(repl, "("+img[1]+img[2]+")", "("+img[2]+")")
				}
			}
			def.Definition = strings.NewReplacer(repl...).Replace(def.Definition)
		}

		def.Definition = strings.TrimSpace(def.Definition)

		dict = append(dict, def)
	}

	sort.Slice(dict, func(i, j int) bool {
		return dict[i].Title < dict[j].Title
	})

	return dict, nil
}

func unmarshalStrictFrontMatter(buf []byte, v interface{}) (content []byte, err error) {
	spl := bytes.SplitN(buf, []byte{'-', '-', '-'}, 3)
	for _, b := range spl[0] {
		if !unicode.IsSpace(rune(b)) {
			return buf, nil
		}
	}
	return spl[2], yaml.UnmarshalStrict(spl[1], v)
}
