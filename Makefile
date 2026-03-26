PKGS        := $(shell go list ./... 2> /dev/null | grep -v '/vendor')
LOCALS      := $(shell find . -type f -name '*.go' -not -path "./vendor*/*")
EXAMPLES    := $(shell ls -1 examples)

.EXPORT_ALL_VARIABLES:
GO111MODULE  = on

all: fmt deps test build

fmt:
	@go fmt ./...
	go generate -x ./...

deps:
	@which peg || go get github.com/pointlander/peg
	go get ./...
	-go mod tidy
	go vet ./...

test: fmt deps
	@go test ./...

bin:
	@mkdir $(@)
bin/friendscript:
	@go build -o $(@) ./cmd/friendscript/

build: bin bin/friendscript
	@which friendscript 2> /dev/null && cp -v bin/friendscript `which friendscript` || true

examples: $(EXAMPLES)

$(EXAMPLES):
	go build -o bin/example-$(basename $(@)) examples/$(@)/*.go

.PHONY: build bin/friendscript examples $(EXAMPLES)
