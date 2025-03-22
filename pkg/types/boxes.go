package types

// Attention, this file is generated. Manual changes get lost with the next
// run of the code generation.
// created by yacg (template: golang_types.mako v1.1.0)

import (
)


/* Model to describe the input of block diagrams
*/
type Boxes struct {

    // Title of the document
    Title string  `yaml:"title"`

    Boxes Layout  `yaml:"boxes"`

    // Map of formats available to be used in the boxes
    Formats map[string]Format  `yaml:"formats,omitempty"`

    // Minimum margin between boxes
    MinBoxMargin *int  `yaml:"minBoxMargin,omitempty"`

    // Minimum margin between connectors
    MinConnectorMargin *int  `yaml:"minConnectorMargin,omitempty"`
}

func NewBoxes() *Boxes {
        return &Boxes{
            Boxes: *NewLayout(),
            Formats: make(map[string]Format, 0),
        }
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

    Vertical []Layout  `yaml:"vertical,omitempty"`

    Horizontal []Layout  `yaml:"horizontal,omitempty"`

    // Tags to annotate the box, tags are used to format and filter
    Tags []string  `yaml:"tags,omitempty"`

    // List of connections to other boxes
    Connections []Connection  `yaml:"connections,omitempty"`
}

func NewLayout() *Layout {
        return &Layout{
            Vertical: make([]Layout, 0),
            Horizontal: make([]Layout, 0),
            Tags: make([]string, 0),
            Connections: make([]Connection, 0),
        }
}








type Format struct {

    FontCaption *FontDef  `yaml:"fontCaption,omitempty"`

    FontText1 *FontDef  `yaml:"fontText1,omitempty"`

    FontText2 *FontDef  `yaml:"fontText2,omitempty"`

    Border *LineDef  `yaml:"border,omitempty"`

    Fill *FillDef  `yaml:"fill,omitempty"`
}









type Connection struct {

    // box id of the destination
    DestId *string  `yaml:"destId,omitempty"`

    // Arrow at the source box
    SourceArrow *bool  `yaml:"sourceArrow,omitempty"`

    // Arrow at the destination box
    DestArrow *bool  `yaml:"destArrow,omitempty"`

    // Tags to annotate the connection, tags are used to format
    Tags []string  `yaml:"tags,omitempty"`
}

func NewConnection() *Connection {
        return &Connection{
            Tags: make([]string, 0),
        }
}




