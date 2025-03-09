package types_test

import (
	"testing"

	"github.com/okieoth/draw.chart.things/pkg/types"
	"github.com/stretchr/testify/assert"
)

type DummyDimensionCalculator struct {
	captionWidth  int32
	captionHeight int32
	text1Width    int32
	text1Height   int32
	text2Width    int32
	text2Height   int32
}

func (d *DummyDimensionCalculator) CaptionDimensions(txt string) (width, height int32) {
	return d.captionWidth, d.captionHeight
}

func (d *DummyDimensionCalculator) Text1Dimensions(txt string) (width, height int32) {
	return d.text1Width, d.text1Height
}

func (d *DummyDimensionCalculator) Text2Dimensions(txt string) (width, height int32) {
	return d.text2Width, d.text2Height
}

func NewDummyDimensionCalculator(captionWidth, captionHeight, text1Width, text1Height, text2Width, text2Height int32) *DummyDimensionCalculator {
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
		expectedHeight int32
		expectedWidth  int32
	}{
		{
			layout: types.Layout{
				Caption: "test1",
			},
			expectedHeight: 60,
			expectedWidth:  110,
		},
		{
			layout: types.Layout{
				Caption: "test2",
				Text1:   "test2-text1",
			},
			expectedHeight: 75,
			expectedWidth:  130,
		},
		{
			layout: types.Layout{
				Caption: "test3",
				Text1:   "test3-text1",
				Text2:   "test3-text2",
			},
			expectedHeight: 90,
			expectedWidth:  130,
		},
		{
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
			expectedHeight: 155,
			expectedWidth:  350,
		},
	}

	dc := NewDummyDimensionCalculator(100, 50, 120, 10, 80, 10)
	emptyFormats := map[string]types.BoxFormat{}
	for _, test := range tests {
		le := types.ExpInitLayoutElement(&test.layout, emptyFormats)
		le.InitDimensions(dc)
		assert.Equal(t, test.expectedHeight, le.Height)
		assert.Equal(t, test.expectedWidth, le.Width)
	}
}
