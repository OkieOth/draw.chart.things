package types

func InitFontDef(l *FontDef, defaultFont string, defaultSize int, defaultBold, defaultItalic bool, spaceTop int) FontDef {
	var f FontDef
	typeNormal := FontDefTypeEnum_normal
	typeItalic := FontDefTypeEnum_italic
	weightNormal := FontDefWeightEnum_normal
	weightBold := FontDefWeightEnum_bold
	alignedLeft := FontDefAlignedEnum_left

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
		if l.Aligned != nil {
			f.Aligned = l.Aligned
		} else {
			f.Aligned = &alignedLeft
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
		f.Aligned = &alignedLeft
		f.SpaceTop = spaceTop
		f.SpaceBottom = 0
		f.MaxLenBeforeBreak = 90
	}
	return f
}
