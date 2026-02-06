package boxes

import (
	"fmt"

	"github.com/okieoth/draw.chart.things/pkg/types"
)

func (doc *BoxesDocument) IncludeComments(c types.TextDimensionCalculator) error {
	doc.collectCommentsFromLayout(&doc.Boxes, c)
	doc.collectCommentsFromConnections(c)
	doc.CommentMarkerRadius = (doc.CommentMarkerRadius / 2) + 4
	// adjust the document height, based on the comments
	if len(doc.Comments) > 0 {
		currentY := doc.Height + types.GlobalPadding
		for i := range doc.Comments {
			currentY += doc.Comments[i].TextHeight
			currentY += types.GlobalPadding

		}
		doc.Height = currentY + doc.GlobalPadding
	}
	return nil
}

func (doc *BoxesDocument) collectCommentsFromConnections(dimensionsCalc types.TextDimensionCalculator) {
CONNECTIONS:
	for i := range doc.Connections {
		c := doc.Connections[i]
		if c.Comment != nil {
			label := doc.newLabel(c.Comment.Label)
			for li := range doc.HorizontalLines {
				// search for the start line of the connection in the horizontal lines
				l := doc.HorizontalLines[li]
				if l.ConnectionIndex == c.ConnectionIndex && l.IsStart {
					diff, changed := absInt2(l.EndX - l.StartX)
					var x int
					if changed {
						x = l.EndX - (diff / 2)
					} else {
						x = l.StartX + (diff / 2)
					}
					y := l.StartY

					c := doc.newCommentContainer(c.Comment.Text, label, c.Comment.Format, x, y, false, dimensionsCalc)
					doc.Comments = append(doc.Comments, c)
					continue CONNECTIONS
				}
			}
			for li := range doc.VerticalLines {
				// search for the start line of the connection in the vertical lines
				l := doc.VerticalLines[li]
				if l.ConnectionIndex == c.ConnectionIndex && l.IsStart {
					x := l.StartX

					diff, changed := absInt2(l.EndY - l.StartY)
					var y int
					if changed {
						y = l.EndY - (diff / 2)
					} else {
						y = l.StartY + (diff / 2)
					}
					c := doc.newCommentContainer(c.Comment.Text, label, c.Comment.Format, x, y, false, dimensionsCalc)
					doc.Comments = append(doc.Comments, c)
					continue CONNECTIONS
				}
			}
		}
	}
}

func (doc *BoxesDocument) GetCommentFormat(format *string) CommentFormat {
	var boxFormat *BoxFormat
	if format != nil {
		if f, ok := doc.Formats[*format]; ok {
			boxFormat = &f
		}
	}
	if boxFormat == nil {
		if f, ok := doc.Formats["defaultComment"]; ok {
			boxFormat = &f
		}
	}
	if boxFormat == nil {
		if f, ok := doc.Formats["default"]; ok {
			boxFormat = &f
		}
	}
	if boxFormat == nil {
		return CommentFormat{
			FontText:   types.InitFontDef(nil, "serif", 8, false, false, 10),
			FontMarker: types.InitFontDef(nil, "monospace", 8, false, true, 10),
			Line:       *types.InitLineDef(nil),
			Fill:       *types.InitFillDef(nil),
		}

	} else {
		if boxFormat.Fill == nil {
			boxFormat.Fill = types.InitFillDef(nil)
		}
		if boxFormat.Line == nil {
			boxFormat.Line = types.InitLineDef(nil)
		}
		return CommentFormat{
			FontText:   boxFormat.FontComment,
			FontMarker: boxFormat.FontCommentMarker,
			Line:       *boxFormat.Line,
			Fill:       *boxFormat.Fill,
		}

	}
}

func (doc *BoxesDocument) newCommentContainer(
	text, label string,
	format *string,
	x, y int,
	moveMarker bool,
	dimensionsCalc types.TextDimensionCalculator) CommentContainer {
	if moveMarker {
		diff := doc.GlobalPadding + 2
		x, y = x+diff, y+diff
	}

	f := doc.GetCommentFormat(format)
	tw, th := dimensionsCalc.DimensionsWithMaxWidth(text, &f.FontText, doc.Width)
	mw, mh := dimensionsCalc.DimensionsWithMaxWidth(label, &f.FontMarker, doc.Width)

	if mw > doc.CommentMarkerRadius {
		doc.CommentMarkerRadius = mw
	}
	if mh > doc.CommentMarkerRadius {
		doc.CommentMarkerRadius = mh
	}

	return CommentContainer{
		Text:             text,
		Label:            label,
		Format:           f,
		MarkerX:          x,
		MarkerY:          y,
		TextWidth:        tw,
		TextHeight:       th,
		MarkerTextWidth:  mw,
		MarkerTextHeight: mh,
	}
}

func (doc *BoxesDocument) newLabel(label *string) string {
	if label != nil {
		return *label
	} else {
		doc.CommentCurrent += 1
		return fmt.Sprintf("%d", doc.CommentCurrent)
	}
}

func (doc *BoxesDocument) collectCommentsFromLayout(l *LayoutElement, dimensionsCalc types.TextDimensionCalculator) {
	if l.Comment != nil {
		label := doc.newLabel(l.Comment.Label)
		c := doc.newCommentContainer(l.Comment.Text, label, l.Comment.Format, l.X, l.Y, true, dimensionsCalc)
		doc.Comments = append(doc.Comments, c)
	}
	doc.collectCommentsFromLayoutCont(l.Horizontal, dimensionsCalc)
	doc.collectCommentsFromLayoutCont(l.Vertical, dimensionsCalc)
}

func (doc *BoxesDocument) collectCommentsFromLayoutCont(cont *LayoutElemContainer, dimensionsCalc types.TextDimensionCalculator) {
	if cont != nil {
		for i := range cont.Elems {
			doc.collectCommentsFromLayout(&cont.Elems[i], dimensionsCalc)
		}
	}
}
