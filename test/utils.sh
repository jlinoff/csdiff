#!/bin/bash
#
# Utilities for setup scripts.
#

# ================================================================
# Functions
# ================================================================
# Print an info message to stdout.
function utilsPrefix() {
    printf "%-28s %s %5s: " "$(date +'%Y-%m-%d %H:%M:%S %z %Z')" "$1" ${BASH_LINENO[$2]}
    shift ; shift
    echo "$*"
}

# Print an info message to stdout.
function utilsInfo() {
    utilsPrefix INFO 1 "$*"
}

# Print an error message to stderr and exit.
function utilsErr() {
    utilsPrefix ERROR 1 "$*" >&2
    exit 1
}

# Print an error message to stderr.
function utilsErrNoExit() {
    utilsPrefix ERROR 1 "$*" >&2
}

# Decorate a command and exit if the return code is not zero.
function utilsExec() {
    local Cmd="$*"
    echo
    utilsPrefix INFO 1 "cmd.cmd=$Cmd"
    eval "$Cmd"
    local Status=$?
    utilsPrefix INFO 1 "cmd.code=$Status"
    if (( $Status )) ; then
        utilsPrefix INFO 1 "cmd.status=FAILED"
        exit 1
    else
        utilsPrefix INFO 1 "cmd.status=PASSED"
    fi
}

# Decorate a command, do not exit on error.
function utilsExecNoExit() {
    local Cmd="$*"
    echo
    utilsPrefix INFO 1 "cmd.cmd=$Cmd"
    eval "$Cmd"
    local Status=$?
    utilsPrefix INFO 1 "cmd.code=$Status"
    if (( $Status )) ; then
        utilsPrefix INFO 1 "cmd.status=FAILED"
        return 1
    else
        utilsPrefix INFO 1 "cmd.status=PASSED"
        return 0
    fi
}

# Banner.
function utilsBanner() {
    echo
    echo "# ================================================================"
    echo "# $*"
    echo "# ================================================================"
}

