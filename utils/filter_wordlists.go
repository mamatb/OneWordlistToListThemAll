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
	wlSuffix = "-unix_print_32max_nohash"
)

type filterWordlistResult struct {
	wordlist            string
	filterWordlistError error
}

func filterWordlist(wordlist string) error {
	wlExt := filepath.Ext(wordlist)
	wlName := strings.TrimSuffix(wordlist, wlExt)
	var wlScanner *bufio.Scanner
	if wlFile, err := os.Open(wlName + wlExt); err != nil {
		return err
	} else {
		defer wlFile.Close()
		wlScanner = bufio.NewScanner(wlFile)
		wlScanner.Buffer(make([]byte, scanTokenSize), scanTokenSize)
	}
	var wlWriter *bufio.Writer
	if wlFile, err := os.Create(wlName + wlSuffix + wlExt); err != nil {
		return err
	} else {
		defer wlFile.Close()
		wlWriter = bufio.NewWriter(wlFile)
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
		password := wlScanner.Text()
		if rePrint.MatchString(password) && !reHash.MatchString(password) &&
			len(password) <= 32 {
			if _, err := wlWriter.WriteString(password); err != nil {
				return err
			}
			if _, err := wlWriter.WriteRune('\n'); err != nil {
				return err
			}
		}
	}
	return wlScanner.Err()
}

func filterWordlistWorker(jobs chan string, results chan filterWordlistResult) {
	for job := range jobs {
		results <- filterWordlistResult{
			wordlist:            job,
			filterWordlistError: filterWordlist(job),
		}
	}
}

func main() {
	log.SetFlags(0)
	log.SetOutput(os.Stderr)
	jobsN, workersN := len(os.Args)-1, runtime.NumCPU()
	jobs, results := make(chan string, jobsN), make(chan filterWordlistResult, jobsN)
	for range workersN {
		go filterWordlistWorker(jobs, results)
	}
	for _, wordlist := range os.Args[1:] {
		jobs <- wordlist
	}
	close(jobs)
	for range jobsN {
		result := <-results
		if result.filterWordlistError != nil {
			slog.Error("filterWordlist(wordlist)",
				"error", result.filterWordlistError,
				"wordlist", result.wordlist,
			)
		}
	}
	close(results)
}
