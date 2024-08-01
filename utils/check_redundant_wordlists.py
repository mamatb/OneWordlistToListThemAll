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


import multiprocessing
import os


def is_redundant(wordlist_small: str, wordlist_big: str) -> None:
    """Checks if all lines in wordlist_small are included in wordlist_big.

    Args:
        wordlist_small: sorted wordlist, smaller than wordlist_big.
        wordlist_big: sorted wordlist, bigger than wordlist_small.

    Returns:
        None.
    """
    with open(wordlist_small) as wl_small, open(wordlist_big) as wl_big:
        line_wl_small = wl_small.readline()
        line_wl_big = wl_big.readline()
        while len(line_wl_small) > 0 and len(line_wl_big) > 0:
            sline_wl_small = line_wl_small.strip('\n')
            sline_wl_big = line_wl_big.strip('\n')
            if sline_wl_small > sline_wl_big:
                line_wl_big = wl_big.readline()
            elif sline_wl_small == sline_wl_big:
                line_wl_small = wl_small.readline()
                line_wl_big = wl_big.readline()
            else:
                break
        if len(line_wl_small) == 0:
            print(f'all lines in {wordlist_small} are included in {wordlist_big}')


def main() -> None:  # pylint: disable=C0116
    wordlists_sorted = [wl for wl in os.listdir() if wl.endswith('.txt')]
    wordlists_sorted.sort(key=os.path.getsize)
    for wordlist_small_index, wordlist_small in enumerate(wordlists_sorted):
        for wordlist_big in wordlists_sorted[wordlist_small_index + 1:]:
            multiprocessing.Process(
                target=is_redundant,
                args=(wordlist_small, wordlist_big),
            ).start()
    for child in multiprocessing.active_children():
        child.join()


if __name__ == '__main__':
    main()
