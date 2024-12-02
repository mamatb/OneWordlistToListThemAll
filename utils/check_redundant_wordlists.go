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

func errCheck(err error) {
	if err != nil {
		panic(err)
	}
}

func readDirExt(name string, ext string) ([]fs.DirEntry, error) {
	entries, err := os.ReadDir(name)
	var entriesExt []fs.DirEntry
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ext {
			entriesExt = append(entriesExt, entry)
		}
	}
	return entriesExt, err
}

func isRedundant(wlSmallName string, wlBigName string) bool {
	wlSmall, err := os.Open(wlSmallName)
	errCheck(err)
	defer wlSmall.Close()
	wlBig, err := os.Open(wlBigName)
	errCheck(err)
	defer wlBig.Close()
	wlSmallScanner, wlBigScanner := bufio.NewScanner(wlSmall), bufio.NewScanner(wlBig)
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
	wordlists, err := readDirExt(".", ".txt")
	errCheck(err)
	slices.SortFunc(wordlists, func(a fs.DirEntry, b fs.DirEntry) int {
		aInfo, err := a.Info()
		errCheck(err)
		bInfo, err := b.Info()
		errCheck(err)
		return int(aInfo.Size() - bInfo.Size())
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
