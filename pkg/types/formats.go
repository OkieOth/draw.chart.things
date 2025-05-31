package types

// Attention, this file is generated. Manual changes get lost with the next
// run of the code generation.
// created by yacg (template: golang_types.mako v1.1.0)

import (
    "encoding/json"
    "errors"
    "fmt"
)


/* Defines the font a text
*/
type FontDef struct {

    Size int  `yaml:"size"`

    Font string  `yaml:"font"`

    Type *FontDefTypeEnum  `yaml:"type,omitempty"`

    Weight *FontDefWeightEnum  `yaml:"weight,omitempty"`

    // Line height of the box
    LineHeight float32  `yaml:"lineHeight"`

    Color string  `yaml:"color"`

    Aligned *FontDefAlignedEnum  `yaml:"aligned,omitempty"`

    SpaceTop int  `yaml:"spaceTop"`

    SpaceBottom int  `yaml:"spaceBottom"`

    // Maximum length of the text before it breaks
    MaxLenBeforeBreak int  `yaml:"maxLenBeforeBreak"`

    Anchor FontDefAnchorEnum  `yaml:"anchor"`
}





type FontDefTypeEnum int64

const (
    FontDefTypeEnum_normal FontDefTypeEnum = iota
        FontDefTypeEnum_italic
)

func (s FontDefTypeEnum) String() string {
	switch s {
	case FontDefTypeEnum_normal:
		return "normal"
	case FontDefTypeEnum_italic:
		return "italic"
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
)

func (s FontDefWeightEnum) String() string {
	switch s {
	case FontDefWeightEnum_normal:
		return "normal"
	case FontDefWeightEnum_bold:
		return "bold"
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
    default:
		msg := fmt.Sprintf("invalid value for DDDDomainType: %s", value)
		return errors.New(msg)
    }
    return nil
}




type FontDefAlignedEnum int64

const (
    FontDefAlignedEnum_left FontDefAlignedEnum = iota
        FontDefAlignedEnum_center
        FontDefAlignedEnum_right
)

func (s FontDefAlignedEnum) String() string {
	switch s {
	case FontDefAlignedEnum_left:
		return "left"
	case FontDefAlignedEnum_center:
		return "center"
	case FontDefAlignedEnum_right:
		return "right"
    default:
        return "???"
	}
}

func (s FontDefAlignedEnum) MarshalJSON() ([]byte, error) {
    return json.Marshal(s.String())
}

func (s *FontDefAlignedEnum) UnmarshalJSON(data []byte) error {
    var value string
    if err := json.Unmarshal(data, &value); err != nil {
        return err
    }

    switch value {
    case "left":
        *s = FontDefAlignedEnum_left 
    case "center":
        *s = FontDefAlignedEnum_center 
    case "right":
        *s = FontDefAlignedEnum_right 
    default:
		msg := fmt.Sprintf("invalid value for DDDDomainType: %s", value)
		return errors.New(msg)
    }
    return nil
}




type FontDefAnchorEnum int64

const (
    FontDefAnchorEnum_middle FontDefAnchorEnum = iota
        FontDefAnchorEnum_left
        FontDefAnchorEnum_right
)

func (s FontDefAnchorEnum) String() string {
	switch s {
	case FontDefAnchorEnum_middle:
		return "middle"
	case FontDefAnchorEnum_left:
		return "left"
	case FontDefAnchorEnum_right:
		return "right"
    default:
        return "???"
	}
}

func (s FontDefAnchorEnum) MarshalJSON() ([]byte, error) {
    return json.Marshal(s.String())
}

func (s *FontDefAnchorEnum) UnmarshalJSON(data []byte) error {
    var value string
    if err := json.Unmarshal(data, &value); err != nil {
        return err
    }

    switch value {
    case "middle":
        *s = FontDefAnchorEnum_middle 
    case "left":
        *s = FontDefAnchorEnum_left 
    case "right":
        *s = FontDefAnchorEnum_right 
    default:
		msg := fmt.Sprintf("invalid value for DDDDomainType: %s", value)
		return errors.New(msg)
    }
    return nil
}





/* Defines how the border of the box looks like
*/
type LineDef struct {

    Width *int  `yaml:"width,omitempty"`

    Style *LineDefStyleEnum  `yaml:"style,omitempty"`

    Color *string  `yaml:"color,omitempty"`

    Opacity *float64  `yaml:"opacity,omitempty"`
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





/* Defines the fill of the box
*/
type FillDef struct {

    Color *string  `yaml:"color,omitempty"`

    Opacity *float64  `yaml:"opacity,omitempty"`
}





