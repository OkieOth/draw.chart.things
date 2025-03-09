package types

func initLayoutElemArray(l []Layout) []LayoutElement {
	var res = make([]LayoutElement, 0)
	for _, elem := range l {
		res = append(res, initLayoutElement(&elem))
	}
	return res
}

// func initFontDef(l *FontDef) FontDef {
// 	return FontDef{
// 		Family: l.Family,
// 		Size:   l.Size,
// 		Style:  l.Style,
// 		Color:  l.Color,
// 	}
// }

// func initBoxFormat(l *Layout) BoxFormat {
// 	var border *LineDef
// 	var fill *FillDef

// 	return BoxFormat{
// 		Padding:       5,
// 		LineHeight:    1.5,
// 		FontCaption:   initFontDef(&l.FontCaption),
// 		CaptionBefore: 0,
// 		FontText1:     *fontText1,
// 		Text1Before:   30,
// 		FontText2:     *fontText2,
// 		Text2Before:   30,
// 		Border:        *border,
// 		Fill:          *fill,
// 	}
// }

func initLayoutElement(l *Layout) LayoutElement {
	return LayoutElement{
		Id:         l.Id,
		Caption:    l.Caption,
		Text1:      l.Text1,
		Text2:      l.Text2,
		Vertical:   initLayoutElemArray(l.Vertical),
		Horizontal: initLayoutElemArray(l.Horizontal),
		//BoxFormat:  initBoxFormat(l),
	}
}

func DocumentFromBoxes(b *Boxes) *BoxesDocument {
	doc := NewBoxesDocument()
	doc.Title = b.Title
	doc.Boxes = initLayoutElement(&b.Boxes)
	return doc
}
