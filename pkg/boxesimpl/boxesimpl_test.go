package boxesimpl_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"testing"

	"github.com/okieoth/draw.chart.things/pkg/boxesimpl"
	"github.com/okieoth/draw.chart.things/pkg/svgdrawing"
	"github.com/okieoth/draw.chart.things/pkg/types"
	"github.com/okieoth/draw.chart.things/pkg/types/boxes"
)

type DummyDimensionCalculator struct {
	width  int
	height int
}

func (d *DummyDimensionCalculator) Dimensions(txt string, format *types.FontDef) (width, height int) {
	return d.width, d.height
}

func (d *DummyDimensionCalculator) SplitTxt(txt string, format *types.FontDef) (width, height int, lines []types.TextAndDimensions) {
	return 0, 0, make([]types.TextAndDimensions, 0)
}

func NewDummyDimensionCalculator(width, height int) *DummyDimensionCalculator {
	return &DummyDimensionCalculator{
		width:  width,
		height: height,
	}
}

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
		b, err := types.LoadInputFromFile[boxes.Boxes](test.inputFile)
		require.Nil(t, err)
		doc, err := boxesimpl.InitialLayoutBoxes(b, textDimensionCalulator)
		require.Nil(t, err)
		require.NotNil(t, doc)
	}

}

// func tested types.InitDimensions
func TestInitDimensions(t *testing.T) {
	tests := []struct {
		layout         boxes.Layout
		expectedHeight int
		expectedWidth  int
	}{
		{
			layout: boxes.Layout{
				Caption: "test1",
			},
			expectedHeight: 72,
			expectedWidth:  120,
		},
		{
			// extends "test1" with an additional text1
			layout: boxes.Layout{
				Caption: "test2",
				Text1:   "test2-text1",
			},
			expectedHeight: 128,
			expectedWidth:  120,
		},
		{
			// extends "test2" with an additional text2
			layout: boxes.Layout{
				Caption: "test3",
				Text1:   "test3-text1",
				Text2:   "test3-text2",
			},
			expectedHeight: 192,
			expectedWidth:  120,
		},
		{
			// basic vertical layout test
			layout: boxes.Layout{
				Caption: "test4",
				Text1:   "test4-text1",
				Text2:   "test4-text2",
				Vertical: []boxes.Layout{
					{
						Caption: "test4-V1",
					},
					{
						Caption: "test4-V2",
					},
					{
						Caption: "test3-V3",
					},
				},
			},
			expectedHeight: 433,
			expectedWidth:  120,
		},
		{
			// basic horizontal layout test
			layout: boxes.Layout{
				Caption: "test5",
				Text1:   "test5-text1",
				Text2:   "test5-text2",
				Horizontal: []boxes.Layout{
					{
						Caption: "test5-V1",
					},
					{
						Caption: "test5-V2",
					},
					{
						Caption: "test5-V3",
					},
				},
			},
			expectedHeight: 269,
			expectedWidth:  390,
		},
	}

	dc := NewDummyDimensionCalculator(100, 50)
	emptyFormats := map[string]boxes.BoxFormat{}
	for _, test := range tests {
		le := boxesimpl.ExpInitLayoutElement(&test.layout, emptyFormats)
		le.InitDimensions(dc, 5, 10)
		assert.Equal(t, test.expectedHeight, le.Height)
		assert.Equal(t, test.expectedWidth, le.Width)
	}
}
