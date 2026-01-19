package boxes

// Attention, this file is generated. Manual changes get lost with the next
// run of the code generation.
// created by yacg (template: golang_types.mako v1.1.0)

import (
    "github.com/okieoth/draw.chart.things/pkg/types"
)


/* Internal model to describe the layout if the boxes chart
*/
type BoxesDocument struct {

    // Title of the document
    Title string  `yaml:"title"`

    Boxes LayoutElement  `yaml:"boxes"`

    // Height of the document
    Height int  `yaml:"height"`

    // Width of the document
    Width int  `yaml:"width"`

    Connections []ConnectionElem  `yaml:"connections,omitempty"`

    ConnectedElems []string  `yaml:"connectedElems,omitempty"`

    // Padding used as default over the whole diagram
    GlobalPadding int  `yaml:"globalPadding"`

    // Minimum margin between boxes
    MinBoxMargin int  `yaml:"minBoxMargin"`

    // Minimum margin between connectors
    MinConnectorMargin int  `yaml:"minConnectorMargin"`

    // Map of formats available to be used in the boxes
    Formats map[string]BoxFormat  `yaml:"formats,omitempty"`

    // optional list of images used in the generated graphic
    Images []types.ImageDef  `yaml:"images,omitempty"`

    // Vertical roads available to connect boxes in a vertical way
    VerticalRoads []ConnectionLine  `yaml:"verticalRoads,omitempty"`

    // Horizontal roads available to connect boxes in a horizontal way
    HorizontalRoads []ConnectionLine  `yaml:"horizontalRoads,omitempty"`

    // helper structure to take the node points for the possible connection graph
    ConnectionNodes []ConnectionNode  `yaml:"connectionNodes,omitempty"`

    // helper structure for resolving the collisions
    HorizontalLines []ConnectionLine  `yaml:"horizontalLines,omitempty"`

    // helper structure for resolving the collisions
    VerticalLines []ConnectionLine  `yaml:"verticalLines,omitempty"`
}

func NewBoxesDocument() *BoxesDocument {
    return &BoxesDocument{
        Boxes: *NewLayoutElement(),
        Connections: make([]ConnectionElem, 0),
        ConnectedElems: make([]string, 0),
        Formats: make(map[string]BoxFormat, 0),
        Images: make([]types.ImageDef, 0),
        VerticalRoads: make([]ConnectionLine, 0),
        HorizontalRoads: make([]ConnectionLine, 0),
        ConnectionNodes: make([]ConnectionNode, 0),
        HorizontalLines: make([]ConnectionLine, 0),
        VerticalLines: make([]ConnectionLine, 0),
    }
}

func CopyBoxesDocument(src *BoxesDocument) *BoxesDocument {
    if src == nil {
        return nil
    }
    var ret BoxesDocument
    ret.Title = src.Title
    ret.Boxes = *CopyLayoutElement(&src.Boxes)
    ret.Height = src.Height
    ret.Width = src.Width
    ret.Connections = make([]ConnectionElem, 0)
    for _, e := range src.Connections {
        ret.Connections = append(ret.Connections, e)
    }
    ret.ConnectedElems = make([]string, 0)
    for _, e := range src.ConnectedElems {
        ret.ConnectedElems = append(ret.ConnectedElems, e)
    }
    ret.GlobalPadding = src.GlobalPadding
    ret.MinBoxMargin = src.MinBoxMargin
    ret.MinConnectorMargin = src.MinConnectorMargin
    ret.Formats = make(map[string]BoxFormat, 0)
    for k, v := range src.Formats {
        ret.Formats[k] = v
    }
    ret.Images = make([]types.ImageDef, 0)
    for _, e := range src.Images {
        ret.Images = append(ret.Images, e)
    }
    ret.VerticalRoads = make([]ConnectionLine, 0)
    for _, e := range src.VerticalRoads {
        ret.VerticalRoads = append(ret.VerticalRoads, e)
    }
    ret.HorizontalRoads = make([]ConnectionLine, 0)
    for _, e := range src.HorizontalRoads {
        ret.HorizontalRoads = append(ret.HorizontalRoads, e)
    }
    ret.ConnectionNodes = make([]ConnectionNode, 0)
    for _, e := range src.ConnectionNodes {
        ret.ConnectionNodes = append(ret.ConnectionNodes, e)
    }
    ret.HorizontalLines = make([]ConnectionLine, 0)
    for _, e := range src.HorizontalLines {
        ret.HorizontalLines = append(ret.HorizontalLines, e)
    }
    ret.VerticalLines = make([]ConnectionLine, 0)
    for _, e := range src.VerticalLines {
        ret.VerticalLines = append(ret.VerticalLines, e)
    }

    return &ret
}





type LayoutElement struct {

    // unique identifier of that entry
    Id string  `yaml:"id"`

    // Some kind of the main text
    Caption string  `yaml:"caption"`

    // First additional text
    Text1 string  `yaml:"text1"`

    // Second additional text
    Text2 string  `yaml:"text2"`

    Vertical *LayoutElemContainer  `yaml:"vertical,omitempty"`

    Horizontal *LayoutElemContainer  `yaml:"horizontal,omitempty"`

    // X position of the element
    X int  `yaml:"x"`

    // Y position of the element
    Y int  `yaml:"y"`

    // Width of the element
    Width int  `yaml:"width"`

    // Height of the element
    Height int  `yaml:"height"`

    // X position of the center of the element
    CenterX int  `yaml:"centerX"`

    // X position of the center of the element
    CenterY int  `yaml:"centerY"`

    Format *BoxFormat  `yaml:"format,omitempty"`

    Connections []LayoutElemConnection  `yaml:"connections,omitempty"`

    // Y position of the left side of the element to start the connection
    LeftYToStart *int  `yaml:"leftYToStart,omitempty"`

    // Y position of the right side of the element to start the connection
    RightYToStart *int  `yaml:"rightYToStart,omitempty"`

    // X position of the top side of the element to start the connection
    TopXToStart *int  `yaml:"topXToStart,omitempty"`

    // X position of the bottom side of the element to start the connection
    BottomXToStart *int  `yaml:"bottomXToStart,omitempty"`

    // X position where the text would start
    XTextBox *int  `yaml:"xTextBox,omitempty"`

    // Y position where the text would start
    YTextBox *int  `yaml:"yTextBox,omitempty"`

    // Width of the text area
    WidthTextBox *int  `yaml:"widthTextBox,omitempty"`

    // Height of the text area
    HeightTextBox *int  `yaml:"heightTextBox,omitempty"`

    // Tags to annotate the box, tags are used to format and filter
    Tags []string  `yaml:"tags,omitempty"`
}

func NewLayoutElement() *LayoutElement {
    return &LayoutElement{
        Vertical: NewLayoutElemContainer(),
        Horizontal: NewLayoutElemContainer(),
        Connections: make([]LayoutElemConnection, 0),
        Tags: make([]string, 0),
    }
}

func CopyLayoutElement(src *LayoutElement) *LayoutElement {
    if src == nil {
        return nil
    }
    var ret LayoutElement
    ret.Id = src.Id
    ret.Caption = src.Caption
    ret.Text1 = src.Text1
    ret.Text2 = src.Text2
    ret.Vertical = CopyLayoutElemContainer(src.Vertical)
    ret.Horizontal = CopyLayoutElemContainer(src.Horizontal)
    ret.X = src.X
    ret.Y = src.Y
    ret.Width = src.Width
    ret.Height = src.Height
    ret.CenterX = src.CenterX
    ret.CenterY = src.CenterY
    ret.Format = CopyBoxFormat(src.Format)
    ret.Connections = make([]LayoutElemConnection, 0)
    for _, e := range src.Connections {
        ret.Connections = append(ret.Connections, e)
    }
    ret.LeftYToStart = src.LeftYToStart
    ret.RightYToStart = src.RightYToStart
    ret.TopXToStart = src.TopXToStart
    ret.BottomXToStart = src.BottomXToStart
    ret.XTextBox = src.XTextBox
    ret.YTextBox = src.YTextBox
    ret.WidthTextBox = src.WidthTextBox
    ret.HeightTextBox = src.HeightTextBox
    ret.Tags = make([]string, 0)
    for _, e := range src.Tags {
        ret.Tags = append(ret.Tags, e)
    }

    return &ret
}





type ConnectionElem struct {

    // ID of the box where the connector starts
    From *string  `yaml:"from,omitempty"`

    // ID of the box where the connector ends
    To *string  `yaml:"to,omitempty"`

    // Arrow at the source box
    SourceArrow *bool  `yaml:"sourceArrow,omitempty"`

    // Arrow at the destination box
    DestArrow *bool  `yaml:"destArrow,omitempty"`

    Format *types.LineDef  `yaml:"format,omitempty"`

    Parts []ConnectionLine  `yaml:"parts,omitempty"`

    // index of this connection, in the boxes_document object
    ConnectionIndex int  `yaml:"connectionIndex"`
}

func NewConnectionElem() *ConnectionElem {
    return &ConnectionElem{
        Parts: make([]ConnectionLine, 0),
    }
}

func CopyConnectionElem(src *ConnectionElem) *ConnectionElem {
    if src == nil {
        return nil
    }
    var ret ConnectionElem
    ret.From = src.From
    ret.To = src.To
    ret.SourceArrow = src.SourceArrow
    ret.DestArrow = src.DestArrow
    ret.Format = types.CopyLineDef(src.Format)
    ret.Parts = make([]ConnectionLine, 0)
    for _, e := range src.Parts {
        ret.Parts = append(ret.Parts, e)
    }
    ret.ConnectionIndex = src.ConnectionIndex

    return &ret
}








type BoxFormat struct {

    // Padding of the box
    Padding int  `yaml:"padding"`

    FontCaption types.FontDef  `yaml:"fontCaption"`

    FontText1 types.FontDef  `yaml:"fontText1"`

    FontText2 types.FontDef  `yaml:"fontText2"`

    Line *types.LineDef  `yaml:"line,omitempty"`

    // radius of the box corners in pixel
    CornerRadius *int  `yaml:"cornerRadius,omitempty"`

    Fill *types.FillDef  `yaml:"fill,omitempty"`

    // ID of an image to use
    Image *string  `yaml:"image,omitempty"`

    // Minimum margin between boxes
    MinBoxMargin int  `yaml:"minBoxMargin"`

    // optional fixed width that will be applied on the box
    FixedWidth *int  `yaml:"fixedWidth,omitempty"`

    // optional fixed height that will be applied on the box
    FixedHeight *int  `yaml:"fixedHeight,omitempty"`

    // If true, the text will be displayed vertically
    VerticalTxt bool  `yaml:"verticalTxt"`
}


func CopyBoxFormat(src *BoxFormat) *BoxFormat {
    if src == nil {
        return nil
    }
    var ret BoxFormat
    ret.Padding = src.Padding
    ret.FontCaption = *types.CopyFontDef(&src.FontCaption)
    ret.FontText1 = *types.CopyFontDef(&src.FontText1)
    ret.FontText2 = *types.CopyFontDef(&src.FontText2)
    ret.Line = types.CopyLineDef(src.Line)
    ret.CornerRadius = src.CornerRadius
    ret.Fill = types.CopyFillDef(src.Fill)
    ret.Image = src.Image
    ret.MinBoxMargin = src.MinBoxMargin
    ret.FixedWidth = src.FixedWidth
    ret.FixedHeight = src.FixedHeight
    ret.VerticalTxt = src.VerticalTxt

    return &ret
}





type ConnectionLine struct {

    StartX int  `yaml:"startX"`

    StartY int  `yaml:"startY"`

    EndX int  `yaml:"endX"`

    EndY int  `yaml:"endY"`

    // in case the line is connected to the start layout element, then here is its id
    SrcLayoutId *string  `yaml:"srcLayoutId,omitempty"`

    // in case the line is connected to a end layout element, then here is its id
    DestLayoutId *string  `yaml:"destLayoutId,omitempty"`

    // index of the connection, in the boxes_document object, where this line belongs too
    ConnectionIndex int  `yaml:"connectionIndex"`
}


func CopyConnectionLine(src *ConnectionLine) *ConnectionLine {
    if src == nil {
        return nil
    }
    var ret ConnectionLine
    ret.StartX = src.StartX
    ret.StartY = src.StartY
    ret.EndX = src.EndX
    ret.EndY = src.EndY
    ret.SrcLayoutId = src.SrcLayoutId
    ret.DestLayoutId = src.DestLayoutId
    ret.ConnectionIndex = src.ConnectionIndex

    return &ret
}





/* helper type for calculating the connections between elements
*/
type ConnectionNode struct {

    // X position of the element
    X int  `yaml:"x"`

    // Y position of the element
    Y int  `yaml:"y"`

    // in case the edge ends in a layout element, it takes the ID
    NodeId *string  `yaml:"nodeId,omitempty"`

    // optional box id, only on connection nodes that are the entry points to real box connections
    BoxId *string  `yaml:"boxId,omitempty"`

    Edges []ConnectionEdge  `yaml:"edges,omitempty"`
}

func NewConnectionNode() *ConnectionNode {
    return &ConnectionNode{
        Edges: make([]ConnectionEdge, 0),
    }
}

func CopyConnectionNode(src *ConnectionNode) *ConnectionNode {
    if src == nil {
        return nil
    }
    var ret ConnectionNode
    ret.X = src.X
    ret.Y = src.Y
    ret.NodeId = src.NodeId
    ret.BoxId = src.BoxId
    ret.Edges = make([]ConnectionEdge, 0)
    for _, e := range src.Edges {
        ret.Edges = append(ret.Edges, e)
    }

    return &ret
}





type LayoutElemContainer struct {

    // X position of the element
    X int  `yaml:"x"`

    // Y position of the element
    Y int  `yaml:"y"`

    // Width of the container
    Width int  `yaml:"width"`

    // Height of the container
    Height int  `yaml:"height"`

    Elems []LayoutElement  `yaml:"elems,omitempty"`
}

func NewLayoutElemContainer() *LayoutElemContainer {
    return &LayoutElemContainer{
        Elems: make([]LayoutElement, 0),
    }
}

func CopyLayoutElemContainer(src *LayoutElemContainer) *LayoutElemContainer {
    if src == nil {
        return nil
    }
    var ret LayoutElemContainer
    ret.X = src.X
    ret.Y = src.Y
    ret.Width = src.Width
    ret.Height = src.Height
    ret.Elems = make([]LayoutElement, 0)
    for _, e := range src.Elems {
        ret.Elems = append(ret.Elems, e)
    }

    return &ret
}





type LayoutElemConnection struct {

    // box id of the destination
    DestId string  `yaml:"destId"`

    // Arrow at the source box
    SourceArrow bool  `yaml:"sourceArrow"`

    // Arrow at the destination box
    DestArrow bool  `yaml:"destArrow"`

    Format *types.LineDef  `yaml:"format,omitempty"`

    // Tags to annotate the connection, tags are used to format
    Tags []string  `yaml:"tags,omitempty"`
}

func NewLayoutElemConnection() *LayoutElemConnection {
    return &LayoutElemConnection{
        Tags: make([]string, 0),
    }
}

func CopyLayoutElemConnection(src *LayoutElemConnection) *LayoutElemConnection {
    if src == nil {
        return nil
    }
    var ret LayoutElemConnection
    ret.DestId = src.DestId
    ret.SourceArrow = src.SourceArrow
    ret.DestArrow = src.DestArrow
    ret.Format = types.CopyLineDef(src.Format)
    ret.Tags = make([]string, 0)
    for _, e := range src.Tags {
        ret.Tags = append(ret.Tags, e)
    }

    return &ret
}





/* edge type to store edges on a connection node
*/
type ConnectionEdge struct {

    // X position of the element
    X int  `yaml:"x"`

    // Y position of the element
    Y int  `yaml:"y"`

    // either the box ID where the edge ends or the ID connection node
    DestNodeId *string  `yaml:"destNodeId,omitempty"`

    // weight of the connection, based on the distance
    Weight int  `yaml:"weight"`
}


func CopyConnectionEdge(src *ConnectionEdge) *ConnectionEdge {
    if src == nil {
        return nil
    }
    var ret ConnectionEdge
    ret.X = src.X
    ret.Y = src.Y
    ret.DestNodeId = src.DestNodeId
    ret.Weight = src.Weight

    return &ret
}




