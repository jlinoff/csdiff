// Package termcolors is the ANSI terminal color tools.
//
// You call it with a list of colorMap attributes and it returns a string that
// will implement them. Here is a simple example that prints some text in
// bold red.
//
//    clear, _ := termcolors.ParseColorExpr("clear")
//    boldRed, _ := termcolors.ParseColorExpr("red,bold")
//
//    fmt.Print(boldRed)
//    fmt.Println("This will be red!")
//    fmt.Print(clear)
//
// You can also use 256 color mode. Here is a simple example that
// shows that.
//
//    clear, _ := termcolors.ParseColorExpr("clear")
//    boldRed, _ := termcolors.ParseColorExpr("fg256[9],bold")
//
//    fmt.Print(boldRed)
//    fmt.Println("This will be red!")
//    fmt.Print(clear)
//
// The color map keys are not case sensitive so RED, red and Red are all the
// same.
//
// License: The MIT License (MIT)
// Copyright (c) 2017 Joe Linoff
package termcolors

import (
	"fmt"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

// ColorMapType is a unique type for color map constants.
type ColorMapType int

// colorMapType constants
const (
	None    ColorMapType = 0
	Mode8     = 1
	Mode256Fg = 2
	Mode256Bg = 3
)

// ColorMapData defines each item in the color map.
type ColorMapData struct {
	Mode  ColorMapType
	Value int
}

// ColorMap is the map of colors and attributes (like bold)
// that are available.
var ColorMap = map[string]ColorMapData{
	// Attributes.
	"blink":              ColorMapData{Mode8,  5},
	"bold":               ColorMapData{Mode8,  1},
	"clear":              ColorMapData{Mode8,  0},
	"dim":                ColorMapData{Mode8,  2},
	"hidden":             ColorMapData{Mode8,  6},
	"init":               ColorMapData{Mode8,  0},
	"italics":            ColorMapData{Mode8,  3},
	"reset":              ColorMapData{Mode8,  0},
	"resetblink":         ColorMapData{Mode8, 25},
	"resetbold":          ColorMapData{Mode8, 21},
	"resetdim":           ColorMapData{Mode8, 22},
	"resethidden":        ColorMapData{Mode8, 26},
	"resetitalics":       ColorMapData{Mode8, 23},
	"resetreverse":       ColorMapData{Mode8, 27},
	"resetstrikethrough": ColorMapData{Mode8, 29},
	"resetunderline":     ColorMapData{Mode8, 24},
	"reverse":            ColorMapData{Mode8,  7},
	"setblink":           ColorMapData{Mode8,  5},
	"setbold":            ColorMapData{Mode8,  1},
	"setdim":             ColorMapData{Mode8,  2},
	"sethidden":          ColorMapData{Mode8,  6},
	"setitalics":         ColorMapData{Mode8,  3},
	"setreverse":         ColorMapData{Mode8,  7},
	"setunderline":       ColorMapData{Mode8,  4},
	"underline":          ColorMapData{Mode8,  4},

	// Standard foreground.
	"black":       ColorMapData{Mode8, 30},  // same as Mode256Fg 0
	"blue":        ColorMapData{Mode8, 34},  // same as Mode25Fg, 4
	"cyan":        ColorMapData{Mode8, 36},  // same as Mode25Fg, 6
	"fgblack":     ColorMapData{Mode8, 30},  // same as Mode256Fg 0
	"fgblue":      ColorMapData{Mode8, 34},  // same as Mode25Fg, 4
	"fgcyan":      ColorMapData{Mode8, 36},  // same as Mode25Fg, 6
	"fgdefault":   ColorMapData{Mode8, 39},
	"fggreen":     ColorMapData{Mode8, 32},  // same as Mode25Fg, 2
	"fglightgray": ColorMapData{Mode8, 37},  // same as Mode25Fg, 7
	"fglightgrey": ColorMapData{Mode8, 37},  // same as Mode25Fg, 7
	"fgmagenta":   ColorMapData{Mode8, 35},  // same as Mode25Fg, 5
	"fgred":       ColorMapData{Mode8, 31},  // same as Mode25Fg, 1
	"fgyellow":    ColorMapData{Mode8, 33},  // same as Mode25Fg, 3
	"green":       ColorMapData{Mode8, 32},  // same as Mode25Fg, 2
	"lightgray":   ColorMapData{Mode8, 37},  // same as Mode25Fg, 7
	"lightgrey":   ColorMapData{Mode8, 37},  // same as Mode25Fg, 7
	"magenta":     ColorMapData{Mode8, 35},  // same as Mode25Fg, 5
	"red":         ColorMapData{Mode8, 31},  // same as Mode25Fg, 1
	"yellow":      ColorMapData{Mode8, 33},  // same as Mode25Fg, 3

	// High intensity foreground.
	"brightblue":      ColorMapData{Mode8, 94},  // same as Mode25Fg, 12
	"brightcyan":      ColorMapData{Mode8, 96},  // same as Mode25Fg, 14
	"brightgreen":     ColorMapData{Mode8, 92},  // same as Mode25Fg, 10
	"brightmagenta":   ColorMapData{Mode8, 95},  // same as Mode25Fg, 13
	"brightred":       ColorMapData{Mode8, 91},  // same as Mode25Fg, 9
	"brightwhite":     ColorMapData{Mode8, 97},  // same as Mode25Fg, 15
	"brightyellow":    ColorMapData{Mode8, 93},  // same as Mode25Fg, 11
	"darkgray":        ColorMapData{Mode8, 90},  // same as Mode25Fg, 8
	"darkgrey":        ColorMapData{Mode8, 90},  // same as Mode25Fg, 8
	"fgbrightblue":    ColorMapData{Mode8, 94},  // same as Mode25Fg, 12
	"fgbrightcyan":    ColorMapData{Mode8, 96},  // same as Mode25Fg, 14
	"fgbrightgreen":   ColorMapData{Mode8, 92},  // same as Mode25Fg, 10
	"fgbrightmagenta": ColorMapData{Mode8, 95},  // same as Mode25Fg, 13
	"fgbrightred":     ColorMapData{Mode8, 91},  // same as Mode25Fg, 9
	"fgbrightwhite":   ColorMapData{Mode8, 97},  // same as Mode25Fg, 15
	"fgbrightyellow":  ColorMapData{Mode8, 93},  // same as Mode25Fg, 11
	"fgdarkgray":      ColorMapData{Mode8, 90},  // same as Mode25Fg, 8
	"fgdarkgrey":      ColorMapData{Mode8, 90},  // same as Mode25Fg, 8

	// Standard background.
	"bgblack":     ColorMapData{Mode8, 40},  // same as Mode256Bg 0
	"bgred":       ColorMapData{Mode8, 41},  // same as Mode256Bg 1
	"bggreen":     ColorMapData{Mode8, 42},  // same as Mode256Bg 2
	"bgyellow":    ColorMapData{Mode8, 43},  // same as Mode256Bg 3
	"bgblue":      ColorMapData{Mode8, 44},  // same as Mode256Bg 4
	"bgmagenta":   ColorMapData{Mode8, 45},  // same as Mode256Bg 5
	"bgcyan":      ColorMapData{Mode8, 46},  // same as Mode256Bg 6
	"bglightgrey": ColorMapData{Mode8, 47},  // same as Mode256Bg 7
	"bglightgray": ColorMapData{Mode8, 47},  // same as Mode256Bg 7
	"bgdefault":   ColorMapData{Mode8, 49},

	// High intensity background.
	"bgdarkgrey":      ColorMapData{Mode8, 100},   // same as Mode256Bg 8
	"bgdarkgray":      ColorMapData{Mode8, 100},  // same as Mode256Bg 8
	"bgbrightred":     ColorMapData{Mode8, 101},  // same as Mode256Bg 9
	"bgbrightgreen":   ColorMapData{Mode8, 102},  // same as Mode256Bg 10
	"bgbrightyellow":  ColorMapData{Mode8, 103},  // same as Mode256Bg 11
	"bgbrightblue":    ColorMapData{Mode8, 104},  // same as Mode256Bg 12
	"bgbrightmagenta": ColorMapData{Mode8, 105},  // same as Mode256Bg 13
	"bgbrightcyan":    ColorMapData{Mode8, 106},  // same as Mode256Bg 14
	"bgbrightwhite":   ColorMapData{Mode8, 107},  // same as Mode256Bg 15

	// Mode 256.
	"fg256": ColorMapData{Mode256Fg, 0},
	"bg256": ColorMapData{Mode256Bg, 0},
}

// ParseColorExpr parses a color expression into the appropriate ANSI
// escape sequences.
//
// The color expression is a comma separated list of keywords from the
// color map. The keywords are not case sensitive so fgred can be input
// as "fgred", "FGRED", "fgRed" and so on.
//
// The mode256 entries are treated as arrays so to get foreground mode
// 256 color 9 (bright red), you would specify: fg256[9]. Similarly to
// get background mode 256 color 252 you would specify bg256[252].
//
// Here are some examples:
//     fgred        --> ESC[31m
//     fgred,bold   --> ESC[31;1m
//     fg256[1]     --> ESC[38;5;1m
func ParseColorExpr(expr string) (result string, err error) {
	// lambda to parse the expression into records for later processing.
	parse := func(expr string) (recs []ColorMapData, err error){
		recs = []ColorMapData{}
		toks := strings.Split(expr, ",")
		for _, tok := range toks {
			token := strings.TrimSpace(strings.ToLower(tok))
			key := token
			val := -1
			if strings.Index(key, "[") >= 0 {
				// Parse fg256[3] --> fg256 and 3
				xs := strings.Split(key, "[")
				key = xs[0]
				vss := strings.Split(xs[1], "]")
				vs := strings.TrimSpace(vss[0])
				v, e := strconv.Atoi(vs)
				if e != nil {
					err = fmt.Errorf("parsing failed for %v: %v", expr, e)
					return
				}
				val = v
			}
			rec, ok := ColorMap[key]
			if ok == false {
				err = fmt.Errorf("parsing failed for %v: %v not found", expr, key)
				return
			}
			if val >= 0 {
				rec.Value = val
			}
			recs = append(recs, rec)
		}
		return
	}

	// lambda to build the return string from the parsed expression.
	build := func(recs []ColorMapData) (result string, err error) {
		mode := None
		result = ""
		first := true
		for i, rec := range recs {
			if rec.Mode != mode {
				first = true
				mode = rec.Mode
				if i > 0 {
					result += "m"
				}
				switch mode {
					case Mode8: // ESC[...m
					result += fmt.Sprintf("\x1b[")
					case Mode256Fg: // ESC[38;5;#m
					result += fmt.Sprintf("\x1b[38;5;")
					case Mode256Bg: // ESC[48;5;#m
					result += fmt.Sprintf("\x1b[48;5;")
				default:
					break
				}
			}
			if first == false {
				result += ";"
			}
			result += fmt.Sprintf("%d", rec.Value)
			first = false
		}
		result += "m"
		return
	}

	// parse the expression and build the return string.
	err = nil
  recs := []ColorMapData{}
	recs, err = parse(expr)
	if err != nil {
		return
	}
	result, err = build(recs)
	if err != nil {
		return
	}

	return
}

// TermInfoType stores attributes about the terminal.
type TermInfoType struct {
	Rows         uint16
	Cols         uint16
	Xpixel       uint16
	Ypixel       uint16
}

// GetTermInfo returns the height and width of the terminal.
func GetTermInfo() (ti TermInfoType) {
	r, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(&ti)))

	if int(r) == -1 {
		panic(errno)
	}

  return
}

// Print256ColorTables prints out the 256 color tables for foreground and
// background colors.
func Print256ColorTables() {
	// Map when to use reverse video.
	m := map[int]bool{
		0: true, 1: true, 2: true, 3: true, 4: true, 5: true, 6: true, 8: true,
		12: true, 16: true, 17: true, 18: true, 19: true, 20: true, 21: true,
		52: true, 53: true, 54: true, 55: true, 56: true, 57: true,
		232: true, 233: true, 234: true, 235: true,
		236: true, 237: true, 238: true, 239: true,
		240: true, 241: true, 242: true, 243: true}

	// ================================================================
	// 8 color mode - background
	// ================================================================
	fmt.Print("\n")
	fmt.Println("8 Color Mode - Background (ESC[40m .. ESC[47m)")
	fmt.Print("   ")
	fmt.Printf("\x1b[%d;37;1m %-7s \x1b[0m", 40, "Black")
	fmt.Printf("\x1b[%d;37;1m %-7s \x1b[0m", 41, "Red")
	fmt.Printf("\x1b[%d;37;1m %-7s \x1b[0m", 42, "Green")
	fmt.Printf("\x1b[%d;37;1m %-7s \x1b[0m", 43, "Yellow")
	fmt.Printf("\x1b[%d;37;1m %-7s \x1b[0m", 44, "Blue")
	fmt.Printf("\x1b[%d;37;1m %-7s \x1b[0m", 45, "Magenta")
	fmt.Printf("\x1b[%d;37;1m %-7s \x1b[0m", 46, "Cyan")
	fmt.Printf("\x1b[%d;30;1m %-7s \x1b[0m", 47, "White")
	fmt.Print("\n")

	// ================================================================
	// 8 color mode - foreground
	// ================================================================
	fmt.Print("\n")
	fmt.Println("8 Color Mode - Foreground (ESC[30m .. ESC[37m)")
	fmt.Print("   ")
	fmt.Printf("\x1b[%d;47;1m %-7s \x1b[0m", 30, "Black")
	fmt.Printf("\x1b[%d;47;1m %-7s \x1b[0m", 31, "Red")
	fmt.Printf("\x1b[%d;47;1m %-7s \x1b[0m", 32, "Green")
	fmt.Printf("\x1b[%d;47;1m %-7s \x1b[0m", 33, "Yellow")
	fmt.Printf("\x1b[%d;47;1m %-7s \x1b[0m", 34, "Blue")
	fmt.Printf("\x1b[%d;47;1m %-7s \x1b[0m", 35, "Magenta")
	fmt.Printf("\x1b[%d;47;1m %-7s \x1b[0m", 36, "Cyan")
	fmt.Printf("\x1b[%d;40;1m %-7s \x1b[0m", 37, "White")
	fmt.Print("\n")

	// ================================================================
	// 256 color mode - background
	// ================================================================
	fmt.Print("\n")
	fmt.Print("256 Color Mode - Background (ESC[48;5;Nm)")
	for i := 0; i < 256; i++ {
		if (i % 16) == 0 {
			fmt.Printf("\n   ")
		}
		v, ok := m[i]
		if ok && v {
			fmt.Print("\x1b[37m")
		} else {
			fmt.Print("\x1b[30m")
		}
		fmt.Printf("\x1b[48;5;%dm %3d \x1b[0m", i, i)
	}
	fmt.Print("\n")

	// ================================================================
	// 256 color mode - foreground
	// ================================================================
	fmt.Print("\n")
	fmt.Print("256 Color Mode - Foreground (ESC[38;5;Nm)")
	for i := 0; i < 256; i++ {
		if (i % 16) == 0 {
			fmt.Printf("\n   ")
		}
		v, ok := m[i]
		if ok && v {
			fmt.Print("\x1b[47;1m")
		} else {
			fmt.Print("\x1b[40;1m")
		}
		fmt.Printf("\x1b[38;5;%dm %3d \x1b[0m", i, i)
	}
	fmt.Print("\n")
	fmt.Print("\n")
}
