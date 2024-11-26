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


def is_redundant(wl_small_path: str, wl_big_path: str) -> tuple[str, str, bool]:
    """Checks if a small sorted wordlist is redundant with a big sorted wordlist.

    Args:
        wl_small_path: path of the small sorted wordlist.
        wl_big_path: path of the big sorted wordlist.

    Returns:
        whether the small sorted wordlist is redundant with the big sorted wordlist.
    """
    with open(wl_small_path) as wl_small, open(wl_big_path) as wl_big:
        wl_small_line = wl_small.readline()
        wl_big_line = wl_big.readline()
        while len(wl_small_line) > 0 and len(wl_big_line) > 0:
            wl_small_sline = wl_small_line.removesuffix('\n')
            wl_big_sline = wl_big_line.removesuffix('\n')
            if wl_small_sline > wl_big_sline:
                wl_big_line = wl_big.readline()
            elif wl_small_sline == wl_big_sline:
                wl_small_line = wl_small.readline()
                wl_big_line = wl_big.readline()
            else:
                break
        return (wl_small_path, wl_big_path, len(wl_small_line) == 0)


def main() -> None:  # pylint: disable=C0116
    wordlists_sorted = sorted(
        [wordlist for wordlist in os.listdir() if wordlist.endswith('.txt')],
        key=os.path.getsize,
    )
    with multiprocessing.Pool() as workers, multiprocessing.Manager() as shared:
        jobs_n, results = 0, shared.Queue()
        for wordlist_small_index, wordlist_small in enumerate(wordlists_sorted):
            for wordlist_big in wordlists_sorted[wordlist_small_index + 1:]:
                workers.apply_async(
                    is_redundant,
                    (wordlist_small, wordlist_big),
                    callback=results.put,
                )
                jobs_n += 1
        for _ in range(jobs_n):
            result = results.get()
            if result[-1]:
                print(f'{result[0]} is redundant with {result[1]}')


if __name__ == '__main__':
    main()
