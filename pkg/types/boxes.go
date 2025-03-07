package types

// Attention, this file is generated. Manual changes get lost with the next
// run of the code generation.
// created by yacg (template: golang_types.mako v1.1.0)

import (
    "encoding/json"
    "errors"
    "fmt"
)


/* Model to describe the input of block diagrams
*/
type Boxes struct {

    // Title of the document
    Title *string  `json:"title,omitempty"`

    Boxes Layout  `json:"boxes"`

    DefaultFormat *Format  `json:"defaultFormat,omitempty"`

    Formats map[string]Format  `json:"formats,omitempty"`

    // Minimum margin between boxes
    MinBoxMargin *int32  `json:"minBoxMargin,omitempty"`

    // Minimum margin between connectors
    MinConnectorMargin *int32  `json:"minConnectorMargin,omitempty"`
}

func NewBoxes() *Boxes {
        return &Boxes{
            Boxes: *NewLayout(),
            Formats: make(map[string]Format, 0),
        }
}





type Layout struct {

    // unique identifier of that entry
    Id *string  `json:"id,omitempty"`

    // Some kind of the main text
    Caption *string  `json:"caption,omitempty"`

    // First additional text
    Text1 *string  `json:"text1,omitempty"`

    // Second additional text
    Text2 *string  `json:"text2,omitempty"`

    Vertical []Layout  `json:"vertical,omitempty"`

    Horizontal []Layout  `json:"horizontal,omitempty"`

    // Tags to annotate the box, tags are used to format and filter
    Tags []string  `json:"tags,omitempty"`

    // List of connections to other boxes
    Connections []Connection  `json:"connections,omitempty"`
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

    FontCaption *FontDef  `json:"fontCaption,omitempty"`

    FontText1 *FontDef  `json:"fontText1,omitempty"`

    FontText2 *FontDef  `json:"fontText2,omitempty"`

    Border *LineDef  `json:"border,omitempty"`

    Fill *FillDef  `json:"fill,omitempty"`
}












type Connection struct {

    // box id of the destination
    DestId *string  `json:"destId,omitempty"`

    // Arrow at the source box
    SourceArrow *bool  `json:"sourceArrow,omitempty"`

    // Arrow at the destination box
    DestArrow *bool  `json:"destArrow,omitempty"`

    // Tags to annotate the connection, tags are used to format
    Tags []string  `json:"tags,omitempty"`
}

func NewConnection() *Connection {
        return &Connection{
            Tags: make([]string, 0),
        }
}





/* Defines the font a text
*/
type FontDef struct {

    Size *int32  `json:"size,omitempty"`

    Font *string  `json:"font,omitempty"`

    Type *FontDefTypeEnum  `json:"type,omitempty"`

    Weight *FontDefWeightEnum  `json:"weight,omitempty"`

    Color *string  `json:"color,omitempty"`

    Alligned *FontDefAllignedEnum  `json:"alligned,omitempty"`

    SpaceTop *int32  `json:"spaceTop,omitempty"`

    SpaceBottom *int32  `json:"spaceBottom,omitempty"`
}






/* Defines how the border of the box looks like
*/
type LineDef struct {

    Width *int32  `json:"width,omitempty"`

    Color *string  `json:"color,omitempty"`

    Opacity *float64  `json:"opacity,omitempty"`
}






/* Defines the fill of the box
*/
type FillDef struct {

    Color *string  `json:"color,omitempty"`

    Opacity *float64  `json:"opacity,omitempty"`
}





type FontDefTypeEnum int64

const (
    FontDefTypeEnum_normal FontDefTypeEnum = iota
        FontDefTypeEnum_italic
        FontDefTypeEnum_oblique
)

func (s FontDefTypeEnum) String() string {
	switch s {
	case FontDefTypeEnum_normal:
		return "normal"
	case FontDefTypeEnum_italic:
		return "italic"
	case FontDefTypeEnum_oblique:
		return "oblique"
    default:
        return "???"
	}
}

func (s FontDefTypeEnum) MarshalJSON() ([]byte, error) {
    return json.Marshal(s.String())
}

func (s *FontDefTypeEnum) UnmarshalJSON(data []byte) error {
    var value string
    if err := json.Unmarshal(data, &value); err != nil {
        return err
    }

    switch value {
    case "normal":
        *s = FontDefTypeEnum_normal 
    case "italic":
        *s = FontDefTypeEnum_italic 
    case "oblique":
        *s = FontDefTypeEnum_oblique 
    default:
		msg := fmt.Sprintf("invalid value for DDDDomainType: %s", value)
		return errors.New(msg)
    }
    return nil
}




type FontDefWeightEnum int64

const (
    FontDefWeightEnum_normal FontDefWeightEnum = iota
        FontDefWeightEnum_bold
        FontDefWeightEnum_bolder
)

func (s FontDefWeightEnum) String() string {
	switch s {
	case FontDefWeightEnum_normal:
		return "normal"
	case FontDefWeightEnum_bold:
		return "bold"
	case FontDefWeightEnum_bolder:
		return "bolder"
    default:
        return "???"
	}
}

func (s FontDefWeightEnum) MarshalJSON() ([]byte, error) {
    return json.Marshal(s.String())
}

func (s *FontDefWeightEnum) UnmarshalJSON(data []byte) error {
    var value string
    if err := json.Unmarshal(data, &value); err != nil {
        return err
    }

    switch value {
    case "normal":
        *s = FontDefWeightEnum_normal 
    case "bold":
        *s = FontDefWeightEnum_bold 
    case "bolder":
        *s = FontDefWeightEnum_bolder 
    default:
		msg := fmt.Sprintf("invalid value for DDDDomainType: %s", value)
		return errors.New(msg)
    }
    return nil
}




type FontDefAllignedEnum int64

const (
    FontDefAllignedEnum_left FontDefAllignedEnum = iota
        FontDefAllignedEnum_center
        FontDefAllignedEnum_right
)

func (s FontDefAllignedEnum) String() string {
	switch s {
	case FontDefAllignedEnum_left:
		return "left"
	case FontDefAllignedEnum_center:
		return "center"
	case FontDefAllignedEnum_right:
		return "right"
    default:
        return "???"
	}
}

func (s FontDefAllignedEnum) MarshalJSON() ([]byte, error) {
    return json.Marshal(s.String())
}

func (s *FontDefAllignedEnum) UnmarshalJSON(data []byte) error {
    var value string
    if err := json.Unmarshal(data, &value); err != nil {
        return err
    }

    switch value {
    case "left":
        *s = FontDefAllignedEnum_left 
    case "center":
        *s = FontDefAllignedEnum_center 
    case "right":
        *s = FontDefAllignedEnum_right 
    default:
		msg := fmt.Sprintf("invalid value for DDDDomainType: %s", value)
		return errors.New(msg)
    }
    return nil
}




