package gantt

// Attention, this file is generated. Manual changes get lost with the next
// run of the code generation.
// created by yacg (template: golang_types.mako v1.1.0)

import (
    "time"
    "github.com/okieoth/draw.chart.things/pkg/types"
)


/* Model to describe the input of gantt diagrams
*/
type GanttDocument struct {

    // Title of the document
    Title string  `yaml:"title"`

    // Width of the group name column, used to calculate the x position of the boxes
    GroupNameWidth *int  `yaml:"groupNameWidth,omitempty"`

    // Height of the document
    Height *int  `yaml:"height,omitempty"`

    // Start date of the document, used to calculate the x position of the boxes
    StartDate *time.Time  `yaml:"startDate,omitempty"`

    // End date of the document, used to calculate the x position of the boxes
    EndDate *time.Time  `yaml:"endDate,omitempty"`

    // Width of the document
    Width *int  `yaml:"width,omitempty"`

    // Padding used as default over the whole diagram
    GlobalPadding int  `yaml:"globalPadding"`

    // Minimum margin between boxes
    MinBoxMargin int  `yaml:"minBoxMargin"`

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

func CopyGanttDocument(src *GanttDocument) *GanttDocument {
    if src == nil {
        return nil
    }
    var ret GanttDocument
    ret.Title = src.Title
    ret.GroupNameWidth = src.GroupNameWidth
    ret.Height = src.Height
    ret.StartDate = src.StartDate
    ret.EndDate = src.EndDate
    ret.Width = src.Width
    ret.GlobalPadding = src.GlobalPadding
    ret.MinBoxMargin = src.MinBoxMargin
    ret.Groups = make([]DocGanttGroup, 0)
    for _, e := range src.Groups {
        ret.Groups = append(ret.Groups, e)
    }
    ret.Events = make([]DocGanttEvent, 0)
    for _, e := range src.Events {
        ret.Events = append(ret.Events, e)
    }
    ret.Formats = make(map[string]DocGanttFormat, 0)
    for k, v := range src.Formats {
        ret.Formats[k] = v
    }

    return &ret
}





type DocGanttGroup struct {

    // Text to name the group
    Name string  `yaml:"name"`

    // Start date when the group is active
    Start *time.Time  `yaml:"start,omitempty"`

    // Start date when the group is active
    End *time.Time  `yaml:"end,omitempty"`

    // List of entries in the group
    Entries []DocGanttEntry  `yaml:"entries,omitempty"`

    // Height of the group, used to calculate the y position of the boxes
    Height int  `yaml:"height"`

    Format *DocGanttFormat  `yaml:"format,omitempty"`
}

func NewDocGanttGroup() *DocGanttGroup {
    return &DocGanttGroup{
        Entries: make([]DocGanttEntry, 0),
    }
}

func CopyDocGanttGroup(src *DocGanttGroup) *DocGanttGroup {
    if src == nil {
        return nil
    }
    var ret DocGanttGroup
    ret.Name = src.Name
    ret.Start = src.Start
    ret.End = src.End
    ret.Entries = make([]DocGanttEntry, 0)
    for _, e := range src.Entries {
        ret.Entries = append(ret.Entries, e)
    }
    ret.Height = src.Height
    ret.Format = CopyDocGanttFormat(src.Format)

    return &ret
}





type DocGanttEvent struct {

    // Date of the event
    Date time.Time  `yaml:"date"`

    // Text to display for the event
    Text string  `yaml:"text"`

    // Description of the event
    Description *string  `yaml:"description,omitempty"`

    // List of references to entries in groups that are affected by the event
    EntryRefs []DocEntryRef  `yaml:"entryRefs,omitempty"`
}

func NewDocGanttEvent() *DocGanttEvent {
    return &DocGanttEvent{
        EntryRefs: make([]DocEntryRef, 0),
    }
}

func CopyDocGanttEvent(src *DocGanttEvent) *DocGanttEvent {
    if src == nil {
        return nil
    }
    var ret DocGanttEvent
    ret.Date = src.Date
    ret.Text = src.Text
    ret.Description = src.Description
    ret.EntryRefs = make([]DocEntryRef, 0)
    for _, e := range src.EntryRefs {
        ret.EntryRefs = append(ret.EntryRefs, e)
    }

    return &ret
}








type DocGanttFormat struct {

    Font *types.FontDef  `yaml:"font,omitempty"`

    GroupFont *types.FontDef  `yaml:"groupFont,omitempty"`

    EntryFont *types.FontDef  `yaml:"entryFont,omitempty"`

    EventFont *types.FontDef  `yaml:"eventFont,omitempty"`

    EntryFill *types.FillDef  `yaml:"entryFill,omitempty"`
}


func CopyDocGanttFormat(src *DocGanttFormat) *DocGanttFormat {
    if src == nil {
        return nil
    }
    var ret DocGanttFormat
    ret.Font = types.CopyFontDef(src.Font)
    ret.GroupFont = types.CopyFontDef(src.GroupFont)
    ret.EntryFont = types.CopyFontDef(src.EntryFont)
    ret.EventFont = types.CopyFontDef(src.EventFont)
    ret.EntryFill = types.CopyFillDef(src.EntryFill)

    return &ret
}





type DocEntryRef struct {

    // Name of the group where the start of the entry depends on
    GroupRef *string  `yaml:"groupRef,omitempty"`

    // Name of the entry, that end defines the start of this entry
    EntryRef *string  `yaml:"entryRef,omitempty"`
}


func CopyDocEntryRef(src *DocEntryRef) *DocEntryRef {
    if src == nil {
        return nil
    }
    var ret DocEntryRef
    ret.GroupRef = src.GroupRef
    ret.EntryRef = src.EntryRef

    return &ret
}





type DocGanttEntry struct {

    // X position of the entry, used to calculate the position in the gantt chart
    X int  `yaml:"x"`

    // Y position of the entry, used to calculate the position in the gantt chart
    Y int  `yaml:"y"`

    // Text to name the entry
    Name string  `yaml:"name"`

    // Start date when the entry is active
    Start *time.Time  `yaml:"start,omitempty"`

    // End date when the entry is active
    End *time.Time  `yaml:"end,omitempty"`

    // Duration of the entry in days
    Duration *int  `yaml:"duration,omitempty"`

    // Description of the entry
    Description *string  `yaml:"description,omitempty"`

    // List of references to entries in other groups
    References []DocEntryRef  `yaml:"references,omitempty"`

    Format *DocGanttFormat  `yaml:"format,omitempty"`
}

func NewDocGanttEntry() *DocGanttEntry {
    return &DocGanttEntry{
        References: make([]DocEntryRef, 0),
    }
}

func CopyDocGanttEntry(src *DocGanttEntry) *DocGanttEntry {
    if src == nil {
        return nil
    }
    var ret DocGanttEntry
    ret.X = src.X
    ret.Y = src.Y
    ret.Name = src.Name
    ret.Start = src.Start
    ret.End = src.End
    ret.Duration = src.Duration
    ret.Description = src.Description
    ret.References = make([]DocEntryRef, 0)
    for _, e := range src.References {
        ret.References = append(ret.References, e)
    }
    ret.Format = CopyDocGanttFormat(src.Format)

    return &ret
}




