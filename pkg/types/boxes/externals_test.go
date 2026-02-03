package boxes_test

import (
	"testing"

	"github.com/okieoth/draw.chart.things/pkg/types"
	"github.com/okieoth/draw.chart.things/pkg/types/boxes"
	"github.com/stretchr/testify/require"
)

func TestLoadExternalFormats(t *testing.T) {
	inputFormat := "../../../resources/examples_boxes/ext_formats.yaml"
	inputLayout := "../../../resources/examples_boxes/ext_complex_horizontal_connected_pics.yaml"
	inputLayout2 := "../../../resources/examples_boxes/complex_horizontal_connected_pics.yaml"

	additionalFormats, err := types.LoadInputFromFile[boxes.BoxesFileMixings](inputFormat)
	require.Nil(t, err)
	require.NotNil(t, additionalFormats)

	require.Len(t, additionalFormats.Formats, 0)
	require.Len(t, additionalFormats.Images, 1)

	img, ok := additionalFormats.Images["smilie_01_43"]
	require.True(t, ok)

	b, err := types.LoadInputFromFile[boxes.Boxes](inputLayout)
	require.Nil(t, err)
	require.NotNil(t, b)

	require.NotNil(t, img.Base64)
	require.Nil(t, img.Base64Src)

	b2, err := types.LoadInputFromFile[boxes.Boxes](inputLayout2)
	require.Nil(t, err)
	require.NotNil(t, b2)

	require.NotEqual(t, len(b2.Images), len(b.Formats))

	b.MixinThings(*additionalFormats)

	require.Equal(t, len(b2.Images), len(b.Images))
}

func TestLoadExternalConnections(t *testing.T) {
	input := "../../../resources/examples_boxes/ext_complex_horizontal_connected_pics.yaml"
	inputConnections := "../../../resources/examples_boxes/ext_connections.yaml"

	b, err := types.LoadInputFromFile[boxes.Boxes](input)
	require.Nil(t, err)
	require.NotNil(t, b)

	c, err := types.LoadInputFromFile[boxes.BoxesFileMixings](inputConnections)
	require.Nil(t, err)
	require.NotNil(t, c)

	// r4_1
	require.Len(t, b.Boxes.Horizontal[0].Vertical[0].Connections, 2)
	// r5_2
	require.Len(t, b.Boxes.Horizontal[1].Vertical[1].Connections, 0)

	b.MixinThings(*c)

	// r4_1
	require.Len(t, b.Boxes.Horizontal[0].Vertical[0].Connections, 4)
	// r5_2
	require.Len(t, b.Boxes.Horizontal[1].Vertical[1].Connections, 1)

	cl, ok := b.Formats["connLines"]
	require.True(t, ok)
	require.NotNil(t, cl.Line)
	require.Equal(t, 2.0, *cl.Line.Width)
	require.Equal(t, "pink", *cl.Line.Color)
}

func TestLoadExternalConnections2(t *testing.T) {
	input := "../../../resources/examples_boxes/ext_complex_horizontal_connected_pics.yaml"
	inputConnections := "../../../resources/examples_boxes/ext_connections2.yaml"

	b, err := types.LoadInputFromFile[boxes.Boxes](input)
	require.Nil(t, err)
	require.NotNil(t, b)

	c, err := types.LoadInputFromFile[boxes.BoxesFileMixings](inputConnections)
	require.Nil(t, err)
	require.NotNil(t, c)

	// r4_1
	require.Len(t, b.Boxes.Horizontal[0].Vertical[0].Connections, 2)
	// r5_2
	require.Len(t, b.Boxes.Horizontal[1].Vertical[1].Connections, 0)

	require.Len(t, b.Boxes.Horizontal[1].Vertical[1].Vertical, 0)
	require.Len(t, b.Boxes.Horizontal[1].Vertical[1].Horizontal, 0)
	require.Len(t, b.Boxes.Horizontal[2].Vertical[0].Vertical, 0)
	require.Len(t, b.Boxes.Horizontal[2].Vertical[0].Horizontal, 0)

	b.MixinThings(*c)

	// r4_1
	require.Len(t, b.Boxes.Horizontal[0].Vertical[0].Connections, 4)
	// r5_2
	require.Len(t, b.Boxes.Horizontal[1].Vertical[1].Vertical, 2)
	require.Len(t, b.Boxes.Horizontal[1].Vertical[1].Horizontal, 0)
	require.Len(t, b.Boxes.Horizontal[2].Vertical[0].Vertical, 0)
	require.Len(t, b.Boxes.Horizontal[2].Vertical[0].Horizontal, 2)
}
