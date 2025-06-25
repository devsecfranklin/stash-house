BINARY	:= gfu
PKGS	:= $(shell go list ./... |grep -v /vendor)

# version checks git tags. if any errors, they are discarded, and the version is
# read from the VERSION file
VERSION	:= $(shell git describe --tags 2> /dev/null)
ifeq ($(VERSION),)	
	VERSION := $(shell cat VERSION)
endif
BUILD		:= $(shell date -Iminutes)

#ldflags for versioning info: https://stackoverflow.com/questions/11354518/application-auto-build-versioning
#-s flag info: https://blog.spiralscout.com/using-w-and-s-flags-in-golang-97ae59b50e26
LDFLAGS	:= -ldflags "-s -X main.version=${VERSION} -X main.build=${BUILD}"

build:
	go build ${LDFLAGS} -o ${BINARY}

.PHONY: install
install:
	go install ${LDFLAGS}

.PHONY: uninstall
uninstall:
	go clean -i

.PHONY: clean
clean:
	go clean
