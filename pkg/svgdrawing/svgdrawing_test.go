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
