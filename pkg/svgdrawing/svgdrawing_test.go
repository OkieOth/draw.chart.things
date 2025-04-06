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
		// {
		// 	name:      "Single line sans-serif",
		// 	inputText: "Hello World",
		// 	fontDef: &types.FontDef{
		// 		Font:              "sans-serif",
		// 		Size:              12,
		// 		MaxLenBeforeBreak: 100,
		// 	},
		// 	expectedWidth:  57,
		// 	expectedHeight: 14,
		// 	expectedLines: []svgdrawing.TextAndDimensions{
		// 		{Text: "Hello World", Width: 57, Height: 14},
		// 	},
		// },
		{
			name:      "Multi line sans-serif",
			inputText: "Hello World Hello World",
			fontDef: &types.FontDef{
				Font:              "sans-serif",
				Size:              12,
				MaxLenBeforeBreak: 100,
			},
			expectedWidth:  88,
			expectedHeight: 28,
			expectedLines: []svgdrawing.TextAndDimensions{
				{Text: "Hello World Hello", Width: 88, Height: 14},
				{Text: "World", Width: 26, Height: 14},
			},
		},
		// {
		// 	name:      "Multi-line sans-serif",
		// 	inputText: "This is a long text that should wrap into multiple lines",
		// 	fontDef: &types.FontDef{
		// 		Font:              "sans-serif",
		// 		Size:              12,
		// 		MaxLenBeforeBreak: 50,
		// 	},
		// 	expectedWidth:  50,
		// 	expectedHeight: 42,
		// 	expectedLines: []svgdrawing.TextAndDimensions{
		// 		{Text: "This is a long", Width: 50, Height: 14},
		// 		{Text: "text that should", Width: 50, Height: 14},
		// 		{Text: "wrap into multiple", Width: 50, Height: 14},
		// 	},
		// },
		// {
		// 	name:      "Single line monospace",
		// 	inputText: "Monospace",
		// 	fontDef: &types.FontDef{
		// 		Font:              "monospace",
		// 		Size:              10,
		// 		MaxLenBeforeBreak: 80,
		// 	},
		// 	expectedWidth:  60,
		// 	expectedHeight: 12,
		// 	expectedLines: []svgdrawing.TextAndDimensions{
		// 		{Text: "Monospace", Width: 60, Height: 12},
		// 	},
		// },
		// {
		// 	name:      "Multi-line serif",
		// 	inputText: "Serif font with multiple lines",
		// 	fontDef: &types.FontDef{
		// 		Font:              "serif",
		// 		Size:              14,
		// 		MaxLenBeforeBreak: 70,
		// 	},
		// 	expectedWidth:  70,
		// 	expectedHeight: 42,
		// 	expectedLines: []svgdrawing.TextAndDimensions{
		// 		{Text: "Serif font with", Width: 70, Height: 14},
		// 		{Text: "multiple lines", Width: 70, Height: 14},
		// 	},
		// },
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
