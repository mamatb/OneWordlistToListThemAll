# OneWordlistToListThemAll

What?
-----

OneWordlistToListThemAll is a huge mix of password wordlists, proven to be pretty useful to provide some quick hits when cracking several hashes. Feel free to hit me up if any link in here no longer works.

How?
----

Just filtering and mixing.

1. Make sure the source wordlists are not using DOS/Windows line breaks (CR + LF). No need to look for Mac line breaks as they switched from CR to LF long time ago.
```bash
dos2unix --force --newfile "${WORDLIST}.txt" "${WORDLIST}-unix.txt"
```
2. Sort each wordlist and remove duplicates, using version sort just makes more sense to me.
```bash
sort --unique --version-sort --output="${WORDLIST}-unix_sort.txt" "${WORDLIST}-unix.txt"
```
3. Get rid of entries containing non-ascii or non-visible characters (except for the space). I'm aware of the built-in POSIX character class `[:graph:]`, but have decided to keep the space in the charset.
```bash
LC_ALL='C' grep --text --perl-regexp '^([\x20-\x7E])*$' "${WORDLIST}-unix_sort.txt" > "${WORDLIST}-unix_sort_graph.txt"
```
4. Remove all entries longer than 63 characters. As OneWordlistToListThemAll aims to provide some quick hits, it does not make much sense trying passwords that long.
```bash
sed --regexp-extended '/.{64,}/d' "${WORDLIST}-unix_sort_graph.txt" > "${WORDLIST}-unix_sort_graph_under64.txt"
```
5. Generate OneWordlistToListThemAll.
```bash
cat *unix_sort_graph_under64.txt > 'OneWordlistToListThemAll.tmp'
sort --unique --version-sort --output='OneWordlistToListThemAll.txt' 'OneWordlistToListThemAll.tmp'
```
Wordlists
-------

name | size | source
--- | --- | ---
Adeptus Mechanicus | 1.6 GB | `.dic.7z` files at [adeptus-mechanicus.com](https://www.adeptus-mechanicus.com/codex/hashpass/)
Breach Compilation | 3.9 GB | magnet link from a [public gist](https://gist.github.com/scottlinux/9a3b11257ac575e4f71de811322ce6b3)
CrackStation | 12.3 GB | [CrackStation's wordlist](https://crackstation.net/files/crackstation.txt.gz)
EAP | 0.4 GB | passwords grabbed while testing WPA/WPA2-MGT fake APs + [HashcatRulesEngine](https://github.com/llamasoft/HashcatRulesEngine) using [OneRuleToRuleThemAll](https://github.com/NotSoSecure/password_cracking_rules/blob/master/OneRuleToRuleThemAll.rule)
Hashes.org | 14.5 GB | `Hashes.org` at [weakpass.com](https://weakpass.com/wordlist/1931)
Hashkiller.io | 2.8 GB | [Hashkiller.io's wordlist](https://hashkiller.io/downloads/hashkiller-dict-2020-01-26.7z)
Have I Been Pwned | 6.6 GB | `Have I Been Pwned` leaks at [hashes.org](https://temp.hashes.org/leaks.php) (V1 - V6)
Kaonashi | 9.4 GB | [mega link](https://mega.nz/#!nWJXzYzS!P1G8HDiMxq5wFaxeWGWx334Wp9wByj5kMEGLZkVX694) from [Kaonashi's repo](https://github.com/kaonashi-passwords/Kaonashi)
Probable Wordlists | 21.1 GB | biggest file included in [Probable-Wordlists' torrent](https://github.com/berzerk0/Probable-Wordlists/tree/master/Real-Passwords/Real-Password-Rev-2-Torrents)
Weakpass | 100.7 GB | `weakpass_3a` at [weakpass.com](https://weakpass.com/wordlist/1948)
 | | 
[OneWordlistToListThemAll](https://mega.nz/file/5UhlnQ5b#XOeSN-vz2c5O_B0oOxetgl2_2qwaXLUbTTSFRMcTb6k) | 103.7 GB | N/A

Acknowledgements
----------------

I'd like to thank the authors of the source wordlists. As stated before, this repo is just a bunch of filtering and mixing other people's work.
