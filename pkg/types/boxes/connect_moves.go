package boxes

func (doc *BoxesDocument) moveBoxContHorizontal(cont *LayoutElemContainer, startX, offset int) {
	if cont == nil {
		return
	}
	if cont.X < startX && (cont.X+cont.Width) > startX {
		cont.Width += offset
	} else if cont.X >= startX {
		cont.X += offset
	}
	for i := range len(cont.Elems) {
		doc.moveBoxHorizontal(&cont.Elems[i], startX, offset)
	}
}

func (doc *BoxesDocument) moveBoxHorizontal(box *LayoutElement, startX, offset int) {
	if box.X < startX && (box.X+box.Width) > startX {
		box.Width += offset
	} else if box.X >= startX {
		box.X += offset
	}
	doc.moveBoxContHorizontal(box.Horizontal, startX, offset)
	doc.moveBoxContHorizontal(box.Vertical, startX, offset)
}

func (doc *BoxesDocument) StretchAndMoveHorizontal(startX, offset int) {
	doc.moveBoxHorizontal(&doc.Boxes, startX, offset)
}

func (doc *BoxesDocument) moveBoxContVertical(cont *LayoutElemContainer, startY, offset int) {
	if cont == nil {
		return
	}
	if cont.Y < startY && (cont.Y+cont.Height) > startY {
		cont.Height += offset
	} else if cont.Y >= startY {
		cont.Y += offset
	}
	for i := range len(cont.Elems) {
		doc.moveBoxVertical(&cont.Elems[i], startY, offset)
	}
}

func (doc *BoxesDocument) moveBoxVertical(box *LayoutElement, startY, offset int) {
	if box.Y < startY && (box.Y+box.Height) > startY {
		box.Height += offset
	} else if box.Y > startY {
		box.Y += offset
	}
	doc.moveBoxContVertical(box.Horizontal, startY, offset)
	doc.moveBoxContVertical(box.Vertical, startY, offset)
}

func (doc *BoxesDocument) StretchAndMoveVertical(startY, offset int) {
	doc.moveBoxVertical(&doc.Boxes, startY, offset)
}
