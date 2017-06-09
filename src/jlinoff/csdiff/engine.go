// Diff engine.
// License: The MIT License (MIT)
// Copyright (c) 2017 Joe Linoff
package main

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// Summary information.
type diffSummaryType struct {
	NumLeftLines       int
	NumRightLines      int
	NumLeftOnlyLines   int
	NumRightOnlyLines  int
	NumLinesMatch      int
	NumLinesDiff       int
	NumLeftCharsDiff   int
	NumLeftCharsMatch  int
	NumRightCharsMatch int
	NumRightCharsDiff  int
}

// Run the diff using the longest common subsequence, do not print anything.
func diffInit(opts options) (seq1, seq2 []string, mp [][]int) {
	seq1 = filter(opts, readLines(opts.File1))
	seq2 = filter(opts, readLines(opts.File2))

	// To support options like ignore whitespace or ignore case,
	// the lines must be modified before the LCS operation.
	lcs := longestCommonSubsequence(seq1, seq2)

	// Find the match points.
	// First entry is the line number in l1 and the second is the line number in l2.
	mp = [][]int{}
	for i := 0; i < len(lcs); i++ {
		r := []int{-1, -1}
		mp = append(mp, r)
	}

	// lambda seq initializer (DRY)
	initSeq := func(seq []string, k int) {
		for i, j := 0, 0; i < len(seq) && j < len(lcs); i++ {
			if seq[i] == lcs[j] {
				mp[j][k] = i
				j++
			}
		}
	}
	initSeq(seq1, 0) // Initialize for the first sequence.
	initSeq(seq2, 1) // Initialize the second sequence.
	return
}

// filter normalizes the output for comparisons using regular expressions.
func filter(opts options, lines []string) []string {
	newLines := []string{}
	if len(opts.Replacements) > 0 {
		for _, line := range lines {
			for _, rep := range opts.Replacements {
				line = rep.Pattern.ReplaceAllString(line, rep.Replacement)
			}
			newLines = append(newLines, line)
		}
	} else {
		newLines = lines
	}
	return newLines
}

// diff prints out the diffs in separate sections.
// It is a better choice for longer lines.
// Always suppress, ignore the -s option.
func (sum *diffSummaryType) diff(opts options, seq1, seq2 []string, mp [][]int) {
	// lambda to print the lines in the interval between
	// match points.
	printInterval := func(i1 *int, n1 int, i2 *int, n2 int) {
		x1 := *i1
		x2 := *i2
		if x1 >= n1 && x2 >= n2 {
			return
		}

		// Print the header:
		//  line numbers for the first file followed by line numbers for the
		//  second file.
		y1 := x1
		y2 := x2
		for x1 < n1 || x2 < n2 {
			// Advance the minimum.
			p1 := x1 < n1
			p2 := x2 < n2

			if p1 {
				x1++
			}
			if p2 {
				x2++
			}

			// Update the summary data.
			if p1 && p2 {
				sum.NumLinesDiff++
			} else if p1 {
				sum.NumLeftOnlyLines++
			} else if p2 {
				sum.NumRightOnlyLines++
			}
		}

		printDiff := func(a, b int) {
			if (b - a) == 0 {
				fmt.Printf("%v", a)
			} else {
				fmt.Printf("%v,%v", a, b)
			}
		}
		printDiff(y1+1, x1)
		fmt.Printf("c")
		printDiff(y2+1, x2)
		fmt.Println("")

		// Print the left diffs.
		x1 = *i1
		x2 = *i2
		for x1 < n1 || x2 < n2 {
			// Advance the minimum.
			p1 := x1 < n1
			p2 := x2 < n2
			refa := []bool{}

			if p1 && p2 {
				refa, _ = mapCommonSubStrings(seq1[x1], seq2[x2])
				if opts.Summary {
					for _, match := range refa {
						if match == false {
							sum.NumLeftCharsDiff++
						} else {
							sum.NumLeftCharsMatch++
						}
					}

					for _, match := range refa {
						if match == false {
							sum.NumRightCharsDiff++
						} else {
							sum.NumRightCharsMatch++
						}
					}
				}
			}

			// Print the left.
			if p1 {
				printSymbol(opts, "< ")
				printLine(opts, -1, seq1[x1], -1, false, refa, p2)
				fmt.Println("") // new line
				x1++
			}

			if p2 {
				x2++
			}
		}

		// Print the right diffs.
		first := true
		x1 = *i1
		x2 = *i2
		for x1 < n1 || x2 < n2 {
			// Advance the minimum.
			p1 := x1 < n1
			p2 := x2 < n2
			refb := []bool{}

			if p1 && p2 {
				_, refb = mapCommonSubStrings(seq1[x1], seq2[x2])
			}

			if p1 {
				x1++
			}

			// Print the right.
			if p2 {
				if first {
					first = false
					fmt.Println("---")
				}
				printSymbol(opts, "> ")
				printLine(opts, -1, seq2[x2], -1, false, refb, p1)
				fmt.Println("") // new line
				x2++
			}
		}

		*i1 = x1
		*i2 = x2
	}

	// Print all of the diffs.
	sum.NumLeftLines = len(seq1)
	sum.NumRightLines = len(seq2)
	info("num a=matchpoints: %v", len(mp))
	i1 := 0
	i2 := 0
	for i := 0; i < len(mp); i++ {
		n1 := mp[i][0]
		n2 := mp[i][1]
		printInterval(&i1, n1, &i2, n2)
		sum.NumLinesMatch++
		i1++
		i2++
	}
	printInterval(&i1, len(seq1), &i2, len(seq2))
	fmt.Println("")
}

// sdiff prints out the side by side diff.
// It uses the long common substring recursively to get the smallest set of
// differences.
func (sum *diffSummaryType) sdiff(opts options, seq1, seq2 []string, mp [][]int) {
	// Adjust the width for each side.
	// Define the formats.
	width := (opts.Width - 2) / 2
	if (opts.Width % 2) == 0 {
		width-- // adjust slightly to make it fit
	}

	fmt1 := fmt.Sprintf("%%6s %%-%ds", width-7) // left line + line num
	fmt2 := fmt.Sprintf("%%-%ds", width-7)      // left line only

	// Print the header.
	fmt.Println("")
	fmt.Printf("%6s ", "")
	fmt.Printf(fmt2, trunc(opts.File1, width-7))
	fmt.Printf("   ")
	fmt.Printf("%6s ", "")
	fmt.Printf("%v", trunc(opts.File2, width-7))
	fmt.Println("")

	// lambda to print the lines in the interval between
	// match points.
	printInterval := func(i1 *int, n1 int, i2 *int, n2 int) {
		for *i1 < n1 || *i2 < n2 {
			// Advance the minimum.
			p1 := *i1 < n1
			p2 := *i2 < n2
			refa := []bool{}
			refb := []bool{}

			if p1 && p2 {
				refa, refb = mapCommonSubStrings(seq1[*i1], seq2[*i2])

				// update the summary data
				if opts.Summary {
					for _, match := range refa {
						if match == false {
							sum.NumLeftCharsDiff++
						} else {
							sum.NumLeftCharsMatch++
						}
					}
					for _, match := range refb {
						if match == false {
							sum.NumRightCharsDiff++
						} else {
							sum.NumRightCharsMatch++
						}
					}
				}
			}

			// Print the left.
			if p1 {
				printLine(opts, *i1+1, seq1[*i1], width, true, refa, p2)
				*i1++
			} else {
				fmt.Printf(fmt1, "", "")
			}

			// Print the separator
			if p1 && p2 {
				printSymbol(opts, " | ") // change the line to match
			} else if p1 {
				printSymbol(opts, " < ") // only in the left
			} else if p2 {
				printSymbol(opts, " > ") // only in the right
			} else {
				printSymbol(opts, " * ")
			}

			// Print the right.
			if p2 {
				printLine(opts, *i2+1, seq2[*i2], width, false, refb, p1)
				*i2++
			} else {
				fmt.Printf("")
			}

			fmt.Println("") // new line

			// Update the summary data.
			if p1 && p2 {
				sum.NumLinesDiff++
			} else if p1 {
				sum.NumLeftOnlyLines++
			} else if p2 {
				sum.NumRightOnlyLines++
			}
		}
	}

	// Print the diffs and the matching lines.
	i1 := 0
	i2 := 0
	for i := 0; i < len(mp); i++ {
		n1 := mp[i][0]
		n2 := mp[i][1]
		printInterval(&i1, n1, &i2, n2)
		sum.NumLinesMatch++
		if opts.Suppress == false {
			// If suppression is off, print the matches.
			// Since they match we can use seq1 for everything.
			printLine(opts, n1+1, seq1[n1], width, true, []bool{}, true)
			if opts.Colorize == true {
				fmt.Print(opts.Colors.Symbol)
			}
			fmt.Print("   ")
			if opts.Colorize == true {
				fmt.Print(opts.Colors.Reset)
			}
			printLine(opts, n2+1, seq2[n2], width, false, []bool{}, true)
			fmt.Println("")
		}
		i1++
		i2++
	}
	printInterval(&i1, len(seq1), &i2, len(seq2))
	fmt.Println("")
}

// func printSymbol prints the symbol with the colorization.
func printSymbol(opts options, sym string) {
	if opts.Colorize == true {
		fmt.Print(opts.Colors.Symbol)
	}
	fmt.Printf("%v", sym)
	if opts.Colorize == true {
		fmt.Print(opts.Colors.Reset)
	}
}

// printLine
// left - true if left, false if right
// refs - map of character diffs
// both - both lines have values
func printLine(opts options, lineNum int, line string, width int, left bool, ref []bool, both bool) {
	if lineNum > 0 {
		fmt.Printf("%6d ", lineNum)
	}
	w := width - 7
	s := trunc(line, w)
	nr := 0 // num runes

	// doColor was added so that this could be used by matching lines without
	// overhead.
	if opts.Colorize {
		// The user specified the -c option, use the color map.
		if len(ref) > 0 {
			// Two lines, each have diffs.
			fmt.Print(opts.Colors.Reset)
			i := 0
			for i < len(s) && (nr < w || w < 1) {
				if (nr > 0 && ref[nr] != ref[nr-1]) || (nr == 0) {
					// Check the difference map (r) to see if we
					// need to colorize.
					if ref[nr] == false {
						fmt.Print(opts.Colors.Reset)
						fmt.Print(opts.Colors.CharsDiff)
					} else {
						fmt.Print(opts.Colors.Reset)
						fmt.Print(opts.Colors.CharsMatch)
					}
				}
				rv, width := utf8.DecodeRuneInString(s[i:])
				fmt.Printf("%c", rv)
				i += width
				nr++
			}
		} else {
			if both == true {
				// both lines match
				fmt.Print(opts.Colors.LinesMatch)
			} else if left == true {
				// only the left line
				fmt.Print(opts.Colors.LeftLineOnly)
			} else { // left is false
				// only the right line
				fmt.Print(opts.Colors.RightLineOnly)
			}
			nr = len(s)
			fmt.Printf("%v", s)
		}
		fmt.Print(opts.Colors.Reset)
	} else {
		nr = len(s)
		fmt.Printf("%v", s)
	}

	// Pad if this is the left side.
	if left {
		for nr < w {
			fmt.Printf(" ")
			nr++
		}
	}
}

// mapCommonSubStrings - maps common sub strings.
// The entry is true if they match or false otherwise.
func mapCommonSubStrings(a, b string) (refa, refb []bool) {
	refa = make([]bool, len(a))
	refb = make([]bool, len(b))

	type fctType func(fctType, string, string, int, int, int)
	fct := func(f fctType, as, bs string, offa, offb int, depth int) {
		if len(a) < 1 || len(b) < 1 {
			return
		}

		x := longestCommonSubstring(as, bs)
		if len(x) < 1 {
			return
		}

		pa := strings.Index(as, x)
		pb := strings.Index(bs, x)
		ia := pa
		ib := pb
		for i := 0; i < len(x); i++ {
			refa[ia+offa] = true
			refb[ib+offb] = true
			ia++
			ib++
		}

		// Left recursive all does not require index updates.
		f(f, as[:pa], bs[:pb], offa, offb, depth+1)
		// Right recursion requires updating the offsets.
		f(f, as[ia:], bs[ib:], offa+ia, offb+ib, depth+1)
	}

	fct(fct, a, b, 0, 0, 0)
	return
}

// longestCommonSubsequence
// Find the longest common subsequence between two strings.
func longestCommonSubsequence(seq1 []string, seq2 []string) []string {
	seq1Len := len(seq1)
	seq2Len := len(seq2)
	matrix := [][]int{}

	// initialize to -1
	// make sure it is 1 larger than necessary to
	// handle the case where the last entry in each row/column
	// mismatches
	for i := 0; i <= seq1Len; i++ {
		row := []int{}
		for j := 0; j <= seq2Len; j++ {
			row = append(row, -1)
		}
		matrix = append(matrix, row)
	}

	// Create the matrix of the longest subsequence lengths.
	// lambda to get the max of two integers.
	maxi := func(a, b int) int {
		if a >= b {
			return a
		}
		return b
	}
	for i := 0; i < seq1Len; i++ {
		for j := 0; j < seq2Len; j++ {
			if i == 0 || j == 0 {
				matrix[i][j] = 0
			} else if seq1[i-1] == seq2[j-1] {
				matrix[i][j] = 1 + matrix[i-1][j-1]
			} else {
				matrix[i][j] = maxi(matrix[i-1][j], matrix[i][j-1])
			}
		}
	}

	// Backtrack from the lower, right corner to get the characters.
	i := seq1Len
	j := seq2Len
	rlcs := []string{} // reversed version of the LCS
	for i > 0 && j > 0 {
		if seq1[i-1] == seq2[j-1] {
			// If the current char in seq1 and seq2, then
			// it is part of the LCS.
			rlcs = append(rlcs, seq1[i-1])
			i--
			j--
		} else if matrix[i-1][j] > matrix[i][j-1] {
			// They aren't the same, look at the previous entry in the row,
			// if it is the larger value, choose it.
			i--
		} else {
			j--
		}
	}

	// lambda function to reverse a list of strings.
	// simular rlcs[::-1]
	reverse := func(list []string) (result []string) {
		for i := len(list) - 1; i >= 0; i-- {
			result = append(result, list[i])
		}
		return
	}

	lcs := reverse(rlcs)
	return lcs
}

// longestCommonSubstring find the longest common substring between two
// two strings.
// CITATION: https://rosettacode.org/wiki/Longest_Common_Substring#Go
func longestCommonSubstring(a, b string) (out string) {
	lengths := make([]int, len(a)*len(b))
	greatestLength := 0
	for i, x := range a {
		for j, y := range b {
			if x == y {
				if i == 0 || j == 0 {
					lengths[i*len(b)+j] = 1
				} else {
					lengths[i*len(b)+j] = lengths[(i-1)*len(b)+j-1] + 1
				}
				if lengths[i*len(b)+j] > greatestLength {
					greatestLength = lengths[i*len(b)+j]
					out = a[i-greatestLength+1 : i+1]
				}
			}
		}
	}
	return
}
