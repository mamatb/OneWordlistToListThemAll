#!/usr/bin/env python3

import multiprocessing
from os import path


def is_redundant(wordlist_1:str, wordlist_2:str) -> None:
    """Check if all lines of smaller sorted wordlist are in bigger sorted wordlist."""
    with open(wordlist_1) as wl_smaller, open(wordlist_2) as wl_bigger:
        line_wl_smaller = wl_smaller.readline()
        line_wl_bigger = wl_bigger.readline()
        while len(line_wl_smaller) > 0 and len(line_wl_bigger) > 0:
            sline_wl_smaller = line_wl_smaller.strip('\n')
            sline_wl_bigger = line_wl_bigger.strip('\n')
            if sline_wl_smaller > sline_wl_bigger:
                line_wl_bigger = wl_bigger.readline()
            elif sline_wl_smaller == sline_wl_bigger:
                line_wl_smaller = wl_smaller.readline()
                line_wl_bigger = wl_bigger.readline()
            else:
                break
        if len(line_wl_smaller) == 0:
            print(f'all lines in {wordlist_1} are already included in {wordlist_2}')


def main() -> None:  # pylint: disable=C0116
    wordlists_sorted = [
        'adeptus_mechanicus-unix_graph_32max_sort.txt',
        'antipublic-unix_graph_32max_sort.txt',
        'breach_compilation-unix_graph_32max_sort.txt',
        'crackstation-unix_graph_32max_sort.txt',
        'cyclone-unix_graph_32max_sort.txt',
        'dna-unix_graph_32max_sort.txt',
        'HIBP-unix_graph_32max_sort.txt',
        'probable_wordlists-unix_graph_32max_sort.txt',
        'weakpass_3a-unix_graph_32max_sort.txt',
    ]
    wordlists_sorted.sort(key=path.getsize)
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

