package svgdrawing

import (
	"fmt"
	"io"
	"strings"
	"unicode/utf8"

	"github.com/okieoth/draw.chart.things/pkg/svg"
	"github.com/okieoth/draw.chart.things/pkg/types"
)

// Font width factor constants
const (
	MAGIC = 1.1
	// Sans-serif font width factors
	SANSSERIF_NORMAL_WIDTH_FACTOR = 2.6 / 6 * MAGIC
	SANSSERIF_BOLD_WIDTH_FACTOR   = 2.8 / 6 * MAGIC

	// Monospace font width factors
	MONOSPACE_NORMAL_WIDTH_FACTOR = 3.0 / 5 * MAGIC
	MONOSPACE_BOLD_WIDTH_FACTOR   = 3.0 / 5 * MAGIC

	// Serif font width factors
	SERIF_NORMAL_WIDTH_FACTOR = 2.05 / 5 * MAGIC
	SERIF_BOLD_WIDTH_FACTOR   = 2.15 / 5 * MAGIC
)

// Function type for calculating text dimensions
type calcDimensions func(runeCount int, fontSize int, bold bool) (width, height int)

// textAndDimensions represents text and its calculated dimensions
type textAndDimensions struct {
	Text   string
	Width  int
	Height int
}

// SvgTextDimensionCalculator calculates text dimensions for SVG rendering
type SvgTextDimensionCalculator struct {
	dimensionCalculators map[string]calcDimensions
}

// NewSvgTextDimensionCalculator creates a new text dimension calculator
func NewSvgTextDimensionCalculator() *SvgTextDimensionCalculator {
	calc := &SvgTextDimensionCalculator{
		dimensionCalculators: make(map[string]calcDimensions),
	}

	calc.dimensionCalculators["serif"] = serifDimensions
	calc.dimensionCalculators["sans-serif"] = sansserifDimensions
	calc.dimensionCalculators["monospace"] = monospaceDimensions

	return calc
}

// Font dimension calculators
func serifDimensions(runeCount int, fontSize int, bold bool) (width, height int) {
	var factor float32
	if bold {
		factor = float32(fontSize*10) * SERIF_BOLD_WIDTH_FACTOR
	} else {
		factor = float32(fontSize*10) * SERIF_NORMAL_WIDTH_FACTOR
	}
	w := int(factor) * runeCount / 10
	return w, fontSize
}

func sansserifDimensions(runeCount int, fontSize int, bold bool) (width, height int) {
	var factor float32
	if bold {
		factor = float32(fontSize*10) * SANSSERIF_BOLD_WIDTH_FACTOR
	} else {
		factor = float32(fontSize*10) * SANSSERIF_NORMAL_WIDTH_FACTOR
	}
	w := int(factor) * runeCount / 10
	return w, fontSize
}

func monospaceDimensions(runeCount int, fontSize int, bold bool) (width, height int) {
	var factor float32
	if bold {
		factor = float32(fontSize*10) * MONOSPACE_BOLD_WIDTH_FACTOR
	} else {
		factor = float32(fontSize*10) * MONOSPACE_NORMAL_WIDTH_FACTOR
	}
	w := int(factor) * runeCount / 10
	return w, fontSize
}

// SplitTxt splits text and calculates dimensions based on font settings
func (s *SvgTextDimensionCalculator) SplitTxt(txt string, format *types.FontDef) (width, height int, lines []textAndDimensions) {
	fontSize := format.Size
	if fontSize == 0 {
		fontSize = 10
	}

	bold := format.Weight != nil && *format.Weight == types.FontDefWeightEnum_bold
	lineHeight := format.LineHeight
	if lineHeight == 0 {
		lineHeight = 1.2
	}

	runeCount := utf8.RuneCountInString(txt)
	fontType := format.Font
	if fontType == "" {
		fontType = "sans-serif"
	}

	// Get appropriate dimension calculator or default to sans-serif
	dimCalc, exists := s.dimensionCalculators[fontType]
	if !exists {
		dimCalc = s.dimensionCalculators["sans-serif"]
	}

	w, h := dimCalc(runeCount, fontSize, bold)

	if w <= format.MaxLenBeforeBreak {
		lines = []textAndDimensions{{Text: txt, Width: w, Height: h}}
		return w, h, lines
	}

	return splitTxtDimensions(w, runeCount, fontSize, bold, lineHeight, txt, format.MaxLenBeforeBreak, dimCalc)
}

// Dimensions calculates text dimensions
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
	wordsWithWidth := make([]textAndDimensions, 0, len(words))

	// Calculate width for each word
	for _, w := range words {
		rc := utf8.RuneCount([]byte(w))
		width, _ := f(rc, fontSize, bold)

		// Truncate words that are too long
		if width > maxLenBeforeBreak {
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

	// Build output lines
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
			line = ww.Text
			curWidth = ww.Width
		}
	}

	if line != "" {
		lines = appendLineToOutput(line, lh)
	}

	return maxWidth, heightSum, lines
}

// Drawing represents an SVG drawing
type Drawing struct {
	canvas                 *svg.SVG
	txtDimensionCalculator *SvgTextDimensionCalculator
}

// NewDrawing creates a new SVG drawing
func NewDrawing(w io.Writer) *Drawing {
	return &Drawing{
		canvas:                 svg.New(w),
		txtDimensionCalculator: NewSvgTextDimensionCalculator(),
	}
}

// Start initializes the SVG drawing
func (d *Drawing) Start(title string, height, width int) error {
	args := fmt.Sprintf("viewBox=\"0,0,%d,%d\"", width, height)
	d.canvas.Start(width, height, args)
	return nil
}

func (d *Drawing) DrawRaster(width, height, rasterSize int) {
	for x := rasterSize; x < width; x += rasterSize {
		for y := rasterSize; y < height; y += rasterSize {
			// vertical lines
			d.canvas.Line(x, y-1, x, y+1, "stroke:lightgray;stroke-width:0.5")
			// horizontal lines
			d.canvas.Line(x-1, y, x+1, y, "stroke:lightgray;stroke-width:0.5")
		}
	}
}

// textFormat generates CSS formatting for text
func (d *Drawing) textFormat(fontDef *types.FontDef) string {
	var font string

	switch fontDef.Font {
	case "serif":
		font = "Times New Roman, Times, Georgia, serif"
	case "sans-serif":
		font = "Arial, Helvetica, sans-serif"
	case "monospace":
		font = "Courier New, Courier, Lucida Console, monospace"
	default:
		if fontDef.Font != "" {
			font = fontDef.Font
		} else {
			font = "Arial, Helvetica, sans-serif"
		}
	}

	txtFormat := fmt.Sprintf("text-anchor:middle;font-family:%s;font-size:%dpx;fill:%s",
		font, fontDef.Size, fontDef.Color)

	if fontDef.Weight != nil && *fontDef.Weight == types.FontDefWeightEnum_bold {
		txtFormat += ";font-weight:bold"
	}

	if fontDef.Type != nil && *fontDef.Type == types.FontDefTypeEnum_italic {
		txtFormat += ";font-style:italic"
	}

	return txtFormat
}

// drawText renders text with appropriate positioning and returns the updated y-position
func (d *Drawing) drawText(text string, currentY, x, width int, fontDef *types.FontDef) int {
	if text == "" {
		return currentY
	}

	if fontDef.SpaceTop != 0 {
		currentY += fontDef.SpaceTop
	}

	txtFormat := d.textFormat(fontDef)
	_, _, lines := d.txtDimensionCalculator.SplitTxt(text, fontDef)

	for _, l := range lines {
		yTxt := currentY + fontDef.Size
		xTxt := x + (width / 2)
		d.canvas.Text(xTxt, yTxt, l.Text, txtFormat)
		currentY += l.Height
	}

	if fontDef.SpaceBottom != 0 {
		currentY += fontDef.SpaceBottom
	}

	return currentY
}

// Draw renders a box with text elements
func (d *Drawing) Draw(id, caption, text1, text2 string, x, y, width, height int, format types.BoxFormat) error {
	if format.Fill != nil || format.Border != nil {
		attr := ""

		if format.Fill != nil {
			attr = fmt.Sprintf("fill: %s", *(*format.Fill).Color)
			if (*format.Fill).Opacity != nil {
				attr += fmt.Sprintf(";opacity: %f", *(*format.Fill).Opacity)
			}
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
	currentY = d.drawText(text2, currentY, x, width, &format.FontText2)

	return nil
}

func (d *Drawing) DrawLine(x1, y1, x2, y2 int, format types.LineDef) error {
	color := "black"
	if (format.Color != nil) && (*format.Color != "") {
		color = *format.Color
	}
	width := 1
	if (format.Width != nil) && (*format.Width != 0) {
		width = *format.Width
	}
	d.canvas.Line(x1, y1, x2, y2, fmt.Sprintf("stroke:%s;stroke-width:%d", color, width))
	return nil
}

func (d *Drawing) DrawArrow(x, y, angle int, format types.LineDef) error {
	// TODO
	return nil
}

// Done finalizes the SVG document
func (d *Drawing) Done() error {
	d.canvas.End()
	return nil
}
