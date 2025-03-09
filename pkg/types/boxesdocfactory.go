package types

func initLayoutElemArray(l []Layout, inputFormats map[string]Format) []LayoutElement {
	var res = make([]LayoutElement, 0)
	for _, elem := range l {
		res = append(res, initLayoutElement(&elem, inputFormats))
	}
	return res
}

func initFontDef(l *FontDef) FontDef {
	var f FontDef
	typeNormal := FontDefTypeEnum_normal
	weightNormal := FontDefWeightEnum_normal
	alignedLeft := FontDefAlignedEnum_left

	if l != nil {
		if l.Size != 0 {
			f.Size = l.Size
		} else {
			f.Size = 10
		}
		if l.Type != nil {
			f.Type = l.Type
		} else {
			f.Type = &typeNormal
		}
		if l.Weight != nil {
			f.Weight = l.Weight
		} else {
			f.Weight = &weightNormal
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
	} else {
		f.Size = 10
		f.Type = &typeNormal
		f.Weight = &weightNormal
		f.LineHeight = 1.5
		f.Color = "black"
		f.Aligned = &alignedLeft
		f.SpaceTop = 0
		f.SpaceBottom = 0
	}
	return f
}

func initBoxFormat(l *Layout, f *Format) BoxFormat {
	var border *LineDef
	var fill *FillDef

	var fontCaption *FontDef
	var fontText1 *FontDef
	var fontText2 *FontDef
	if f != nil {
		fontCaption = f.FontCaption
		fontText1 = f.FontText1
		fontText2 = f.FontText2
		border = f.Border
		fill = f.Fill
	}

	return BoxFormat{
		Padding:     5,
		FontCaption: initFontDef(fontCaption),
		FontText1:   initFontDef(fontText1),
		FontText2:   initFontDef(fontText2),
		Border:      border,
		Fill:        fill,
	}
}

func initLayoutElement(l *Layout, inputFormats map[string]Format) LayoutElement {
	var f *Format
	for _, tag := range l.Tags {
		if val, ok := inputFormats[tag]; ok {
			f = &val
			break
		}
	}
	return LayoutElement{
		Id:         l.Id,
		Caption:    l.Caption,
		Text1:      l.Text1,
		Text2:      l.Text2,
		Vertical:   initLayoutElemArray(l.Vertical, inputFormats),
		Horizontal: initLayoutElemArray(l.Horizontal, inputFormats),
		Format:     initBoxFormat(l, f),
	}
}

func DocumentFromBoxes(b *Boxes) *BoxesDocument {
	doc := NewBoxesDocument()
	doc.Title = b.Title
	doc.Boxes = initLayoutElement(&b.Boxes, b.Formats)
	return doc
}
