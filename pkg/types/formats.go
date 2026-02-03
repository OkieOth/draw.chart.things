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


func CopyFontDef(src *FontDef) *FontDef {
    if src == nil {
        return nil
    }
    var ret FontDef
    ret.Size = src.Size
    ret.Font = src.Font
    ret.Type = src.Type
    ret.Weight = src.Weight
    ret.LineHeight = src.LineHeight
    ret.Color = src.Color
    ret.Aligned = src.Aligned
    ret.SpaceTop = src.SpaceTop
    ret.SpaceBottom = src.SpaceBottom
    ret.MaxLenBeforeBreak = src.MaxLenBeforeBreak
    ret.Anchor = src.Anchor

    return &ret
}




type FontDefTypeEnum string

const (
    FontDefTypeEnum_normal FontDefTypeEnum = "normal"
    FontDefTypeEnum_italic FontDefTypeEnum = "italic"
)


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




type FontDefWeightEnum string

const (
    FontDefWeightEnum_normal FontDefWeightEnum = "normal"
    FontDefWeightEnum_bold FontDefWeightEnum = "bold"
)


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




type FontDefAlignedEnum string

const (
    FontDefAlignedEnum_left FontDefAlignedEnum = "left"
    FontDefAlignedEnum_center FontDefAlignedEnum = "center"
    FontDefAlignedEnum_right FontDefAlignedEnum = "right"
)


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




type FontDefAnchorEnum string

const (
    FontDefAnchorEnum_middle FontDefAnchorEnum = "middle"
    FontDefAnchorEnum_left FontDefAnchorEnum = "left"
    FontDefAnchorEnum_right FontDefAnchorEnum = "right"
)


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

    Width *float64  `yaml:"width,omitempty"`

    Style *LineDefStyleEnum  `yaml:"style,omitempty"`

    Color *string  `yaml:"color,omitempty"`

    Opacity *float64  `yaml:"opacity,omitempty"`
}


func CopyLineDef(src *LineDef) *LineDef {
    if src == nil {
        return nil
    }
    var ret LineDef
    ret.Width = src.Width
    ret.Style = src.Style
    ret.Color = src.Color
    ret.Opacity = src.Opacity

    return &ret
}




type LineDefStyleEnum string

const (
    LineDefStyleEnum_solid LineDefStyleEnum = "solid"
    LineDefStyleEnum_dotted LineDefStyleEnum = "dotted"
    LineDefStyleEnum_dashed LineDefStyleEnum = "dashed"
)


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


func CopyFillDef(src *FillDef) *FillDef {
    if src == nil {
        return nil
    }
    var ret FillDef
    ret.Color = src.Color
    ret.Opacity = src.Opacity

    return &ret
}





/* parameters of an image to be displayed in the SVG
*/
type ImageDef struct {

    // some words to explain what this image is about
    Description *string  `yaml:"description,omitempty"`

    // with of the displayed image
    Width int  `yaml:"width"`

    // height of the displayed image
    Height int  `yaml:"height"`

    // base64 string of the image to use
    Base64 *string  `yaml:"base64,omitempty"`

    // distance top and bottom of the image
    MarginTopBottom *int  `yaml:"marginTopBottom,omitempty"`

    // distance left and right of the image
    MarginLeftRight *int  `yaml:"marginLeftRight,omitempty"`

    // file path to a text file that contains the base64 of the png
    Base64Src *string  `yaml:"base64Src,omitempty"`
}


func CopyImageDef(src *ImageDef) *ImageDef {
    if src == nil {
        return nil
    }
    var ret ImageDef
    ret.Description = src.Description
    ret.Width = src.Width
    ret.Height = src.Height
    ret.Base64 = src.Base64
    ret.MarginTopBottom = src.MarginTopBottom
    ret.MarginLeftRight = src.MarginLeftRight
    ret.Base64Src = src.Base64Src

    return &ret
}




