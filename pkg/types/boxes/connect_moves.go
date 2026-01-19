package boxes

func (doc *BoxesDocument) moveBoxContVertical(cont *LayoutElemContainer, startX, offset int) {
	if cont == nil {
		return
	}
	if cont.X < startX && (cont.X+cont.Width) > startX {
		cont.Width += offset
	} else if cont.X >= startX {
		cont.X += offset
	}
	for i := range len(cont.Elems) {
		doc.moveBoxVertical(&cont.Elems[i], startX, offset)
	}
}

func (doc *BoxesDocument) moveBoxVertical(box *LayoutElement, startX, offset int) {
	if box.X < startX && (box.X+box.Width) > startX {
		box.Width += offset
	} else if box.X > startX {
		box.X += offset
	}
	doc.moveBoxContVertical(box.Horizontal, startX, offset)
	doc.moveBoxContVertical(box.Vertical, startX, offset)
}

func (doc *BoxesDocument) MoveBoxesVertical(startX, offset int) {
	doc.moveBoxVertical(&doc.Boxes, startX, offset)
}

func (doc *BoxesDocument) moveBoxContHorizontal(cont *LayoutElemContainer, startY, offset int) {
	if cont == nil {
		return
	}
	if cont.Y < startY && (cont.Y+cont.Height) > startY {
		cont.Height += offset
	} else if cont.Y >= startY {
		cont.Y += offset
	}
	for i := range len(cont.Elems) {
		doc.moveBoxHorizontal(&cont.Elems[i], startY, offset)
	}
}

func (doc *BoxesDocument) moveBoxHorizontal(box *LayoutElement, startY, offset int) {
	if box.Y < startY && (box.Y+box.Height) > startY {
		box.Height += offset
	} else if box.Y > startY {
		box.Y += offset
	}
	doc.moveBoxContHorizontal(box.Horizontal, startY, offset)
	doc.moveBoxContHorizontal(box.Vertical, startY, offset)
}

func (doc *BoxesDocument) MoveBoxesHorizontal(startY, offset int) {
	doc.moveBoxHorizontal(&doc.Boxes, startY, offset)
}
