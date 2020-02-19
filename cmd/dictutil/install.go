package main

import (
	"bufio"
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	_ "github.com/mattn/go-sqlite3"

	"github.com/geek1011/koboutils/kobo"
	"github.com/spf13/pflag"
)

var (
	builtinUpdated = "4.19.14123"
	builtinDict    = map[string]string{
		"de":        "Deutsch",
		"de-en":     "Deutsch - English",
		"en-de":     "English - Deutsch",
		"en":        "English",
		"en-es":     "English - Español",
		"en-fr":     "English - Français",
		"en-it":     "English - Italiano",
		"en-ja":     "English - 日本語（ジーニアス）",
		"en-ja-pgs": "English - 日本語（プログレッシブ）",
		"en-nl":     "English - Nederlands",
		"en-pt":     "English - Português",
		"en-tr":     "English - Türkçe",
		"es-en":     "Español - English",
		"es":        "Español",
		"fr-en":     "Français - English",
		"fr":        "Français",
		"it-en":     "Italiano - English",
		"it":        "Italiano",
		"jaxxdjs":   "日本語",
		"nl":        "Nederlands",
		"pt-en":     "Português - English",
		"pt":        "Português",
	}
	builtinSorted []string
)

func findDevice(root string) (string, string, error) {
	if len(root) == 0 {
		kobos, err := kobo.Find()
		if err != nil {
			return "", "", err
		} else if len(kobos) == 0 {
			return "", "", fmt.Errorf("no devices detected")
		}
		root = kobos[0]
	}

	_, version, _, err := kobo.ParseKoboVersion(root)
	if err != nil {
		return "", "", fmt.Errorf("parse Kobo version file for %#v: %w.\n", root, err)
	}
	return root, version, nil
}

func builtinHelp() {
	fmt.Fprintf(os.Stderr, "Built-in Kobo dictionaries (last updated for %s):\n", builtinUpdated)
	for _, loc := range builtinSorted {
		lbl := builtinDict[loc]
		if loc == "en" {
			fmt.Fprintf(os.Stderr, "  %-40s %s\n", "en (dicthtml.zip)", lbl)
		} else {
			fmt.Fprintf(os.Stderr, "  %-40s %s\n", fmt.Sprintf("%s (dicthtml-%s.zip)", loc, loc), lbl)
		}
	}
}

func builtinInit() {
	for k := range builtinDict {
		builtinSorted = append(builtinSorted, k)
	}
	sort.Strings(builtinSorted)
}

// the stuff above is shared with uninstall

func init() {
	commands = append(commands, &command{Name: "install", Short: "I", Description: "Install a dictzip file", Main: installMain})
	builtinInit()
}

func installMain(args []string, fs *pflag.FlagSet) int {
	fs.SortFlags = false
	root := fs.StringP("kobo", "k", "", "KOBOeReader path (default: automatically detected)")
	locale := fs.StringP("locale", "l", "", "Locale name to use (format: ALPHANUMERIC{2}; translation dictionaries are not supported) (default: detected from filename if in format dicthtml-**.zip)")
	name := fs.StringP("name", "n", "", "Custom additional label for dictionary (ignored when replacing built-in dictionaries)")
	builtin := fs.StringP("builtin", "b", "replace", "How to handle built-in locales [replace = replace and prevent from syncing] [ignore = replace and leave syncing as-is]")
	help := fs.BoolP("help", "h", false, "Show this help text")
	fs.Parse(args[1:])

	if *help || fs.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] dictzip\n\nOptions:\n%s\n", args[0], fs.FlagUsages())
		builtinHelp()
		fmt.Fprintf(os.Stderr, "\nNote:\n  If you are not replacing a built-in dictionary, the 'Enable searches on extra\n  dictionaries patch' must be installed, or you will not be able to select\n  your custom dictionary.\n")
		return 0
	}

	if *builtin != "replace" && *builtin != "ignore" {
		fmt.Fprintf(os.Stderr, "Error: invalid built-in dictionary mode %#v, see --help for more details.\n", *builtin)
		return 2
	}

	df, err := os.Open(fs.Args()[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not open specified dictzip: %v.\n", err)
		return 1
	}
	defer df.Close()

	dfi, err := df.Stat()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not stat specified dictzip: %v.\n", err)
		return 1
	}
	dictSize := dfi.Size()

	dictLocale := *locale
	if len(dictLocale) == 0 {
		m := regexp.MustCompile(`^dicthtml-([a-zA-Z0-9]{2})\.zip$`).FindStringSubmatch(filepath.Base(fs.Args()[0]))
		if len(m) == 0 {
			fmt.Fprintf(os.Stderr, "Error: no locale specified, and dictzip name doesn't include one.\n")
			return 1
		}
		dictLocale = m[1]
	}

	if !regexp.MustCompile(`^[a-zA-Z0-9]{2}$`).MatchString(dictLocale) { // this is a bit on the overly safe side, but there's not much harm in it anyways, and it can be loosened if needed
		fmt.Fprintf(os.Stderr, "Error: invalid locale %#v specified.\n", dictLocale)
		return 1
	}

	var dictFilename string
	if dictLocale == "en" {
		dictFilename = "dicthtml.zip"
	} else {
		dictFilename = "dicthtml-" + dictLocale + ".zip"
	}

	kobopath, version, err := findDevice(*root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not detect a Kobo eReader (you can specify one manually with --kobo): %v.\n", err)
		return 1
	}

	fmt.Printf("Found Kobo eReader at %#v with firmware version %s.\n", kobopath, version)
	if kobo.VersionCompare(version, "4.7.10364") < 0 {
		fmt.Fprintf(os.Stderr, "Error: firmware version too old (v2 dictionaries were only introduced in 4.7.10364).\n")
		return 1
	}

	dictName, dictBuiltin := builtinDict[dictLocale]
	if !dictBuiltin {
		dictName = "Extra:_" + dictLocale
		if len(*name) != 0 {
			dictName += " " + *name
		}
	} else if len(*name) != 0 {
		fmt.Fprintf(os.Stderr, "Warning: Ignoring custom dictionary label for built-in dictionary.\n")
	}

	fmt.Printf("Installing dictzip %#v (locale: %s) as %#v (overwriting_builtin: %t) with label %#v.\n\n", fs.Args()[0], dictLocale, dictFilename, dictBuiltin, dictName)

	// TODO: maybe split these functions out and test them?

	fmt.Printf("Copying dictzip.\n")
	if err := func() error {
		dz := filepath.Join(kobopath, ".kobo", "dict", dictFilename)

		dfo, err := os.OpenFile(dz, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer dfo.Close()

		if _, err := io.Copy(dfo, df); err != nil {
			return err
		}

		if err := dfo.Close(); err != nil {
			return err
		}

		fmt.Printf("  Wrote dictzip to %#v.\n", dz)
		return nil
	}(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: copy dictzip: %v.\n", err)
		return 1
	}

	fmt.Printf("Updating ExtraLocales.\n")
	if dictBuiltin {
		fmt.Printf("  No need; replacing built-in dictionary.\n")
	} else {
		if err := func() error {
			cfg := filepath.Join(kobopath, ".kobo", "Kobo", "Kobo eReader.conf")

			f, err := os.OpenFile(cfg, os.O_RDONLY, 0)
			if err != nil {
				return fmt.Errorf("open config file: %w", err)
			}
			defer f.Close()

			var locales []string
			buf := bytes.NewBuffer(nil)

			fs := bufio.NewScanner(f)
			for fs.Scan() {
				if bytes.HasPrefix(fs.Bytes(), []byte("ExtraLocales=")) {
					for _, loc := range strings.Split(strings.SplitN(fs.Text(), "=", 2)[1], ",") {
						locales = append(locales, strings.TrimSpace(loc))
					}
					continue
				}
				_, _ = buf.Write(fs.Bytes()) // err is always nil
				buf.WriteRune('\n')
			}

			var added bool
			for _, loc := range locales {
				if loc == dictLocale {
					added = true
					break
				}
			}
			if added {
				fmt.Printf("  Locale %#v already added to ExtraLocales.\n", dictLocale)
				return nil
			}

			fmt.Printf("  Adding locale %#v to ExtraLocales.\n", dictLocale)
			locales = append(locales, dictLocale)
			sort.Strings(locales)

			buf.WriteString("\n[ApplicationPreferences]\n") // this will get merged by Qt
			buf.WriteString("ExtraLocales=" + strings.Join(locales, ","))

			f.Close()

			fo, err := os.OpenFile(cfg+".tmp", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
			if err != nil {
				return fmt.Errorf("open new config file: %w", err)
			}
			defer os.Remove(cfg + ".tmp")
			defer fo.Close()

			if _, err := fo.Write(buf.Bytes()); err != nil {
				return fmt.Errorf("write new config file: %w", err)
			}

			if err := fo.Sync(); err != nil {
				return fmt.Errorf("write new config file: %w", err)
			}

			if err := fo.Close(); err != nil {
				return fmt.Errorf("write new config file: %w", err)
			}

			if err := os.Rename(cfg+".tmp", cfg); err != nil {
				return fmt.Errorf("rename new config file: %w", err)
			}

			return nil
		}(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: update ExtraLocales: %v.\n", err)
			return 1
		}
	}

	fmt.Printf("Updating database.\n")
	if err := func() error {
		db, err := sql.Open("sqlite3", filepath.Join(kobopath, ".kobo", "KoboReader.sqlite"))
		if err != nil {
			return fmt.Errorf("open database: %w", err)
		}
		defer db.Close()

		rSuffix := "-" + dictLocale
		rName := dictName
		rInstalled := "true"
		rSize := dictSize
		rIsSynced := "false"
		if dictBuiltin && *builtin == "ignore" {
			rIsSynced = "true"
		}

		if _, err := db.Exec("INSERT OR REPLACE INTO Dictionary (Suffix, Name, Installed, Size, IsSynced) VALUES (?, ?, ?, ?, ?)", rSuffix, rName, rInstalled, rSize, rIsSynced); err != nil {
			return fmt.Errorf("update database: %w", err)
		}
		fmt.Printf("  Added row to database: Suffix=%#v Name=%#v Installed=%#v Size=%#v IsSynced=%#v.\n", rSuffix, rName, rInstalled, rSize, rIsSynced)

		if err := db.Close(); err != nil {
			return fmt.Errorf("close database: %w", err)
		}

		return nil
	}(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: update database: %v.\n", err)
		return 1
	}

	fmt.Printf("\nSuccessfully installed dictzip %#v to Kobo %#v.\n", fs.Args()[0], kobopath)

	return 0
}
