package types_test

import (
	"testing"

	"github.com/okieoth/draw.chart.things/pkg/boxesimpl"
	"github.com/okieoth/draw.chart.things/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type DummyDimensionCalculator struct {
	captionWidth  int
	captionHeight int
	text1Width    int
	text1Height   int
	text2Width    int
	text2Height   int
}

func (d *DummyDimensionCalculator) CaptionDimensions(txt string, format *types.FontDef) (width, height int) {
	return d.captionWidth, d.captionHeight
}

func (d *DummyDimensionCalculator) Text1Dimensions(txt string, format *types.FontDef) (width, height int) {
	return d.text1Width, d.text1Height
}

func (d *DummyDimensionCalculator) Text2Dimensions(txt string, format *types.FontDef) (width, height int) {
	return d.text2Width, d.text2Height
}

func NewDummyDimensionCalculator(captionWidth, captionHeight, text1Width, text1Height, text2Width, text2Height int) *DummyDimensionCalculator {
	return &DummyDimensionCalculator{
		captionWidth:  captionWidth,
		captionHeight: captionHeight,
		text1Width:    text1Width,
		text1Height:   text1Height,
		text2Width:    text2Width,
		text2Height:   text2Height,
	}
}

// func tested types.InitDimensions
func TestInitDimensions(t *testing.T) {
	tests := []struct {
		layout         types.Layout
		expectedHeight int
		expectedWidth  int
	}{
		{
			layout: types.Layout{
				Caption: "test1",
			},
			expectedHeight: 60,
			expectedWidth:  110,
		},
		{
			// extends "test1" with an additional text1
			layout: types.Layout{
				Caption: "test2",
				Text1:   "test2-text1",
			},
			expectedHeight: 75,
			expectedWidth:  130,
		},
		{
			// extends "test2" with an additional text2
			layout: types.Layout{
				Caption: "test3",
				Text1:   "test3-text1",
				Text2:   "test3-text2",
			},
			expectedHeight: 90,
			expectedWidth:  130,
		},
		{
			// basic vertical layout test
			layout: types.Layout{
				Caption: "test4",
				Text1:   "test4-text1",
				Text2:   "test4-text2",
				Vertical: []types.Layout{
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
			expectedHeight: 295,
			expectedWidth:  130,
		},
		{
			// basic horizontal layout test
			layout: types.Layout{
				Caption: "test5",
				Text1:   "test5-text1",
				Text2:   "test5-text2",
				Horizontal: []types.Layout{
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
			expectedHeight: 150,
			expectedWidth:  350,
		},
	}

	dc := NewDummyDimensionCalculator(100, 50, 120, 10, 80, 10)
	emptyFormats := map[string]types.BoxFormat{}
	for _, test := range tests {
		le := types.ExpInitLayoutElement(&test.layout, emptyFormats)
		le.InitDimensions(dc, 5, 10)
		assert.Equal(t, test.expectedHeight, le.Height)
		assert.Equal(t, test.expectedWidth, le.Width)
	}
}

func checkLayoutElement(t *testing.T, le *types.LayoutElement, initX, initY int) {
	require.GreaterOrEqual(t, le.X, initX)
	require.GreaterOrEqual(t, le.Y, initY)
	if le.Vertical != nil {
		require.GreaterOrEqual(t, le.Vertical.X, le.X)
		require.GreaterOrEqual(t, le.Vertical.Y, le.Y)
		for _, v := range le.Vertical.Elems {
			checkLayoutElement(t, &v, le.X, le.Y)
		}
	}
	if le.Horizontal != nil {
		require.GreaterOrEqual(t, le.Horizontal.X, le.X)
		require.GreaterOrEqual(t, le.Horizontal.Y, le.Y)
		for _, h := range le.Horizontal.Elems {
			checkLayoutElement(t, &h, le.X, le.Y)
		}
	}
}

func TestCenteredCoordinates(t *testing.T) {
	tests := []struct {
		inputFile  string
		outputFile string
	}{
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
	}
	dc := NewDummyDimensionCalculator(100, 50, 120, 10, 80, 10)
	for _, test := range tests {
		b, err := types.LoadInputFromFile[types.Boxes](test.inputFile)
		require.Nil(t, err)
		doc, err := boxesimpl.InitialLayoutBoxes(b, dc)
		require.Nil(t, err)
		checkLayoutElement(t, &doc.Boxes, 0, 0)
	}

}
