.PHONY: test build

VERSION = $(shell grep "const Version =" cmd/sub/version.go | grep "const Version =" | sed -e 's-.*= "--' -e 's-".*--')

build:
	go build -o build/draw -ldflags "-s -w" cmd/main.go

generate-all:
	bash -c scripts/generateAll.sh

test:
	go test ./... && echo ":)" || echo ":-/"
