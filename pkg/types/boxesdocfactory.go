package types

func initLayoutElemArray(l []Layout) []LayoutElement {
	var res = make([]LayoutElement, 0)
	for _, elem := range l {
		res = append(res, initLayoutElement(&elem))
	}
	return res
}

func initLayoutElement(l *Layout) LayoutElement {
	return LayoutElement{
		Id:         l.Id,
		Caption:    l.Caption,
		Text1:      l.Text1,
		Text2:      l.Text2,
		Vertical:   initLayoutElemArray(l.Vertical),
		Horizontal: initLayoutElemArray(l.Horizontal),
	}
}

func DocumentFromBoxes(b *Boxes) *BoxesDocument {
	doc := NewBoxesDocument()
	doc.Title = b.Title
	doc.Boxes = initLayoutElement(&b.Boxes)
	return doc
}
