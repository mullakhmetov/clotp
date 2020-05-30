.PHONY: \
	build
	test

build:
	go build -o clotp

test:
	go test -race -mod=vendor -timeout=60s -count 1 ./...