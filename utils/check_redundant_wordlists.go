package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"slices"
)

type isRedundantJob struct {
	wlSmallName string
	wlBigName   string
}

type isRedundantResult struct {
	wlSmallName string
	wlBigName   string
	isRedundant bool
}

func readDirExt(name string, ext string) ([]fs.DirEntry, error) {
	var entriesExt []fs.DirEntry
	entries, err := os.ReadDir(name)
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ext {
			entriesExt = append(entriesExt, entry)
		}
	}
	return entriesExt, err
}

func isRedundant(wlSmallName string, wlBigName string) bool {
	var wlSmallScanner, wlBigScanner *bufio.Scanner
	if wlSmall, err := os.Open(wlSmallName); err == nil {
		defer wlSmall.Close()
		wlSmallScanner = bufio.NewScanner(wlSmall)
	} else {
		panic(err)
	}
	if wlBig, err := os.Open(wlBigName); err == nil {
		defer wlBig.Close()
		wlBigScanner = bufio.NewScanner(wlBig)
	} else {
		panic(err)
	}
	wlSmallScan, wlBigScan := wlSmallScanner.Scan(), wlBigScanner.Scan()
	for wlSmallScan && wlBigScan {
		switch bytes.Compare(wlSmallScanner.Bytes(), wlBigScanner.Bytes()) {
		case 1:
			wlBigScan = wlBigScanner.Scan()
		case 0:
			wlSmallScan, wlBigScan = wlSmallScanner.Scan(), wlBigScanner.Scan()
		case -1:
			wlBigScan = false
		}
	}
	return !wlSmallScan
}

func isRedundantWorker(jobs chan isRedundantJob, results chan isRedundantResult) {
	for job := range jobs {
		results <- isRedundantResult{
			wlSmallName: job.wlSmallName,
			wlBigName:   job.wlBigName,
			isRedundant: isRedundant(job.wlSmallName, job.wlBigName),
		}
	}
}

func main() {
	var wordlists []fs.DirEntry
	if entries, err := readDirExt(".", ".txt"); err == nil {
		wordlists = entries
	} else {
		panic(err)
	}
	slices.SortFunc(wordlists, func(a fs.DirEntry, b fs.DirEntry) int {
		if aInfo, aErr := a.Info(); aErr == nil {
			if bInfo, bErr := b.Info(); bErr == nil {
				return int(aInfo.Size() - bInfo.Size())
			} else {
				panic(bErr)
			}
		} else {
			panic(aErr)
		}
	})
	jobsN, workersN := len(wordlists)*(len(wordlists)-1)/2, runtime.NumCPU()
	jobs, results := make(chan isRedundantJob, jobsN), make(chan isRedundantResult, jobsN)
	for i := 0; i < workersN; i++ {
		go isRedundantWorker(jobs, results)
	}
	for wlSmallIndex, wlSmall := range wordlists {
		for _, wlBig := range wordlists[wlSmallIndex+1:] {
			jobs <- isRedundantJob{
				wlSmallName: wlSmall.Name(),
				wlBigName:   wlBig.Name(),
			}
		}
	}
	close(jobs)
	for i := 0; i < jobsN; i++ {
		result := <-results
		if result.isRedundant {
			fmt.Println(result.wlSmallName + " is redundant with " + result.wlBigName)
		}
	}
	close(results)
}
