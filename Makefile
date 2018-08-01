.PHONY: build

PKGS=`go list ./... 2> /dev/null | grep -v '/vendor'`
LOCALS=`find . -type f -name '*.go' -not -path "./vendor*/*"`


all: fmt deps test

fmt:
	@go list golang.org/x/tools/cmd/goimports || go get golang.org/x/tools/cmd/goimports
	goimports -w $(LOCALS)
	go generate -x ./...

deps:
	@go list github.com/pointlander/peg || go get github.com/pointlander/peg
	go get .
	go vet .
	dep ensure

test: fmt deps
	go test $(PKGS)