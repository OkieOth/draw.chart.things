package types

import (
	"fmt"
	"strings"
)

type BoxesDrawing interface {
	Start(title string, height, width int) error
	Draw(id, caption, text1, text2 string, x, y, width, height int, format BoxFormat) error
	DrawLine(x1, y1, x2, y2 int, format LineDef) error
	DrawArrow(x, y, angle int, format LineDef) error
	Done() error
}

func initLayoutElemContainer(l []Layout, inputFormats map[string]BoxFormat) *LayoutElemContainer {
	if len(l) == 0 {
		return nil
	}
	var ret LayoutElemContainer
	ret.Elems = make([]LayoutElement, 0)
	for _, elem := range l {
		ret.Elems = append(ret.Elems, initLayoutElement(&elem, inputFormats))
	}
	return &ret
}

func initConnections(l []Connection) []LayoutElemConnection {
	ret := make([]LayoutElemConnection, 0)
	for _, elem := range l {
		var conn LayoutElemConnection
		conn.DestId = elem.DestId
		conn.SourceArrow = elem.SourceArrow
		conn.DestArrow = elem.DestArrow
		conn.Tags = elem.Tags
		ret = append(ret, conn)
	}
	return ret
}

func InitFontDef(l *FontDef, defaultFont string, defaultSize int, defaultBold, defaultItalic bool, spaceTop int) FontDef {
	var f FontDef
	typeNormal := FontDefTypeEnum_normal
	typeItalic := FontDefTypeEnum_italic
	weightNormal := FontDefWeightEnum_normal
	weightBold := FontDefWeightEnum_bold
	alignedLeft := FontDefAlignedEnum_left

	if l != nil {
		if l.Font != "" {
			f.Font = l.Font
		} else {
			f.Font = defaultFont
		}
		if l.Size != 0 {
			f.Size = l.Size
		} else {
			f.Size = defaultSize
		}
		if l.Type != nil {
			f.Type = l.Type
		} else {
			if defaultItalic {
				f.Type = &typeItalic
			} else {
				f.Type = &typeNormal
			}
		}
		if l.Weight != nil {
			f.Weight = l.Weight
		} else {
			if defaultBold {
				f.Weight = &weightBold
			} else {
				f.Weight = &weightNormal
			}
		}
		if l.LineHeight != 0 {
			f.LineHeight = l.LineHeight
		} else {
			f.LineHeight = 1.5
		}
		if l.Color != "" {
			f.Color = l.Color
		} else {
			f.Color = "black"
		}
		if l.Aligned != nil {
			f.Aligned = l.Aligned
		} else {
			f.Aligned = &alignedLeft
		}
		f.SpaceTop = l.SpaceTop
		if f.SpaceTop == 0 {
			f.SpaceTop = spaceTop
		}
		f.SpaceBottom = l.SpaceBottom
		if f.MaxLenBeforeBreak != 0 {
			f.MaxLenBeforeBreak = l.MaxLenBeforeBreak
		} else {
			f.MaxLenBeforeBreak = 90
		}
	} else {
		f.Size = defaultSize
		if defaultItalic {
			f.Type = &typeItalic
		} else {
			f.Type = &typeNormal
		}
		f.Font = defaultFont
		if defaultBold {
			f.Weight = &weightBold
		} else {
			f.Weight = &weightNormal
		}
		f.LineHeight = 1.5
		f.Color = "black"
		f.Aligned = &alignedLeft
		f.SpaceTop = spaceTop
		f.SpaceBottom = 0
		f.MaxLenBeforeBreak = 90
	}
	return f
}

func initBoxFormat(f *Format) BoxFormat {
	var border *LineDef
	var fill *FillDef

	var fontCaption *FontDef
	var fontText1 *FontDef
	var fontText2 *FontDef
	var verticalTxt bool
	padding := GlobalPadding
	boxMargin := GlobalMinBoxMargin
	if f != nil {
		fontCaption = f.FontCaption
		fontText1 = f.FontText1
		fontText2 = f.FontText2
		border = f.Border
		fill = f.Fill
		if f.Padding != nil {
			padding = *f.Padding
		}
		if f.BoxMargin != nil {
			boxMargin = *f.BoxMargin
		}
		if f.VerticalTxt != nil {
			verticalTxt = *f.VerticalTxt
		}
	}

	return BoxFormat{
		Padding:      padding,
		MinBoxMargin: boxMargin,
		FontCaption:  InitFontDef(fontCaption, "sans-serif", 10, true, false, 0),
		FontText1:    InitFontDef(fontText1, "serif", 8, false, false, 10),
		FontText2:    InitFontDef(fontText2, "monospace", 8, false, true, 10),
		Border:       border,
		Fill:         fill,
		VerticalTxt:  verticalTxt,
	}
}

func getDefaultFormat() BoxFormat {
	return BoxFormat{
		Padding:      GlobalPadding,
		MinBoxMargin: GlobalMinBoxMargin,
		FontCaption:  InitFontDef(nil, "sans-serif", 10, true, false, 0),
		FontText1:    InitFontDef(nil, "serif", 8, false, false, 10),
		FontText2:    InitFontDef(nil, "monospace", 8, false, true, 10),
	}
}

func initFormats(inputFormat map[string]Format) map[string]BoxFormat {
	var res = make(map[string]BoxFormat)
	for key, elem := range inputFormat {
		res[key] = initBoxFormat(&elem)
	}
	if _, hasDefault := res["default"]; !hasDefault {
		res["default"] = getDefaultFormat()
	}
	return res
}

func initLayoutElement(l *Layout, inputFormats map[string]BoxFormat) LayoutElement {
	var f *BoxFormat
	for _, tag := range l.Tags {
		if val, ok := inputFormats[tag]; ok {
			f = &val
			break
		}
	}
	if (f == nil) && (l.Caption != "" || l.Text1 != "" || l.Text2 != "") {
		formatKey := "default"
		if l.Format != nil {
			formatKey = *l.Format
		}
		for key, format := range inputFormats {
			if key == formatKey {
				f = &format
				break
			}
		}
		if f == nil {
			d := getDefaultFormat()
			f = &d
		}
	}
	if l.Id == "" {
		if l.Caption != "" {
			l.Id = strings.ToLower(l.Caption)
		}
	}
	return LayoutElement{
		Id:          l.Id,
		Caption:     l.Caption,
		Text1:       l.Text1,
		Text2:       l.Text2,
		Vertical:    initLayoutElemContainer(l.Vertical, inputFormats),
		Horizontal:  initLayoutElemContainer(l.Horizontal, inputFormats),
		Format:      f,
		Connections: initConnections(l.Connections),
	}
}

func DocumentFromBoxes(b *Boxes) *BoxesDocument {
	doc := NewBoxesDocument()
	doc.Title = b.Title
	doc.Formats = initFormats(b.Formats)
	doc.Boxes = initLayoutElement(&b.Boxes, doc.Formats)
	if doc.MinBoxMargin == 0 {
		doc.MinBoxMargin = GlobalMinBoxMargin
	}
	if doc.MinConnectorMargin == 0 {
		doc.MinConnectorMargin = GlobalMinConnectorMargin
	}
	if doc.GlobalPadding == 0 {
		doc.GlobalPadding = GlobalPadding
	}
	return doc
}

func (d *BoxesDocument) DrawBoxes(drawingImpl BoxesDrawing) error {
	return d.Boxes.Draw(drawingImpl)
}

func (doc *BoxesDocument) drawStartPositionsImpl(drawingImpl *BoxesDrawing, elem *LayoutElement, f *LineDef) {
	if doc.ShouldHandle(elem) {
		(*drawingImpl).DrawLine(*elem.TopXToStart, elem.Y, *elem.TopXToStart, elem.Y-RasterSize, *f)
		(*drawingImpl).DrawLine(*elem.BottomXToStart, elem.Y+elem.Height, *elem.BottomXToStart, elem.Y+elem.Height+RasterSize, *f)
		(*drawingImpl).DrawLine(elem.X, *elem.LeftYToStart, elem.X-RasterSize, *elem.LeftYToStart, *f)
		(*drawingImpl).DrawLine(elem.X+elem.Width, *elem.RightYToStart, elem.X+elem.Width+RasterSize, *elem.RightYToStart, *f)
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

func (d *BoxesDocument) DrawStartPositions(drawingImpl BoxesDrawing) {
	w := 2
	b := "blue"
	f := LineDef{
		Width: &w,
		Color: &b,
	}
	d.InitStartPositions()
	d.drawStartPositionsImpl(&drawingImpl, &d.Boxes, &f)
}

func (d *BoxesDocument) DrawRoads(drawingImpl BoxesDrawing) {
	w := 1
	b := "lightgray"
	f := LineDef{
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

func (d *BoxesDocument) DrawConnections(drawingImpl BoxesDrawing) error {
	b := "black"
	w := 1
	format := LineDef{
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

func (b *LayoutElement) Draw(drawing BoxesDrawing) error {
	if b.Format != nil {
		if err := drawing.Draw(b.Id, b.Caption, b.Text1, b.Text2, b.X, b.Y, b.Width, b.Height, *b.Format); err != nil {
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
