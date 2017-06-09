// Program to diff to files to find mismatches using the
// longest common subsequence (LCS) algorithm for line differences
// and longest common substring for character differences.
// See the help (-h) for more information.
// License: The MIT License (MIT)
// Copyright (c) 2017 Joe Linoff
package main

import "fmt"

var version = "0.4.3"

func main() {
	opts := getopts()
	seq1, seq2, mp := diffInit(opts)

	sum := diffSummaryType{}
	if opts.SideBySide {
		sum.sdiff(opts, seq1, seq2, mp)
	} else {
		sum.diff(opts, seq1, seq2, mp)
	}
	if opts.Summary {
		printSummary(sum)
	}
}

// printSummary prints the diff summary.
func printSummary(sum diffSummaryType) {
	fct := func(key string, val int) {
		fmt.Printf("%-30s : %6d\n", key, val)
	}

	fct("summary: NumLinesMatch", sum.NumLinesMatch)
	fct("summary: NumLinesDiff", sum.NumLinesDiff)
	fct("summary: NumLeftLines", sum.NumLeftLines)
	fct("summary: NumLeftOnlyLines", sum.NumLeftOnlyLines)
	fct("summary: NumLeftCharsDiff", sum.NumLeftCharsDiff)
	fct("summary: NumLeftCharsMatch", sum.NumLeftCharsMatch)
	fct("summary: NumRightLines", sum.NumRightLines)
	fct("summary: NumRightOnlyLines", sum.NumRightOnlyLines)
	fct("summary: NumRightCharsDiff", sum.NumRightCharsDiff)
	fct("summary: NumRightCharsMatch", sum.NumRightCharsMatch)
}
