.PHONY: test build

VERSION = $(shell grep "const Version =" cmd/sub/version.go | grep "const Version =" | sed -e 's-.*= "--' -e 's-".*--')

build:
	go build -o build/draw -ldflags "-s -w" cmd/main.go

build-docker:
	docker build -f Dockerfile.release -t ghcr.io/okieoth/draw.chart.things:$(VERSION) .

generate-all:
	bash -c scripts/generateAll.sh

test:
	go test ./... && echo ":)" || echo ":-/"
