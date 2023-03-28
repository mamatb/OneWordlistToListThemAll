# OneWordlistToListThemAll

* [What?](#what)
* [How?](#how)
* [Wordlists](#wordlists)
* [Acknowledgements](#acknowledgements)

## What? <a name="what" />

OneWordlistToListThemAll is a huge mix of password wordlists, proven to be pretty useful to provide some quick hits when cracking several hashes. Feel free to hit me up if any link in here no longer works.

## How? <a name="how" />

Just filtering and mixing.

1. Make sure the source wordlists are not using DOS/Windows line breaks (CR + LF). No need to look for Mac line breaks as they switched from CR to LF long time ago.
```bash
dos2unix --force --newfile "${WORDLIST}.txt" "${WORDLIST}-unix.txt"
```
2. Get rid of passwords containing non-ascii or non-visible characters (except for the space). I'm aware of the built-in POSIX character class `[:graph:]`, but have decided to keep the space in the charset.
```bash
LC_ALL='C' grep --text --perl-regexp '^([\x20-\x7E])*$' "${WORDLIST}-unix.txt" > "${WORDLIST}-unix_graph.txt"
```
3. Remove all passwords longer than 32 characters. As OneWordlistToListThemAll aims to provide some quick hits, it does not make much sense trying passwords that long.
```bash
sed --regexp-extended '/.{33,}/d' "${WORDLIST}-unix_graph.txt" > "${WORDLIST}-unix_graph_32max.txt"
```
4. Sort each wordlist and remove duplicates.
```bash
sort --unique --output="${WORDLIST}-unix_graph_32max_sort.txt" "${WORDLIST}-unix_graph_32max.txt"
```
5. Generate OneWordlistToListThemAll.
```bash
cat *unix_graph_32max_sort.txt > 'OneWordlistToListThemAll.tmp'
sort --unique --output='OneWordlistToListThemAll.txt' 'OneWordlistToListThemAll.tmp'
```
## Wordlists <a name="wordlists" />

name | size (post filtering) | source
-- | -- | --
Adeptus Mechanicus | 1.6 GB | `.dic.7z` files from [adeptus-mechanicus.com](https://www.adeptus-mechanicus.com/codex/hashpass/)
Breach Compilation | 3.8 GB | magnet link from [GitHub Gist](https://gist.github.com/scottlinux/9a3b11257ac575e4f71de811322ce6b3)
CrackStation | 12.1 GB | [CrackStation's wordlist](https://crackstation.net/files/crackstation.txt.gz)
Cyclone | 6.2 GB | anonfiles link from [cyclone's repo](https://github.com/cyclone-github/wordlist/tree/master/cyclone_hk_v2)
EAP | 7 KB | passwords grabbed while testing WPA/WPA2-MGT fake APs
Hashes.org | 14.2 GB | `.7z` file from [weakpass.com](https://weakpass.com/wordlist/1931)
Hashkiller.io | 2.8 GB | [Hashkiller.io's wordlist](https://hashkiller.io/download)
Have I Been Pwned | 6.2 GB | leaks from [hashes.org](https://temp.hashes.org/leaks.php) (HIBP V1 - V6)
Kaonashi | 9.4 GB | MEGA link from [Kaonashi's repo](https://github.com/kaonashi-passwords/Kaonashi/tree/master/wordlists)
Password DNA | 80.5 KB | `.dict` file from [unix-ninja's post](https://www.unix-ninja.com/p/Password_DNA)
Probable Wordlists | 21 GB | `.torrent` file from [Probable-Wordlists' repo](https://github.com/berzerk0/Probable-Wordlists/tree/master/Real-Passwords/Real-Password-Rev-2-Torrents)
RockYou | 136.2 MB | `.tar.gz` file from [SecLists' repo](https://github.com/danielmiessler/SecLists/tree/master/Passwords/Leaked-Databases)
RockYou 2021 | 90.6 GB | `.7z` file from [weakpass.com](https://weakpass.com/wordlist/1943)
Weakpass | 100 GB | `.7z` file from [weakpass.com](https://weakpass.com/wordlist/1948)
 | | 
OneWordlistToListThemAll ([OneWordlistToListThemAll.7z.001](https://anonfiles.com/Lfv2b7h8z4/OneWordlistToListThemAll_7z_001) + [OneWordlistToListThemAll.7z.002](https://anonfiles.com/3dq8bah2z5/OneWordlistToListThemAll_7z_002)) | 102.6 GB | N/A
OneWordlistToListThemAll WPA-PSK, at least 8 characters per password ([OneWordlistToListThemAll_WPA-PSK.7z.001](https://anonfiles.com/lcZcg6hbz4/OneWordlistToListThemAll_WPA_PSK_7z_001) + [OneWordlistToListThemAll_WPA-PSK.7z.002](https://anonfiles.com/9dNeg3hfz9/OneWordlistToListThemAll_WPA_PSK_7z_002)) | 94.1 GB | N/A

## Acknowledgements <a name="acknowledgements" />

I'd like to thank the authors of the source wordlists. As stated before, this repo is just a bunch of filtering and mixing other people's work.
