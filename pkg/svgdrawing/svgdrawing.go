package svgdrawing

import (
	"fmt"
	"io"
	"strings"
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

type textAndDimensions struct {
	text   string
	width  int
	height int
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

type calcDimensions func(runeCount int, fontSize int, bold bool) (width, height int)

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

func (s *SvgTextDimensionCalculator) splitTxt(txt string, format *types.FontDef) (width, height int, lines []textAndDimensions) {
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
	lineHeight := float32(1.2)
	if format.LineHeight != 0 {
		lineHeight = format.LineHeight
	}
	runeCount := utf8.RuneCountInString(txt)

	var w, h int
	switch format.Font {
	case "serif":
		if w, h = serifDimensions(runeCount, fontSize, bold); w <= format.MaxLenBeforeBreak {
			w, h, lines = splitTxtDimensions(w, runeCount, fontSize, bold, lineHeight, txt, format.MaxLenBeforeBreak, serifDimensions)
		}
	case "sans-serif":
		if w, h = sansserifDimensions(runeCount, fontSize, bold); w <= format.MaxLenBeforeBreak {
			w, h, lines = splitTxtDimensions(w, runeCount, fontSize, bold, lineHeight, txt, format.MaxLenBeforeBreak, sansserifDimensions)
		}
	case "monospace":
		if w, h = monospaceDimensions(runeCount, fontSize, bold); w <= format.MaxLenBeforeBreak {
			w, h, lines = splitTxtDimensions(w, runeCount, fontSize, bold, lineHeight, txt, format.MaxLenBeforeBreak, monospaceDimensions)
		}
	default:
		if w, h = sansserifDimensions(runeCount, fontSize, bold); w <= format.MaxLenBeforeBreak {
			w, h, lines = splitTxtDimensions(w, runeCount, fontSize, bold, lineHeight, txt, format.MaxLenBeforeBreak, sansserifDimensions)
		}
	}
	return w, h, lines
}

func (s *SvgTextDimensionCalculator) Dimensions(txt string, format *types.FontDef) (width, height int) {
	w, h, _ := s.splitTxt(txt, format)
	return w, h
}

func splitTxtDimensions(
	originalWidth, runeCount, fontSize int,
	bold bool, lineHeight float32,
	txt string,
	maxLenBeforeBreak int,
	f calcDimensions) (width, height int, lines []textAndDimensions) {
	words := strings.Fields(txt)
	wordsWithWidth := make([]textAndDimensions, 0)
	for _, w := range words {
		rc := utf8.RuneCount([]byte(w))
		width, _ := f(rc, fontSize, bold)
		if width > maxLenBeforeBreak {
			// the included word is truncated :D
			for j := rc - 1; j >= 0; j-- {
				w2, _ := f(j, fontSize, bold)
				if w2 <= maxLenBeforeBreak {
					w = string([]rune(w)[:j])
					width = w2
					w = w[:rc+1]
					break
				}
			}
		}
		wordsWithWidth = append(wordsWithWidth, textAndDimensions{text: w, width: width})
	}
	// build the output lines ...
	curWidth := 0
	var line string
	var maxWidth int
	lines = make([]textAndDimensions, 0)
	lh := int(float32(fontSize) * lineHeight)
	heightSum := 0
	appendLineToOutput := func(txt string, width, height int) []textAndDimensions {
		lines = append(lines, textAndDimensions{text: txt, width: width, height: height})
		heightSum += lh
		lw, _ := f(utf8.RuneCount([]byte(txt)), fontSize, bold)
		if lw > maxWidth {
			maxWidth = lw
		}
		return lines
	}
	for _, ww := range wordsWithWidth {
		if curWidth+ww.width < maxLenBeforeBreak {
			if line != "" {
				line += " "
			}
			line += ww.text
			curWidth += ww.width
		} else {
			lines = appendLineToOutput(line, curWidth, lh)
			lines = append(lines, textAndDimensions{text: line, width: curWidth, height: lh})
			heightSum += lh
			lw, _ := f(utf8.RuneCount([]byte(line)), fontSize, bold)
			if lw > maxWidth {
				maxWidth = lw
			}
			line = ww.text
			curWidth = ww.width
		}
	}
	if line != "" {
		lines = appendLineToOutput(line, curWidth, lh)
	}
	return maxWidth, heightSum, lines
}

type Drawing struct {
	canvas                 *svg.SVG
	txtDimensionCalculator *SvgTextDimensionCalculator
}

func NewDrawing(w io.Writer) *Drawing {
	return &Drawing{
		canvas:                 svg.New(w),
		txtDimensionCalculator: NewSvgTextDimensionCalculator(),
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
	if caption != "" {
		_, _, lines := d.txtDimensionCalculator.splitTxt(caption, &format.FontCaption)
		yTxt := y + format.Padding
		for _, l := range lines {
			xTxt := x + (width-l.width)/2
			txtFormat := "" // TODO
			d.canvas.Text(xTxt, yTxt, l.text, txtFormat)
			y += l.height
		}
	}
	return nil
}

func (d *Drawing) Done() error {
	d.canvas.End()
	return nil
}
