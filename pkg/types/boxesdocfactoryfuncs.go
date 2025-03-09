package types

type TextDimensionCalculator interface {
	CaptionDimensions(txt string) (width, height int32)
	Text1Dimensions(txt string) (width, height int32)
	Text2Dimensions(txt string) (width, height int32)
}

func (l *LayoutElement) InitDimensions(c TextDimensionCalculator) {
	initVertical := func() {
		if len(l.Vertical) > 0 {
			var h, w int32
			for _, sub := range l.Vertical {
				if h > 0 {
					h += l.Format.MinBoxMargin
				}
				sub.InitDimensions(c)
				h += sub.Height
				if sub.Width > w {
					w = sub.Width
				}
			}
			l.Height += h
			if w > l.Width {
				l.Width = w
			}
		}
	}
	initHorizontal := func() {
		if len(l.Horizontal) > 0 {
			var h, w int32
			for _, sub := range l.Horizontal {
				if w > 0 {
					w += l.Format.MinBoxMargin
				}
				sub.InitDimensions(c)
				w += sub.Width
				if sub.Height > h {
					h = sub.Height
				}
			}
			l.Height += h
			if l.Width < w {
				l.Width = w
			}
		}
	}

	var cW, cH, t1W, t1H, t2W, t2H int32
	l.Height = (2 * l.Format.Padding)
	if l.Caption != "" {
		cW, cH = c.CaptionDimensions(l.Caption)
		l.Height += cH + l.Format.FontCaption.SpaceTop + l.Format.FontCaption.SpaceBottom
	}
	if l.Text1 != "" {
		t1W, t1H = c.Text1Dimensions(l.Text1)
		l.Height += t1H + l.Format.Padding + l.Format.FontText1.SpaceTop + l.Format.FontText1.SpaceBottom
	}
	if l.Text2 != "" {
		t2W, t2H = c.Text2Dimensions(l.Text2)
		l.Height += t2H + l.Format.Padding + l.Format.FontText2.SpaceTop + l.Format.FontText2.SpaceBottom
	}
	l.Width = max(cW, max(t1W, t2W)) + (2 * l.Format.Padding)

	initVertical()
	initHorizontal()
}
