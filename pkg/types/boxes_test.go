package types_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/okieoth/draw.chart.things/pkg/types"
)

type testFunc func(t *testing.T, b *types.Boxes)

func checkLayout(t *testing.T, l *types.Layout, id string, horizontalLen, verticalLen int) {
	assert.Len(t, l.Horizontal, horizontalLen)
	assert.Len(t, l.Vertical, verticalLen)
	assert.Equal(t, id, l.Id)
}

func TestLoadBoxes(t *testing.T) {
	tests := []struct {
		fileName string
		verify   testFunc
	}{{
		fileName: "../../resources/examples_boxes/simple_box.yaml",
		verify: func(t *testing.T, b *types.Boxes) {
			assert.NotNil(t, b)
			checkLayout(t, &b.Boxes, "main", 0, 0)

			assert.Equal(t, "I am a simple box", b.Boxes.Caption)
			assert.Equal(t, "The first example layout", b.Boxes.Text1)
			assert.Equal(t, "This is a simple box layout. It is used to demonstrate the basic layout features.", b.Boxes.Text2)
			assert.Len(t, b.Boxes.Tags, 2)
			assert.Equal(t, "simple", b.Boxes.Tags[0])
			assert.Equal(t, "test", b.Boxes.Tags[1])

			defaultFormat, defFormatExist := b.Formats["default"]

			assert.True(t, defFormatExist)
			assert.NotNil(t, defaultFormat.Border)
			assert.Equal(t, "black", *defaultFormat.Border.Color)
			assert.Equal(t, 1, *defaultFormat.Border.Width)
			assert.NotNil(t, defaultFormat.Fill)
			assert.Equal(t, "lightgreen", *defaultFormat.Fill.Color)
			assert.Nil(t, defaultFormat.FontCaption)
			assert.Nil(t, defaultFormat.FontText1)
			assert.Nil(t, defaultFormat.FontText2)

			assert.Len(t, b.Formats, 1)
		},
	},
		{
			fileName: "../../resources/examples_boxes/simple_diamond.yaml",
			verify: func(t *testing.T, b *types.Boxes) {
				assert.NotNil(t, b)
				assert.Len(t, b.Boxes.Horizontal, 0)
				assert.Len(t, b.Boxes.Vertical, 3)
				checkLayout(t, &b.Boxes.Vertical[0], "r1_1", 0, 0)

				assert.Len(t, b.Boxes.Vertical[1].Horizontal, 3)
				checkLayout(t, &b.Boxes.Vertical[1].Horizontal[0], "r2_1", 0, 0)
				checkLayout(t, &b.Boxes.Vertical[1].Horizontal[1], "r2_2", 0, 0)
				checkLayout(t, &b.Boxes.Vertical[1].Horizontal[2], "r2_3", 0, 0)

				checkLayout(t, &b.Boxes.Vertical[2], "r3_1", 0, 0)
			},
		}}

	for _, test := range tests {
		b, err := types.LoadInputFromFile[types.Boxes](test.fileName)
		assert.Nil(t, err)
		test.verify(t, b)
	}
}
