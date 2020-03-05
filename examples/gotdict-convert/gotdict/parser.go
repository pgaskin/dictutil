// Package gotdict parses GOTDict (https://github.com/wjdp/gotdict).
package gotdict

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

// Dict represents the Dict.
type Dict []*Def

// Def represents a definition.
type Def struct {
	// Title is the main title of the definition (it may contain spaces) (i.e. Tyrion Lannister).
	Title string
	// Terms are other forms of the title which should be recognized.
	Terms []string
	// Type is the record type. Currently, not many entries have one.
	Type Type
	// Images contains referenced image files.
	Images map[string][]byte
	// Definition contains the Markdown definition.
	Definition string
}

// Type is a Dict record type.
type Type string

const (
	// TypeUnknown is used for definitions without a type set (i.e. before types were used).
	TypeUnknown Type = ""
	// TypeCharacter is a character (e.g. Jon, Tyrion).
	TypeCharacter Type = "character"
	// TypeHouse is a house (e.g. Lannister, Stark).
	TypeHouse Type = "house"
	// TypeEvent is an event in time.
	TypeEvent Type = "event"
	// TypeCity is a city.
	TypeCity Type = "city"
	// TypeLocation is a location (e.g. King's Landing).
	TypeLocation Type = "location"
	// TypeRiver is a river.
	TypeRiver Type = "river"
	// TypeShip is a ship.
	TypeShip Type = "ship"
	// TypeWord is an uncommon or ASOIAF-specific word.
	TypeWord Type = "word"
)

// Parse parses the Dict. If imgdir is an empty string, images are removed. If
// imgref is true, image paths are set to the full filepath rather than reading
// the images to memory.
func Parse(defdir, imgdir string, imgref bool) (Dict, error) {
	var dict Dict

	fis, err := ioutil.ReadDir(defdir)
	if err != nil {
		return nil, err
	}

	seen := map[string]*Def{}
	for _, fi := range fis {
		if filepath.Ext(fi.Name()) != ".mdd" {
			continue
		}

		buf, err := ioutil.ReadFile(filepath.Join(defdir, fi.Name()))
		if err != nil {
			return nil, err
		}

		var obj struct {
			Title string   `yaml:"title"`
			Terms []string `yaml:"terms"`
			Type  Type     `yaml:"type"`
		}

		md, err := unmarshalStrictFrontMatter(buf, &obj)
		if err != nil {
			return nil, fmt.Errorf("parse %s frontmatter: %w", fi.Name(), err)
		} else if obj.Title == "" {
			return nil, fmt.Errorf("parse %s frontmatter: title not set", fi.Name())
		}

		def := &Def{}

		obj.Title = strings.TrimSpace(obj.Title)
		if odef, ok := seen[obj.Title]; ok {
			return nil, fmt.Errorf("parse %s: already seen %#v in other def %#v", fi.Name(), def.Title, odef)
		}
		seen[obj.Title] = def
		def.Title = obj.Title

		for _, term := range obj.Terms {
			term = strings.TrimSpace(term)
			if odef, ok := seen[term]; ok && term != "Jon Umber" { // it's usually a mistake to have duplicate terms (but remember that dictgen will handle them fine)
				return nil, fmt.Errorf("parse %s: already seen term %#v in other def %#v", fi.Name(), term, odef)
			}
			seen[term] = def
			def.Terms = append(def.Terms, term)
		}

		def.Type = Type(strings.TrimSpace(string(obj.Type)))
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
