package types

// Attention, this file is generated. Manual changes get lost with the next
// run of the code generation.
// created by yacg (template: golang_types.mako v1.1.0)

import (
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

    // Padding used as default over the whole diagram
    GlobalPadding int  `yaml:"globalPadding"`

    // Minimum margin between boxes
    MinBoxMargin int  `yaml:"minBoxMargin"`

    // Minimum margin between connectors
    MinConnectorMargin int  `yaml:"minConnectorMargin"`

    // Map of formats available to be used in the boxes
    Formats map[string]BoxFormat  `yaml:"formats,omitempty"`

    // Vertical roads available to connect boxes in a vertical way
    VerticalRoads []ConnectionLine  `yaml:"verticalRoads,omitempty"`

    // Horizontal roads available to connect boxes in a horizontal way
    HorizontalRoads []ConnectionLine  `yaml:"horizontalRoads,omitempty"`
}

func NewBoxesDocument() *BoxesDocument {
        return &BoxesDocument{
            Boxes: *NewLayoutElement(),
            Connections: make([]ConnectionElem, 0),
            Formats: make(map[string]BoxFormat, 0),
            VerticalRoads: make([]ConnectionLine, 0),
            HorizontalRoads: make([]ConnectionLine, 0),
        }
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





type ConnectionElem struct {

    // Reference to the box where the connector starts
    From *LayoutElement  `yaml:"from,omitempty"`

    // Reference to the box where the connector ends
    To *LayoutElement  `yaml:"to,omitempty"`

    // Arrow at the source box
    SourceArrow *bool  `yaml:"sourceArrow,omitempty"`

    // Arrow at the destination box
    DestArrow *bool  `yaml:"destArrow,omitempty"`

    Format *LineDef  `yaml:"format,omitempty"`

    Parts []ConnectionLine  `yaml:"parts,omitempty"`
}

func NewConnectionElem() *ConnectionElem {
        return &ConnectionElem{
            From: NewLayoutElement(),
            To: NewLayoutElement(),
            Parts: make([]ConnectionLine, 0),
        }
}








type BoxFormat struct {

    // Padding of the box
    Padding int  `yaml:"padding"`

    FontCaption FontDef  `yaml:"fontCaption"`

    FontText1 FontDef  `yaml:"fontText1"`

    FontText2 FontDef  `yaml:"fontText2"`

    Border *LineDef  `yaml:"border,omitempty"`

    Fill *FillDef  `yaml:"fill,omitempty"`

    // Minimum margin between boxes
    MinBoxMargin int  `yaml:"minBoxMargin"`
}






type ConnectionLine struct {

    StartX int  `yaml:"startX"`

    StartY int  `yaml:"startY"`

    EndX int  `yaml:"endX"`

    EndY int  `yaml:"endY"`

    // If the line is moved out by a box with a collision
    MovedOut bool  `yaml:"movedOut"`
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





type LayoutElemConnection struct {

    // box id of the destination
    DestId string  `yaml:"destId"`

    // Arrow at the source box
    SourceArrow bool  `yaml:"sourceArrow"`

    // Arrow at the destination box
    DestArrow bool  `yaml:"destArrow"`

    // Tags to annotate the connection, tags are used to format
    Tags []string  `yaml:"tags,omitempty"`
}

func NewLayoutElemConnection() *LayoutElemConnection {
        return &LayoutElemConnection{
            Tags: make([]string, 0),
        }
}




