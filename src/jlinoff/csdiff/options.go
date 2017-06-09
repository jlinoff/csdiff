// Process the command line options.
// License: The MIT License (MIT)
// Copyright (c) 2017 Joe Linoff
package main

import (
	"fmt"
	"jlinoff/termcolors"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

// colorsType are the colors used in colorize mode.
type colorsType struct {
	CharsMatch    string
	CharsDiff     string
	LinesMatch    string
	LeftLineOnly  string
	RightLineOnly string
	Symbol        string // |, <, >
	Reset         string
}

// replaceType is the list of replacement patterns that are used for filtering.
type replaceType struct {
	Pattern     *regexp.Regexp
	Replacement string
}

type options struct {
	File1        string
	File2        string
	Suppress     bool
	Width        int
	Colorize     bool
	Colors       colorsType
	SideBySide   bool
	Summary      bool
	Replacements []replaceType
}

func getopts() (opts options) {
	// lambda to get the next argument on the command line.
	nextArg := func(idx *int, o string) (arg string) {
		*idx++
		if *idx < len(os.Args) {
			arg = os.Args[*idx]
		} else {
			log.Fatalf("ERROR: missing argumnent for option '%s'", o)
		}
		return
	}

	// lambda to get the next argument on the command line.
	nextArgN := func(idx *int, o string, n int) (arg string) {
		*idx++
		if *idx < len(os.Args) {
			arg = os.Args[*idx]
		} else {
			log.Fatalf("ERROR: missing argumnent %d for option '%s'", n, o)
		}
		return
	}

	// lambda to get a range in an interval
	nextArgInt := func(idx *int, o string, min int, max int) (arg int) {
		a := nextArg(idx, o)
		arg = 0
		if v, e := strconv.Atoi(a); e == nil {
			if v < min {
				log.Fatalf("ERROR: '%v' too small, minimum accepted value is %v", o, min)
			} else if v > max {
				log.Fatalf("ERROR: '%v' too large, maximum value accepted is %v", o, max)
			}
			arg = v
		} else {
			log.Fatalf("ERROR: '%v' expected a number in the range [%v..%v]", o, min, max)
		}
		return
	}

	// Initialize the colors.
	// Use background to differentiate so that the user can space differences.
	reset, _ := termcolors.ParseColorExpr("clear")
	def, _ := termcolors.ParseColorExpr("bgLightGrey")
	symdef, _ := termcolors.ParseColorExpr("red,bold")
	ct := colorsType{
		CharsMatch:    reset,
		CharsDiff:     def,
		LinesMatch:    reset,
		LeftLineOnly:  def,
		RightLineOnly: def,
		Symbol:        symdef,
		Reset:         reset,
	}

	// Initialize the options structure.
	opts = options{
		Width:        int(termcolors.GetTermInfo().Cols),
		Colorize:     true,
		Colors:       ct,
		SideBySide:   true,
		Replacements: []replaceType{},
	}

	// Process the CLI arguments.
	for i := 1; i < len(os.Args); i++ {
		opt := os.Args[i]
		switch opt {
		case "--256":
			termcolors.Print256ColorTables()
			os.Exit(0)
		case "-h", "--help":
			help()
		case "-c", "--color-map":
			cm := nextArg(&i, opt)
			getColorMap(opt, cm, &opts)
		case "--clear":
			clear, _ := termcolors.ParseColorExpr("clear")
			opts.Colors = colorsType{
				CharsMatch:    clear,
				CharsDiff:     clear,
				LinesMatch:    clear,
				LeftLineOnly:  clear,
				RightLineOnly: clear,
				Symbol:        clear,
			}
		case "--config":
			config := nextArg(&i, opt)
			readConfig(opt, config, &opts)
		case "-d", "--diff":
			opts.SideBySide = false
		case "-n", "--no-colorize":
			opts.Colorize = false
		case "-r", "--replace":
			p := nextArgN(&i, opt, 1)
			r := nextArgN(&i, opt, 2)
			rp, e := regexp.Compile(p)
			if e != nil {
				log.Fatalf("invalid regular expression '%v' for %v", p, opt)
			}
			replace := replaceType{Pattern: rp, Replacement: r}
			opts.Replacements = append(opts.Replacements, replace)
		case "-s", "--suppress-common-lines":
			opts.Suppress = true
		case "--summary":
			opts.Summary = true
		case "-w", "--width":
			opts.Width = nextArgInt(&i, opt, 8, 100000)
		case "-V", "--version":
			b := filepath.Base(os.Args[0])
			fmt.Printf("%v v%v\n", b, version)
			os.Exit(0)
		default:
			if len(opts.File1) == 0 {
				opts.File1 = opt
				if _, err := os.Stat(opt); os.IsNotExist(err) {
					log.Fatalf("file does not exist: '%v'", opt)
				}
			} else if len(opts.File2) == 0 {
				opts.File2 = opt
				if _, err := os.Stat(opt); os.IsNotExist(err) {
					log.Fatalf("file does not exist: '%v'", opt)
				}
			} else {
				log.Fatalf("too many arguments specified")
			}
		}
	}
	return
}

// getColorMap gets the color map argument.
// this is quite complex.
func getColorMap(opt string, cms string, opts *options) {
	// The format is:
	//  -c <target>=attr[,attr][;<target>=attr[,attr]]

	lines := strings.Split(cms, ";")
	for _, cm := range lines {
		toks := strings.SplitN(cm, "=", 2)
		if len(toks) < 2 {
			log.Fatalf("invalid argument for '%v', expected <fld>=<values>: %v", opt, cm)
		}

		key := strings.TrimSpace(toks[0]) // for file parsing
		seq, err := termcolors.ParseColorExpr(toks[1])
		if err != nil {
			log.Fatalf("invalid key value '%v' for '%v', see help (-h): %v", key, opt, err)
		}
		switch strings.ToLower(key) {
		case "charsmatch", "cm":
			opts.Colors.CharsMatch = seq
		case "charsdiff", "cd":
			opts.Colors.CharsDiff = seq
		case "linesmatch", "lm":
			opts.Colors.LinesMatch = seq
		case "leftlineonly", "left", "llo":
			opts.Colors.LeftLineOnly = seq
		case "rightlineonly", "right", "rlo":
			opts.Colors.RightLineOnly = seq
		case "symbol", "sym", "s":
			opts.Colors.Symbol = seq
		default:
			log.Fatalf("invalid key value '%v' for '%v', see help (-h)", key, opt)
		}
	}
}

// readConfig reads the config file.
func readConfig(opt string, config string, opts *options) {
	lines := readLines(config)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		r, _ := utf8.DecodeRuneInString(line[0:])
		if string(r) == "#" {
			continue
		}
		getColorMap(opt, line, opts)
	}
}
