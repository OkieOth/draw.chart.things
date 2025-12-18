//go:build js && wasm
// +build js,wasm

package main

import (
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
func createSvg(boxesYaml string, defaultDepth int, filter []string) string {
	var boxes boxes.Boxes
	if err := y.Unmarshal([]byte(boxesYaml), &boxes); err != nil {
		return unknownSvg
	}

	ret := boxesimpl.DrawBoxesFiltered(boxes, defaultDepth, filter)
	if ret.ErrorMsg != "" {
		return unknownSvg
	}
	return ret.SVG
}

// JS wrapper: exposes getSvg to JavaScript
func createSvgWrapper(this js.Value, args []js.Value) interface{} {
	if len(args) < 3 {
		return "error: expected (string, number, string[])"
	}
	input := args[0].String()
	depth := args[1].Int()
	jsArray := args[2]
	if jsArray.Type() != js.TypeObject {
		return "error: filter must be an array"
	}
	length := jsArray.Length()
	filter := make([]string, 0, length)
	for i := 0; i < length; i++ {
		val := jsArray.Index(i)
		if val.Type() == js.TypeString {
			filter = append(filter, val.String())
		}
	}
	return createSvg(input, depth, filter)
}

func main() {
	// Expose the function to JS as `getSvg`
	js.Global().Set("createSvg", js.FuncOf(createSvgWrapper))

	// Keep Go running
	select {}
}
