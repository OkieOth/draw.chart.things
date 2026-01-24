package boxesimpl

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/okieoth/draw.chart.things/pkg/types"
	"github.com/okieoth/draw.chart.things/pkg/types/boxes"
)

func initLayoutElemContainer(l []boxes.Layout, inputFormats map[string]boxes.BoxFormat, connectedElems *[]string, hideTextsForParents bool, definedImages map[string]types.ImageDef) *boxes.LayoutElemContainer {
	if len(l) == 0 {
		return nil
	}
	var ret boxes.LayoutElemContainer
	ret.Elems = make([]boxes.LayoutElement, 0)
	for _, elem := range l {
		ret.Elems = append(ret.Elems, initLayoutElement(&elem, inputFormats, connectedElems, hideTextsForParents, definedImages))
	}
	return &ret
}

func initConnections(l []boxes.Connection, inputFormats map[string]boxes.BoxFormat) []boxes.LayoutElemConnection {
	ret := make([]boxes.LayoutElemConnection, 0)
	defaultConnectionFormat, ok := inputFormats["defaultConnection"]
	for _, elem := range l {
		var conn boxes.LayoutElemConnection
		conn.DestId = elem.DestId
		conn.SourceArrow = elem.SourceArrow
		conn.DestArrow = elem.DestArrow
		conn.Tags = elem.Tags
		if elem.Format != nil {
			if formatInst, ok := inputFormats[*elem.Format]; ok {
				conn.Format = formatInst.Line
			}
		} else if ok {
			conn.Format = defaultConnectionFormat.Line
		}
		ret = append(ret, conn)
	}
	return ret
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
	var fixedHeight, fixedWidth, cornerRadius *int
	if f != nil {
		fontCaption = f.FontCaption
		fontText1 = f.FontText1
		fontText2 = f.FontText2
		border = f.Line
		fill = f.Fill
		if f.Padding != nil {
			padding = *f.Padding
		} else {
			f.Padding = &padding
		}
		if f.BoxMargin != nil {
			boxMargin = *f.BoxMargin
		} else {
			f.BoxMargin = &boxMargin
		}
		if f.VerticalTxt != nil {
			verticalTxt = *f.VerticalTxt
		}
		fixedHeight = f.FixedHeight
		fixedWidth = f.FixedWidth
		cornerRadius = f.CornerRadius
	}

	return boxes.BoxFormat{
		Padding:      padding,
		MinBoxMargin: boxMargin,
		FontCaption:  types.InitFontDef(fontCaption, "sans-serif", 10, true, false, 0),
		FontText1:    types.InitFontDef(fontText1, "serif", 8, false, false, 10),
		FontText2:    types.InitFontDef(fontText2, "monospace", 8, false, true, 10),
		Line:         border,
		Fill:         fill,
		FixedWidth:   fixedWidth,
		FixedHeight:  fixedHeight,
		VerticalTxt:  verticalTxt,
		CornerRadius: cornerRadius,
	}
}

func getDefaultFormat() boxes.BoxFormat {
	return boxes.BoxFormat{
		Padding:      types.GlobalPadding,
		MinBoxMargin: types.GlobalMinBoxMargin,
		FontCaption:  types.InitFontDef(nil, "sans-serif", 10, true, false, 0),
		FontText1:    types.InitFontDef(nil, "serif", 8, false, false, 10),
		FontText2:    types.InitFontDef(nil, "monospace", 8, false, true, 10),
		Line:         types.InitLineDef(nil),
		Fill:         nil,
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
	if _, hasDefault := res["defaultConnection"]; !hasDefault {
		res["defaultConnection"] = getDefaultFormat()
	}
	return res
}

func initImage(l *boxes.Layout, definedImages map[string]types.ImageDef) *boxes.ImageContainer {
	if l.Image == nil {
		return nil
	}
	imageDef, ok := definedImages[*l.Image]
	if !ok {
		return nil
	}
	marginTopBottom := 5
	if imageDef.MarginTopBottom != nil {
		marginTopBottom = *imageDef.MarginTopBottom
	}
	marginLeftRight := 2
	if imageDef.MarginLeftRight != nil {
		marginLeftRight = *imageDef.MarginLeftRight
	}
	img := boxes.ImageContainer{
		Width:           imageDef.Width,
		Height:          imageDef.Height,
		ImgId:           *l.Image,
		MarginTopBottom: marginTopBottom,
		MarginLeftRight: marginLeftRight,
	}
	return &img
}

func initLayoutElement(l *boxes.Layout, inputFormats map[string]boxes.BoxFormat, connectedIds *[]string, hideTextsForParents bool, definedImages map[string]types.ImageDef) boxes.LayoutElement {
	var f *boxes.BoxFormat
	// for _, tag := range l.Tags {
	// 	if val, ok := inputFormats[tag]; ok {
	// 		f = &val
	// 		break
	// 	}
	// }
	if l.Format != nil {
		if val, ok := inputFormats[*l.Format]; ok {
			f = &val
		}
	}
	if len(l.Connections) > 0 && l.Id != "" {
		if !slices.Contains(*connectedIds, l.Id) {
			*connectedIds = append(*connectedIds, l.Id)
		}
		for _, c := range l.Connections {
			if !slices.Contains(*connectedIds, c.DestId) {
				*connectedIds = append(*connectedIds, c.DestId)
			}
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
	var text1, text2 string
	if !hideTextsForParents || (len(l.Vertical) == 0 && len(l.Horizontal) == 0) {
		text1 = l.Text1
		text2 = l.Text2
	}
	return boxes.LayoutElement{
		Id:          l.Id,
		Caption:     l.Caption,
		Text1:       text1,
		Text2:       text2,
		Image:       initImage(l, definedImages), // TODO
		Vertical:    initLayoutElemContainer(l.Vertical, inputFormats, connectedIds, hideTextsForParents, definedImages),
		Horizontal:  initLayoutElemContainer(l.Horizontal, inputFormats, connectedIds, hideTextsForParents, definedImages),
		Format:      f,
		Connections: initConnections(l.Connections, inputFormats),
	}
}

func initExternalImages(doc *boxes.BoxesDocument) error {
	for key, img := range doc.Images {
		if img.Base64 == nil && img.Base64Src == nil {
			return fmt.Errorf("Missing 'base64' or 'base64Src' attribute for imageDef id=%s", key)
		}
		if img.Base64Src != nil {
			// load the base64 string from the file given by the attrib
			bytes, err := os.ReadFile(*img.Base64Src)
			if err != nil {
				return fmt.Errorf("Error while reading content of 'base64Src' (%s) for imageDef id=%s", *img.Base64Src, key)
			}
			base64Str := string(bytes)
			img.Base64 = &base64Str
			doc.Images[key] = img
		}
	}
	return nil
}

func DocumentFromBoxes(b *boxes.Boxes) (*boxes.BoxesDocument, error) {
	doc := boxes.NewBoxesDocument()
	doc.Title = b.Title
	doc.Formats = initFormats(b.Formats)
	doc.Images = b.Images
	if b.MinBoxMargin != nil {
		doc.MinBoxMargin = *b.MinBoxMargin
	}
	if b.MinConnectorMargin != nil {
		doc.MinConnectorMargin = *b.MinConnectorMargin
	}
	if b.GlobalPadding != nil {
		doc.GlobalPadding = *b.GlobalPadding
	}
	if b.LineDist != nil {
		doc.LineDist = *b.LineDist
	}
	if err := initExternalImages(doc); err != nil {
		return nil, err
	}

	doc.Boxes = initLayoutElement(&b.Boxes, doc.Formats, &doc.ConnectedElems, b.HideTextsForParents, b.Images)
	if doc.MinBoxMargin == 0 {
		doc.MinBoxMargin = types.GlobalMinBoxMargin
	}
	if doc.MinConnectorMargin == 0 {
		doc.MinConnectorMargin = types.GlobalMinConnectorMargin
	}
	if doc.GlobalPadding == 0 {
		doc.GlobalPadding = types.GlobalPadding
	}
	if doc.LineDist == 0 {
		doc.LineDist = types.LineDist
	}
	return doc, nil
}
