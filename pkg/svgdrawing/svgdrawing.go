package svgdrawing

import (
	"github.com/ajstarks/svgo"
	"github.com/okieoth/draw.chart.things/pkg/types"
	"io"
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
	d.canvas.Start(width, height)
	return nil
}

func (d *Drawing) Draw(caption, text1, text2 string, x, y, width, height int, format types.BoxFormat) error {
	return nil
}

func (d *Drawing) Done() error {
	d.canvas.End()
	return nil
}
