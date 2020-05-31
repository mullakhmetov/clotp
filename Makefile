.PHONY: \
	build \
	test \
	vet \
	fmt \
	fmtcheck

SRCS = $(shell git ls-files '*.go')
PKGS = $(shell go list ./... | grep -v /vendor/)

build:
	go build -o clotp

lint:
	golangci-lint run ./...

test:
	go test -race -mod=vendor -timeout=60s -count 1 ./...

vet:
	go vet $(PKGS) || exit;

fmt:
	gofmt -w $(SRCS)

fmtcheck:
	@ $(foreach file,$(SRCS),gofmt -s -l $(file);)