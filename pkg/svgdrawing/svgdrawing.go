package svgdrawing

import (
	"fmt"
	"io"
	"unicode/utf8"

	"github.com/okieoth/draw.chart.things/pkg/svg"
	"github.com/okieoth/draw.chart.things/pkg/types"
)

const SANSSERIF_NORMAL_WIDTH_FACTOR = 2.6 / 6
const SANSSERIF_BOLD_WIDTH_FACTOR = 2.8 / 6

const MONOSPACE_NORMAL_WIDTH_FACTOR = 3 / 5
const MONOSPACE_BOLD_WIDTH_FACTOR = 3 / 5

type SvgTextDimensionCalculator struct {
}

func NewSvgTextDimensionCalculator() *SvgTextDimensionCalculator {
	return &SvgTextDimensionCalculator{}
}

func getFormat(format *types.FontDef) *types.FontDef {
	if format != nil {
		return format
	}
	f := types.InitFontDef(nil)
	return &f
}

func serifDimensions(runeCount int, fontSize int, bold bool) (width, height int) {
	var factor float32
	if bold {
		factor = float32(fontSize*10) * 2.15 / 5
	} else {
		factor = float32(fontSize*10) * 2.05 / 5
	}
	w := int(factor) * runeCount / 10
	return w, fontSize
}

func sansserifDimensions(runeCount int, fontSize int, bold bool) (width, height int) {
	var factor float32
	if bold {
		factor = float32(fontSize*10) * 2.8 / 6
	} else {
		factor = float32(fontSize*10) * 2.6 / 6
	}
	w := int(factor) * runeCount / 10
	return w, fontSize
}
func monospaceDimensions(runeCount int, fontSize int, bold bool) (width, height int) {
	var factor float32
	if bold {
		factor = float32(fontSize*10) * 3 / 5
	} else {
		factor = float32(fontSize*10) * 3 / 5
	}
	w := int(factor) * runeCount / 10
	return w, fontSize
}

func (s *SvgTextDimensionCalculator) Dimensions(txt string, format *types.FontDef) (width, height int) {
	if format == nil {
		f := types.InitFontDef(nil)
		format = &f
	}
	fontSize := format.Size
	if fontSize == 0 {
		fontSize = 10
	}
	bold := false
	if format.Weight != nil && *format.Weight == types.FontDefWeightEnum_bold {
		bold = true
	}
	runeCount := utf8.RuneCountInString(txt)

	switch format.Font {
	case "serif":
		return serifDimensions(runeCount, fontSize, bold)
	case "sans-serif":
		return sansserifDimensions(runeCount, fontSize, bold)
	case "monospace":
		return monospaceDimensions(runeCount, fontSize, bold)
	default:
		return sansserifDimensions(runeCount, fontSize, bold)
	}
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
