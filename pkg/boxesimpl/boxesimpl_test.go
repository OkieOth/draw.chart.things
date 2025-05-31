package boxesimpl_test

import (
	"github.com/stretchr/testify/require"

	"testing"

	"github.com/okieoth/draw.chart.things/pkg/boxesimpl"
	"github.com/okieoth/draw.chart.things/pkg/svgdrawing"
	"github.com/okieoth/draw.chart.things/pkg/types"
)

func TestDrawBoxesFromFile(t *testing.T) {
	tests := []struct {
		inputFile string
	}{
		// {
		// 	inputFile:  "../../resources/examples_boxes/simple_box.yaml",
		// },
		{
			inputFile: "../../resources/examples_boxes/simple_diamond.yaml",
		}}

	textDimensionCalulator := svgdrawing.NewSvgTextDimensionCalculator()

	for _, test := range tests {
		b, err := types.LoadInputFromFile[types.Boxes](test.inputFile)
		require.Nil(t, err)
		doc, err := boxesimpl.InitialLayoutBoxes(b, textDimensionCalulator)
		require.Nil(t, err)
		require.NotNil(t, doc)
	}

}
