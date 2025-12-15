package main

import (
	"bufio"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

const (
	scanTokenSize = 1024 * bufio.MaxScanTokenSize
	wlSuffix      = "-unix_print_32max_nohash"
)

type filterWlJob struct {
	wlNameIn  string
	wlNameOut string
}

type filterWlRes struct {
	wlNameIn    string
	wlNameOut   string
	filterWlErr error
}

func filterWordlist(wlNameIn string, wlNameOut string) error {
	var wlScanner *bufio.Scanner
	if wlFileIn, err := os.Open(wlNameIn); err != nil {
		return err
	} else {
		defer wlFileIn.Close()
		wlScanner = bufio.NewScanner(wlFileIn)
		wlScanner.Buffer(make([]byte, scanTokenSize), scanTokenSize)
	}
	var wlWriter *bufio.Writer
	if wlFileOut, err := os.Create(wlNameOut); err != nil {
		return err
	} else {
		defer wlFileOut.Close()
		wlWriter = bufio.NewWriter(wlFileOut)
		defer wlWriter.Flush()
	}
	rePrint, err := regexp.Compile(`^[[:print:]]*$`)
	if err != nil {
		return err
	}
	reHash, err := regexp.Compile(`^[[:xdigit:]]{32}$`)
	if err != nil {
		return err
	}
	for wlScanner.Scan() {
		word := wlScanner.Text()
		if rePrint.MatchString(word) && len(word) <= 32 && !reHash.MatchString(word) {
			if _, err := wlWriter.WriteString(word); err != nil {
				return err
			}
			if _, err := wlWriter.WriteRune('\n'); err != nil {
				return err
			}
		}
	}
	return wlScanner.Err()
}

func filterWlWorker(jobs chan filterWlJob, results chan filterWlRes) {
	for job := range jobs {
		results <- filterWlRes{
			wlNameIn:    job.wlNameIn,
			wlNameOut:   job.wlNameOut,
			filterWlErr: filterWordlist(job.wlNameIn, job.wlNameOut),
		}
	}
}

func main() {
	log.SetFlags(0)
	log.SetOutput(os.Stderr)
	jobsN, workersN := len(os.Args)-1, runtime.NumCPU()
	jobs, results := make(chan filterWlJob, jobsN), make(chan filterWlRes, jobsN)
	for range workersN {
		go filterWlWorker(jobs, results)
	}
	for _, wlName := range os.Args[1:] {
		wlExt := filepath.Ext(wlName)
		wlStem := strings.TrimSuffix(wlName, wlExt)
		jobs <- filterWlJob{
			wlNameIn:  wlName,
			wlNameOut: wlStem + wlSuffix + wlExt,
		}
	}
	close(jobs)
	for range jobsN {
		result := <-results
		if result.filterWlErr != nil {
			slog.Error("filterWordlist(wlNameIn, wlNameOut)",
				"error", result.filterWlErr,
				"wlNameInt", result.wlNameIn,
				"wlNameOut", result.wlNameOut,
			)
		}
	}
	close(results)
}
