package boxes

// Attention, this file is generated. Manual changes get lost with the next
// run of the code generation.
// created by yacg (template: golang_types.mako v1.1.0)

import (
    "github.com/okieoth/draw.chart.things/pkg/types"
)


/* Model to inject additional connections in a boxes layout definition
*/
type AdditionalConnections struct {

    // dictionary of connection objects
    Connections map[string]ConnectionCont  `yaml:"connections,omitempty"`

    Formats map[string]Format  `yaml:"formats,omitempty"`
}

func NewAdditionalConnections() *AdditionalConnections {
    return &AdditionalConnections{
        Connections: make(map[string]ConnectionCont, 0),
        Formats: make(map[string]Format, 0),
    }
}

func CopyAdditionalConnections(src *AdditionalConnections) *AdditionalConnections {
    if src == nil {
        return nil
    }
    var ret AdditionalConnections
    ret.Connections = make(map[string]ConnectionCont, 0)
    for k, v := range src.Connections {
        ret.Connections[k] = v
    }
    ret.Formats = make(map[string]Format, 0)
    for k, v := range src.Formats {
        ret.Formats[k] = v
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








/* Model to inject additional formats in a boxes layout definition
*/
type AdditionalFormats struct {

    Formats map[string]Format  `yaml:"formats,omitempty"`

    // optional list of images used in the generated graphic
    Images map[string]types.ImageDef  `yaml:"images,omitempty"`
}

func NewAdditionalFormats() *AdditionalFormats {
    return &AdditionalFormats{
        Formats: make(map[string]Format, 0),
        Images: make(map[string]types.ImageDef, 0),
    }
}

func CopyAdditionalFormats(src *AdditionalFormats) *AdditionalFormats {
    if src == nil {
        return nil
    }
    var ret AdditionalFormats
    ret.Formats = make(map[string]Format, 0)
    for k, v := range src.Formats {
        ret.Formats[k] = v
    }
    ret.Images = make(map[string]types.ImageDef, 0)
    for k, v := range src.Images {
        ret.Images[k] = v
    }

    return &ret
}










