package main

import (
	"bufio"
	"bytes"
	"io/fs"
	"log"
	"log/slog"
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
	wlSmallName      string
	wlBigName        string
	isRedundantBool  bool
	isRedundantError error
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

func isRedundant(wlSmallName string, wlBigName string) (bool, error) {
	var wlSmallScanner, wlBigScanner *bufio.Scanner
	if wlSmall, err := os.Open(wlSmallName); err == nil {
		defer wlSmall.Close()
		wlSmallScanner = bufio.NewScanner(wlSmall)
	} else {
		return false, err
	}
	if wlBig, err := os.Open(wlBigName); err == nil {
		defer wlBig.Close()
		wlBigScanner = bufio.NewScanner(wlBig)
	} else {
		return false, err
	}
	wlSmallScan, wlBigScan := wlSmallScanner.Scan(), wlBigScanner.Scan()
	for wlSmallScan && wlBigScan {
		switch bytes.Compare(wlSmallScanner.Bytes(), wlBigScanner.Bytes()) {
		case 1:
			wlBigScan = wlBigScanner.Scan()
		case 0:
			wlSmallScan, wlBigScan = wlSmallScanner.Scan(), wlBigScanner.Scan()
		case -1:
			return false, nil
		}
	}
	return !wlSmallScan, nil
}

func isRedundantWorker(jobs chan isRedundantJob, results chan isRedundantResult) {
	for job := range jobs {
		isRedundantBool, isRedundantError := isRedundant(job.wlSmallName, job.wlBigName)
		results <- isRedundantResult{
			wlSmallName:      job.wlSmallName,
			wlBigName:        job.wlBigName,
			isRedundantBool:  isRedundantBool,
			isRedundantError: isRedundantError,
		}
	}
}

func main() {
	const CWD, EXT string = ".", ".txt"
	log.SetFlags(0)
	log.SetOutput(os.Stderr)
	wordlists, err := readDirExt(CWD, EXT)
	if err != nil {
		slog.Error("readDirExt(name, ext)",
			"error", err,
			"name", CWD,
			"ext", EXT,
		)
	}
	slices.SortFunc(wordlists, func(a fs.DirEntry, b fs.DirEntry) int {
		var aSize, bSize int64
		if aInfo, err := a.Info(); err == nil {
			aSize = aInfo.Size()
		} else {
			return 0
		}
		if bInfo, err := b.Info(); err == nil {
			bSize = bInfo.Size()
		} else {
			return 0
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
		if result.isRedundantError != nil {
			slog.Error("isRedundant(wlSmallName, wlBigName)",
				"error", result.isRedundantError,
				"wlSmallName", result.wlSmallName,
				"wlBigName", result.wlBigName,
			)
		} else if result.isRedundantBool {
			slog.Info("isRedundant(wlSmallName, wlBigName) returned true",
				"wlSmallName", result.wlSmallName,
				"wlBigName", result.wlBigName,
			)
		}
	}
	close(results)
}
