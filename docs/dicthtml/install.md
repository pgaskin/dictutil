---
layout: default
title: Installing custom dictionaries
parent: dicthtml
---

# Installing custom dictionaries
Sideloading custom dictionaries is easy, but slightly finicky.

## Using dictutil
You can easily install dictionaries using dictutil. First, if you are not replacing a built-in dictionary, enable and install the **Enable searches on extra dictionaries** [patch](https://pgaskin.net/kobopatch-patches). Then, follow the [instructions for using the install command](../dictutil/install.html).

You can uninstall custom dictionaries (including reverting overwritten built-in ones) using the [uninstall command](../dictutil/uninstall.html).

## Manual installation
1. Enable and install the **Enable searches on extra dictionaries** [patch](https://pgaskin.net/kobopatch-patches).
2. Copy the dictionary to `KOBOeReader/.kobo/dict/dicthtml-LOCALE.zip`, where **LOCALE** is a string consisting of 2 lowercase alphanumeric characters. It does not have to be a valid locale.
3. If using a a firmware version 4.20.14601 or newer, mark the file as read-only (in Windows Explorer, or `chmod 444 dicthtml-LOCALE.zip`) to prevent nickel from overwriting it during the sync process.
4. If using a firmware version older than 4.20.14601, open `KOBOeReader/.kobo/KoboReader.sqlite` in a SQLite3 editor, and add a row to the Dictionary table with the following values:
    - **Suffix:** `-LOCALE`, where **LOCALE** is the locale code you chose earlier. This is used when constructing filenames.
    - **Name:** `Extra:_LOCALE LABEL`, where **LOCALE** is the locale code you chose earlier, and **LABEL** is a custom label (it can have spaces in it).
    - **Installed:** `true`. This one is self-explanatory.
    - **Size:** `SIZE`, where *SIZE* is the size of the dictzip in bytes. This is displayed in the dictionary settings, but is unused otherwise, so it's fine if it isn't accurate as long as it is a valid number. For built-in dictionaries with `IsSynced` set, it is used to check for updates.
    - **IsSynced:** `false`. This is used to see if the sync process should attempt to sync the specified dictionary. If true, the `Size` column is checked against the expected size of the latest version (from the dictionary download server), and if it does not match, the new dictionary is downloaded over it.
5. Open `KOBOeReader/.kobo/Kobo/Kobo eReader.conf`, and add a line like `ExtraLocales=LOCALE` in the `ApplicationPreferences` section. If it already exists, add your locale code to it and keep the items separated by a comma and a space (e.g. `ExtraLocales=a1, a2`).
6. Eject your eReader and test the dictionary.
    - If the dictionary is unselectable, ensure you followed the steps correctly, especially regarding the locale codes.
    - If the dictionary says that the word wasn't found, or just acts unusually in general, ensure the dictionary file is valid.

## About locale names and patches
The reason why the patch is required is due to a bug in the firmware. When you choose an entry from the dictionary dropdown, it tries to find a locale name matching it (which it uses to construct the filename for the dicthtml). Kobo has a hard-coded list of supported built-in locales, and supports adding extra ones using the **ApplicationPreferences->ExtraLocales** config file option (a comma separated list of locale codes). These locales have an automatically generated name of "Extra: LOCALE".

But, this is where the bug occurs. To support translation dictionaries, the dictionary selector will split the name by spaces, and only check against the first element. This is perfectly fine for one-word locale names (i.e. all the built-in ones) For custom locales, it will try to match **Extra:**, which doesn't exist, so it will default to the English dictionary. Thus, to fix this, the "Extra: " prefix used for the custom locales needs to be changed to one without a space. The patch replaces the space with an underscore. This bug does have one benefit though: since only stuff before the first space is considered, you can have a custom label after it.

## Alternative method
It is also possible to install custom dictionaries by replacing an existing built-in installed dictionary in `KOBOeReader/.kobo/dict`. To prevent it from being overwritten during a sync, set the `IsSynced` column to `false` for it in the DB on firmware versions older than 4.20.14601, otherwise, mark it read-only.

## About changes in firmware 4.20.14601

In short:

- **Same:** Nickel will still attempt to sync all dictionaries, including sideloaded ones, unless IsSynced is false.
- **New:** IsSynced can't be changed anymore due to the dictionary table being removed.
- **New:** Nickel will avoid overwriting dictionary files if they are marked read-only, and will instead write `"dicthtml-LOCALE" marked as read-only.. skipping` to the log in the `sync` category.
- **Same:** Nickel still generates locale names by default with `Extra: LOCALE`.
- **New:** Nickel doesn't read the dictionary table anymore, so the name in it is ignored. In addition, entries in the table won't change anything even if it is still present.
- **New:** The built-in dictionaries are hard-coded, rather than writing them to the db during migrations and reading from it at runtime.
- **Same:** Nickel still has the bug where the locale splitting is messed up, so the `Extra: LOCALE` names are inherently broken.
- **Same:** The matching can be fixed by replacing `Extra: ` with `Extra:_` (or anything not containing Unicode whitespace).
- **New:** The database doesn't need to be changed anymore in addition to the patch, as the names are generated dynamically using the same string.
- **Therefore:** If the dictionary table is present, it can safely be removed.
- **Therefore:** The steps required to install custom dictionaries are now (note that these have already been incorporated into the instructions above, they are just here for convenience):
  - Copy the dictzip and mark it read-only.
  - Add it to ExtraLocales if it is not a built-in locale.
  - Use the patch to replace `Extra: ` in libnickel with any other string (same length or shorter with a null byte at the end), but does not contain a space (` `).

See [#49](https://github.com/geek1011/kobopatch-patches/issues/49) for more information.
