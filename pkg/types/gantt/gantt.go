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
type Gantt struct {

    // Title of the document
    Title string  `yaml:"title"`

    // List of groups to be displayed in the gantt chart
    Groups []Group  `yaml:"groups,omitempty"`

    // List of events to be displayed in the gantt chart
    Events []Event  `yaml:"events,omitempty"`

    // Map of formats available to be used in the gantt chart
    Formats map[string]GanttFormat  `yaml:"formats,omitempty"`
}

func NewGantt() *Gantt {
    return &Gantt{
        Groups: make([]Group, 0),
        Events: make([]Event, 0),
        Formats: make(map[string]GanttFormat, 0),
    }
}

func CopyGantt(src *Gantt) *Gantt {
    if src == nil {
        return nil
    }
    var ret Gantt
    ret.Title = src.Title
    ret.Groups = make([]Group, 0)
    for _, e := range src.Groups {
        ret.Groups = append(ret.Groups, e)
    }
    ret.Events = make([]Event, 0)
    for _, e := range src.Events {
        ret.Events = append(ret.Events, e)
    }
    ret.Formats = make(map[string]GanttFormat, 0)
    for k, v := range src.Formats {
        ret.Formats[k] = v
    }

    return &ret
}





type Group struct {

    // Text to name the group
    Name string  `yaml:"name"`

    // Start date when the group is active
    Start *time.Time  `yaml:"start,omitempty"`

    // Start date when the group is active
    End *time.Time  `yaml:"end,omitempty"`

    // List of entries in the group
    Entries []Entry  `yaml:"entries,omitempty"`

    // Optional reference to the format to be used for the group
    Format *string  `yaml:"format,omitempty"`
}

func NewGroup() *Group {
    return &Group{
        Entries: make([]Entry, 0),
    }
}

func CopyGroup(src *Group) *Group {
    if src == nil {
        return nil
    }
    var ret Group
    ret.Name = src.Name
    ret.Start = src.Start
    ret.End = src.End
    ret.Entries = make([]Entry, 0)
    for _, e := range src.Entries {
        ret.Entries = append(ret.Entries, e)
    }
    ret.Format = src.Format

    return &ret
}





/* an event that is displayed in the gantt chart and has references to specific
    group entries
*/
type Event struct {

    // Date of the event
    Date time.Time  `yaml:"date"`

    // Text to display for the event
    Text string  `yaml:"text"`

    // Description of the event
    Description *string  `yaml:"description,omitempty"`

    // List of references to entries in groups that are affected by the event
    EntryRefs []EntryRef  `yaml:"entryRefs,omitempty"`
}

func NewEvent() *Event {
    return &Event{
        EntryRefs: make([]EntryRef, 0),
    }
}

func CopyEvent(src *Event) *Event {
    if src == nil {
        return nil
    }
    var ret Event
    ret.Date = src.Date
    ret.Text = src.Text
    ret.Description = src.Description
    ret.EntryRefs = make([]EntryRef, 0)
    for _, e := range src.EntryRefs {
        ret.EntryRefs = append(ret.EntryRefs, e)
    }

    return &ret
}








type GanttFormat struct {

    Font *types.FontDef  `yaml:"font,omitempty"`

    GroupFont *types.FontDef  `yaml:"groupFont,omitempty"`

    EntryFont *types.FontDef  `yaml:"entryFont,omitempty"`

    EventFont *types.FontDef  `yaml:"eventFont,omitempty"`

    EntryFill *types.FillDef  `yaml:"entryFill,omitempty"`
}


func CopyGanttFormat(src *GanttFormat) *GanttFormat {
    if src == nil {
        return nil
    }
    var ret GanttFormat
    ret.Font = types.CopyFontDef(src.Font)
    ret.GroupFont = types.CopyFontDef(src.GroupFont)
    ret.EntryFont = types.CopyFontDef(src.EntryFont)
    ret.EventFont = types.CopyFontDef(src.EventFont)
    ret.EntryFill = types.CopyFillDef(src.EntryFill)

    return &ret
}





type Entry struct {

    // Text to name the entry
    Name string  `yaml:"name"`

    // Start date when the entry is active
    Start *time.Time  `yaml:"start,omitempty"`

    // Start, relative to the start of another entry from a group
    StartsAfter *RelativeStart  `yaml:"startsAfter,omitempty"`

    // End date when the entry is active
    End *time.Time  `yaml:"end,omitempty"`

    // Duration of the entry in days
    Duration *int  `yaml:"duration,omitempty"`

    // Description of the entry
    Description *string  `yaml:"description,omitempty"`

    // List of references to entries in other groups
    References []EntryRef  `yaml:"references,omitempty"`

    // Optional reference to the format to be used for the entry
    Format *string  `yaml:"format,omitempty"`
}

func NewEntry() *Entry {
    return &Entry{
        References: make([]EntryRef, 0),
    }
}

func CopyEntry(src *Entry) *Entry {
    if src == nil {
        return nil
    }
    var ret Entry
    ret.Name = src.Name
    ret.Start = src.Start
    ret.StartsAfter = CopyRelativeStart(src.StartsAfter)
    ret.End = src.End
    ret.Duration = src.Duration
    ret.Description = src.Description
    ret.References = make([]EntryRef, 0)
    for _, e := range src.References {
        ret.References = append(ret.References, e)
    }
    ret.Format = src.Format

    return &ret
}





type RelativeStart struct {

    // Optional, name of the group where the start of the entry depends on
    GroupRef *string  `yaml:"groupRef,omitempty"`

    // Name of the entry, that end defines the start of this entry
    EntryRef string  `yaml:"entryRef"`
}


func CopyRelativeStart(src *RelativeStart) *RelativeStart {
    if src == nil {
        return nil
    }
    var ret RelativeStart
    ret.GroupRef = src.GroupRef
    ret.EntryRef = src.EntryRef

    return &ret
}





type EntryRef struct {

    // Name of the group where the start of the entry depends on
    GroupRef *string  `yaml:"groupRef,omitempty"`

    // Name of the entry, that end defines the start of this entry
    EntryRef *string  `yaml:"entryRef,omitempty"`
}


func CopyEntryRef(src *EntryRef) *EntryRef {
    if src == nil {
        return nil
    }
    var ret EntryRef
    ret.GroupRef = src.GroupRef
    ret.EntryRef = src.EntryRef

    return &ret
}




