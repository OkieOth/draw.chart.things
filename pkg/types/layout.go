package types

type TextDimensionCalculator interface {
	Dimensions(txt string, format *FontDef) (width, height int)
}

func (l *LayoutElement) incrementX(xOffset int) {
	l.X += xOffset
	if l.Vertical != nil {
		l.Vertical.X += xOffset
		for i := 0; i < len(l.Vertical.Elems); i++ {
			sub := &l.Vertical.Elems[i]
			sub.incrementX(xOffset)
		}
	}
	if l.Horizontal != nil {
		l.Horizontal.X += xOffset
		for i := 0; i < len(l.Horizontal.Elems); i++ {
			sub := &l.Horizontal.Elems[i]
			sub.incrementX(xOffset)
		}
	}
}

func (l *LayoutElement) incrementY(yOffset int) {
	l.Y += yOffset
	if l.Vertical != nil {
		l.Vertical.Y += yOffset
		for i := 0; i < len(l.Vertical.Elems); i++ {
			sub := &l.Vertical.Elems[i]
			sub.incrementY(yOffset)
		}
	}
	if l.Horizontal != nil {
		l.Horizontal.Y += yOffset
		for i := 0; i < len(l.Horizontal.Elems); i++ {
			sub := &l.Horizontal.Elems[i]
			sub.incrementY(yOffset)
		}
	}
}

func (l *LayoutElement) centerHorizontalElems() {
	if l.Horizontal != nil {
		cont := l.Horizontal
		contYOffset := (l.Height - cont.Height) / 2
		contXOffset := (l.Width - cont.Width) / 2
		cont.Y = l.Y + contYOffset
		cont.X = l.X + contXOffset
		for i := 0; i < len(cont.Elems); i++ {
			sub := &cont.Elems[i]
			offsetY := ((cont.Height - sub.Height) / 2)
			sub.incrementY(offsetY)
			sub.incrementX(contXOffset)
			//sub.centerHorizontalElems()
			sub.Center()
		}
	}
}

func (l *LayoutElement) centerVerticalElems() {
	if l.Vertical != nil {
		cont := l.Vertical
		contYOffset := (l.Y + l.Height - cont.Y - cont.Height) / 2
		contXOffset := (l.Width - cont.Width) / 2
		cont.Y = l.Y + contYOffset
		cont.X = l.X + contXOffset
		for i := 0; i < len(cont.Elems); i++ {
			sub := &cont.Elems[i]
			sub.incrementY(contYOffset)
			offsetX := contXOffset + ((cont.Width - sub.Width) / 2)
			sub.incrementX(offsetX)
			//sub.centerVerticalElems()
			sub.Center()
		}
	}
}

func (l *LayoutElement) AreOnTheSameVerticalLevel(otherElem *LayoutElement) bool {
	return ((l.CenterY > otherElem.Y) && (l.CenterY < otherElem.Y+otherElem.Height)) ||
		((otherElem.CenterY > l.Y) && (otherElem.CenterY < l.Y+l.Height))
}

func (l *LayoutElement) initVertical(c TextDimensionCalculator, yInnerOffset, defaultPadding, defaultBoxMargin int) {
	if l.Vertical != nil && len(l.Vertical.Elems) > 0 {
		curX := l.X
		l.Vertical.X = curX
		curY := l.Y + yInnerOffset
		l.Vertical.Y = curY
		var h, w int
		var hasChilds bool
		for i := 0; i < len(l.Vertical.Elems); i++ {
			sub := &l.Vertical.Elems[i]
			if (sub.Horizontal != nil && len(sub.Horizontal.Elems) > 0) || (sub.Vertical != nil && len(sub.Vertical.Elems) > 0) {
				hasChilds = true
			}
			if h > 0 {
				h += defaultBoxMargin
			}
			sub.X = curX
			sub.Y = curY
			sub.InitDimensions(c, defaultPadding, defaultBoxMargin)
			curY += (sub.Height + defaultBoxMargin)
			h += sub.Height
			if sub.Width > w {
				w = sub.Width
			}
			if sub.Width > l.Vertical.Width {
				l.Vertical.Width = sub.Width
			}
		}
		if !hasChilds {
			for i := 0; i < len(l.Vertical.Elems); i++ {
				sub := &l.Vertical.Elems[i]
				sub.Width = w
			}
		}

		l.Vertical.Height = h + defaultPadding
		l.Height += l.Vertical.Height
		if w > l.Width {
			l.Width = w
		}
	}
}

func (l *LayoutElement) initHorizontal(c TextDimensionCalculator, yInnerOffset, defaultPadding, defaultBoxMargin int) {
	if l.Horizontal != nil && len(l.Horizontal.Elems) > 0 {
		curX := l.X
		l.Horizontal.X = curX
		curY := l.Y + yInnerOffset
		l.Horizontal.Y = curY
		var h, w int
		var hasChilds bool
		for i := 0; i < len(l.Horizontal.Elems); i++ {
			sub := &l.Horizontal.Elems[i]
			if (sub.Horizontal != nil && len(sub.Horizontal.Elems) > 0) || (sub.Vertical != nil && len(sub.Vertical.Elems) > 0) {
				hasChilds = true
			}
			if w > 0 {
				w += defaultBoxMargin
			}
			sub.X = curX
			sub.Y = curY
			sub.InitDimensions(c, defaultPadding, defaultBoxMargin)
			curX += (sub.Width + defaultBoxMargin)
			w += sub.Width
			if sub.Height > h {
				h = sub.Height
			}
			if sub.Height > l.Horizontal.Height {
				l.Horizontal.Height = sub.Height
			}
		}
		if !hasChilds {
			for i := 0; i < len(l.Horizontal.Elems); i++ {
				sub := &l.Horizontal.Elems[i]
				sub.Height = h
			}
		}

		l.Horizontal.Height = h + defaultPadding
		l.Height += l.Horizontal.Height
		l.Horizontal.Width = w

		if l.Width < w {
			l.Width = w
		}
	}
}

func (l *LayoutElement) InitDimensions(c TextDimensionCalculator, defaultPadding, defaultBoxMargin int) {
	var cW, cH, t1W, t1H, t2W, t2H int
	//var yCaptionOffset, yText1Offset, yText2Offset, yInnerOffset int
	var yInnerOffset int
	if l.Caption != "" || l.Text1 != "" || l.Text2 != "" {
		l.Height = (2 * l.Format.Padding)
		if l.Caption != "" {
			cW, cH = c.Dimensions(l.Caption, &l.Format.FontCaption)
			l.Height += cH + l.Format.FontCaption.SpaceTop + l.Format.FontCaption.SpaceBottom
		}
		if l.Text1 != "" {
			p := l.Format.Padding
			if l.Format.FontText1.SpaceTop > 0 {
				p = l.Format.FontText1.SpaceTop
			}
			t1W, t1H = c.Dimensions(l.Text1, &l.Format.FontText1)
			l.Height += t1H + p + l.Format.FontText1.SpaceBottom
		}
		if l.Text2 != "" {
			p := l.Format.Padding
			if l.Format.FontText2.SpaceTop > 0 {
				p = l.Format.FontText2.SpaceTop
			}
			t2W, t2H = c.Dimensions(l.Text2, &l.Format.FontText2)
			l.Height += t2H + p + l.Format.FontText2.SpaceBottom
		}
		//yInnerOffset = l.Format.Padding + max(yCaptionOffset, max(yText1Offset, yText2Offset))
		//yInnerOffset = l.Format.Padding + l.Height
		yInnerOffset = l.Height
		l.Width = max(cW, max(t1W, t2W)) + (2 * l.Format.Padding)
	}
	l.initVertical(c, yInnerOffset, defaultPadding, defaultBoxMargin)
	l.initHorizontal(c, yInnerOffset, defaultPadding, defaultBoxMargin)
}

func (l *LayoutElement) Center() {
	l.CenterX = l.X + (l.Width / 2)
	l.CenterY = l.Y + (l.Height / 2)
	l.centerHorizontalElems()
	l.centerVerticalElems()
}
