package types

type TextDimensionCalculator interface {
	CaptionDimensions(txt string) (width, height int32)
	Text1Dimensions(txt string) (width, height int32)
	Text2Dimensions(txt string) (width, height int32)
}

func (l *LayoutElement) InitDimensions(c TextDimensionCalculator) {
	// TODO
	// l.Width, l.Height = td(l.Caption)
	// for _, v := range l.Vertical {
	// 	v.InitDimensions(td)
	// }
	// for _, h := range l.Horizontal {
	// 	h.InitDimensions(td)
	// }
}
