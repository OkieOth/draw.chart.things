package boxes

import (
	"fmt"
	"github.com/okieoth/draw.chart.things/pkg/types"
)

func (d *BoxesDocument) DrawBoxes(drawingImpl types.Drawing) error {
	return d.Boxes.Draw(drawingImpl)
}

func (doc *BoxesDocument) drawStartPositionsImpl(drawingImpl *types.Drawing, elem *LayoutElement, f *types.LineDef) {
	if doc.ShouldHandle(elem) {
		(*drawingImpl).DrawLine(*elem.TopXToStart, elem.Y, *elem.TopXToStart, elem.Y-types.RasterSize, *f)
		(*drawingImpl).DrawLine(*elem.BottomXToStart, elem.Y+elem.Height, *elem.BottomXToStart, elem.Y+elem.Height+types.RasterSize, *f)
		(*drawingImpl).DrawLine(elem.X, *elem.LeftYToStart, elem.X-types.RasterSize, *elem.LeftYToStart, *f)
		(*drawingImpl).DrawLine(elem.X+elem.Width, *elem.RightYToStart, elem.X+elem.Width+types.RasterSize, *elem.RightYToStart, *f)
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
	w := 2
	b := "blue"
	f := types.LineDef{
		Width: &w,
		Color: &b,
	}
	d.InitStartPositions()
	d.drawStartPositionsImpl(&drawingImpl, &d.Boxes, &f)
}

func (d *BoxesDocument) DrawRoads(drawingImpl types.Drawing) {
	w := 1
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

func (d *BoxesDocument) DrawConnections(drawingImpl types.Drawing) error {
	b := "black"
	w := 1
	format := types.LineDef{
		Width: &w,
		Color: &b,
	}

	for _, elem := range d.Connections {
		// iterate over the connections of the document
		for _, l := range elem.Parts {
			// drawing the connection lines
			drawingImpl.DrawLine(l.StartX, l.StartY, l.EndX, l.EndY, format)
		}

	}
	return nil
}

func (doc *BoxesDocument) AdjustDocHeight(le *LayoutElement, currentMax int) int {
	if le != &doc.Boxes {
		if le.Y+le.Height > currentMax {
			currentMax = le.Y + le.Height
		}
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
			FontCaption: b.Format.FontCaption,
			FontText1:   b.Format.FontText1,
			FontText2:   b.Format.FontText2,
			Padding:     b.Format.Padding,
			Border:      b.Format.Border,
			Fill:        b.Format.Fill,
			VerticalTxt: b.Format.VerticalTxt,
		}
		if err := drawing.DrawRectWithText(b.Id, b.Caption, b.Text1, b.Text2, b.X, b.Y, b.Width, b.Height, f); err != nil {
			return fmt.Errorf("Error drawing element %s: %w", b.Id, err)
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
