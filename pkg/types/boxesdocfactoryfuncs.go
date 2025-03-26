package types

type TextDimensionCalculator interface {
	CaptionDimensions(txt string) (width, height int)
	Text1Dimensions(txt string) (width, height int)
	Text2Dimensions(txt string) (width, height int)
}

func (l *LayoutElement) centerHorizontal() {
	if l.Horizontal != nil {
		for i := 0; i < len(l.Horizontal.Elems); i++ {
			sub := &l.Horizontal.Elems[i]
			sub.Y = l.Y + ((l.Horizontal.Height - sub.Height) / 2)
			sub.centerHorizontal()
		}
	}
}

func (l *LayoutElement) centerVertical() {
	if l.Vertical != nil {
		for i := 0; i < len(l.Vertical.Elems); i++ {
			sub := &l.Vertical.Elems[i]
			sub.X = l.X + ((l.Vertical.Width - sub.Width) / 2)
			sub.centerVertical()
		}
	}
}

func (l *LayoutElement) initVertical(c TextDimensionCalculator, yInnerOffset, defaultPadding, defaultBoxMargin int) {
	if l.Vertical != nil && len(l.Vertical.Elems) > 0 {
		curX := l.X
		curY := l.Y + yInnerOffset
		var h, w int
		for i := 0; i < len(l.Vertical.Elems); i++ {
			sub := &l.Vertical.Elems[i]
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
		l.Vertical.Height = h + defaultPadding
		l.Height += l.Vertical.Height
		if w > l.Width {
			l.Width = w
		}
		l.centerVertical()
	}
}

func (l *LayoutElement) initHorizontal(c TextDimensionCalculator, yInnerOffset, defaultPadding, defaultBoxMargin int) {
	if l.Horizontal != nil && len(l.Horizontal.Elems) > 0 {
		curX := l.X
		curY := l.Y + yInnerOffset
		var h, w int
		for i := 0; i < len(l.Horizontal.Elems); i++ {
			sub := &l.Horizontal.Elems[i]
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
		l.Height += h
		l.Horizontal.Width = w

		if l.Width < w {
			l.Width = w
		}
		l.centerHorizontal()
	}
}

func (l *LayoutElement) InitDimensions(c TextDimensionCalculator, defaultPadding, defaultBoxMargin int) {
	var cW, cH, t1W, t1H, t2W, t2H int
	var yCaptionOffset, yText1Offset, yText2Offset, yInnerOffset int
	if l.Caption != "" || l.Text1 != "" || l.Text2 != "" {
		l.Height = (2 * l.Format.Padding)
		if l.Caption != "" {
			yCaptionOffset = l.Format.FontCaption.SpaceTop + l.Format.Padding
			cW, cH = c.CaptionDimensions(l.Caption)
			l.Height += cH + l.Format.FontCaption.SpaceTop + l.Format.FontCaption.SpaceBottom
		}
		if l.Text1 != "" {
			yText1Offset = yCaptionOffset + l.Format.Padding + l.Format.FontText1.SpaceTop
			t1W, t1H = c.Text1Dimensions(l.Text1)
			l.Height += t1H + l.Format.Padding + l.Format.FontText1.SpaceTop + l.Format.FontText1.SpaceBottom
		}
		if l.Text2 != "" {
			yText2Offset = yText1Offset + l.Format.Padding + l.Format.FontText2.SpaceTop
			t2W, t2H = c.Text2Dimensions(l.Text2)
			l.Height += t2H + l.Format.Padding + l.Format.FontText2.SpaceTop + l.Format.FontText2.SpaceBottom
		}
		yInnerOffset = l.Format.Padding + max(yCaptionOffset, max(yText1Offset, yText2Offset))
		l.Width = max(cW, max(t1W, t2W)) + (2 * l.Format.Padding)
	}
	l.initVertical(c, yInnerOffset, defaultPadding, defaultBoxMargin)
	l.initHorizontal(c, yInnerOffset, defaultPadding, defaultBoxMargin)
}
