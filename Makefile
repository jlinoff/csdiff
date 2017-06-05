#
# Make csdiff for different platforms.
# Assumes that gnu make is present.
#
# Just type "make" to build, install and test. Install will install
# csdiff in a platform dependent directory under bin.
#
# If you are on a Mac and want to build the linux version, make sure
# that you have Docker installed and type the following command.
#
#   make "linux"
#
# You can use the make "linux" command on linux as well if you have
# docker but do not have go installed.
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

# Special case, build linux version on Mac using
# docker. If you are on linux, just type "make".
.PHONY: linux
linux: 
	@docker images
	@docker images goco | grep latest
	@if ! docker images goco | grep latest >/dev/null ; then \
		echo "INFO: creating goco docker image" ; \
		echo "INFO: docker build --build-arg gover=1.8.3 -f Dockerfile -t goco:1.8.3 -t goco:latest ." ; \
		docker build --build-arg gover=1.8.3 -f Dockerfile -t goco:1.8.3 -t goco:latest . ; \
	fi
	docker run -it --rm -v $$(pwd):/opt/go/project goco make
	@echo "INFO: created bin/Linux-x86_64/csdiff"
	@echo "INFO: done"
