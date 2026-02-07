package types

func InitFontDef(l *FontDef, defaultFont string, defaultSize int, defaultBold, defaultItalic bool, spaceTop int) FontDef {
	var f FontDef
	typeNormal := FontDefTypeEnum_normal
	typeItalic := FontDefTypeEnum_italic
	weightNormal := FontDefWeightEnum_normal
	weightBold := FontDefWeightEnum_bold
	alignedCenter := FontDefAlignedEnum_center

	if l != nil {
		if l.Font != "" {
			f.Font = l.Font
		} else {
			f.Font = defaultFont
		}
		if l.Size != 0 {
			f.Size = l.Size
		} else {
			f.Size = defaultSize
		}
		if l.Type != nil {
			f.Type = l.Type
		} else {
			if defaultItalic {
				f.Type = &typeItalic
			} else {
				f.Type = &typeNormal
			}
		}
		if l.Weight != nil {
			f.Weight = l.Weight
		} else {
			if defaultBold {
				f.Weight = &weightBold
			} else {
				f.Weight = &weightNormal
			}
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
		if l.Anchor != "" {
			f.Anchor = l.Anchor
		} else {
			f.Anchor = FontDefAnchorEnum_middle
		}
		if l.Aligned != nil {
			f.Aligned = l.Aligned
		} else {
			f.Aligned = &alignedCenter
		}
		f.SpaceTop = l.SpaceTop
		if f.SpaceTop == 0 {
			f.SpaceTop = spaceTop
		}
		f.SpaceBottom = l.SpaceBottom
		if l.MaxLenBeforeBreak != 0 {
			f.MaxLenBeforeBreak = l.MaxLenBeforeBreak
		} else {
			f.MaxLenBeforeBreak = 90
		}
	} else {
		f.Size = defaultSize
		if defaultItalic {
			f.Type = &typeItalic
		} else {
			f.Type = &typeNormal
		}
		f.Font = defaultFont
		if defaultBold {
			f.Weight = &weightBold
		} else {
			f.Weight = &weightNormal
		}
		f.LineHeight = 1.5
		f.Color = "black"
		f.Aligned = &alignedCenter
		f.Anchor = FontDefAnchorEnum_middle
		f.SpaceTop = spaceTop
		f.SpaceBottom = 0
		f.MaxLenBeforeBreak = 90
	}
	return f
}

func InitLineDef(l *LineDef) *LineDef {
	wDef := 1.0
	sDef := LineDefStyleEnum_solid
	cDef := "black"
	oDef := 1.0

	w := &wDef
	o := &oDef
	s := &sDef
	c := &cDef

	if l != nil {
		if l.Width != nil {
			w = l.Width
		}
		if l.Style != nil {
			s = l.Style
		}
		if l.Opacity != nil {
			o = l.Opacity
		}
		if l.Color != nil {
			c = l.Color
		}
	}

	return &LineDef{
		Width:   w,
		Style:   s,
		Opacity: o,
		Color:   c,
	}

}

func InitFillDef(l *FillDef) *FillDef {
	cDef := "#efefef"
	oDef := 1.0

	o := &oDef
	c := &cDef

	if l != nil {
		if l.Opacity != nil {
			o = l.Opacity
		}
		if l.Color != nil {
			c = l.Color
		}
	}

	return &FillDef{
		Opacity: o,
		Color:   c,
	}

}
