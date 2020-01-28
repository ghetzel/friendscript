PKGS        := $(shell go list ./... 2> /dev/null | grep -v '/vendor')
LOCALS      := $(shell find . -type f -name '*.go' -not -path "./vendor*/*")
EXAMPLES    := $(shell ls -1 examples)

.EXPORT_ALL_VARIABLES:
GO111MODULE  = on

all: fmt deps test

fmt:
	gofmt -w $(LOCALS)
	go generate -x ./...

deps:
	@go list github.com/pointlander/peg || go get github.com/pointlander/peg
	go get ./...
	-go mod tidy
	go vet ./...

test: fmt deps
	go test ./...

examples: $(EXAMPLES)

$(EXAMPLES):
	go build -o bin/example-$(basename $(@)) examples/$(@)/*.go

.PHONY: build examples $(EXAMPLES)