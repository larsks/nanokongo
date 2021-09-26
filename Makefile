SRCS = $(wildcard *.go) \
       $(wildcard */*.go) \
       actions/keycodes.go

GENERATED = $(EXENAME) \
	    actions/keycodes.go \
	    actions/keycodes.go.txt

all: nanokongo

test:
	golangci-lint run

nanokongo: $(SRCS)
	go build

actions/keycodes.go: actions/keycodes.go.txt
	go generate ./...

actions/keycodes.go.txt:
	curl -o $@ -sfL https://raw.githubusercontent.com/bendahl/uinput/master/keycodes.go

clean:
	rm -f $(GENERATED)
