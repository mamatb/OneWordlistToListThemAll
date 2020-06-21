# OneWordlistToListThemAll

What?
-----

OneWordlistToListThemAll is a huge mix of password wordlists, proven to be pretty useful to provide some quick hits when cracking several hashes. Feel free to hit me up if any mega link in here no longer works.

How?
----

Just filtering and mixing.

1. Make sure the source wordlists are not using DOS/Windows line breaks (CR + LF). No need to look for Mac line breaks as they switched from CR to LF long time ago.
```
dos2unix --force --newfile "${WORDLIST}.txt" "${WORDLIST}-unix.txt"
```
2. Sort each wordlist and remove duplicates, using version sort just makes more sense to me.
```
sort --unique --version-sort --output="${WORDLIST}-unix_sort.txt" "${WORDLIST}-unix.txt"
```
3. Get rid of entries containing non-ascii or non-visible characters (except for the space). I'm aware of the built-in POSIX character class `[:graph:]`, but have decided to keep the space in the charset.
```
LC_ALL='C' grep --text --perl-regexp '^([\x20-\x7E])*$' "${WORDLIST}-unix_sort.txt" > "${WORDLIST}-unix_sort_graph.txt"
```
4. Remove all entries longer than 63 characters. As OneWordlistToListThemAll aims to provide some quick hits, it does not make much sense trying passwords that long.
```
sed --regexp-extended '/.{64,}/d' "${WORDLIST}-unix_sort_graph.txt" > "${WORDLIST}-unix_sort_graph_under64.txt"
```
5. Generate OneWordlistToListThemAll.
```
cat *.txt > 'OneWordlistToListThemAll.tmp'
sort --unique --version-sort --output='OneWordlistToListThemAll.txt' 'OneWordlistToListThemAll.tmp'
```
Sources
-------

wordlist name with mega link | decompressed file size | original source
--- | --- | ---
[adeptus_mechanicus](https://mega.nz/file/wR5FhKCS#PsdGoH44-ofBCSQAKyxURAjX7ttL6KqO34KuSCW80XE) | 1.6 GB | `.dic.7z` files at [adeptus-mechanicus.com](https://www.adeptus-mechanicus.com/codex/hashpass/)
[breach_compilation](https://mega.nz/file/RB4TjCjS#QV8u4vFUGYNswB-xIQt9udKrJ2nC6am2FVCWZWM5xbk) | 3.9 GB | magnet link from a [public gist](https://gist.github.com/scottlinux/9a3b11257ac575e4f71de811322ce6b3)
[crackstation](https://mega.nz/file/5N5BVaiB#DT2fLFRdeHtKjYSNv1X7BOtJLEqyYPvo3_V5hNosSZo) | 12.3 GB | [CrackStation's wordlist](https://crackstation.net/files/crackstation.txt.gz)
[DCHTPass](https://mega.nz/file/YEh10AbS#WwEEgT-26IKmOD53TzNdLMQ0ossv0sw7Qsr6R6ZkXLU) | 23.9 GB | `DCHTPassv1.0.txt` at [weakpass.com](https://weakpass.com/wordlist/1257)
[eap](https://mega.nz/file/JApA1CJS#NJOyMFDO4S61KraiPvaz_eKKsUIsjTsYgiyvg7FYCS8) | 0.3 GB | passwords grabbed while testing WPA/WPA2-MGT fake APs + [HashcatRulesEngine](https://github.com/llamasoft/HashcatRulesEngine) using [OneRuleToRuleThemAll](https://github.com/NotSoSecure/password_cracking_rules/blob/master/OneRuleToRuleThemAll.rule)
[hashesorg2019](https://mega.nz/file/9FwFzAYA#daRmuI84P9UOKKTGdZ4xaJLiXy4ze13w-i4LibljxBk) | 13.2 GB | `hashesorg2019` at [weakpass.com](https://weakpass.com/wordlist/1851)
[kaonashi](https://mega.nz/file/JcwzEKiL#A6dXWlaMZepq9abRmcUHL9LyZOX2F97uo-DVL-6tNck) | 9.4 GB | [mega link](https://mega.nz/#!nWJXzYzS!P1G8HDiMxq5wFaxeWGWx334Wp9wByj5kMEGLZkVX694) from [Kaonashi's repo](https://github.com/kaonashi-passwords/Kaonashi)
[probable_wordlists](https://mega.nz/file/II5zXKhZ#mJaNjRiJbqagRPX36Uj0pG7P--73x7SoQt_axy5UZjw) | 21.1 GB | biggest file included in [Probable-Wordlists' torrent](https://github.com/berzerk0/Probable-Wordlists/tree/master/Real-Passwords/Real-Password-Rev-2-Torrents)
[weakpass_2a](https://mega.nz/file/gZoFiAKT#qIbE2JJbtDIkbEjXqxjRISFtBcXQo11h1Vl1GpeQKC8) | 89.6 GB | `weakpass_2a` at [weakpass.com](https://weakpass.com/wordlist/1919)
 | | 
[OneWordlistToListThemAll](https://mega.nz/file/0U5nQS4Y#UrpqxFWvOntGrsgOeZtWRwt3ZhiG5tRMqddciOx-MR0) | 107.8 GB | nope

Acknowledgements
----------------

I'd like to thank the authors of the source wordlists. As stated before, this repo is just a bunch of filtering and mixing other people's work.
