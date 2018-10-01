.PHONY: build

PKGS        := $(shell go list ./... 2> /dev/null | grep -v '/vendor')
LOCALS      := $(shell find . -type f -name '*.go' -not -path "./vendor*/*")

.EXPORT_ALL_VARIABLES:
GO111MODULE  = on

all: fmt deps test

fmt:
	gofmt -w $(LOCALS)
	go generate -x ./...

deps:
	@go list github.com/pointlander/peg || go get github.com/pointlander/peg
	go get ./...
	go vet ./...

test: fmt deps
	go test ./...