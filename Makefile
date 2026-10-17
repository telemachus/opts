.DEFAULT_GOAL := test

fmt:
	golangci-lint fmt --no-config -Egofumpt

staticcheck:
	staticcheck ./...

revive:
	revive -config revive.toml ./...

golangci:
	golangci-lint run

lint: staticcheck revive golangci

build: testr lint
	go build .

test:
	go test -shuffle on .

testr:
	go test -race -shuffle on .

testv:
	go test -shuffle on -v .

clean:
	go clean

.PHONY: fmt staticcheck revive golangci lint
.PHONY: test testv testr
.PHONY: build clean
