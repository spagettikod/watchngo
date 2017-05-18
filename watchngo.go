package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// EventType are the types of events that can occur for a file
type EventType int

const (
	// Created is the event type used when a file was created
	Created EventType = iota

	// Modified is the event type sent when a file was modified
	Modified

	// Deleted is the event type sent when a file was deleted
	Deleted
)

// FileEvent is sent on a channel when monitored directory changes.
type FileEvent struct {
	AbsPath   string
	Type      EventType
	Timestamp time.Time
}

func watch(dir string) chan FileEvent {
	mc := make(chan FileEvent)

	abs, err := filepath.Abs(dir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	go process(abs, mc)
	return mc
}

func process(dir string, events chan FileEvent) {
	prev := make(map[string]time.Time)
	for {
		timestamp := time.Now()
		scanned, err := scan(dir)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		// find new and modified files
		for k, v := range scanned {
			if p, exists := prev[k]; !exists || p != v {
				prev[k] = v
				if !exists {
					events <- FileEvent{AbsPath: k, Type: Created, Timestamp: timestamp}
				} else {
					events <- FileEvent{AbsPath: k, Type: Modified, Timestamp: timestamp}
				}
			}
		}
		// find removed files or directories
		for k := range prev {
			if _, exists := scanned[k]; !exists {
				delete(prev, k)
				events <- FileEvent{AbsPath: k, Type: Deleted, Timestamp: timestamp}
			}
		}
		time.Sleep(time.Second * 1)
	}
}

// scan a directory (recursively) and return a map of absolute file names and their modified timestamps
func scan(dir string) (files map[string]time.Time, err error) {
	files = make(map[string]time.Time)

	// open file in 'dir' defer closing it
	f, err := os.Open(dir)
	if err != nil {
		return //files, err
	}
	defer func() {
		if err = f.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	// get directory statistics
	stat, err := f.Stat()
	if err != nil {
		return files, err
	}

	// return an error if file in 'dir' is not an directory
	if !stat.IsDir() {
		return files, fmt.Errorf("%v is not a directory", dir)
	}

	// fetch all FileInfo in the directory
	fis, err := f.Readdir(-1)
	if err != nil {
		return files, err
	}

	// loop through all FileInfo
	for _, fi := range fis {
		abs := filepath.Join(dir, fi.Name())
		files[abs] = fi.ModTime()

		// do a recursive call if file is a directory
		if fi.IsDir() {
			r, err := scan(abs)
			if err != nil {
				return files, err
			}

			// add files found in the recursive call to the final list
			for k, v := range r {
				files[k] = v
			}
			continue
		}
	}

	return files, nil
}

func usage() {
	fmt.Println("usage: watchngo [flags] DIRECTORY COMMAND")
	fmt.Println("")
	fmt.Println("flags:")
	fmt.Println("  -v	verbose, watchngo outputmore output than what comes from stdout and stderr")
	fmt.Println("")
	fmt.Println("Watches DIRECTORY and runs COMMAND when changes are detected.")
	os.Exit(1)
}

func doLog(msg string, verbose *bool) {
	if *verbose {
		log.Println(msg)
	}
}

func main() {
	verbose := flag.Bool("v", false, "verbose")
	flag.Parse()
	args := flag.Args()

	if len(args) < 2 {
		usage()
	}
	dir := flag.Arg(0)
	command := flag.Arg(1)

	absDir, err := filepath.Abs(dir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	mod := watch(absDir)

	if *verbose {
		log.Printf("watchngo: watching directory %v", absDir)
	}

	var prevTimestamp time.Time
	for e := range mod {
		doLog(fmt.Sprintf("watchngo: change detected on file %v", e.AbsPath), verbose)
		if prevTimestamp != e.Timestamp {
			prevTimestamp = e.Timestamp
			doLog(fmt.Sprintf("watchngo: executing %v", command), verbose)
			output, err := exec.Command("sh", "-c", command).Output()
			if err != nil {
				if xerr, ok := err.(*exec.ExitError); ok {
					fmt.Fprint(os.Stderr, string(xerr.Stderr))
				} else {
					fmt.Fprint(os.Stderr, err.Error())
				}
				continue
			}
			if len(output) > 0 {
				fmt.Fprint(os.Stdout, string(output))
			}
			doLog(fmt.Sprintf("watchngo: finished %v", command), verbose)
		} else {
			doLog("watchngo: skipping execution, multiple files changed at the same time", verbose)
		}
	}
}
