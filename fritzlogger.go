package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func findinLogs(logline string) int {
	var found bool = true
	var newlogs []string

	totallog, err := os.OpenFile("fritz.logs", os.O_RDONLY|os.O_CREATE, 0644)
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

	ttwr, _ := os.OpenFile("fritz.logs", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer ttwr.Close()

	for _, line := range newlogs {
		ttwr.WriteString(line + "\n")
	}

	// fmt.Printf("%v log lines found.\n", newlines)

	return newlines
}

func main() {
	localpath, _ := os.Getwd()
	if len(os.Args) > 1 {
		if os.Args[1] != "" {
			localpath = os.Args[1]
		}
	}

	files, err := ioutil.ReadDir(localpath)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if !file.IsDir() {
			if filepath.Ext(file.Name()) == ".log" {
				logfile, err := os.Open(localpath + "/" + file.Name())
				if err != nil {
					fmt.Println(err)
				}
				defer logfile.Close()

				linesscanner := bufio.NewScanner(logfile)
				for linesscanner.Scan() {
					findinLogs(linesscanner.Text())
				}
				os.Remove(localpath + "/" + file.Name())
			}
		}
	}
}
