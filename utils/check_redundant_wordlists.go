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
	"strings"
)

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

func isRedundant(wlSmallName string, wlBigName string) string {
	wlSmall, err := os.Open(wlSmallName)
	errCheck(err)
	defer wlSmall.Close()
	wlBig, err := os.Open(wlBigName)
	errCheck(err)
	defer wlBig.Close()
	wlSmallScanner, wlBigScanner := bufio.NewScanner(wlSmall), bufio.NewScanner(wlBig)
	wlSmallScan, wlBigScan := wlSmallScanner.Scan(), wlBigScanner.Scan()
	wlSmallLine, wlBigLine := wlSmallScanner.Bytes(), wlBigScanner.Bytes()
	for wlSmallScan && wlBigScan {
		switch bytes.Compare(wlSmallLine, wlBigLine) {
		case 1:
			wlBigScan = wlBigScanner.Scan()
			wlBigLine = wlBigScanner.Bytes()
		case 0:
			wlSmallScan, wlBigScan = wlSmallScanner.Scan(), wlBigScanner.Scan()
			wlSmallLine, wlBigLine = wlSmallScanner.Bytes(), wlBigScanner.Bytes()
		case -1:
			wlBigScan = false
		}
	}
	if wlSmallScan {
		return wlSmallName + " is not redundant with " + wlBigName
	} else {
		return wlSmallName + " is redundant with " + wlBigName
	}
}

func isRedundantWorker(jobs chan []string, results chan string) {
	for wlNames := range jobs {
		results <- isRedundant(wlNames[0], wlNames[1])
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
	jobs, results := make(chan []string, jobsN), make(chan string, jobsN)
	for i := 0; i < workersN; i++ {
		go isRedundantWorker(jobs, results)
	}
	for wlSmallIndex, wlSmall := range wordlists {
		for _, wlBig := range wordlists[wlSmallIndex+1:] {
			jobs <- []string{wlSmall.Name(), wlBig.Name()}
		}
	}
	close(jobs)
	for i := 0; i < jobsN; i++ {
		result := <-results
		if strings.Contains(result, "is redundant") {
			fmt.Println(result)
		}
	}
	close(results)
}
