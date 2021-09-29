PKG = github.com/larsks/nanokongo
EXENAME=nanokongo-$(shell go env GOOS)-$(shell go env GOARCH)

SRCS = $(wildcard *.go) \
       $(wildcard */*.go) \
       actions/keycodes.go

GENERATED = build/$(EXENAME) \
	    actions/keycodes.go \
	    actions/keycodes.go.txt

VERSION = $(shell git describe --tags --exact-match 2> /dev/null || echo unknown)
COMMIT = $(shell git rev-parse --short=10 HEAD)
DATE = $(shell date -u +"%Y-%m-%dT%H:%M:%S")

GOLDFLAGS = \
	    -X '$(PKG)/version.BuildVersion=$(VERSION)' \
	    -X '$(PKG)/version.BuildRef=$(COMMIT)' \
	    -X '$(PKG)/version.BuildDate=$(DATE)'

all: build/$(EXENAME)

test:
	golangci-lint run

build/$(EXENAME): $(SRCS)
	go build -o $@ -ldflags "$(GOLDFLAGS)"

actions/keycodes.go: actions/keycodes.go.txt
	go generate ./...

actions/keycodes.go.txt:
	curl -o $@ -sfL https://raw.githubusercontent.com/bendahl/uinput/master/keycodes.go

clean:
	rm -f $(GENERATED)
