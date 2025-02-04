.DEFAULT_GOAL := test

PREFIX := $(HOME)/local/gitmirror

fmt:
	golangci-lint run --disable-all --no-config -Egofmt --fix
	golangci-lint run --disable-all --no-config -Egofumpt --fix

lint: fmt
	staticcheck .
	revive -config revive.toml .
	golangci-lint run

golangci: fmt
	golangci-lint run

staticcheck: fmt
	staticcheck .

revive: fmt
	revive -config revive.toml ./...

test:
	go test -shuffle on

testv:
	go test -shuffle on -v

testr:
	go test -race -shuffle on

build: lint testr
	go build .

install: build
	go install .

clean:
	go clean -i -r -cache

.PHONY: fmt lint build install test testv testr clean
