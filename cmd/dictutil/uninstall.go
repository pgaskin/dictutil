package main

import (
	"bufio"
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/geek1011/koboutils/v2/kobo"
	"github.com/spf13/pflag"
)

func init() {
	commands = append(commands, &command{Name: "uninstall", Short: "U", Description: "Uninstall a dictzip file", Main: uninstallMain})
}

func uninstallMain(args []string, fs *pflag.FlagSet) int {
	fs.SortFlags = false
	root := fs.StringP("kobo", "k", "", "KOBOeReader path (default: automatically detected)")
	builtin := fs.StringP("builtin", "b", "normal", "How to handle built-in locales [normal = uninstall the same way as the UI] [delete = completely delete the entry (doesn't have any effect on 4.20.14601+)] [restore = download the original dictionary from Kobo again]")
	help := fs.BoolP("help", "h", false, "Show this help text")
	fs.Parse(args[1:])

	if *help || fs.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] locale\n\nOptions:\n%s\n", args[0], fs.FlagUsages())
		builtinHelp()
		return 0
	}

	if *builtin != "normal" && *builtin != "delete" && *builtin != "restore" {
		fmt.Fprintf(os.Stderr, "Error: invalid built-in dictionary mode %#v, see --help for more details.\n", *builtin)
		return 2
	}

	kobopath, version, err := findDevice(*root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not detect a Kobo eReader (you can specify one manually with --kobo): %v.\n", err)
		return 1
	}

	fmt.Printf("Found Kobo eReader at %s with firmware version %s.\n", kobopath, version)
	if kobo.VersionCompare(version, "4.7.10364") < 0 {
		fmt.Fprintf(os.Stderr, "Error: firmware version too old (v2 dictionaries were only introduced in 4.7.10364).\n")
		return 1
	}
	newMethod := kobo.VersionCompare(version, "4.20.14601") >= 0 // https://github.com/geek1011/kobopatch-patches/issues/49

	var dictPath, dictLocale string
	if dictLocale = strings.TrimLeft(fs.Args()[0], "-"); dictLocale == "en" {
		dictPath = filepath.Join(kobopath, ".kobo", "dict", "dicthtml.zip")
	} else if regexp.MustCompile(`^[a-zA-Z0-9-]+$`).MatchString(dictLocale) {
		dictPath = filepath.Join(kobopath, ".kobo", "dict", "dicthtml-"+dictLocale+".zip")
	} else {
		fmt.Fprintf(os.Stderr, "Error: invalid locale name.\n")
		return 1
	}
	dictSuffix := "-" + dictLocale
	_, dictBuiltin := builtinDict[dictLocale]

	fmt.Printf("Uninstalling dictionary %#v (locale: %s).\n\n", dictPath, dictLocale)

	fmt.Printf("Updating database.\n")
	if err := func() error {
		db, err := sql.Open("sqlite3", filepath.Join(kobopath, ".kobo", "KoboReader.sqlite"))
		if err != nil {
			return fmt.Errorf("open database: %w", err)
		}
		defer db.Close()

		if exists, err := func() (bool, error) {
			res, err := db.Query(`SELECT name FROM sqlite_master WHERE type="table" AND name="Dictionary";`)
			if err != nil {
				return false, fmt.Errorf("check dictionary table: %w", err)
			}
			defer res.Close()

			if !res.Next() { // if no rows are returned, there was an error or the table didn't exist
				if err := res.Err(); err != nil {
					return false, fmt.Errorf("check dictionary table: %w", err)
				}
				return false, nil
			}
			return true, nil
		}(); err != nil {
			return fmt.Errorf("check dictionary table: %w", err)
		} else if exists {
			if newMethod {
				fmt.Printf("  Note: the dictionary table is unnecessary and inconsequential in firmware 4.20.14601+ and can be safely removed.\n")
			}
		} else {
			if newMethod {
				// show a message to prevent confusion
				fmt.Printf("  No need to update dictionary table on 4.20.14601, skipping.\n")
				return nil
			} else {
				return fmt.Errorf("check dictionary table: not found, and version < 4.20.14123")
			}
		}

		if !dictBuiltin || *builtin == "delete" {
			if res, err := db.Exec("DELETE FROM Dictionary WHERE Suffix = ?", dictSuffix); err != nil {
				return fmt.Errorf("delete row from database: %w", err)
			} else if ra, _ := res.RowsAffected(); ra == 0 {
				fmt.Printf("  Row already removed from database (suffix=%s).\n", dictSuffix)
			} else {
				fmt.Printf("  Removed row from database (suffix=%s).\n", dictSuffix)
			}
		}

		if dictBuiltin && *builtin == "normal" {
			if _, err := db.Exec("UPDATE Dictionary SET Installed = ? WHERE Suffix = ?", "false", dictSuffix); err != nil {
				return fmt.Errorf("update row in database: %w", err)
			} else {
				fmt.Printf("  Set IsInstalled to false in database for built-in dictionary (suffix=%s).\n", dictSuffix)
			}
		}

		if dictBuiltin && *builtin == "restore" {
			if _, err := db.Exec("UPDATE Dictionary SET Installed = ? WHERE Suffix = ?", "true", dictSuffix); err != nil {
				return fmt.Errorf("update row in database: %w", err)
			} else {
				fmt.Printf("  Set IsInstalled to true in database for built-in dictionary (suffix=%s).\n", dictSuffix)
			}
		}

		if err := db.Close(); err != nil {
			return fmt.Errorf("close database: %w", err)
		}

		return nil
	}(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: update database: %v.\n", err)
		return 1
	}

	fmt.Printf("Updating ExtraLocales.\n")
	if dictBuiltin {
		fmt.Printf("  No need; built-in dictionary.\n")
	} else {
		if err := func() error {
			cfg := filepath.Join(kobopath, ".kobo", "Kobo", "Kobo eReader.conf")

			f, err := os.OpenFile(cfg, os.O_RDONLY, 0)
			if err != nil {
				return fmt.Errorf("open config file: %w", err)
			}
			defer f.Close()

			var locales []string
			var filtered bool
			buf := bytes.NewBuffer(nil)

			fs := bufio.NewScanner(f)
			for fs.Scan() {
				if bytes.HasPrefix(fs.Bytes(), []byte("ExtraLocales=")) {
					for _, loc := range strings.Split(strings.SplitN(fs.Text(), "=", 2)[1], ",") {
						loc = strings.TrimSpace(loc)
						if loc == dictLocale {
							filtered = true
						} else {
							locales = append(locales, loc)
						}
					}
					continue
				}
				_, _ = buf.Write(fs.Bytes()) // err is always nil
				buf.WriteRune('\n')
			}

			if !filtered {
				fmt.Printf("  Locale %#v already removed from ExtraLocales.\n", dictLocale)
				return nil
			}

			fmt.Printf("  Removing locale %#v from ExtraLocales.\n", dictLocale)
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

	fmt.Printf("Removing dictzip.\n")
	if err := os.Remove(dictPath); os.IsNotExist(err) { // this will still remove it if it's readonly on Windows (golang/go@2ffb3e5d905b5622204d199128dec06cefd57790)
		fmt.Printf("  Already removed.\n")
	} else if err != nil {
		fmt.Fprintf(os.Stderr, "Error: remove dictzip: %v.\n", err)
		return 1
	} else {
		fmt.Printf("  Removed.\n")
	}

	if *builtin == "restore" {
		// TODO: reconsider whether this belongs in uninstall, as:
		//  - This doesn't update the file size.
		//  - This doesn't ensure there is actually a DB entry for the restored
		//    dict.
		//  - This isn't really uninstalling.
		//  - It might not even belong in dictutil at all because the URLs may
		//    change (and it isn't that hard to manually download a dictionary
		//    to install it with dictutil install)
		url := "https://kbdownload1-a.akamaihd.net/ereader/dictionaries/v2/" + filepath.Base(dictPath)
		fmt.Printf("Restoring original dictionary from %#v.\n", url)

		if err := func() error {
			resp, err := http.Get(url)
			if err != nil {
				return fmt.Errorf("get dictionary: %w", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("get dictionary: response status %s", resp.Status)
			}

			df, err := os.OpenFile(dictPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
			if err != nil {
				return fmt.Errorf("open output dictzip: %w", err)
			}
			defer df.Close()

			if _, err := io.Copy(df, resp.Body); err != nil {
				return fmt.Errorf("write output dictzip: %w", err)
			}

			if err := df.Close(); err != nil {
				return fmt.Errorf("write output dictzip: %w", err)
			}

			return nil
		}(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: download dictionary: %v.\n", err)
			return 1
		}
	}

	fmt.Printf("\nSuccessfully uninstalled dictionary for locale %s.\n", dictLocale)

	return 0
}
