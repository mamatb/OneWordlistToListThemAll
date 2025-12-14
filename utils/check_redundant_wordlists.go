package main

import (
	"bufio"
	"bytes"
	"log"
	"log/slog"
	"os"
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

func isRedundant(wlSmallName string, wlBigName string) (bool, error) {
	var wlSmallScanner, wlBigScanner *bufio.Scanner
	if wlSmall, err := os.Open(wlSmallName); err != nil {
		return false, err
	else {
		defer wlSmall.Close()
		wlSmallScanner = bufio.NewScanner(wlSmall)
	}
	if wlBig, err := os.Open(wlBigName); err != nil {
		return false, err
	} else {
		defer wlBig.Close()
		wlBigScanner = bufio.NewScanner(wlBig)
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
	log.SetFlags(0)
	log.SetOutput(os.Stderr)
	wordlists := os.Args[1:]
	slices.SortFunc(wordlists, func(a string, b string) int {
		var aSize, bSize int64
		if aInfo, err := os.Stat(a); err != nil {
			return 0
		} else {
			aSize = aInfo.Size()
		}
		if bInfo, err := os.Stat(b); err != nil {
			return 0
		} else {
			bSize = bInfo.Size()
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
				wlSmallName: wlSmall,
				wlBigName:   wlBig,
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
			slog.Info("isRedundant(wlSmallName, wlBigName)=true",
				"wlSmallName", result.wlSmallName,
				"wlBigName", result.wlBigName,
			)
		}
	}
	close(results)
}
