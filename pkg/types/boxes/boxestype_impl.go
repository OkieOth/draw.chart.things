package boxes

import (
	"fmt"

	"github.com/okieoth/draw.chart.things/pkg/types"
)

func (d *BoxesDocument) DrawBoxes(drawingImpl types.Drawing) error {
	return d.Boxes.Draw(drawingImpl)
}

func (doc *BoxesDocument) drawStartPositionsImpl(drawingImpl *types.Drawing, elem *LayoutElement, f *types.LineDef) {
	if elem.TopXToStart != nil && elem.BottomXToStart != nil && elem.LeftYToStart != nil && elem.RightYToStart != nil {
		(*drawingImpl).DrawLine(*elem.TopXToStart, elem.Y, *elem.TopXToStart, elem.Y-2, *f)
		(*drawingImpl).DrawLine(*elem.BottomXToStart, elem.Y+elem.Height, *elem.BottomXToStart, elem.Y+elem.Height+2, *f)
		(*drawingImpl).DrawLine(elem.X, *elem.LeftYToStart, elem.X-2, *elem.LeftYToStart, *f)
		(*drawingImpl).DrawLine(elem.X+elem.Width, *elem.RightYToStart, elem.X+elem.Width+2, *elem.RightYToStart, *f)
	}
	if elem.Vertical != nil {
		for i := 0; i < len(elem.Vertical.Elems); i++ {
			doc.drawStartPositionsImpl(drawingImpl, &elem.Vertical.Elems[i], f)
		}
	}
	if elem.Horizontal != nil {
		for i := 0; i < len(elem.Horizontal.Elems); i++ {
			doc.drawStartPositionsImpl(drawingImpl, &elem.Horizontal.Elems[i], f)
		}
	}
}

func (d *BoxesDocument) ShouldHandle(elem *LayoutElement) bool {
	if elem == &d.Boxes {
		return false
	}
	if elem.Caption == "" && elem.Text1 == "" && elem.Text2 == "" && elem.Id == "" {
		return false
	}
	return true
}

func (d *BoxesDocument) DrawStartPositions(drawingImpl types.Drawing) {
	w := 2.0
	b := "blue"
	f := types.LineDef{
		Width: &w,
		Color: &b,
	}
	d.drawStartPositionsImpl(&drawingImpl, &d.Boxes, &f)
}

func (d *BoxesDocument) DrawConnectionNodes(drawingImpl types.Drawing) {
	w1 := 4.0
	b1 := "grey"
	w2 := 4.0
	b2 := "purple"
	w3 := 2.0
	b3 := "pink"
	w4 := 0.5
	f1 := types.LineDef{
		Width: &w1,
		Color: &b1,
	}
	f2 := types.LineDef{
		Width: &w2,
		Color: &b2,
	}
	f3 := types.LineDef{
		Width: &w3,
		Color: &b3,
	}
	f4 := types.LineDef{
		Width: &w4,
		Color: &b1,
	}
	_ = f3
	for _, n := range d.ConnectionNodes {
		if n.NodeId != nil && *n.NodeId != "" {
			drawingImpl.DrawLine(n.X, n.Y-2, n.X, n.Y+2, f2)
			drawingImpl.DrawLine(n.X-2, n.Y, n.X+2, n.Y, f2)
		} else {
			drawingImpl.DrawLine(n.X, n.Y-2, n.X, n.Y+2, f1)
			drawingImpl.DrawLine(n.X-2, n.Y, n.X+2, n.Y, f1)
		}
		for _, e := range n.Edges {
			nx, ny, ex, ey := n.X, n.Y, e.X, e.Y
			if n.X == e.X {
				// vertical line
				if n.Y > e.Y {
					// line up
					nx = nx - 2
				} else {
					// line down
					nx = nx + 2
				}
				ex = nx
			} else {
				// horizontal line
				if n.X > e.X {
					// line to left
					ny = ny - 2
				} else {
					// line to right
					ny = ny + 2
				}
				ey = ny
			}

			if e.DestNodeId != nil && e.Weight != 0 {
				drawingImpl.DrawLine(nx, ny, ex, ey, f3)
			} else {
				//fmt.Println("DEBUG: boxId=", *n.BoxId, ", id=", *n.NodeId, "edgeCount=", len(n.Edges), ", destNodeId=", e.DestNodeId)
				drawingImpl.DrawLine(nx, ny, ex, ey, f4)
			}
		}
	}
}

func (d *BoxesDocument) DrawRoads(drawingImpl types.Drawing) {
	w := 1.0
	b := "lightgray"
	f := types.LineDef{
		Width: &w,
		Color: &b,
	}
	for _, r := range d.VerticalRoads {
		drawingImpl.DrawLine(r.StartX, r.StartY, r.EndX, r.EndY, f)
	}
	for _, r := range d.HorizontalRoads {
		drawingImpl.DrawLine(r.StartX, r.StartY, r.EndX, r.EndY, f)
	}
}

func (d *BoxesDocument) adjustLineToWidth(startX, startY, endX, endY, offset int, firstLine, lastLine bool) (int, int, int, int) {
	if startX == endX {
		// vertical line
		if startY < endY {
			// top/down
			if !firstLine {
				startY += -offset
			}
			if !lastLine {
				endY += offset
			}
		} else {
			//botton/up
			if !firstLine {
				endY += -offset
			}
			if !lastLine {
				startY += offset
			}
		}
	} else {
		// horizontal line
		if startX < endX {
			// left/right
			if !firstLine {
				startX += -offset
			}
			if !lastLine {
				endX += offset
			}
		} else {
			// right/left
			if !firstLine {
				startX += offset
			}
			if !lastLine {
				endX += -offset
			}
		}
	}
	return startX, startY, endX, endY
}

func (d *BoxesDocument) DrawConnections(drawingImpl types.Drawing) error {
	// DEBUG - start
	// b := "blue"
	// w := 2.0
	// format := types.LineDef{
	// 	Width: &w,
	// 	Color: &b,
	// }

	// for _, elem := range d.Connections {
	// 	// iterate over the connections of the document
	// 	lineFormat := format
	// 	if elem.Format != nil {
	// 		lineFormat = *elem.Format
	// 	}
	// 	offset := int(*lineFormat.Width / 2)
	// 	lastPos := len(elem.Parts) - 1
	// 	for i, l := range elem.Parts {
	// 		// drawing the connection lines
	// 		x1, y1, x2, y2 := d.adjustLineToWidth(l.StartX, l.StartY, l.EndX, l.EndY, offset, i == 0, i == lastPos)
	// 		// DEBUG - Start
	// 		if x1 == x2 {
	// 			// vertical line
	// 			y1 += 3
	// 			y2 -= 3
	// 		} else {
	// 			// horizontal line
	// 			x1 += 3
	// 			x2 -= 3
	// 		}
	// 		// DEBUG - End
	// 		drawingImpl.DrawLine(x1, y1, x2, y2, lineFormat)
	// 	}
	// }
	// DEBUG - end
	d.DrawMovedConnectionLines(drawingImpl)
	return nil
}

func (d *BoxesDocument) DrawMovedConnectionLines(drawingImpl types.Drawing) {
	b := "black"
	w := 1.0
	s := types.LineDefStyleEnum_dashed
	format := types.LineDef{
		Width: &w,
		Color: &b,
		Style: &s,
	}
	for _, l := range d.HorizontalLines {
		lineFormat := format
		if l.Format != nil {
			lineFormat = *l.Format
		}
		offset := int(*lineFormat.Width / 2)
		x1, y1, x2, y2 := d.adjustLineToWidth(l.StartX, l.StartY, l.EndX, l.EndY, offset, l.IsStart, l.IsEnd)
		drawingImpl.DrawLine(x1, y1, x2, y2, lineFormat)
	}
	for _, l := range d.VerticalLines {
		lineFormat := format
		if l.Format != nil {
			lineFormat = *l.Format
		}
		drawingImpl.DrawLine(l.StartX, l.StartY, l.EndX, l.EndY, lineFormat)
	}
}

func (doc *BoxesDocument) AdjustDocHeight(le *LayoutElement, currentMax int) int {
	if le != &doc.Boxes {
		if le.Y+le.Height > currentMax {
			currentMax = le.Y + le.Height
		}
	} else {
		currentMax = doc.Boxes.Height
	}
	if le.Vertical != nil {
		for _, elem := range le.Vertical.Elems {
			currentMax = doc.AdjustDocHeight(&elem, currentMax)
		}
	}
	if le.Horizontal != nil {
		for _, elem := range le.Horizontal.Elems {
			currentMax = doc.AdjustDocHeight(&elem, currentMax)
		}
	}
	return currentMax
}

func (b *LayoutElement) Draw(drawing types.Drawing) error {
	if b.Format != nil {
		f := types.RectWithTextFormat{
			FontCaption:  b.Format.FontCaption,
			FontText1:    b.Format.FontText1,
			FontText2:    b.Format.FontText2,
			Padding:      b.Format.Padding,
			Border:       b.Format.Line,
			Fill:         b.Format.Fill,
			VerticalTxt:  b.Format.VerticalTxt,
			CornerRadius: b.Format.CornerRadius,
		}
		textYOffset := 0
		if b.Image != nil {
			textYOffset = (b.Image.Y - b.Y) + b.Image.Height + b.Image.MarginTopBottom
		}
		if err := drawing.DrawRectWithText(b.Id, b.Caption, b.Text1, b.Text2, b.X, b.Y, b.Width, b.Height, textYOffset, f); err != nil {
			return fmt.Errorf("Error drawing element %s: %w", b.Id, err)
		}
	}
	if b.Image != nil {
		if err := drawing.DrawPng(b.Image.X, b.Image.Y, b.Image.ImgId); err != nil {
			return fmt.Errorf("Error drawing image %s: %w", b.Id, err)
		}
	}
	if b.Vertical != nil {
		for _, elem := range b.Vertical.Elems {
			if err := elem.Draw(drawing); err != nil {
				return err
			}
		}
	}
	if b.Horizontal != nil {
		for _, elem := range b.Horizontal.Elems {
			if err := elem.Draw(drawing); err != nil {
				return err
			}
		}
	}
	return nil
}

func (b *LayoutElement) DrawTextBoxes(drawing types.Drawing) error {
	if b.Format != nil {
		c := "black"
		w := 2.0
		o := 0.4
		f := types.LineDef{
			Color:   &c,
			Width:   &w,
			Opacity: &o,
		}
		if b.XTextBox != nil {
			if err := drawing.DrawSolidRect(*b.XTextBox, *b.YTextBox, *b.WidthTextBox, *b.HeightTextBox, f); err != nil {
				return fmt.Errorf("Error drawing text boxes element %s: %w", b.Id, err)
			}
		}
	}
	if b.Vertical != nil {
		for _, elem := range b.Vertical.Elems {
			if err := elem.DrawTextBoxes(drawing); err != nil {
				return err
			}
		}
	}
	if b.Horizontal != nil {
		for _, elem := range b.Horizontal.Elems {
			if err := elem.DrawTextBoxes(drawing); err != nil {
				return err
			}
		}
	}
	return nil
}
