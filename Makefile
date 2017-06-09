#
# Make csdiff for different platforms.
# Assumes that gnu make is present.
#
# Just type "make" to build, install and test. Install will install
# csdiff in a platform dependent directory under bin.
#
# If you want to build for all platforms, type "make rel".
#
# If you want to use the docker image for linux, type: "make docker-linux".
#
GOOS = $(shell uname -s | tr '[:upper:]' '[:lower:]')

all: $(GOOS) test

# Make the release executables.
rel: allplats rel/csdiff-darwin_amd64.zip \
	rel/csdiff-linux_amd64.tar.gz rel/csdiff-linux_amd64.zip

build: $(GOOS)

clean:
	find . -type f -name '*~' -delete
	rm -rf bin pkg rel

# Make different platforms.
.PHONY: allplats darwin linux windows
allplats: darwin linux

darwin:
	@go version
	GOOS=$@ GOARCH=amd64 GOPATH=$$(pwd) go install -pkgdir $$(pwd)/pkg jlinoff/termcolors
	GOOS=$@ GOARCH=amd64 GOPATH=$$(pwd) GOBIN=$$(pwd)/bin/$@_amd64 go install jlinoff/csdiff

linux:
	@go version
	GOOS=$@ GOARCH=amd64 GOPATH=$$(pwd) go install -pkgdir $$(pwd)/pkg jlinoff/termcolors
	GOOS=$@ GOARCH=amd64 GOPATH=$$(pwd) go install -pkgdir $$(pwd)/pkg jlinoff/csdiff

windows:
	@go version
	GOOS=$@ GOARCH=amd64 GOPATH=$$(pwd) go install -pkgdir $$(pwd)/pkg jlinoff/termcolors
	GOOS=$@ GOARCH=amd64 GOPATH=$$(pwd) go install -pkgdir $$(pwd)/pkg jlinoff/csdiff

# Package the executable into a zip file for release.
rel/csdiff-%.zip: bin/%/csdiff
	@[ ! -d $(dir $@) ] && mkdir -p $(dir $@) || true
	@rm -f $@
	cp $< $(dir $@)/csdiff
	@chmod a+x $(dir $@)/csdiff
	@cd $(dir $@) && zip -r $(notdir $@) csdiff
	@rm -f $(dir $@)/csdiff
	@unzip -l $@

# Package the executable into a tar file for release.
rel/csdiff-%.tar.gz: bin/%/csdiff
	@[ ! -d $(dir $@) ] && mkdir -p $(dir $@) || true
	@rm -f $@
	@rm -f $(dir $@)/csdiff
	cp $< $(dir $@)/csdiff
	@cd $(dir $@) && tar zcvf $(notdir $@) csdiff
	@chmod a+x $(dir $@)/csdiff
	@rm -f $(dir $@)/csdiff
	@tar ztvf $@

.PHONY: edit
edit:
	GOPATH=$$(pwd) open -a /opt/atom/latest/Atom.app/Contents/MacOS/Atom .

.PHONY: test
test:
	@cd test && make

# Special case, build linux version on Mac using
# docker. If you are on linux, just type "make".
.PHONY: docker-linux
docker-linux: 
	@docker images
	@docker images goco | grep latest
	@if ! docker images goco | grep latest >/dev/null ; then \
	echo "INFO: creating goco docker image" ; \
	echo "INFO: docker build --build-arg gover=1.8.3 -f Dockerfile -t goco:1.8.3 -t goco:latest ." ; \
	docker build --build-arg gover=1.8.3 -f Dockerfile -t goco:1.8.3 -t goco:latest . ; \
	fi
	docker run -it --rm -v $$(pwd):/opt/go/project goco make
	@echo "INFO: done"
