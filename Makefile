.PHONY: test build

VERSION = $(shell grep "const Version =" cmd/sub/version.go | grep "const Version =" | sed -e 's-.*= "--' -e 's-".*--')
CURRENT_DIR = $(shell pwd)

build:
	go build -o build/draw -ldflags "-s -w" cmd/main.go

build-docker:
	docker build -f Dockerfile.release -t ghcr.io/okieoth/draw.chart.things:$(VERSION) .

generate-all:
	bash -c scripts/generateAll.sh

open-simple-in-browser:
	firefox file://$(CURRENT_DIR)/temp/TestSimpleSvg_box.svg \
		file://$(CURRENT_DIR)/temp/TestSimpleSvg_box_nested.svg \
		file://$(CURRENT_DIR)/temp/TestSimpleSvg_box_nested2.svg \
		file://$(CURRENT_DIR)/temp/TestSimpleSvg_box_nested3.svg \
		file://$(CURRENT_DIR)/temp/TestSimpleSvg_box_nested4.svg \
		file://$(CURRENT_DIR)/temp/TestSimpleSvg_box_nested5.svg \
		file://$(CURRENT_DIR)/temp/TestSimpleSvg_hdiamond_nestedx.svg \
		file://$(CURRENT_DIR)/temp/TestSimpleSvg_hdiamond_nestedx2.svg &

open-complex-in-browser:
	firefox file://$(CURRENT_DIR)/temp/TestSimpleSvg_diamond.svg \
		file://$(CURRENT_DIR)/temp/TestSimpleSvg_hdiamond.svg \
		file://$(CURRENT_DIR)/temp/TestSimpleSvg_hdiamond_nestedx.svg \
		file://$(CURRENT_DIR)/temp/TestSimpleSvg_hdiamond_nestedx2.svg \
		file://$(CURRENT_DIR)/temp/TestSimpleSvg_ccomplex.svg \
		file://$(CURRENT_DIR)/temp/TestSimpleSvg_hcomplex.svg \
		file://$(CURRENT_DIR)/temp/TestSimpleSvg_vcomplex.svg &

open-connected-in-browser:
	firefox file://$(CURRENT_DIR)/temp/long_horizontal_01.svg \
		file://$(CURRENT_DIR)/temp/long_horizontal_02.svg \
		file://$(CURRENT_DIR)/temp/long_vertical_01.svg \
		file://$(CURRENT_DIR)/temp/long_vertical_02.svg \
		file://$(CURRENT_DIR)/temp/TestSimpleSvg_hcomplex_connected_01.svg \
		file://$(CURRENT_DIR)/temp/TestSimpleSvg_hcomplex_connected_02.svg \
		file://$(CURRENT_DIR)/temp/TestSimpleSvg_hcomplex_connected_03.svg \
		file://$(CURRENT_DIR)/temp/TestSimpleSvg_hcomplex_connected_04.svg &


test:
	go test ./... && echo ":)" || echo ":-/"
