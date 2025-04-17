package svgdrawing_test

import (
	"os"

	"github.com/stretchr/testify/require"

	"testing"

	"github.com/okieoth/draw.chart.things/pkg/boxesimpl"
	"github.com/okieoth/draw.chart.things/pkg/svgdrawing"
	"github.com/okieoth/draw.chart.things/pkg/types"
)

func TestSimpleSvg(t *testing.T) {
	tests := []struct {
		inputFile  string
		outputFile string
	}{
		{
			inputFile:  "../../resources/examples/simple_box.yaml",
			outputFile: "../../temp/TestSimpleSvg_box.svg",
		},
		{
			inputFile:  "../../resources/examples/simple_box_nested.yaml",
			outputFile: "../../temp/TestSimpleSvg_box_nested.svg",
		},
		{
			inputFile:  "../../resources/examples/simple_box_nested2.yaml",
			outputFile: "../../temp/TestSimpleSvg_box_nested2.svg",
		},
		{
			inputFile:  "../../resources/examples/simple_box_nested3.yaml",
			outputFile: "../../temp/TestSimpleSvg_box_nested3.svg",
		},
		{
			inputFile:  "../../resources/examples/simple_box_nested4.yaml",
			outputFile: "../../temp/TestSimpleSvg_box_nested4.svg",
		},
		{
			inputFile:  "../../resources/examples/simple_box_nested5.yaml",
			outputFile: "../../temp/TestSimpleSvg_box_nested5.svg",
		},
		{
			inputFile:  "../../resources/examples/simple_diamond.yaml",
			outputFile: "../../temp/TestSimpleSvg_diamond.svg",
		},
		{
			inputFile:  "../../resources/examples/horizontal_diamond.yaml",
			outputFile: "../../temp/TestSimpleSvg_hdiamond.svg",
		},
		{
			inputFile:  "../../resources/examples/complex_vertical.yaml",
			outputFile: "../../temp/TestSimpleSvg_vcomplex.svg",
		},
		{
			inputFile:  "../../resources/examples/complex_horizontal.yaml",
			outputFile: "../../temp/TestSimpleSvg_hcomplex.svg",
		},
		{
			inputFile:  "../../resources/examples/complex_complex.yaml",
			outputFile: "../../temp/TestSimpleSvg_ccomplex.svg",
		},
		{
			inputFile:  "../../resources/examples/horizontal_nested_diamond.yaml",
			outputFile: "../../temp/TestSimpleSvg_hdiamond_nestedx.svg",
		},
		{
			inputFile:  "../../resources/examples/horizontal_nested_diamond2.yaml",
			outputFile: "../../temp/TestSimpleSvg_hdiamond_nestedx2.svg",
		},
	}

	textDimensionCalulator := svgdrawing.NewSvgTextDimensionCalculator()

	for _, test := range tests {
		b, err := types.LoadInputFromFile[types.Boxes](test.inputFile)
		require.Nil(t, err)
		doc, err := boxesimpl.InitialLayoutBoxes(b, textDimensionCalulator)
		require.Nil(t, err)
		output, err := os.Create(test.outputFile)
		require.Nil(t, err)
		svgdrawing := svgdrawing.NewDrawing(output)
		doc.DrawBoxes(svgdrawing)
		svgdrawing.Done()
		output.Close()
	}
}

func TestSplitTxt(t *testing.T) {
	textDimensionCalculator := svgdrawing.NewSvgTextDimensionCalculator()
	tests := []struct {
		name           string
		inputText      string
		fontDef        *types.FontDef
		expectedWidth  int
		expectedHeight int
		expectedLines  []svgdrawing.TextAndDimensions
	}{
		{
			name:      "Single line sans-serif",
			inputText: "Hello World",
			fontDef: &types.FontDef{
				Font:              "sans-serif",
				Size:              12,
				MaxLenBeforeBreak: 100,
			},
			expectedWidth:  62,
			expectedHeight: 12,
			expectedLines: []svgdrawing.TextAndDimensions{
				{Text: "Hello World", Width: 62, Height: 12},
			},
		},
		{
			name:      "Multi line sans-serif",
			inputText: "Hello World Hello World",
			fontDef: &types.FontDef{
				Font:              "sans-serif",
				Size:              12,
				MaxLenBeforeBreak: 100,
			},
			expectedWidth:  96,
			expectedHeight: 28,
			expectedLines: []svgdrawing.TextAndDimensions{
				{Text: "Hello World Hello", Width: 96, Height: 14},
				{Text: "World", Width: 28, Height: 14},
			},
		},
		{
			name:      "Multi-line sans-serif",
			inputText: "This is a long text that should wrap into multiple lines",
			fontDef: &types.FontDef{
				Font:              "sans-serif",
				Size:              12,
				MaxLenBeforeBreak: 51,
			},
			expectedWidth:  51,
			expectedHeight: 98,
			expectedLines: []svgdrawing.TextAndDimensions{
				{Text: "This is a", Width: 51, Height: 14},
				{Text: "long text", Width: 51, Height: 14},
				{Text: "that", Width: 22, Height: 14},
				{Text: "should", Width: 34, Height: 14},
				{Text: "wrap into", Width: 51, Height: 14},
				{Text: "multiple", Width: 45, Height: 14},
				{Text: "lines", Width: 28, Height: 14},
			},
		},
		{
			name:      "Single line monospace",
			inputText: "Monospace",
			fontDef: &types.FontDef{
				Font:              "monospace",
				Size:              10,
				MaxLenBeforeBreak: 80,
			},
			expectedWidth:  59,
			expectedHeight: 10,
			expectedLines: []svgdrawing.TextAndDimensions{
				{Text: "Monospace", Width: 59, Height: 10},
			},
		},
		{
			name:      "Multi-line serif",
			inputText: "Serif font with multiple lines",
			fontDef: &types.FontDef{
				Font:              "serif",
				Size:              14,
				MaxLenBeforeBreak: 100,
			},
			expectedWidth:  94,
			expectedHeight: 32,
			expectedLines: []svgdrawing.TextAndDimensions{
				{Text: "Serif font with", Width: 94, Height: 16},
				{Text: "multiple lines", Width: 88, Height: 16},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			width, height, lines := textDimensionCalculator.SplitTxt(tt.inputText, tt.fontDef)
			require.Equal(t, tt.expectedWidth, width)
			require.Equal(t, tt.expectedHeight, height)
			for i, line := range lines {
				require.Equal(t, tt.expectedLines[i].Text, line.Text)
				require.Equal(t, tt.expectedLines[i].Width, line.Width)
				require.Equal(t, tt.expectedLines[i].Height, line.Height)
			}
		})
	}
}
