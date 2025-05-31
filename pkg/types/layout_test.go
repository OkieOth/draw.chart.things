package types_test

import (
	"testing"

	"github.com/okieoth/draw.chart.things/pkg/boxesimpl"
	"github.com/okieoth/draw.chart.things/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type DummyDimensionCalculator struct {
	width  int
	height int
}

func (d *DummyDimensionCalculator) Dimensions(txt string, format *types.FontDef) (width, height int) {
	return d.width, d.height
}

func NewDummyDimensionCalculator(width, height int) *DummyDimensionCalculator {
	return &DummyDimensionCalculator{
		width:  width,
		height: height,
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
			expectedHeight: 72,
			expectedWidth:  120,
		},
		{
			// extends "test1" with an additional text1
			layout: types.Layout{
				Caption: "test2",
				Text1:   "test2-text1",
			},
			expectedHeight: 128,
			expectedWidth:  120,
		},
		{
			// extends "test2" with an additional text2
			layout: types.Layout{
				Caption: "test3",
				Text1:   "test3-text1",
				Text2:   "test3-text2",
			},
			expectedHeight: 192,
			expectedWidth:  120,
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
			expectedHeight: 433,
			expectedWidth:  120,
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
			expectedHeight: 269,
			expectedWidth:  390,
		},
	}

	dc := NewDummyDimensionCalculator(100, 50)
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
	require.Greater(t, le.CenterX, 0)
	require.Greater(t, le.CenterY, 0)
	require.Equal(t, le.CenterX, le.X+(le.Width/2))
	require.Equal(t, le.CenterY, le.Y+(le.Height/2))
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
		inputFile string
	}{
		{
			inputFile: "../../resources/examples/simple_diamond.yaml",
		},
		{
			inputFile: "../../resources/examples/horizontal_diamond.yaml",
		},
		{
			inputFile: "../../resources/examples/complex_vertical.yaml",
		},
		{
			inputFile: "../../resources/examples/complex_horizontal.yaml",
		},
	}
	dc := NewDummyDimensionCalculator(100, 50)
	for _, test := range tests {
		b, err := types.LoadInputFromFile[types.Boxes](test.inputFile)
		require.Nil(t, err)
		doc, err := boxesimpl.InitialLayoutBoxes(b, dc)
		require.Nil(t, err)
		checkLayoutElement(t, &doc.Boxes, 0, 0)
	}

}

func TestAreOnTheSameVerticalLevel(t *testing.T) {
	tests := []struct {
		inputFile string
		checkFunc func(t *testing.T, doc *types.BoxesDocument)
	}{
		{
			inputFile: "../../resources/examples/simple_diamond.yaml",
			checkFunc: func(t *testing.T, doc *types.BoxesDocument) {
				require.NotNil(t, doc)
				e1 := doc.Boxes.Vertical.Elems[0]
				e2 := doc.Boxes.Vertical.Elems[1].Horizontal.Elems[0]
				e3 := doc.Boxes.Vertical.Elems[1].Horizontal.Elems[1]
				e4 := doc.Boxes.Vertical.Elems[1].Horizontal.Elems[2]
				e5 := doc.Boxes.Vertical.Elems[2]
				require.False(t, e1.AreOnTheSameVerticalLevel(&e2))
				require.False(t, e1.AreOnTheSameVerticalLevel(&e3))
				require.False(t, e1.AreOnTheSameVerticalLevel(&e4))
				require.False(t, e1.AreOnTheSameVerticalLevel(&e5))

				require.True(t, e2.AreOnTheSameVerticalLevel(&e3))
				require.True(t, e2.AreOnTheSameVerticalLevel(&e4))
				require.True(t, e3.AreOnTheSameVerticalLevel(&e4))
				require.True(t, e4.AreOnTheSameVerticalLevel(&e2))
				require.True(t, e4.AreOnTheSameVerticalLevel(&e3))
				require.False(t, e5.AreOnTheSameVerticalLevel(&e2))
				require.False(t, e5.AreOnTheSameVerticalLevel(&e3))
				require.False(t, e5.AreOnTheSameVerticalLevel(&e4))
			},
		},
		{
			inputFile: "../../resources/examples/horizontal_diamond.yaml",
			checkFunc: func(t *testing.T, doc *types.BoxesDocument) {
				require.NotNil(t, doc)
				e1 := doc.Boxes.Horizontal.Elems[0]
				e2 := doc.Boxes.Horizontal.Elems[1].Vertical.Elems[0]
				e3 := doc.Boxes.Horizontal.Elems[1].Vertical.Elems[1]
				e4 := doc.Boxes.Horizontal.Elems[1].Vertical.Elems[2]
				e5 := doc.Boxes.Horizontal.Elems[2]
				require.False(t, e1.AreOnTheSameVerticalLevel(&e2))
				require.True(t, e1.AreOnTheSameVerticalLevel(&e3))
				require.True(t, e3.AreOnTheSameVerticalLevel(&e1))
				require.False(t, e1.AreOnTheSameVerticalLevel(&e4))
				require.False(t, e4.AreOnTheSameVerticalLevel(&e1))
				require.True(t, e3.AreOnTheSameVerticalLevel(&e5))
			},
		},
	}
	dc := NewDummyDimensionCalculator(100, 50)
	for _, test := range tests {
		b, err := types.LoadInputFromFile[types.Boxes](test.inputFile)
		require.Nil(t, err)
		doc, err := boxesimpl.InitialLayoutBoxes(b, dc)
		require.Nil(t, err)
		test.checkFunc(t, doc)
	}

}
