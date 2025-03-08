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
    Title *string  `yaml:"title,omitempty"`

    Boxes LayoutElement  `yaml:"boxes"`

    // Height of the document
    Height int32  `yaml:"height"`

    // Width of the document
    Width int32  `yaml:"width"`

    Connections []ConnectionElem  `yaml:"connections,omitempty"`
}

func NewBoxesDocument() *BoxesDocument {
        return &BoxesDocument{
            Boxes: *NewLayoutElement(),
            Connections: make([]ConnectionElem, 0),
        }
}





type LayoutElement struct {

    // unique identifier of that entry
    Id string  `yaml:"id"`

    // Some kind of the main text
    Caption *string  `yaml:"caption,omitempty"`

    // First additional text
    Text1 *string  `yaml:"text1,omitempty"`

    // Second additional text
    Text2 *string  `yaml:"text2,omitempty"`

    Vertical []LayoutElement  `yaml:"vertical,omitempty"`

    Horizontal []LayoutElement  `yaml:"horizontal,omitempty"`

    // X position of the element
    X int32  `yaml:"x"`

    // Y position of the element
    Y int32  `yaml:"y"`

    // Width of the element
    Width int32  `yaml:"width"`

    // Height of the element
    Height int32  `yaml:"height"`

    Format Format  `yaml:"format"`

    // Tags to annotate the box, tags are used to format and filter
    Tags []string  `yaml:"tags,omitempty"`
}

func NewLayoutElement() *LayoutElement {
        return &LayoutElement{
            Vertical: make([]LayoutElement, 0),
            Horizontal: make([]LayoutElement, 0),
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

    Points []ConnectionPoint  `yaml:"points,omitempty"`
}

func NewConnectionElem() *ConnectionElem {
        return &ConnectionElem{
            From: NewLayoutElement(),
            To: NewLayoutElement(),
            Points: make([]ConnectionPoint, 0),
        }
}








type ConnectionPoint struct {

    // X position of the point
    X *int32  `yaml:"x,omitempty"`

    // Y position of the point
    Y *int32  `yaml:"y,omitempty"`
}





