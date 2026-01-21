//go:build js && wasm
// +build js,wasm

package main

import (
	"errors"
	"fmt"
	"syscall/js"

	"github.com/okieoth/draw.chart.things/pkg/boxesimpl"
	"github.com/okieoth/draw.chart.things/pkg/types/boxes"

	y "gopkg.in/yaml.v3"
)

const unknownSvg string = `<svg xmlns="http://www.w3.org/2000/svg" width="800" height="400" viewBox="0 0 800 400">
	<rect x="150" y="100" width="500" height="200" fill="#eeeeee" stroke="#555555" stroke-width="2" />
	<text x="400" y="200" font-family="sans-serif" font-size="20" fill="#333333" text-anchor="middle" dominant-baseline="middle">unknown content</text>
</svg>
`

// getSVG implements the requested signature
func createSvg(boxesYaml string, defaultDepth int, expanded, blacklisted []string, debug bool) string {
	var boxes boxes.Boxes
	if err := y.Unmarshal([]byte(boxesYaml), &boxes); err != nil {
		return unknownSvg
	}

	ret := boxesimpl.DrawBoxesFiltered(boxes, defaultDepth, expanded, blacklisted, debug)
	if ret.ErrorMsg != "" {
		return unknownSvg
	}
	return ret.SVG
}

func createSvgExt(boxesYaml string, additionalFormats, additionalConnections []string, defaultDepth int, expanded, blacklisted []string, debug bool) string {
	var b boxes.Boxes
	if err := y.Unmarshal([]byte(boxesYaml), &b); err != nil {
		fmt.Printf("error while unmarshalling boxes layout: %v", err)
		return unknownSvg
	}

	for i, c := range additionalConnections {
		var extConnections boxes.AdditionalConnections
		if err := y.Unmarshal([]byte(c), &extConnections); err != nil {
			fmt.Printf("error while unmarshalling external connections (%d): %v", i, err)
			return unknownSvg
		}
		b.MixinConnections(extConnections)
	}
	for i, c := range additionalFormats {
		var extFormats boxes.AdditionalFormats
		if err := y.Unmarshal([]byte(c), &extFormats); err != nil {
			fmt.Printf("error while unmarshalling external formats (%d): %v", i, err)
			return unknownSvg
		}
		b.MixinFormats(extFormats)
	}

	ret := boxesimpl.DrawBoxesFiltered(b, defaultDepth, expanded, blacklisted, debug)
	if ret.ErrorMsg != "" {
		return unknownSvg
	}
	return ret.SVG
}

func getArrayFromJsValue(args []js.Value, index int) ([]string, error) {
	jsArray := args[index]
	if jsArray.Type() != js.TypeObject {
		return []string{}, errors.New("")
	}
	length := jsArray.Length()
	ret := make([]string, 0, length)
	for i := range length {
		val := jsArray.Index(i)
		if val.Type() == js.TypeString {
			ret = append(ret, val.String())
		}
	}
	return ret, nil
}

// JS wrapper: exposes getSvg to JavaScript
func createSvgWrapper(this js.Value, args []js.Value) interface{} {
	if len(args) < 5 {
		return "error: expected (string, number, string[], string[], bool)"
	}
	input := args[0].String()
	depth := args[1].Int()
	debug := args[4].Bool()
	expanded, err := getArrayFromJsValue(args, 2)
	if err != nil {
		return "error: expanded must be an array"
	}
	blacklisted, err := getArrayFromJsValue(args, 3)
	if err != nil {
		return "error: blacklisted must be an array"
	}
	return createSvg(input, depth, expanded, blacklisted, debug)
}

func createSvgExtWrapper(this js.Value, args []js.Value) interface{} {
	if len(args) < 7 {
		return "error: expected (string, string[], string[], number, string[], string[], bool)"
	}
	input := args[0].String()
	additionalFormats, err := getArrayFromJsValue(args, 1)
	if err != nil {
		return "error: additionalFormats needs to be an array"
	}
	additionalConnections, err := getArrayFromJsValue(args, 2)
	if err != nil {
		return "error: additionalConnections needs to be an array"
	}
	expanded, err := getArrayFromJsValue(args, 4)
	if err != nil {
		return "error: expanded must be an array"
	}
	blacklisted, err := getArrayFromJsValue(args, 5)
	if err != nil {
		return "error: blacklisted must be an array"
	}
	depth := args[3].Int()
	debug := args[6].Bool()
	return createSvgExt(input, additionalFormats, additionalConnections, depth, expanded, blacklisted, debug)
}

func main() {
	// Expose the function to JS as `getSvg`
	js.Global().Set("createSvg", js.FuncOf(createSvgWrapper))
	js.Global().Set("createSvgExt", js.FuncOf(createSvgExtWrapper))

	// Keep Go running
	select {}
}
