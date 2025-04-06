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
	Text   string
	Width  int
	Height int
}

func NewSvgTextDimensionCalculator() *SvgTextDimensionCalculator {
	return &SvgTextDimensionCalculator{}
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

func (s *SvgTextDimensionCalculator) SplitTxt(txt string, format *types.FontDef) (width, height int, lines []textAndDimensions) {
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
			lines = []textAndDimensions{{Text: txt, Width: w, Height: h}}
		} else {
			w, h, lines = splitTxtDimensions(w, runeCount, fontSize, bold, lineHeight, txt, format.MaxLenBeforeBreak, serifDimensions)
		}
	case "sans-serif":
		if w, h = sansserifDimensions(runeCount, fontSize, bold); w <= format.MaxLenBeforeBreak {
			lines = []textAndDimensions{{Text: txt, Width: w, Height: h}}
		} else {
			w, h, lines = splitTxtDimensions(w, runeCount, fontSize, bold, lineHeight, txt, format.MaxLenBeforeBreak, sansserifDimensions)
		}
	case "monospace":
		if w, h = monospaceDimensions(runeCount, fontSize, bold); w <= format.MaxLenBeforeBreak {
			lines = []textAndDimensions{{Text: txt, Width: w, Height: h}}
		} else {
			w, h, lines = splitTxtDimensions(w, runeCount, fontSize, bold, lineHeight, txt, format.MaxLenBeforeBreak, monospaceDimensions)
		}
	default:
		if w, h = sansserifDimensions(runeCount, fontSize, bold); w <= format.MaxLenBeforeBreak {
			lines = []textAndDimensions{{Text: txt, Width: w, Height: h}}
		} else {
			w, h, lines = splitTxtDimensions(w, runeCount, fontSize, bold, lineHeight, txt, format.MaxLenBeforeBreak, sansserifDimensions)
		}
	}
	return w, h, lines
}

func (s *SvgTextDimensionCalculator) Dimensions(txt string, format *types.FontDef) (width, height int) {
	w, h, _ := s.SplitTxt(txt, format)
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
					break
				}
			}
		}
		wordsWithWidth = append(wordsWithWidth, textAndDimensions{Text: w, Width: width})
	}
	// build the output lines ...
	curWidth := 0
	var line string
	var maxWidth int
	lines = make([]textAndDimensions, 0)
	lh := int(float32(fontSize) * lineHeight)
	heightSum := 0
	appendLineToOutput := func(txt string, height int) []textAndDimensions {
		lw, _ := f(utf8.RuneCount([]byte(txt)), fontSize, bold)
		if lw > maxWidth {
			maxWidth = lw
		}
		lines = append(lines, textAndDimensions{Text: txt, Width: lw, Height: height})
		heightSum += lh
		return lines
	}
	for _, ww := range wordsWithWidth {
		if curWidth+ww.Width < maxLenBeforeBreak {
			if line != "" {
				line += " "
			}
			line += ww.Text
			curWidth += ww.Width
		} else {
			lines = appendLineToOutput(line, lh)
			lw, _ := f(utf8.RuneCount([]byte(line)), fontSize, bold)
			if lw > maxWidth {
				maxWidth = lw
			}
			line = ww.Text
			curWidth = ww.Width
		}
	}
	if line != "" {
		lines = appendLineToOutput(line, lh)
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

func (d *Drawing) textFormat(fontDev *types.FontDef) string {
	var font string
	switch fontDev.Font {
	case "serif":
		font = "Times New Roman, Times, Georgia, serif"
	case "sans-serif":
		font = "Arial, Helvetica, sans-serif"
	case "monospace":
		font = "Courier New, Courier, Lucida Console, monospace"
	default:
		if fontDev.Font != "" {
			font = fontDev.Font
		} else {
			font = "Arial, Helvetica, sans-serif;"
		}
	}
	txtFormat := fmt.Sprintf("text-anchor:start;font-family:%s;font-size:%dpx;fill:%s", font, fontDev.Size, fontDev.Color)
	if fontDev.Weight != nil && *fontDev.Weight == types.FontDefWeightEnum_bold {
		txtFormat += ";font-weight:bold"
	}
	if fontDev.Type != nil && *fontDev.Type == types.FontDefTypeEnum_italic {
		txtFormat += ";font-style:italic"
	}
	return txtFormat
}

func (d *Drawing) drawText(text string, currentY, x, width int, fontDev *types.FontDef) int {
	if text != "" {
		if fontDev.SpaceTop != 0 {
			currentY += fontDev.SpaceTop
		}
		_, _, lines := d.txtDimensionCalculator.SplitTxt(text, fontDev)
		for _, l := range lines {
			yTxt := currentY + fontDev.Size
			xTxt := x + (width-l.Width)/2
			txtFormat := d.textFormat(fontDev)
			d.canvas.Text(xTxt, yTxt, l.Text, txtFormat)
			currentY += l.Height
		}
		if fontDev.SpaceBottom != 0 {
			currentY += fontDev.SpaceBottom
		}
	}
	return currentY
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
			attr += fmt.Sprintf("stroke: %s;stroke-Width: %d", c, w)
		}
		d.canvas.RectWithId(id, x, y, width, height, attr)
	}
	currentY := y + format.Padding
	currentY = d.drawText(caption, currentY, x, width, &format.FontCaption)
	currentY = d.drawText(text1, currentY, x, width, &format.FontText1)
	d.drawText(text2, currentY, x, width, &format.FontText2)
	return nil
}

func (d *Drawing) Done() error {
	d.canvas.End()
	return nil
}
