package types

type TextDimensionCalculator interface {
	CaptionDimensions(txt string) (width, height int)
	Text1Dimensions(txt string) (width, height int)
	Text2Dimensions(txt string) (width, height int)
}

func (l *LayoutElement) initVertical(c TextDimensionCalculator, yInnerOffset, defaultPadding, defaultBoxMargin int) {
	if len(l.Vertical) > 0 {
		curX := l.X
		curY := l.Y + yInnerOffset
		var h, w int
		for i := 0; i < len(l.Vertical); i++ {
			sub := &l.Vertical[i]
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
		}
		l.Height += h + defaultPadding
		if w > l.Width {
			l.Width = w
		}
	}
}

func (l *LayoutElement) initHorizontal(c TextDimensionCalculator, yInnerOffset, defaultPadding, defaultBoxMargin int) {
	if len(l.Horizontal) > 0 {
		curX := l.X + defaultPadding
		curY := l.Y + yInnerOffset
		var h, w int
		for i := 0; i < len(l.Horizontal); i++ {
			sub := &l.Horizontal[i]
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
		}
		l.Height += h + defaultPadding
		if l.Width < w {
			l.Width = w
		}
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
