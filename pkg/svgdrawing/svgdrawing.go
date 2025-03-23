package svgdrawing

import (
	"fmt"
	"io"

	"github.com/okieoth/draw.chart.things/pkg/types"
	svg "github.com/okieoth/svgo"
)

type SvgTextDimensionCalculator struct {
}

func NewSvgTextDimensionCalculator() *SvgTextDimensionCalculator {
	return &SvgTextDimensionCalculator{}
}

func (s *SvgTextDimensionCalculator) CaptionDimensions(txt string) (width, height int) {
	return 100, 50 // TODO - implement this
}

func (s *SvgTextDimensionCalculator) Text1Dimensions(txt string) (width, height int) {
	return 100, 50 // TODO - implement this
}

func (s *SvgTextDimensionCalculator) Text2Dimensions(txt string) (width, height int) {
	return 100, 50 // TODO - implement this
}

type Drawing struct {
	canvas *svg.SVG
}

func NewDrawing(w io.Writer) *Drawing {
	return &Drawing{
		canvas: svg.New(w),
	}
}

func (d *Drawing) Start(title string, height, width int) error {
	args := fmt.Sprintf("viewBox=\"0,0,%d,%d\"", width, height)
	d.canvas.Start(width, height, args)
	return nil
}

func (d *Drawing) Draw(id, caption, text1, text2 string, x, y, width, height int, format types.BoxFormat) error {

	if format.Fill != nil || format.Border != nil {
		attr := ""
		if format.Fill != nil {
			attr = fmt.Sprintf("fill: %s", *(*format.Fill).Color)
		}
		if format.Border != nil {
			if attr != "" {
				attr += ";"
			}
			w := 1
			if (*format.Border).Width != nil {
				w = *(*format.Border).Width
			}
			c := "black"
			if (*format.Border).Color != nil {
				c = *(*format.Border).Color
			}
			attr += fmt.Sprintf("stroke: %s;stroke-width: %d", c, w)
		}
		d.canvas.RectWithId(id, x, y, width, height, attr)
	}
	return nil
}

func (d *Drawing) Done() error {
	d.canvas.End()
	return nil
}
