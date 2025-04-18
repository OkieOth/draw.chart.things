package types

import (
	"fmt"
)

type BoxesDrawing interface {
	Start(title string, height, width int) error
	Draw(id, caption, text1, text2 string, x, y, width, height int, format BoxFormat) error
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
	}

	return BoxFormat{
		Padding:      padding,
		MinBoxMargin: boxMargin,
		FontCaption:  InitFontDef(fontCaption, "sans-serif", 10, true, false, 0),
		FontText1:    InitFontDef(fontText1, "serif", 8, false, false, 10),
		FontText2:    InitFontDef(fontText2, "monospace", 8, false, true, 10),
		Border:       border,
		Fill:         fill,
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
		for key, format := range inputFormats {
			if key == "default" {
				f = &format
				break
			}
		}
		if f == nil {
			d := getDefaultFormat()
			f = &d
		}
	}
	return LayoutElement{
		Id:         l.Id,
		Caption:    l.Caption,
		Text1:      l.Text1,
		Text2:      l.Text2,
		Vertical:   initLayoutElemContainer(l.Vertical, inputFormats),
		Horizontal: initLayoutElemContainer(l.Horizontal, inputFormats),
		Format:     f,
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
		doc.MinConnectorMargin = GlobalMinBoxMargin
	}
	if doc.GlobalPadding == 0 {
		doc.GlobalPadding = GlobalPadding
	}
	return doc
}

func (d *BoxesDocument) DrawBoxes(drawingImpl BoxesDrawing) error {
	if err := drawingImpl.Start(d.Title, d.Height, d.Width); err != nil {
		return fmt.Errorf("Error starting drawing: %w", err)
	}
	return d.Boxes.Draw(drawingImpl)
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
