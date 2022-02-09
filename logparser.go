package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/schollz/progressbar/v3"
)

func getLogFiles(localpath string) []string {
	fmt.Printf("Listing log files in %s\n", localpath)
	bar := progressbar.NewOptions(-1,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(false),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetDescription("[yellow][1/2][reset] Indexing Log Files..."))

	bar.Add(1)

	files, err := ioutil.ReadDir(localpath)
	if err != nil {
		log.Fatal(err)
	}

	var logfiles []string
	for _, file := range files {
		bar.Add(1)
		_, filename := filepath.Split(file.Name())
		if !file.IsDir() {
			if (filepath.Ext(file.Name()) == ".log") && (filename[0:6] == "fritz-") {
				logfiles = append(logfiles, file.Name())
			}
		}
	}

	fmt.Printf("%v log files found\n", len(logfiles))
	return logfiles
}

func parselogs(remove bool, localpath, mergelogfilename string) {
	var logfiles []string = getLogFiles(localpath)

	bar := progressbar.NewOptions(len(logfiles),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(false),
		progressbar.OptionSetPredictTime(true),
		progressbar.OptionFullWidth(),
		progressbar.OptionShowCount(),
		progressbar.OptionShowIts(),
		progressbar.OptionSetDescription("[yellow][2/2][reset] Parsing Log Files..."))

	if len(logfiles) > 0 {
		for _, file := range logfiles {
			bar.Add(1)
			var addedloglines int
			logfile, err := os.Open(localpath + "/" + file)
			if err != nil {
				fmt.Println(err)
			}
			defer logfile.Close()

			linesscanner := bufio.NewScanner(logfile)
			for linesscanner.Scan() {
				addedloglines = findinLogs(linesscanner.Text(), mergelogfilename)
			}

			if remove {
				os.Remove(localpath + "/" + file)
			}

			if addedloglines > 0 {
				fmt.Printf("%v log lines added.", addedloglines)
			}
		}
	}
}

func findinLogs(logline, mergelogfilename string) int {
	var found bool = true
	var newlogs []string

	totallog, err := os.OpenFile(mergelogfilename, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer totallog.Close()

	var linecount int = 0
	var newlines int = 0

	linesscanner := bufio.NewScanner(totallog)
	for linesscanner.Scan() {
		if linesscanner.Text() == logline {
			found = true
			linecount++
		}
	}

	if !found {
		newlines++
		newlogs = append(newlogs, logline)
	}

	if err := linesscanner.Err(); err != nil {
		fmt.Println(err)
	}

	if linecount == 0 {
		newlogs = append(newlogs, logline)
	}

	ttwr, _ := os.OpenFile(mergelogfilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer ttwr.Close()

	for _, line := range newlogs {
		ttwr.WriteString(line + "\n")
	}

	return newlines
}
