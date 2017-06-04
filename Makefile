#
# Make csdiff for different platforms.
# Assumes that gnu make is present.
#
# Just type "make" to build, install and
# test. Install will install csdiff in
# a platform dependent directory under
# bin.
#
OS = $(shell uname -s)
MACH = $(shell uname -m)
OS_DIR = $(OS)-$(MACH)
BIN_DIR = bin/$(OS_DIR)
PROG = $(BIN_DIR)/csdiff

all: $(PROG) test

build: $(PROG)

clean:
	find . -type f -name '*~' -delete
	rm -rf bin pkg

$(PROG): jlinoff/csdiff

jlinoff/csdiff:
	go version
	GOPATH=$$(pwd) go install jlinoff/termcolors
	GOPATH=$$(pwd) GOBIN=$$(pwd)/$(BIN_DIR) go install $@

.PHONY: edit
edit:
	GOPATH=$$(pwd) open -a /opt/atom/latest/Atom.app/Contents/MacOS/Atom .

.PHONY: test
test:
	@cd test && make
