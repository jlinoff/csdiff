package main

import (
	"fmt"
	"jlinoff/termcolors"
	"os"
	"path/filepath"
	"strings"
)

// help for program usage
//    Diff two files side by side allowing for colorized output using
//    the longest common subsequence algorithm for the line comparisons
//    and a recursive longest common substring algorithm for the
//    character differences between two lines.
func help() {
	f := `
USAGE
   %[1]v [OPTIONS] FILE1 FILE2

DESCRIPTION
    Command line tool that does a side by side diff of two text files
    with regular expression filtering and ANSI terminal colorization.

    It is useful for analyzing text files that have patterns like
    timestamps that can easily be filtered out.

    The file output is side by side. Here is a simple example. The
    width was truncated to 90 characters. Normally the width is the
    size of the terminal window.

        $ cat f1.txt
        start
        Lorem ipsum dolor sit amet, consectetur
        adipiscing elit, sed do eiusmod tempor
        incididunt ut labore et dolore magna
        aliqua. Ut enim ad minim veniam, quis
        nostrud exercitation ullamco laboris
        nisi ut aliquip ex ea commodo consequat.

        $ cat f2.txt
        prefix
        Lorem ipsum dolor sit amet, consectetur
        adipiscing elit, sed do eiusmod tempor
        incididunt ut labore et dolore magna
        aliqua. Ut enim ad minim veniam, quis
        nostrud exercitation ullamco laboris
        infix
        nisi ut aliquip ex ea commodo consequat.
        suffix

        $ %[1]v -w 90 f1.txt f2.txt

               test/file04a.txt                              test/file05a.txt
             1 [47msta[0mr[47mtm[0m                                [31;1m|[0m     1 [47mp[0mr[47mefix[0m
             2 Lorem ipsum dolor sit amet, consect$        2 Lorem ipsum dolor sit amet, consect$
             3 adipiscing elit, sed do eiusmod tem$        3 adipiscing elit, sed do eiusmod tem$
             4 incididunt ut labore et dolore magna        4 incididunt ut labore et dolore magna
             5 aliqua. Ut enim ad minim veniam, qu$        5 aliqua. Ut enim ad minim veniam, qu$
             6 nostrud exercitation ullamco laboris        6 nostrud exercitation ullamco laboris
                                                    [31;1m>[0m      7 [47minfix  [0m
             7 nisi ut aliquip ex ea commodo conse$        8 nisi ut aliquip ex ea commodo conse$
                                                    [31;1m>[0m      9 [47msuffix[0m

    Note that truncated lines have a $ as the last character.

    For comparing files that have time stamps or other regular
    patterns, you can use the -r (--replace) option to replace
    then with a common value. Here is a simple example that
    replaces dates of the form YYYY-MM-DD (ex. 2017-01-02)
    with the string 'YYYY-MM-DD'.

        $ %[1]v -r '\d{4}-\d{2}-\d{2}' 'YYYY-MM-DD' file1 file2

    Date differences will be ignored.

    Another illustration of replacement use would be to mask the
    difference between start and prefix in the earlier example. If you
    specify -r 'start' 'prefix', the first lines will match.

        $ %[1]v -w 90 f1.txt f2.txt

               test/file04a.txt                              test/file05a.txt
             1 prefix                                      1 prefix
             2 Lorem ipsum dolor sit amet, consect$        2 Lorem ipsum dolor sit amet, consect$
             3 adipiscing elit, sed do eiusmod tem$        3 adipiscing elit, sed do eiusmod tem$
             4 incididunt ut labore et dolore magna        4 incididunt ut labore et dolore magna
             5 aliqua. Ut enim ad minim veniam, qu$        5 aliqua. Ut enim ad minim veniam, qu$
             6 nostrud exercitation ullamco laboris        6 nostrud exercitation ullamco laboris
                                                    [31;1m>[0m      7 [47minfix  [0m
             7 nisi ut aliquip ex ea commodo conse$        8 nisi ut aliquip ex ea commodo conse$
                                                    [31;1m>[0m      9 [47msuffix[0m

    As you can see the first line now matches because start was
    replaced by prefix.

OPTIONS
    --256       Print the ANSI terminal 256 color table color
                values for foreground and background and exit.
                This is useful for determing which extended
                colors work for your terminals.

    -c COLOR_VAL, --color-map COLOR_VAL
                Specify a color value for a diff condition.
                The syntax is COND=ATTR1[,[ATTR2[,ATTR3]]].
                Multiple conditions can be specified by semi-colons.
                The available conditions are.

                   CharsMatch     cm    Chars match on both lines.
                   CharsDiff      cd    Chars differ on both lines.
                   LineMatch      lm    Color when both lines match.
                   LeftLineOnly   llo   Only the left line, no right.
                   RightLineOnly  rlo   Only the right line, no left.
                   Symbol         sym   The line diff symbol.

                The conditions are case insensitive so diff could be
                specified as Diff, diff, or d.

                Here is an example that shows the default settings.

                   -c cd=bgLightGrey
                   -c cm=bgDefault
                   -c symbol=bold,fgRed
                   -c llo=bgLightGrey
                   -c rlo=bgLightGrey

                Note you also specify them like this using the
                semi-colon separator.

                   -c 'cd=bgLightGrey;cm=bgDefault;sym=bold,fgRed;llo=bgLightGrey;rlo=bgLightGrey'

                These are the available foreground colors (case
                insensitive):

                   fgDefault

                   fgBlack          [30mXYZ[0m
                   fgBlue           [34mXYZ[0m
                   fgCyan           [36mXYZ[0m
                   fgGreen          [32mXYZ[0m
                   fgMagenta        [35mXYZ[0m
                   fgRed            [31mXYZ[0m
                   fgYellow         [33mXYZ[0m
                   fgWhite          [97mXYZ[0m

                   fgLightBlue      [94mXYZ[0m
                   fgLightCyan      [96mXYZ[0m
                   fgLightGreen     [92mXYZ[0m
                   fgLightGrey      [37mXYZ[0m
                   fgLightMagenta   [95mXYZ[0m
                   fgLightRed       [91mXYZ[0m
                   fgLightYellow    [93mXYZ[0m

                   fgDarkGrey       [90mXYZ[0m

                These are the available background colors (case
                insensitive):

                   bgDefault

                   bgBlack          [40;97mXYZ[0m
                   bgBlue           [44;97mXYZ[0m
                   bgCyan           [46mXYZ[0m
                   bgGreen          [42mXYZ[0m
                   bgMagenta        [45mXYZ[0m
                   bgRed            [41;97mXYZ[0m
                   bgYellow         [43mXYZ[0m

                   bgLightBlue      [104mXYZ[0m
                   bgLightCyan      [106mXYZ[0m
                   bgLightGreen     [102mXYZ[0m
                   bgLightGrey      [47mXYZ[0m
                   bgLightMagenta   [105mXYZ[0m
                   bgLightRed       [101mXYZ[0m
                   bgLightYellow    [103mXYZ[0m

                   bgDarkGrey       [100mXYZ[0m

                The following foreground attribute modifiers are
                available (case insensitive):

                   bold             [1mXYZ[0m
                   dim              [2mXYZ[0m
                   underline        [4mXYZ[0m
                   blink            [5mXYZ[0m
                   reverse          [7mXYZ[0m

                Here is another example of color combinations.

                   -c cd=bold,fgRed
                   -c cm=bold,fgBlue
                   -c s=bold,fgMagenta
                   -c lm=bold,fgGreen
                   -c llo=bold,fgCyan
                   -c rlo=bold,fgCyan

    --clear     Clear the default settings. This is useful when you
                want to create a new color map. It is the same setting
                all color map fields to fgDefault.

    --config FILE
                Read a configuration files with color map data. There
                is one color map specification per line. Blank lines
                and lines that start with # are ignored. White space
                is allowed and the values are case insensitive.

                Here is an example.

                  # This is an example color map file.
                  cd  = bgLightGrey
                  cm  = bold, fgBlue
                  sym = bold, fgMagenta
                  lm  = bold, fgGreen
                  llo = bgLightGrey
                  rlo = bgLightGrey

    -d, --diff  Don't do the side by side diff. Use separate lines.
                This is similar to the standard diff output. This
                option always suppresses common lines.

    -h, --help  This help message.

    -n, --no-color
               Turn off color mode. This option really isn't useful
               because tools like sdiff are much faster. It was only
               made available for testing.

    -r PATTERN REPLACEMENT, --replace PATTERN REPLACEMENT
               Replace regular expression pattern PATTERN with
               REPLACEMENT where PATTERN is a regular expression that
               is recognized by the go regexp package. You can find
               out more information about the syntax here.

                   https://github.com/google/re2/wiki/Syntax

    -s, --suppress
               Suppress common lines.

    -V, --version
               Print the program version and exit.

    -w INT, --width INT
               The width of the output. The default is %[3]v.

EXAMPLES
    # Example 1. help
    $ %[1]v -h

    # Example 2: Diff two files with colors
    $ %[1]v file1 file2

    # Example 3: Diff two files with no colors.
    $ %[1]v -n file1 file2

    # Example 4: Diff two files with custom colors.
    $ %[1]v -c cd=bgYellow,bold,blink,underline,fgRed \
        -c cm=bgGreen,bold,fgBlue \
        -c s=bgBlack,bold,fgRed \
        -c lm=bgBlack,bold,fgYellow \
        file1 file2

    # Example 5: Diff two files, ignore timestamp differences
    #            where the time stamp looks like this:
    #                2017-10-02 23:01:57.2734743
    $ %[1]v -r '\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d+' 'yyyy-mm-dd HH:MM:SS.ssssss'

    # Example 6: Show the extended colors available from the
    #            ANSI 256 color terminal tables.
    $ %[1]v --256

VERSION
    v%[2]v

PROJECT
    https://github.com/jlinoff/csdiff

COPYRIGHT
    Copyright (c) 2017 by Joe Linoff

LICENSE
    MIT Open Source
  `
	f = "\n" + strings.TrimSpace(f) + "\n\n"
	fmt.Printf(f, filepath.Base(os.Args[0]), version, termcolors.MakeTermInfo().Width())
	os.Exit(0)
}
