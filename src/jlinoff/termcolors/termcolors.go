// Terminal color tools. Could be a separate package if it is needed elsewhere.
// License: The MIT License (MIT)
// Copyright (c) 2017 Joe Linoff
package termcolors

import (
	"fmt"
	"syscall"
	"unsafe"
)

//func main() {
//	MakeTermInfo().TestTermInfo()
//}

// TermAnsiAttrType defines the terminal attribute contants type.
type TermAnsiAttrType int

// TermAnsiAttrTypes is a list of attr types.
type TermAnsiAttrTypes []TermAnsiAttrType

// Terminal attribute contants.
const (
	SetBold        TermAnsiAttrType = 1
	SetDim                          = 2
	SetUnderline                    = 4
	SetBlink                        = 5
	SetReverse                      = 7
	SetHidden                       = 8
	Reset                           = 0 // reset everything
	ResSetBold                      = 21
	ResetDim                        = 22
	ResetUnderline                  = 24
	ResetBlink                      = 25
	ResetReverse                    = 27
	ResetHidden                     = 28
	FgDefault                       = 39
	FgBlack                         = 30
	FgRed                           = 31
	FgGreen                         = 32
	FgYellow                        = 33
	FgBlue                          = 34
	FgMagenta                       = 35
	FgCyan                          = 36
	FgLightgray                     = 37
	FgLightgrey                     = 37
	FgDarkgray                      = 90
	FgDarkgrey                      = 90
	FgLightred                      = 91
	FgLightgreen                    = 92
	FgLightyellow                   = 93
	FgLightblue                     = 94
	FgLightmagenta                  = 95
	FgLightcyan                     = 96
	FgWhite                         = 97
	BgDefault                       = 49
	BgBlack                         = 40
	BgRed                           = 41
	BgGreen                         = 42
	BgYellow                        = 43
	BgBlue                          = 44
	BgMagenta                       = 45
	BgCyan                          = 46
	BgLightgray                     = 47
	BgLightgrey                     = 47
	BgDarkgray                      = 100
	BgDarkgrey                      = 100
	BgLightred                      = 101
	BgLightgreen                    = 102
	BgLightyellow                   = 103
	BgLightblue                     = 104
	BgLightmagenta                  = 105
	BgLightcyan                     = 106
	BgWhite                         = 107
)

// TermInfoType stores attributes about the terminal.
type TermInfoType struct {
	Rows         uint16
	Cols         uint16
	Xpixel       uint16
	Ypixel       uint16
	AttrsByName  map[string]TermAnsiAttrType
	AttrsByValue map[TermAnsiAttrType]string
}

// MakeTermInfo creates the terminal information
func MakeTermInfo() (ti TermInfoType) {
	retCode, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(&ti)))

	if int(retCode) == -1 {
		panic(errno)
	}

	// Create string map to the terminal attributes.
	ti.AttrsByName = map[string]TermAnsiAttrType{}
	ti.AttrsByName["SetBold"] = SetBold
	ti.AttrsByName["SetDim"] = SetDim
	ti.AttrsByName["SetUnderline"] = SetUnderline
	ti.AttrsByName["SetBlink"] = SetBlink
	ti.AttrsByName["SetReverse"] = SetReverse
	ti.AttrsByName["SetHidden"] = SetHidden
	ti.AttrsByName["Reset"] = Reset
	ti.AttrsByName["ResSetBold"] = ResSetBold
	ti.AttrsByName["ResetDim"] = ResetDim
	ti.AttrsByName["ResetUnderline"] = ResetUnderline
	ti.AttrsByName["ResetBlink"] = ResetBlink
	ti.AttrsByName["ResetReverse"] = ResetReverse
	ti.AttrsByName["ResetHidden"] = ResetHidden
	ti.AttrsByName["FgDefault"] = FgDefault
	ti.AttrsByName["FgBlack"] = FgBlack
	ti.AttrsByName["FgRed"] = FgRed
	ti.AttrsByName["FgGreen"] = FgGreen
	ti.AttrsByName["FgYellow"] = FgYellow
	ti.AttrsByName["FgBlue"] = FgBlue
	ti.AttrsByName["FgMagenta"] = FgMagenta
	ti.AttrsByName["FgCyan"] = FgCyan
	ti.AttrsByName["FgLightgray"] = FgLightgray
	ti.AttrsByName["FgLightgrey"] = FgLightgrey
	ti.AttrsByName["FgDarkgray"] = FgDarkgray
	ti.AttrsByName["FgDarkgrey"] = FgDarkgrey
	ti.AttrsByName["FgLightred"] = FgLightred
	ti.AttrsByName["FgLightgreen"] = FgLightgreen
	ti.AttrsByName["FgLightyellow"] = FgLightyellow
	ti.AttrsByName["FgLightblue"] = FgLightblue
	ti.AttrsByName["FgLightmagenta"] = FgLightmagenta
	ti.AttrsByName["FgLightcyan"] = FgLightcyan
	ti.AttrsByName["FgWhite"] = FgWhite
	ti.AttrsByName["BgDefault"] = BgDefault
	ti.AttrsByName["BgBlack"] = BgBlack
	ti.AttrsByName["BgRed"] = BgRed
	ti.AttrsByName["BgGreen"] = BgGreen
	ti.AttrsByName["BgYellow"] = BgYellow
	ti.AttrsByName["BgBlue"] = BgBlue
	ti.AttrsByName["BgMagenta"] = BgMagenta
	ti.AttrsByName["BgCyan"] = BgCyan
	ti.AttrsByName["BgLightgray"] = BgLightgray
	ti.AttrsByName["BgLightgrey"] = BgLightgrey
	ti.AttrsByName["BgDarkgray"] = BgDarkgray
	ti.AttrsByName["BgDarkgrey"] = BgDarkgrey
	ti.AttrsByName["BgLightred"] = BgLightred
	ti.AttrsByName["BgLightgreen"] = BgLightgreen
	ti.AttrsByName["BgLightyellow"] = BgLightyellow
	ti.AttrsByName["BgLightblue"] = BgLightblue
	ti.AttrsByName["BgLightmagenta"] = BgLightmagenta
	ti.AttrsByName["BgLightcyan"] = BgLightcyan
	ti.AttrsByName["BgWhite"] = BgWhite

	// Now populate the reverse map.
	ti.AttrsByValue = map[TermAnsiAttrType]string{}
	for v, k := range ti.AttrsByName {
		ti.AttrsByValue[k] = v
	}

	return ti
}

// TestTermInfo shows how to use this package.
// It displays the terminal width, height and a host of information
// about fg, bg and attribute pairs.
func (ti TermInfoType) TestTermInfo() {
	fmt.Printf("width  = %4v\n", ti.Width())
	fmt.Printf("height = %4v\n", ti.Height())

	fgs := []TermAnsiAttrType{FgDefault, FgBlack, FgRed, FgGreen, FgYellow, FgBlue, FgMagenta, FgCyan}
	bgs := []TermAnsiAttrType{BgDefault, BgBlack, BgRed, BgGreen, BgYellow, BgBlue, BgMagenta, BgCyan}
	as := []TermAnsiAttrType{Reset, SetBold, SetDim, SetUnderline, SetBlink, SetReverse}
	n := 0
	for _, bg := range bgs {
		for _, fg := range fgs {
			for _, a1 := range as {
				for _, a2 := range as {
					n++
					fmt.Printf("%5d ", n)
					fmt.Printf("%4d:%-12s  ", fg, ti.AttrsByValue[fg])
					fmt.Printf("%4d:%-12s  ", bg, ti.AttrsByValue[bg])
					fmt.Printf("%4d:%-12s  ", a1, ti.AttrsByValue[a1])
					fmt.Printf("%4d:%-12s  ", a2, ti.AttrsByValue[a2])
					if a1 != a2 {
						ti.Set(bg, a1, a2, fg)
					} else {
						ti.Set(bg, a1, fg)
					}
					fmt.Printf("%v", "Lorem ipsum dolor sit amet")
					ti.Reset()
					fmt.Println("")
				}
			}
		}
	}
}

// Width returns the width of the terminal.
func (ti TermInfoType) Width() uint16 {
	return ti.Cols
}

// Height returns the height of the terminal.
func (ti TermInfoType) Height() uint16 {
	return ti.Rows
}

// Set term characteristics.
// ti.Set(BgGreen, SetBold, FgRed)
// fmt.Println("Bold Red on Green")
// ti.Reset()  // same as ti.Set(Reset)
func (ti TermInfoType) Set(attrs ...TermAnsiAttrType) {
	fmt.Printf("\x1b[")
	for i, attr := range attrs {
		if i > 0 {
			fmt.Printf(";")
		}
		fmt.Printf("%d", attr)
	}
	fmt.Printf("m")
}

// Reset the terminal characteristics.
// It is shorthand for ti.Set(Reset).
func (ti TermInfoType) Reset() {
	ti.Set(Reset)
}
