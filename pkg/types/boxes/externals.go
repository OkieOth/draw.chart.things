package boxes

// Attention, this file is generated. Manual changes get lost with the next
// run of the code generation.
// created by yacg (template: golang_types.mako v1.1.0)

import (
    "github.com/okieoth/draw.chart.things/pkg/types"
)


/* Model to inject additional things in a boxes layout definition
*/
type BoxesFileMixings struct {

    // optional title, that's appended to the original layout title
    Title *string  `yaml:"title,omitempty"`

    // allows to include a version for the layout description
    Version *string  `yaml:"version,omitempty"`

    // Legend definition used in this diagram
    Legend *Legend  `yaml:"legend,omitempty"`

    // dictionary for layout mixins. key of the dictionary is the caption of the box that will take the additional content
    LayoutMixins map[string]LayoutMixin  `yaml:"layoutMixins,omitempty"`

    // dictionary of connection objects
    Connections map[string]ConnectionCont  `yaml:"connections,omitempty"`

    Formats map[string]Format  `yaml:"formats,omitempty"`

    // Set of formats that overwrites the style of boxes, if specific conditions are met
    FormatVariations *FormatVariations  `yaml:"formatVariations,omitempty"`

    // dictionary of comment objects
    Comments map[string]types.Comment  `yaml:"comments,omitempty"`

    // optional map of images used in the generated graphic
    Images map[string]types.ImageDef  `yaml:"images,omitempty"`
}

func NewBoxesFileMixings() *BoxesFileMixings {
    return &BoxesFileMixings{
        Legend: NewLegend(),
        LayoutMixins: make(map[string]LayoutMixin, 0),
        Connections: make(map[string]ConnectionCont, 0),
        Formats: make(map[string]Format, 0),
        FormatVariations: NewFormatVariations(),
        Comments: make(map[string]types.Comment, 0),
        Images: make(map[string]types.ImageDef, 0),
    }
}

func CopyBoxesFileMixings(src *BoxesFileMixings) *BoxesFileMixings {
    if src == nil {
        return nil
    }
    var ret BoxesFileMixings
    ret.Title = src.Title
    ret.Version = src.Version
    ret.Legend = CopyLegend(src.Legend)
    ret.LayoutMixins = make(map[string]LayoutMixin, 0)
    for k, v := range src.LayoutMixins {
        ret.LayoutMixins[k] = v
    }
    ret.Connections = make(map[string]ConnectionCont, 0)
    for k, v := range src.Connections {
        ret.Connections[k] = v
    }
    ret.Formats = make(map[string]Format, 0)
    for k, v := range src.Formats {
        ret.Formats[k] = v
    }
    ret.FormatVariations = CopyFormatVariations(src.FormatVariations)
    ret.Comments = make(map[string]types.Comment, 0)
    for k, v := range src.Comments {
        ret.Comments[k] = v
    }
    ret.Images = make(map[string]types.ImageDef, 0)
    for k, v := range src.Images {
        ret.Images[k] = v
    }

    return &ret
}











type ConnectionCont struct {

    Connections []Connection  `yaml:"connections,omitempty"`
}

func NewConnectionCont() *ConnectionCont {
    return &ConnectionCont{
        Connections: make([]Connection, 0),
    }
}

func CopyConnectionCont(src *ConnectionCont) *ConnectionCont {
    if src == nil {
        return nil
    }
    var ret ConnectionCont
    ret.Connections = make([]Connection, 0)
    for _, e := range src.Connections {
        ret.Connections = append(ret.Connections, e)
    }

    return &ret
}













