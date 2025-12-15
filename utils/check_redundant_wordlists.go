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

type isRedundantRes struct {
	wlSmallName    string
	wlBigName      string
	isRedundantRet bool
	isRedundantErr error
}

func isRedundant(wlSmallName string, wlBigName string) (bool, error) {
	var wlSmallScanner, wlBigScanner *bufio.Scanner
	if wlSmall, err := os.Open(wlSmallName); err != nil {
		return false, err
	} else {
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

func isRedundantWorker(jobs chan isRedundantJob, results chan isRedundantRes) {
	for job := range jobs {
		isRedundantRet, isRedundantErr := isRedundant(job.wlSmallName, job.wlBigName)
		results <- isRedundantRes{
			wlSmallName:    job.wlSmallName,
			wlBigName:      job.wlBigName,
			isRedundantRet: isRedundantRet,
			isRedundantErr: isRedundantErr,
		}
	}
}

func main() {
	log.SetFlags(0)
	log.SetOutput(os.Stderr)
	wlNames := os.Args[1:]
	slices.SortFunc(wlNames, func(a string, b string) int {
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
	jobsN, workersN := len(wlNames)*(len(wlNames)-1)/2, runtime.NumCPU()
	jobs, results := make(chan isRedundantJob, jobsN), make(chan isRedundantRes, jobsN)
	for range workersN {
		go isRedundantWorker(jobs, results)
	}
	for wlSmallIndex, wlSmallName := range wlNames {
		for _, wlBigName := range wlNames[wlSmallIndex+1:] {
			jobs <- isRedundantJob{
				wlSmallName: wlSmallName,
				wlBigName:   wlBigName,
			}
		}
	}
	close(jobs)
	for range jobsN {
		result := <-results
		if result.isRedundantErr != nil {
			slog.Error("isRedundant(wlSmallName, wlBigName)",
				"error", result.isRedundantErr,
				"wlSmallName", result.wlSmallName,
				"wlBigName", result.wlBigName,
			)
		} else if result.isRedundantRet {
			slog.Info("isRedundant(wlSmallName, wlBigName)=true",
				"wlSmallName", result.wlSmallName,
				"wlBigName", result.wlBigName,
			)
		}
	}
	close(results)
}
