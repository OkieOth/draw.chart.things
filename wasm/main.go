//go:build js && wasm
// +build js,wasm

package main

import (
	"errors"
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

func getArrayFromJsValue(args []js.Value, index int) ([]string, error) {
	jsArray := args[index]
	if jsArray.Type() != js.TypeObject {
		return []string{}, errors.New("")
	}
	length := jsArray.Length()
	ret := make([]string, 0, length)
	for i := 0; i < length; i++ {
		val := jsArray.Index(i)
		if val.Type() == js.TypeString {
			ret = append(ret, val.String())
		}
	}
	return ret, nil
}

// JS wrapper: exposes getSvg to JavaScript
func createSvgWrapper(this js.Value, args []js.Value) interface{} {
	if len(args) < 4 {
		return "error: expected (string, number, string[], bool)"
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

func main() {
	// Expose the function to JS as `getSvg`
	js.Global().Set("createSvg", js.FuncOf(createSvgWrapper))

	// Keep Go running
	select {}
}
