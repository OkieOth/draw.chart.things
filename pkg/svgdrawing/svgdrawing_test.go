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
		svgdrawing.DrawRaster(doc.Width, doc.Height, types.RasterSize)
		svgdrawing.Done()
		output.Close()
		_, err = os.Stat(test.outputFile)
		require.Nil(t, err)
	}
}

func TestSvgWithConnections(t *testing.T) {
	type testData struct {
		inputFile  string
		outputFile string
		checkFunc  func(t *testing.T, doc *types.BoxesDocument)
	}
	runTests := func(tests []testData) {
		textDimensionCalulator := svgdrawing.NewSvgTextDimensionCalculator()

		for _, test := range tests {
			b, err := types.LoadInputFromFile[types.Boxes](test.inputFile)
			require.Nil(t, err)
			doc, err := boxesimpl.InitialLayoutBoxes(b, textDimensionCalulator)
			require.Nil(t, err)
			// debug - can help to see all possible connections in the created file
			//doc.ConnectBoxesFull()
			doc.ConnectBoxes()
			output, err := os.Create(test.outputFile)
			require.Nil(t, err)
			svgdrawing := svgdrawing.NewDrawing(output)
			doc.DrawBoxes(svgdrawing)
			doc.DrawConnections(svgdrawing)
			svgdrawing.DrawRaster(doc.Width, doc.Height, types.RasterSize)
			svgdrawing.Done()
			output.Close()
			_, err = os.Stat(test.outputFile)
			require.Nil(t, err)
			// test.checkFunc(t, doc)
		}
	}

	tests := []testData{
		{
			inputFile:  "../../resources/examples/complex_horizontal_connected_01.yaml",
			outputFile: "../../temp/TestSimpleSvg_hcomplex_connected_01.svg",
			checkFunc: func(t *testing.T, doc *types.BoxesDocument) {
				require.Equal(t, 1, len(doc.Connections))
				require.Equal(t, 4, len(doc.Connections[0].Parts))
			},
		},
		{
			inputFile:  "../../resources/examples/complex_horizontal_connected_02.yaml",
			outputFile: "../../temp/TestSimpleSvg_hcomplex_connected_02.svg",
			checkFunc: func(t *testing.T, doc *types.BoxesDocument) {
				require.Equal(t, 1, len(doc.Connections))
				require.Equal(t, 2, len(doc.Connections[0].Parts))
			},
		},
		{
			inputFile:  "../../resources/examples/complex_horizontal_connected_03.yaml",
			outputFile: "../../temp/TestSimpleSvg_hcomplex_connected_03.svg",
			checkFunc: func(t *testing.T, doc *types.BoxesDocument) {
				require.Equal(t, 1, len(doc.Connections))
				require.Equal(t, 2, len(doc.Connections[0].Parts))
			},
		},
		{
			inputFile:  "../../resources/examples/complex_horizontal_connected_04.yaml",
			outputFile: "../../temp/TestSimpleSvg_hcomplex_connected_04.svg",
			checkFunc: func(t *testing.T, doc *types.BoxesDocument) {
				require.Equal(t, 5, len(doc.Connections))
				require.Equal(t, 2, len(doc.Connections[0].Parts))
				require.Equal(t, 2, len(doc.Connections[1].Parts))
				require.Equal(t, 4, len(doc.Connections[2].Parts))
				require.Equal(t, 2, len(doc.Connections[3].Parts))
				require.Equal(t, 2, len(doc.Connections[4].Parts))
				// TODO
			},
		},
		{
			inputFile:  "../../resources/examples/complex_horizontal_connected_05.yaml",
			outputFile: "../../temp/TestSimpleSvg_hcomplex_connected_05.svg",
			checkFunc: func(t *testing.T, doc *types.BoxesDocument) {
				require.Equal(t, 5, len(doc.Connections))
				require.Equal(t, 2, len(doc.Connections[0].Parts))
				require.Equal(t, 2, len(doc.Connections[1].Parts))
				require.Equal(t, 2, len(doc.Connections[2].Parts))
				require.Equal(t, 2, len(doc.Connections[3].Parts))
				require.Equal(t, 4, len(doc.Connections[4].Parts))
				// TODO
			},
		},
		{
			inputFile:  "../../resources/examples/long_horizontal_01.yaml",
			outputFile: "../../temp/long_horizontal_01.svg",
			checkFunc: func(t *testing.T, doc *types.BoxesDocument) {
				require.Equal(t, 6, len(doc.Connections))
				for _, c := range doc.Connections {
					require.Equal(t, 2, len(c.Parts))
				}
			},
		},
		{
			inputFile:  "../../resources/examples/long_horizontal_02.yaml",
			outputFile: "../../temp/long_horizontal_02.svg",
			checkFunc: func(t *testing.T, doc *types.BoxesDocument) {
				require.Equal(t, 6, len(doc.Connections))
				for _, c := range doc.Connections {
					require.Equal(t, 2, len(c.Parts))
				}
			},
		},
		{
			inputFile:  "../../resources/examples/long_vertical_01.yaml",
			outputFile: "../../temp/long_vertical_01.svg",
			checkFunc: func(t *testing.T, doc *types.BoxesDocument) {
				require.Equal(t, 6, len(doc.Connections))
				for i, c := range doc.Connections {
					if i == 2 {
						require.Equal(t, 3, len(c.Parts))
					} else {
						require.Equal(t, 2, len(c.Parts))
					}
				}
			},
		},
		{
			inputFile:  "../../resources/examples/long_vertical_02.yaml",
			outputFile: "../../temp/long_vertical_02.svg",
			checkFunc: func(t *testing.T, doc *types.BoxesDocument) {
				require.Equal(t, 6, len(doc.Connections))
				for i, c := range doc.Connections {
					if i == 2 {
						require.Equal(t, 3, len(c.Parts))
					} else {
						require.Equal(t, 2, len(c.Parts))
					}
				}
			},
		},
		{
			inputFile:  "../../resources/examples/horizontal_nested_diamond2_connected.yaml",
			outputFile: "../../temp/horizontal_nested_diamond2_connected.svg",
			checkFunc: func(t *testing.T, doc *types.BoxesDocument) {
				require.Equal(t, 8, len(doc.Connections))
				// for _, c := range doc.Connections {
				// 	require.Equal(t, 4, len(c.Parts))
				// }
			},
		},
	}
	runTests(tests)
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
