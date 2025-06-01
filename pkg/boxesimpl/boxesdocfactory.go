package boxesimpl

import (
	"strings"

	"github.com/okieoth/draw.chart.things/pkg/types"
	"github.com/okieoth/draw.chart.things/pkg/types/boxes"
)

func initLayoutElemContainer(l []boxes.Layout, inputFormats map[string]boxes.BoxFormat) *boxes.LayoutElemContainer {
	if len(l) == 0 {
		return nil
	}
	var ret boxes.LayoutElemContainer
	ret.Elems = make([]boxes.LayoutElement, 0)
	for _, elem := range l {
		ret.Elems = append(ret.Elems, initLayoutElement(&elem, inputFormats))
	}
	return &ret
}

func initConnections(l []boxes.Connection) []boxes.LayoutElemConnection {
	ret := make([]boxes.LayoutElemConnection, 0)
	for _, elem := range l {
		var conn boxes.LayoutElemConnection
		conn.DestId = elem.DestId
		conn.SourceArrow = elem.SourceArrow
		conn.DestArrow = elem.DestArrow
		conn.Tags = elem.Tags
		ret = append(ret, conn)
	}
	return ret
}

func InitFontDef(l *types.FontDef, defaultFont string, defaultSize int, defaultBold, defaultItalic bool, spaceTop int) types.FontDef {
	var f types.FontDef
	typeNormal := types.FontDefTypeEnum_normal
	typeItalic := types.FontDefTypeEnum_italic
	weightNormal := types.FontDefWeightEnum_normal
	weightBold := types.FontDefWeightEnum_bold
	alignedLeft := types.FontDefAlignedEnum_left

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

func initBoxFormat(f *boxes.Format) boxes.BoxFormat {
	var border *types.LineDef
	var fill *types.FillDef

	var fontCaption *types.FontDef
	var fontText1 *types.FontDef
	var fontText2 *types.FontDef
	var verticalTxt bool
	padding := types.GlobalPadding
	boxMargin := types.GlobalMinBoxMargin
	if f != nil {
		fontCaption = f.FontCaption
		fontText1 = f.FontText1
		fontText2 = f.FontText2
		border = f.Border
		fill = f.Fill
		if f.Padding != nil {
			padding = *f.Padding
		}
		if f.BoxMargin != nil {
			boxMargin = *f.BoxMargin
		}
		if f.VerticalTxt != nil {
			verticalTxt = *f.VerticalTxt
		}
	}

	return boxes.BoxFormat{
		Padding:      padding,
		MinBoxMargin: boxMargin,
		FontCaption:  InitFontDef(fontCaption, "sans-serif", 10, true, false, 0),
		FontText1:    InitFontDef(fontText1, "serif", 8, false, false, 10),
		FontText2:    InitFontDef(fontText2, "monospace", 8, false, true, 10),
		Border:       border,
		Fill:         fill,
		VerticalTxt:  verticalTxt,
	}
}

func getDefaultFormat() boxes.BoxFormat {
	return boxes.BoxFormat{
		Padding:      types.GlobalPadding,
		MinBoxMargin: types.GlobalMinBoxMargin,
		FontCaption:  InitFontDef(nil, "sans-serif", 10, true, false, 0),
		FontText1:    InitFontDef(nil, "serif", 8, false, false, 10),
		FontText2:    InitFontDef(nil, "monospace", 8, false, true, 10),
	}
}

func initFormats(inputFormat map[string]boxes.Format) map[string]boxes.BoxFormat {
	var res = make(map[string]boxes.BoxFormat)
	for key, elem := range inputFormat {
		res[key] = initBoxFormat(&elem)
	}
	if _, hasDefault := res["default"]; !hasDefault {
		res["default"] = getDefaultFormat()
	}
	return res
}

func initLayoutElement(l *boxes.Layout, inputFormats map[string]boxes.BoxFormat) boxes.LayoutElement {
	var f *boxes.BoxFormat
	for _, tag := range l.Tags {
		if val, ok := inputFormats[tag]; ok {
			f = &val
			break
		}
	}
	if (f == nil) && (l.Caption != "" || l.Text1 != "" || l.Text2 != "") {
		formatKey := "default"
		if l.Format != nil {
			formatKey = *l.Format
		}
		for key, format := range inputFormats {
			if key == formatKey {
				f = &format
				break
			}
		}
		if f == nil {
			d := getDefaultFormat()
			f = &d
		}
	}
	if l.Id == "" {
		if l.Caption != "" {
			l.Id = strings.ToLower(l.Caption)
		}
	}
	return boxes.LayoutElement{
		Id:          l.Id,
		Caption:     l.Caption,
		Text1:       l.Text1,
		Text2:       l.Text2,
		Vertical:    initLayoutElemContainer(l.Vertical, inputFormats),
		Horizontal:  initLayoutElemContainer(l.Horizontal, inputFormats),
		Format:      f,
		Connections: initConnections(l.Connections),
	}
}

func DocumentFromBoxes(b *boxes.Boxes) *boxes.BoxesDocument {
	doc := boxes.NewBoxesDocument()
	doc.Title = b.Title
	doc.Formats = initFormats(b.Formats)
	doc.Boxes = initLayoutElement(&b.Boxes, doc.Formats)
	if doc.MinBoxMargin == 0 {
		doc.MinBoxMargin = types.GlobalMinBoxMargin
	}
	if doc.MinConnectorMargin == 0 {
		doc.MinConnectorMargin = types.GlobalMinConnectorMargin
	}
	if doc.GlobalPadding == 0 {
		doc.GlobalPadding = types.GlobalPadding
	}
	return doc
}
