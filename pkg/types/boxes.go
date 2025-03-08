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
    Title *string  `yaml:"title,omitempty"`

    Boxes Layout  `yaml:"boxes"`

    DefaultFormat *Format  `yaml:"defaultFormat,omitempty"`

    Formats map[string]Format  `yaml:"formats,omitempty"`

    // Minimum margin between boxes
    MinBoxMargin *int32  `yaml:"minBoxMargin,omitempty"`

    // Minimum margin between connectors
    MinConnectorMargin *int32  `yaml:"minConnectorMargin,omitempty"`
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





/* Defines the font a text
*/
type FontDef struct {

    Size *int32  `yaml:"size,omitempty"`

    Font *string  `yaml:"font,omitempty"`

    Type *FontDefTypeEnum  `yaml:"type,omitempty"`

    Weight *FontDefWeightEnum  `yaml:"weight,omitempty"`

    Color *string  `yaml:"color,omitempty"`

    Alligned *FontDefAllignedEnum  `yaml:"alligned,omitempty"`

    SpaceTop *int32  `yaml:"spaceTop,omitempty"`

    SpaceBottom *int32  `yaml:"spaceBottom,omitempty"`
}






/* Defines how the border of the box looks like
*/
type LineDef struct {

    Width *int32  `yaml:"width,omitempty"`

    Style *LineDefStyleEnum  `yaml:"style,omitempty"`

    Color *string  `yaml:"color,omitempty"`

    Opacity *float64  `yaml:"opacity,omitempty"`
}






/* Defines the fill of the box
*/
type FillDef struct {

    Color *string  `yaml:"color,omitempty"`

    Opacity *float64  `yaml:"opacity,omitempty"`
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




type LineDefStyleEnum int64

const (
    LineDefStyleEnum_solid LineDefStyleEnum = iota
        LineDefStyleEnum_dotted
        LineDefStyleEnum_dashed
)

func (s LineDefStyleEnum) String() string {
	switch s {
	case LineDefStyleEnum_solid:
		return "solid"
	case LineDefStyleEnum_dotted:
		return "dotted"
	case LineDefStyleEnum_dashed:
		return "dashed"
    default:
        return "???"
	}
}

func (s LineDefStyleEnum) MarshalJSON() ([]byte, error) {
    return json.Marshal(s.String())
}

func (s *LineDefStyleEnum) UnmarshalJSON(data []byte) error {
    var value string
    if err := json.Unmarshal(data, &value); err != nil {
        return err
    }

    switch value {
    case "solid":
        *s = LineDefStyleEnum_solid 
    case "dotted":
        *s = LineDefStyleEnum_dotted 
    case "dashed":
        *s = LineDefStyleEnum_dashed 
    default:
		msg := fmt.Sprintf("invalid value for DDDDomainType: %s", value)
		return errors.New(msg)
    }
    return nil
}




