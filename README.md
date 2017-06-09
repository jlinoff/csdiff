# csdiff
[![Releases](https://img.shields.io/github/release/jlinoff/csdiff.svg?style=flat)](https://github.com/jlinoff/csdiff/releases)

side by side diff tool with colorization and regular expression support for filtering

## Contents
1. [Introduction](#introduction)
2. [Using It](#using)
3. [Colors](#colors)
4. [Command Line Options](#cliopts)
5. [Installation](#installation)
6. [Modification](#modification)
7. [About](#about)
8. [TODO](#todo)

<a name="introduction"></a>
## Introduction
Command line tool that does a side by side diff of two text files with [regular expression](https://github.com/google/re2/wiki/Syntax) filtering
and [ANSI terminal colorization](https://en.wikipedia.org/wiki/ANSI_escape_code).

It is useful for analyzing text files that have patterns like timestamps that can easily be filtered out.

<a name="using"></a>
## Using It
The tool is very similar to sdiff, [wdiff](https://www.gnu.org/software/wdiff/) or diff. You specify two files or substitutions and it outputs the differences between them using ANSI
terminal colorization to highlight the differences. 

You can specify replacement patterns to mask or filter differences that you don't care about,
like line numbers or time stamps. It is very similar to using process substitution in sdiff
but the regular expression patterns are defined by go. Doing it with sdiff would look something
like this (with no colorization).
```bash
$ sdiff \
    <(cat test/td03.txt | sed -E -e 's@[0-9][0-9]:[0-9][0-9]:[0-9][0-9]@HH:MM:SS@g') \
    <(cat test/td04.txt | sed -E -e 's@[0-9][0-9]:[0-9][0-9]:[0-9][0-9]@HH:MM:SS@g' )
```

The default colorization is very simple, a light grey background is used to highlight differences.
A background color was chosen so that it would be easy to see extra spaces.
You can customize the colors.  

Here is an example that shows two files whose content only differ by a time stamp.
```bash
$ bin/Darwin-x86_64/csdiff test/td03.txt test/td04.txt
```
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<img src="https://cloud.githubusercontent.com/assets/2991242/26766790/21bf7818-494d-11e7-88c2-84eea6022a0e.png" alt="example-1">

You can see the differences in the vertical, light grey stripes.

Here is an example that shows how to use the replacement option (-r) to ignore the timestamp differences. As you can see, the go regular expressions are a bit more concise than the previous `sed` example.
```bash
$ bin/Darwin-x86_64/csdiff \
    -r '\d{4}-\d{2}-\d{2}\s+\d{2}:\d{2}:\d{2}' 'yyyy-mm-dd HH:MM:SS' \
    test/td03.txt test/td04.txt
```
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<img src="https://cloud.githubusercontent.com/assets/2991242/26766793/2d0d2530-494d-11e7-849b-a03bec7a1a5c.png" alt="example-2">

Now there are no differences because they were masked.

<a name="colors"></a>
## Colors
You have the option of customizing the output colors based on the type of data by specifying colormaps for different
types of output.

Here is an example of how that is done.

```bash
$ bin/Darwin-x86_64/csdiff -c cd=bold,fgRed \
                           -c cm=bold,fgBlue \
                           -c s=bold,fgMagenta \
                           -c lm=bold,fgGreen \
                           -c llo=bold,fgCyan \
                           -c rlo=bold,fgCyan \
                           test/td01.txt test/td02.txt
```
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<img src="https://cloud.githubusercontent.com/assets/2991242/26766795/32be864a-494d-11e7-9b37-1554c4821494.png" alt="example-3">

Contrast that with the default below and you can see the differences.
```bash
$ bin/Darwin-x86_64/csdiff test/td01.txt test/td02.txt
```
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<img src="https://cloud.githubusercontent.com/assets/2991242/26789827/5df9308e-49c6-11e7-9c3e-4426b5f31f7e.png" alt="example-4">

As you can see, there are many more colors in the first example.

### Specifying Color
Color is specified by defining a color map for colorizable entities. The syntax looks like this.

```
   -c cd=bold,fgRed
   ^  ^ ^^    ^
   |  | ||    +--- Another colorization value or text modifier, you can specify as many as you like,
   |  | |+-------- A colorization value or text modifier (described below).
   |  | +--------- Equals sign that separates the tag from the attributes.
   |  +----------- Tag (described below)
   +-------------- Option that specifies a color map.
```

Note that you can have multiple color maps on the same line by separating them with a semi-colon as show below.
```
   -c 'cd=bold,fgRed;cm=bold,fgBlue'
```

### Colorizable Entities (Tags)
These tags describe the data that can be colored.

| Entity        | Abbreviation | Description |
| ------------- | ------------ | ----------- |
| CharsMatch    | cm  | Color of characters that match on both lines when there are differences. |
| CharsDiff     | cd  | Color of characters that differ on both lines. |
| LineMatch     | lm  | Color of characters when both lines match. |
| LeftLineOnly  | llo | Color of characters when there is no right line. |
| RightLineOnly | rlo | Color of characters when there is no left line. |
| Symbol        | sym | Color of the sdiff symbol in the middle. The symbol is &vert;, &lt;, &gt; or nothing. |

### Symbols
These are the symbols that csdiff inserts between the lines. They cannot be changed.

| Symbol | Description |
| :----: | ----------- |
| &vert; | There are differences between the two lines. |
| &lt;   | There is no matching right line, the left line was inserted. |
| &gt;   | There is no matching left line, the right line was inserted. |
| &nbsp; | There is no difference between the lines. |

### Text (Foreground) Colorization Values
These are the values that you can specify to color text.
It is a subset of the ANSI terminal colors that should work
everywhere.
They are *not* case sensitive.

1. fgDefault
2. fgBlack
3. fgBlue
4. fgCyan
5. fgGreen
6. fgMagenta
7. fgRed
8. fgYellow
9. fgWhite
10. fgLightBlue
11. fgLightCyan
12. fgLightGreen
13. fgLightGrey
14. fgLightMagenta
15. fgLightRed
16. fgLightYellow
17. fgDarkGrey

### Background Colorization Values
These are the values that you can specify for the background.
It is a subset of the ANSI terminal colors that should work
everywhere.
They are *not* case sensitive.

1. bgDefault
2. bgBlack
3. bgBlue
4. bgCyan
5. bgGreen
6. bgMagenta
7. bgRed
8. bgYellow
9. bgLightBlue
10. bgLightCyan
11. bgLightGreen
12. bgLightGrey
13. bgLightMagenta
14. bgLightRed
15. bgLightYellow
16. bgDarkGrey

### Text Modifiers
These are the values that you can specify to change how the text is displayed.
They are *not* case sensitive.

| Tag       | Description |
| --------- | ----------- |
| bold      | Make the text bold. |
| dim       | Make the text dim. |
| underline | Underline the text. |
| blink     | Blink. |
| reverse   | Reverse the foreground and background. |
| reset     | Reset the colors and modifiers. |

### 256 Color ANSI Colors
You can specify a much wider range of colors using the ANSI 256 color mode colors using: `fg256[N]` and `bg256[N]`
where N is a number between in the range [0..255]. To see all of the colors available use the `--256` option.

```bash
$ csdiff --256
```
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<img src="https://user-images.githubusercontent.com/2991242/26957671-6ac182e6-4c7b-11e7-80c3-495d098ad811.png" alt="example-256">

> The `fg256[N]` and `bg256[N]` color options are only available in version 0.5.x or later.

<a name="cliopts"></a>
## Command Line Options
The command line options are best accessed by looking at the inline help because they may change
over time but these are the command line options available in 0.4.x.

| Long Option           | Short Option    | Brief Description |
| --------------------- | --------------- | ----------------- |
| --256                 | NONE            | Print the 256 color ANSI color map values. |
| --color-map COLOR_MAP | --c COLOR_MAP   | Specify a color map for a tag. |
| --clear               | NONE            | Clear the default color map. |
| --config FILE         | NONE            | Specify a color map config file. |
| --help                | -h              | Inline help. |
| --diff                | -d              | Do a traditional diff. Useful for very long lines. |
| --no-color            | -n              | Turn off colorization. Used for testing. |
| --replace PATT REP    | -r PATT REP     | Specify a pattern to replace. Can be specified multiple times. |
| --suppress            | -s              | Suppress common lines. |
| --version             | -V              | Print the program version and exit. |
| --width NUM           | -w NUM          | The width of the output. The default is the width of the terminal. |

<a name="installation"></a>
## Installation
Just download the tar or zip image for your system and extract the executable. You can use it directly
or copy it.

### Linux Example
```bash
$ curl -s -k -L -O https://github.com/jlinoff/csdiff/releases/download/v0.4.2/csdiff-x64.tar.gz
$ tar ztvf csdiff-x64.tar.gz 
-rwxr-xr-x  1 jlinoff 1784424920 2319784 Jun  5 11:02 csdiff
```

### Mac Example
```bash
$ curl -s -k -L -O https://github.com/jlinoff/csdiff/releases/download/v0.4.2/csdiff-mac.zip
$ unzip -l csdiff-mac.zip
Archive:  csdiff-mac.zip
  Length      Date    Time    Name
---------  ---------- -----   ----
  2301184  06-05-2017 11:02   csdiff
---------                     -------
  2301184                     1 file
```

<a name="modification"></a>
## Modification (docker)
Here is how you check out the source to modify it.
If you make changes or implement bug fixes that you think might be useful,
please let me know so that I can incorporate them.

```bash
$ git clone https://github.com/jlinoff/csdiff.git
```

### Dockerfile
Note that I used the following [Dockerfile](https://docs.docker.com/engine/reference/builder/)
to create the linux image because I was working on a Mac.

```
# This docker file creates a go compilation container that can be used
# to cross-compile go for linux on any platform that supports docker.
#
#   $ cd csdiff
#   $ docker run -it --rm -v $(pwd):/opt/go/project goco make
#
# To create it. Note that gover is an ARG that defines the version of
# go that you want to build. This example shows how to build it for go-1.8.3.
#
#   $ docker build --build-arg gover=1.8.3 -f Dockerfile -t goco:1.8.3 -t goco:latest .
FROM centos:latest

RUN yum clean all && yum update -y && yum install -y git make

ARG gover
ENV GO_VERSION=$gover
ENV GOROOT=/opt/go/latest
ENV GOPATH=/opt/go/project
ENV GO_PROG=/opt/go/latest/bin/go

# Setup the volume.
RUN mkdir -p ${GOPATH}
VOLUME ${GOPATH}

# Install go in /opt/go
RUN mkdir -p /opt/go/${GO_VERSION}/dl && \
    cd /opt/go/${GO_VERSION}/dl && \
    curl -k -O -L https://storage.googleapis.com/golang/go${GO_VERSION}.linux-amd64.tar.gz && \
    cd /opt/go/${GO_VERSION} && \
    tar zxf dl/go${GO_VERSION}.linux-amd64.tar.gz && \
    ln -s /opt/go/${GO_VERSION}/go /opt/go/latest && \
    ${GO_PROG} version
    
# Install golint
RUN cd /opt/go/${GO_VERSION}/dl && \
    GOPATH=/opt/go/${GO_VERSION}/dl ${GOROOT}/bin/go get -u github.com/golang/lint/golint && \
    cp bin/* ${GOROOT}/bin

# Wrapper for the go command that makes it
# natural for the user to run something like
# docker run -it --rm -v $(pwd):/opt/go/project goco go build myprog.go
RUN /bin/echo '#!/bin/bash'                           > /opt/go/goco.sh && \
    /bin/echo 'export PATH="${GOROOT}/bin:${PATH}"'  >> /opt/go/goco.sh && \
    /bin/echo 'cd /opt/go/project'                   >> /opt/go/goco.sh && \
    /bin/echo '$*'                                   >> /opt/go/goco.sh && \
    chmod a+rx /opt/go/goco.sh && \
    /opt/go/goco.sh go version

# Run in go environment.
ENTRYPOINT ["/opt/go/goco.sh"]
CMD ["/opt/go/latest/bin/go", "version"]
```

### Using the goco docker image
This is how you use the goco docker image to compile csdiff for linux.
```bash
$ cd csdiff
$ docker run -it --rm -v $(pwd):/opt/go/project goco make
```
The binary will appear in bin/Linux-x86_64/csdiff.

<a name="about"></a>
## About

I originally developed this because I was working on a Mac and FileMerge did not provide line numbers
but then added regular expression filtering which also made it useful on linux.

This tool was written in go-1.8.3 using the [atom](https://atom.io/) and [emacs](https://emacsformacosx.com/)
editors on Mac OSX 10.12.5.

It has been tested on Mac OS X 10.12.5, CentOS 7 and CentOS 6.

Suggestions and improvements are greatly appreciated.

<a name="todo"></a>
