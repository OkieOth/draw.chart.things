package types

// Attention, this file is generated. Manual changes get lost with the next
// run of the code generation.
// created by yacg (template: golang_types.mako v1.1.0)

import (
)


/* Object to define comments in the diagrams
*/
type Comment struct {

    // text of the comment
    Text string  `yaml:"text"`

    // optional number or a short text, displayed in the marker of that comment
    Label *string  `yaml:"label,omitempty"`

    // format to use to render this comment
    Format *string  `yaml:"format,omitempty"`
}


func CopyComment(src *Comment) *Comment {
    if src == nil {
        return nil
    }
    var ret Comment
    ret.Text = src.Text
    ret.Label = src.Label
    ret.Format = src.Format

    return &ret
}




