PWD = $(shell pwd)
ORG_PATH = github.com/mingqing
REPO_PATH = ${ORG_PATH}/toolpub

BIN = toolpub
ARGS = ""
BUILD_OS_TARGETS = "linux"

GOROOT := /opt/17173/go
GOPATH := ${PWD}/gopath

all: clean build

build:
	mkdir build
	mkdir -p ${GOPATH}/src/${ORG_PATH}
	ln -s ${PWD} ${GOPATH}/src/${REPO_PATH}
	GOPATH=$(GOPATH) go build -o build/${BIN}

run: build
	./build/${BIN}

clean:
	rm -rf build
	rm -f ${GOPATH}/src/${REPO_PATH}
