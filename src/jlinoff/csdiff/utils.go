// Utilities.
// License: The MIT License (MIT)
// Copyright (c) 2017 Joe Linoff
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime"
)

// check an error, report it and exit with the callers line number.
func check(e error) {
	if e != nil {
		_, _, lineno, _ := runtime.Caller(1)
		log.Fatalf("ERROR:%v %v", lineno, e)
	}
}

// truncate string
// same as s[:w] in python
func trunc(s string, w int) string {
	if w < 1 || len(s) <= w {
		return s
	}
	n := s[:w-1] + "$" // add $ at the end to show that we truncated
	return n
}

// debug prints a debug message with the callers
// line number.
func debug(f string, a ...interface{}) {
	_, _, lineno, _ := runtime.Caller(1)
	s := fmt.Sprintf(f, a...)
	fmt.Fprint(os.Stderr, "DEBUG:")
	log.Printf(":%v %v\n", lineno, s)
}

// info prints a info message with the callers
// line number.
func info(f string, a ...interface{}) {
	_, _, lineno, _ := runtime.Caller(1)
	s := fmt.Sprintf(f, a...)
	fmt.Fprint(os.Stderr, "INFO:")
	log.Printf(":%v %v\n", lineno, s)
}

// readlines reads lines from a text file
// If no data is available, the lines slice is empty.
func readLines(path string) (lines []string) {
	lines = []string{}

	fp, err := os.Open(path)
	check(err)
	defer fp.Close()

	s := bufio.NewScanner(fp)
	for s.Scan() {
		lines = append(lines, s.Text())
	}
	return
}
