---
layout: default
---

# Installing custom dictionaries
Sideloading custom dictionaries is easy, but slightly finicky.

## Manual installation
1. Enable and install the **Enable searches on extra dictionaries** [patch](https://pgaskin.net/kobopatch-patches).
2. Copy the dictionary to `KOBOeReader/.kobo/dict/dicthtml-LOCALE.zip`, where **LOCALE** is a string consisting of 2 lowercase alphanumeric characters. It does not have to be a valid locale.
3. Open `KOBOeReader/.kobo/KoboReader.sqlite` in a SQLite3 editor, and add a row to the Dictionary table with the following values:
    - **Suffix:** `-LOCALE`, where **LOCALE** is the locale code you chose earlier. This is used when constructing filenames.
    - **Name:** `Extra:_LOCALE LABEL`, where **LOCALE** is the locale code you chose earlier, and **LABEL** is a custom label (it can have spaces in it).
    - **Installed:** `true`. This one is self-explanatory.
    - **Size:** `SIZE`, where *SIZE* is the size of the dictzip in bytes. This is displayed in the dictionary settings, but is unused otherwise, so it's fine if it isn't accurate as long as it is a valid number. For built-in dictionaries, it is used to check for updates.
    - **IsSynced:** `false`. This is used to see if the sync process should attempt to sync the specified dictionary.
4. Open `KOBOeReader/.kobo/Kobo/Kobo eReader.conf`, and add a line like `ExtraLocales=LOCALE` in the `ApplicationPreferences` section. If it already exists, add your locale code to it and keep the items separated by a comma and a space (e.g. `ExtraLocales=a1, a2`).
5. Eject your eReader and test the dictionary.
    - If the dictionary is unselectable, ensure you followed the steps correctly, especially regarding the locale codes.
    - If the dictionary says that the word wasn't found, or just acts unusually in general, ensure the dictionary file is valid.

## Using dictutil
Coming soon.

## About locale names and patches
The reason why the patch is required is due to a bug in the firmware. When you choose an entry from the dictionary dropdown, it tries to find a locale name matching it (which it uses to construct the filename for the dicthtml). Kobo has a hard-coded list of supported built-in locales, and supports adding extra ones using the **ApplicationPreferences->ExtraLocales** config file option (a comma separated list of locale codes). These locales have an automatically generated name of "Extra: LOCALE".

But, this is where the bug occurs. To support translation dictionaries, the dictionary selector will split the name by spaces, and only check against the first element. This is perfectly fine for one-word locale names (i.e. all the built-in ones) For custom locales, it will try to match **Extra:**, which doesn't exist, so it will default to the English dictionary. Thus, to fix this, the "Extra: " prefix used for the custom locales needs to be changed to one without a space. The patch replaces the space with an underscore. This bug does have one benefit though: since only stuff before the first space is considered, you can have a custom label after it.
