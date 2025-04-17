.PHONY: test build

VERSION = $(shell grep "const Version =" cmd/sub/version.go | grep "const Version =" | sed -e 's-.*= "--' -e 's-".*--')
CURRENT_DIR = $(shell pwd)

build:
	go build -o build/draw -ldflags "-s -w" cmd/main.go

build-docker:
	docker build -f Dockerfile.release -t ghcr.io/okieoth/draw.chart.things:$(VERSION) .

generate-all:
	bash -c scripts/generateAll.sh

open-nested-in-browser:
	firefox file://$(CURRENT_DIR)/temp/TestSimpleSvg_box.svg \
		file://$(CURRENT_DIR)/temp/TestSimpleSvg_box_nested.svg \
		file://$(CURRENT_DIR)/temp/TestSimpleSvg_box_nested2.svg \
		file://$(CURRENT_DIR)/temp/TestSimpleSvg_box_nested3.svg \
		file://$(CURRENT_DIR)/temp/TestSimpleSvg_box_nested4.svg \
		file://$(CURRENT_DIR)/temp/TestSimpleSvg_box_nested5.svg &


test:
	go test ./... && echo ":)" || echo ":-/"
