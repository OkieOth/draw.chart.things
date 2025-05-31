package types

// Attention, this file is generated. Manual changes get lost with the next
// run of the code generation.
// created by yacg (template: golang_types.mako v1.1.0)

import (
    "time"
)


/* Model to describe the input of gantt diagrams
*/
type GanttDocument struct {

    // Title of the document
    Title string  `yaml:"title"`

    // Height of the document
    Height *int  `yaml:"height,omitempty"`

    // Start date of the document, used to calculate the x position of the boxes
    StartDate *time.Time  `yaml:"startDate,omitempty"`

    // End date of the document, used to calculate the x position of the boxes
    EndDate *time.Time  `yaml:"endDate,omitempty"`

    // Width of the document
    Width *int  `yaml:"width,omitempty"`

    // Padding used as default over the whole diagram
    GlobalPadding *int  `yaml:"globalPadding,omitempty"`

    // Minimum margin between boxes
    MinBoxMargin *int  `yaml:"minBoxMargin,omitempty"`

    // List of groups to be displayed in the gantt chart
    Groups []DocGanttGroup  `yaml:"groups,omitempty"`

    // List of events to be displayed in the gantt chart
    Events []DocGanttEvent  `yaml:"events,omitempty"`

    // Map of formats available to be used in the boxes
    Formats map[string]DocGanttFormat  `yaml:"formats,omitempty"`
}

func NewGanttDocument() *GanttDocument {
        return &GanttDocument{
            Groups: make([]DocGanttGroup, 0),
            Events: make([]DocGanttEvent, 0),
            Formats: make(map[string]DocGanttFormat, 0),
        }
}





type DocGanttGroup struct {
}






type DocGanttEvent struct {
}









type DocGanttFormat struct {

    Font *FontDef  `yaml:"font,omitempty"`

    GroupFont *FontDef  `yaml:"groupFont,omitempty"`

    EntryFont *FontDef  `yaml:"entryFont,omitempty"`

    EventFont *FontDef  `yaml:"eventFont,omitempty"`

    EntryFill *FillDef  `yaml:"entryFill,omitempty"`
}





