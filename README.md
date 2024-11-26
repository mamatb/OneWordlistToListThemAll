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
LC_ALL='C' dos2unix --force --newfile "${WORDLIST}.txt" "${WORDLIST}-unix.txt"
```
2. Get rid of passwords containing non-ascii or non-visible characters (except for the space).
```bash
LC_ALL='C' grep --text --extended-regexp '^[[:print:]]*$' "${WORDLIST}-unix.txt" > "${WORDLIST}-unix_print.txt"
```
3. Remove all passwords longer than 32 characters. As OneWordlistToListThemAll aims to provide some quick hits, it does not make much sense trying passwords that long.
```bash
LC_ALL='C' grep --text --invert-match --extended-regexp '.{33}' "${WORDLIST}-unix_print.txt" > "${WORDLIST}-unix_print_32max.txt"
```
4. Remove hash-like passwords that may remain uncracked in the source wordlists.
```bash
LC_ALL='C' grep --text --invert-match --extended-regexp '[[:xdigit:]]{32}' "${WORDLIST}-unix_print_32max.txt" > "${WORDLIST}-unix_print_32max_nohash.txt"
```
5. Sort each source wordlist and remove duplicates.
```bash
LC_ALL='C' sort --unique --output="${WORDLIST}-unix_print_32max_nohash_sort.txt" "${WORDLIST}-unix_print_32max_nohash.txt"
```
6. Generate OneWordlistToListThemAll.
```bash
cat *unix_print_32max_nohash_sort.txt > 'OneWordlistToListThemAll.tmp'
LC_ALL='C' sort --unique --output='OneWordlistToListThemAll.txt' 'OneWordlistToListThemAll.tmp'
```
7. Generate OneWordlistToListThemAll WPA-PSK, at least 8 characters per password.
```bash
LC_ALL='C' grep --text --extended-regexp '.{8}' 'OneWordlistToListThemAll.txt' > 'OneWordlistToListThemAll_WPA-PSK.txt'
```

## Wordlists <a name="wordlists" />

name | size (post filtering) | source
-- | -- | --
adeptus mechanicus | 1.6 GB | `.dic.7z` files from [adeptus-mechanicus.com](https://www.adeptus-mechanicus.com/codex/hashpass/)
anti public | 9.3 GB | `.7z` file from [weakpass.com](https://weakpass.com/wordlists/antipublic_breach)
breach compilation | 3.8 GB | `.gz` file from [weakpass.com](https://weakpass.com/wordlists/breachcompilation.txt)
crackstation | 12.1 GB | `.gz` file from [crackstation.net](https://crackstation.net/crackstation-wordlist-password-cracking-dictionary.htm)
cyclone | 6.2 GB | MediaFire link from [cyclone's repo](https://github.com/cyclone-github/wordlist/tree/master/cyclone_hk_v2)
hashes.org | 14.2 GB | `.7z` file from [weakpass.com](https://weakpass.com/wordlists/hashes.org)
hashkiller | 2.8 GB | `.7z` file from [hashkiller.io](https://hashkiller.io/download)
have i been pwned | 6.1 GB | leaks from [hashes.org](https://temp.hashes.org/leaks.php) (HIBP V1 - V6)
kaonashi | 9.4 GB | MEGA link from [Kaonashi's repo](https://github.com/kaonashi-passwords/Kaonashi/tree/master/wordlists)
password dna | 80.5 KB | `.dict` file from [unix-ninja.com](https://www.unix-ninja.com/p/Password_DNA)
probable-wordlists | 21 GB | `.torrent` file from [Probable-Wordlists' repo](https://github.com/berzerk0/Probable-Wordlists/tree/master/Real-Passwords/Real-Password-Rev-2-Torrents)
rockyou 2024 | 94.9 GB | `.zip` file from [DragonJAR](https://djar.co/rockyou2024zip)
weakpass | 83.3 GB | `.7z` file from [weakpass.com](https://weakpass.com/wordlists/weakpass_4a.txt)
 | | 
[OneWordlistToListThemAll](https://drive.google.com/file/d/1cNm5bqDLvK_STuuBaa1bD2eHlUs7FFJB) | 117 GB | N/A
[OneWordlistToListThemAll WPA-PSK](https://drive.google.com/file/d/1FU0kAVo5j3SXX_f9h8RdAeFjhm78xcjw) | 99.8 GB | N/A

## Acknowledgements <a name="acknowledgements" />

I'd like to thank the authors of the source wordlists. As stated before, this repo is just a bunch of filtering and mixing other people's work.
