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
			return false
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
		var aSize, bSize int64
		if aInfo, err := a.Info(); err == nil {
			aSize = aInfo.Size()
		} else {
			panic(err)
		}
		if bInfo, err := b.Info(); err == nil {
			bSize = bInfo.Size()
		} else {
			panic(err)
		}
		return int(aSize - bSize)
	})
	jobsN, workersN := len(wordlists)*(len(wordlists)-1)/2, runtime.NumCPU()
	jobs, results := make(chan isRedundantJob, jobsN), make(chan isRedundantResult, jobsN)
	for range workersN {
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
	for range jobsN {
		result := <-results
		if result.isRedundant {
			fmt.Println(result.wlSmallName + " is redundant with " + result.wlBigName)
		}
	}
	close(results)
}
