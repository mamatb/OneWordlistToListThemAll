#!/usr/bin/env python3

# OneWordlistToListThemAll is a huge mix of password wordlists, proven to be
# pretty useful to provide some quick hits when cracking several hashes
#
# author - mamatb (t.me/m_amatb)
# location - https://github.com/mamatb/OneWordlistToListThemAll
# style guide - https://google.github.io/styleguide/pyguide.html

# TODO
#
# add module docstring
# add tests using pytest
# deal with SIGTERM in child processes


import multiprocessing
import os


def is_redundant(wl_small_path: str, wl_big_path: str) -> None:
    """Checks if all lines in the small wordlist are included in the big wordlist.

    Args:
        wl_small_path: path of the small sorted wordlist.
        wl_big_path: path of the big sorted wordlist.

    Returns:
        None.
    """
    with open(wl_small_path) as wl_small, open(wl_big_path) as wl_big:
        wl_small_line = wl_small.readline()
        wl_big_line = wl_big.readline()
        while len(wl_small_line) > 0 and len(wl_big_line) > 0:
            wl_small_sline = wl_small_line.strip('\n')
            wl_big_sline = wl_big_line.strip('\n')
            if wl_small_sline > wl_big_sline:
                wl_big_line = wl_big.readline()
            elif wl_small_sline == wl_big_sline:
                wl_small_line = wl_small.readline()
                wl_big_line = wl_big.readline()
            else:
                break
        if len(wl_small_line) == 0:
            print(f'all lines in {wl_small_path} are included in {wl_big_path}')


def main() -> None:  # pylint: disable=C0116
    wordlists_sorted = sorted(
        [wordlist for wordlist in os.listdir() if wordlist.endswith('.txt')],
        key=os.path.getsize,
    )
    wordlists_pairs = []
    for wordlist_small_index, wordlist_small in enumerate(wordlists_sorted):
        for wordlist_big in wordlists_sorted[wordlist_small_index + 1:]:
            wordlists_pairs.append((wordlist_small, wordlist_big))
    with multiprocessing.Pool() as is_redundant_pool:
        is_redundant_pool.starmap(is_redundant, wordlists_pairs)


if __name__ == '__main__':
    main()
