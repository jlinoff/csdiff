#!/bin/bash
#
# TODO: create a simple test framework.
#

# ================================================================
# Includes
# ================================================================
Location="$(cd $(dirname $0) && pwd)"
source $Location/utils.sh

# ================================================================
# Main
# ================================================================
OS=$(uname -s)
MACH=$(uname -m)
OS_DIR=${OS}-${MACH}
BIN_DIR=../bin/${OS_DIR}
PROG=${BIN_DIR}/csdiff

# ================================================================
# Tests
# ================================================================
utilsExec ${PROG} -h
utilsExec ${PROG} -V
utilsExec ${PROG} -n td01.txt td02.txt
utilsExec ${PROG} td01.txt td02.txt
utilsExec ${PROG} \
          -c cd=bgYellow,bold,blink,underline,fgRed \
          -c cm=bgGreen,bold,fgblue \
	  -c s=bgBlack,bold,fgRed \
	  -c lm=bgBlack,bold,fgYellow \
          td01.txt td02.txt
utilsExec ${PROG} \
          -c cd=bold,fgRed \
          -c cm=fgDefault \
	  -c s=fgDefault \
          td01.txt td02.txt
utilsExec ${PROG} \
          -c cd=bold,fgRed \
          -c cm=fgDefault \
	  -c s=fgDefault \
          -c lm=bold,fgGreen \
          td01.txt td02.txt
utilsExec ${PROG} \
          -c cd=bold,fgRed \
          -c cm=bold,fgBlue \
	  -c s=bold,fgMagenta \
          -c lm=bold,fgGreen \
          td01.txt td02.txt
utilsExec ${PROG} \
          -c cd=bold,fgRed \
          -c cm=bold,fgBlue \
	  -c s=bold,fgMagenta \
          -c lm=bold,fgGreen \
          -c llo=bold,fgCyan \
          -c rlo=bold,fgCyan \
          td03.txt td04.txt
utilsExec ${PROG} td03.txt td04.txt
utilsExec ${PROG} --config test.conf td03.txt td04.txt
utilsExec ${PROG} --diff td03.txt td04.txt
utilsExec ${PROG} --summary --config test.conf td03.txt td04.txt
utilsExec ${PROG} --summary --diff td03.txt td04.txt
utilsExec ${PROG} td03.txt td04.txt
utilsExec ${PROG} -r "'\d{4}-\d{2}-\d{2}\s+\d{2}:\d{2}:\d{2}'" "'yyyy-mm-dd HH:MM:SS'" td03.txt td04.txt
utilsExec ${PROG} td02.txt td05.txt
utilsExec ${PROG} -d td02.txt td05.txt

utilsInfo "done"
