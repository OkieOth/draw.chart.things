package types

// Attention, this file is generated. Manual changes get lost with the next
// run of the code generation.
// created by yacg (template: golang_types.mako v1.1.0)

import (
    "time"
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





type Group struct {

    // Text to name the group
    Name string  `yaml:"name"`

    // Start date when the group is active
    Start *time.Time  `yaml:"start,omitempty"`

    // Start date when the group is active
    End *time.Time  `yaml:"end,omitempty"`

    // List of entries in the group
    Entries []Entry  `yaml:"entries,omitempty"`
}

func NewGroup() *Group {
        return &Group{
            Entries: make([]Entry, 0),
        }
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








type GanttFormat struct {

    Font *FontDef  `yaml:"font,omitempty"`

    GroupFont *FontDef  `yaml:"groupFont,omitempty"`

    EntryFont *FontDef  `yaml:"entryFont,omitempty"`

    EventFont *FontDef  `yaml:"eventFont,omitempty"`

    EntryFill *FillDef  `yaml:"entryFill,omitempty"`
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

    // List of resources assigned to the entry
    Resources []string  `yaml:"resources,omitempty"`
}

func NewEntry() *Entry {
        return &Entry{
            Resources: make([]string, 0),
        }
}





type RelativeStart struct {

    // Optional, name of the group where the start of the entry depends on
    GroupRef *string  `yaml:"groupRef,omitempty"`

    // Name of the entry, that end defines the start of this entry
    EntryRef string  `yaml:"entryRef"`
}






type EntryRef struct {

    // Name of the group where the start of the entry depends on
    GroupRef *string  `yaml:"groupRef,omitempty"`

    // Name of the entry, that end defines the start of this entry
    EntryRef string  `yaml:"entryRef"`
}





