package types_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/okieoth/draw.chart.things/pkg/types"
)

type testFunc func(t *testing.T, b *types.Boxes)

func TestLoadBoxes(t *testing.T) {
	tests := []struct {
		fileName string
		verify   testFunc
	}{{
		fileName: "../../resources/examples/simple_box.yaml",
		verify: func(t *testing.T, b *types.Boxes) {
			assert.NotNil(t, b)
			assert.NotNil(t, b.Layout)
			assert.Equal(t, "main", *b.Layout.Id)
			assert.Equal(t, "I am a simple box", *b.Layout.Caption)
			assert.Equal(t, "The first example layout", *b.Layout.Text1)
			assert.Equal(t, "This is a simple box layout. It is used to demonstrate the basic layout features.", *b.Layout.Text2)
			assert.Equal(t, 2, len(b.Layout.Tags))
			assert.Equal(t, "simple", b.Layout.Tags[0])
			assert.Equal(t, "test", b.Layout.Tags[1])

			assert.NotNil(t, b.DefaultFormat)
			assert.NotNil(t, b.DefaultFormat.Border)
			assert.Equal(t, "black", *b.DefaultFormat.Border.Color)
			assert.Equal(t, int32(1), *b.DefaultFormat.Border.Width)
			assert.NotNil(t, b.DefaultFormat.Fill)
			assert.Equal(t, "lightgreen", *b.DefaultFormat.Fill.Color)
			assert.Nil(t, b.DefaultFormat.FontCaption)
			assert.Nil(t, b.DefaultFormat.FontText1)
			assert.Nil(t, b.DefaultFormat.FontText2)

			assert.Equal(t, 0, len(b.Formats))
		},
	}}

	for _, test := range tests {
		b, err := types.LoadInputFromFile[types.Boxes](test.fileName)
		assert.Nil(t, err)
		test.verify(t, b)
	}
}
