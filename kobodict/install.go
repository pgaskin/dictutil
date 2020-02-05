package kobodict

// TODO(maybe): func Install(kobopath, dictzip string, locale, label string) error
//       which will:
//       - add locale to config ApplicationPreferences/ExtraLocales
//       - copy and rename dict to .kobo/dict/dicthtml*.zip
//       - add row to db dictionary table ("-{locale}", os.Stat(dictzip).Size(), "Extra:_ {label}", true, true)
//       - prompt about installing patch
