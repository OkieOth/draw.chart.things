package boxes

// Attention, this file is generated. Manual changes get lost with the next
// run of the code generation.
// created by yacg (template: golang_types.mako v1.1.0)

import (
    "github.com/okieoth/draw.chart.things/pkg/types"
)


/* Model to describe the input of block diagrams
*/
type Boxes struct {

    // Title of the document
    Title string  `yaml:"title"`

    Boxes Layout  `yaml:"boxes"`

    // Map of formats available to be used in the boxes
    Formats map[string]Format  `yaml:"formats,omitempty"`

    // optional list of images used in the generated graphic
    Images []types.ImageDef  `yaml:"images,omitempty"`

    // Padding used as default over the whole diagram
    GlobalPadding *int  `yaml:"globalPadding,omitempty"`

    // Minimum margin between boxes
    MinBoxMargin *int  `yaml:"minBoxMargin,omitempty"`

    // Minimum margin between connectors
    MinConnectorMargin *int  `yaml:"minConnectorMargin,omitempty"`
}

func NewBoxes() *Boxes {
    return &Boxes{
        Boxes: *NewLayout(),
        Formats: make(map[string]Format, 0),
        Images: make([]types.ImageDef, 0),
    }
}

func CopyBoxes(src *Boxes) *Boxes {
    if src == nil {
        return nil
    }
    var ret Boxes
    ret.Title = src.Title
    ret.Boxes = *CopyLayout(&src.Boxes)
    ret.Formats = make(map[string]Format, 0)
    for k, v := range src.Formats {
        ret.Formats[k] = v
    }
    ret.Images = make([]types.ImageDef, 0)
    for _, e := range src.Images {
        ret.Images = append(ret.Images, e)
    }
    ret.GlobalPadding = src.GlobalPadding
    ret.MinBoxMargin = src.MinBoxMargin
    ret.MinConnectorMargin = src.MinConnectorMargin

    return &ret
}





type Layout struct {

    // unique identifier of that entry
    Id string  `yaml:"id"`

    // Some kind of the main text
    Caption string  `yaml:"caption"`

    // First additional text
    Text1 string  `yaml:"text1"`

    // Second additional text
    Text2 string  `yaml:"text2"`

    // If set, then the content for 'vertical' attrib is loaded from an external file
    ExtVertical *string  `yaml:"extVertical,omitempty"`

    Vertical []Layout  `yaml:"vertical,omitempty"`

    // If set, then the content for 'horizontal' attrib is loaded from an external file
    ExtHorizontal *string  `yaml:"extHorizontal,omitempty"`

    Horizontal []Layout  `yaml:"horizontal,omitempty"`

    // Tags to annotate the box, tags are used to format and filter
    Tags []string  `yaml:"tags,omitempty"`

    // List of connections to other boxes
    Connections []Connection  `yaml:"connections,omitempty"`

    // reference to the format to use for this box
    Format *string  `yaml:"format,omitempty"`
}

func NewLayout() *Layout {
    return &Layout{
        Vertical: make([]Layout, 0),
        Horizontal: make([]Layout, 0),
        Tags: make([]string, 0),
        Connections: make([]Connection, 0),
    }
}

func CopyLayout(src *Layout) *Layout {
    if src == nil {
        return nil
    }
    var ret Layout
    ret.Id = src.Id
    ret.Caption = src.Caption
    ret.Text1 = src.Text1
    ret.Text2 = src.Text2
    ret.ExtVertical = src.ExtVertical
    ret.Vertical = make([]Layout, 0)
    for _, e := range src.Vertical {
        ret.Vertical = append(ret.Vertical, e)
    }
    ret.ExtHorizontal = src.ExtHorizontal
    ret.Horizontal = make([]Layout, 0)
    for _, e := range src.Horizontal {
        ret.Horizontal = append(ret.Horizontal, e)
    }
    ret.Tags = make([]string, 0)
    for _, e := range src.Tags {
        ret.Tags = append(ret.Tags, e)
    }
    ret.Connections = make([]Connection, 0)
    for _, e := range src.Connections {
        ret.Connections = append(ret.Connections, e)
    }
    ret.Format = src.Format

    return &ret
}








type Format struct {

    // optional fixed width that will be applied on the box
    FixedWidth *int  `yaml:"fixedWidth,omitempty"`

    // optional fixed height that will be applied on the box
    FixedHeight *int  `yaml:"fixedHeight,omitempty"`

    // If true, the text will be displayed vertically
    VerticalTxt *bool  `yaml:"verticalTxt,omitempty"`

    FontCaption *types.FontDef  `yaml:"fontCaption,omitempty"`

    FontText1 *types.FontDef  `yaml:"fontText1,omitempty"`

    FontText2 *types.FontDef  `yaml:"fontText2,omitempty"`

    Line *types.LineDef  `yaml:"line,omitempty"`

    Fill *types.FillDef  `yaml:"fill,omitempty"`

    // Padding used for this format
    Padding *int  `yaml:"padding,omitempty"`

    // Minimum margin between boxes
    BoxMargin *int  `yaml:"boxMargin,omitempty"`

    // radius of the box corners in pixel
    CornerRadius *int  `yaml:"cornerRadius,omitempty"`

    // ID of an image to use
    Image *string  `yaml:"image,omitempty"`
}


func CopyFormat(src *Format) *Format {
    if src == nil {
        return nil
    }
    var ret Format
    ret.FixedWidth = src.FixedWidth
    ret.FixedHeight = src.FixedHeight
    ret.VerticalTxt = src.VerticalTxt
    ret.FontCaption = types.CopyFontDef(src.FontCaption)
    ret.FontText1 = types.CopyFontDef(src.FontText1)
    ret.FontText2 = types.CopyFontDef(src.FontText2)
    ret.Line = types.CopyLineDef(src.Line)
    ret.Fill = types.CopyFillDef(src.Fill)
    ret.Padding = src.Padding
    ret.BoxMargin = src.BoxMargin
    ret.CornerRadius = src.CornerRadius
    ret.Image = src.Image

    return &ret
}








type Connection struct {

    // box id of the destination
    DestId string  `yaml:"destId"`

    // Arrow at the source box
    SourceArrow bool  `yaml:"sourceArrow"`

    // Arrow at the destination box
    DestArrow bool  `yaml:"destArrow"`

    // optional format to style the connection
    Format *string  `yaml:"format,omitempty"`

    // Tags to annotate the connection, tags are used to format
    Tags []string  `yaml:"tags,omitempty"`
}

func NewConnection() *Connection {
    return &Connection{
        Tags: make([]string, 0),
    }
}

func CopyConnection(src *Connection) *Connection {
    if src == nil {
        return nil
    }
    var ret Connection
    ret.DestId = src.DestId
    ret.SourceArrow = src.SourceArrow
    ret.DestArrow = src.DestArrow
    ret.Format = src.Format
    ret.Tags = make([]string, 0)
    for _, e := range src.Tags {
        ret.Tags = append(ret.Tags, e)
    }

    return &ret
}




