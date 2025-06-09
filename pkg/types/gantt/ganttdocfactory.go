package gantt

import (
	"time"

	"github.com/okieoth/draw.chart.things/pkg/types"
)

func (d *GanttDocument) CreateDocGanttGroup(gg *Group, startDate, endDate time.Time) *DocGanttGroup {
	g := NewDocGanttGroup()
	g.Name = gg.Name
	g.Start = gg.Start
	g.End = gg.End
	if gg.Entries != nil {
		for _, entry := range gg.Entries {
			dge := NewDocGanttEntry()
			dge.Name = entry.Name
			dge.Start = entry.Start
			dge.End = entry.End
			if entry.Format != nil {
				if f, found := d.Formats[*entry.Format]; found {
					dge.Format = &f
				} else {
					f = d.Formats["default"]
					dge.Format = &f
				}
			} else {
				f := d.Formats["default"]
				dge.Format = &f
			}
			g.Entries = append(g.Entries, *dge)
		}
	}
	return g
}

func (d *GanttDocument) initFormats(formats map[string]GanttFormat) {
	for name, format := range formats {
		f := DocGanttFormat{
			Font:      format.Font,
			GroupFont: format.GroupFont,
			EntryFont: format.EntryFont,
			EventFont: format.EventFont,
			EntryFill: format.EntryFill,
		}
		d.Formats[name] = f
	}
	if _, exists := d.Formats["default"]; !exists {
		font := types.InitFontDef(nil, "sans-serif", 10, false, false, 0)
		font.MaxLenBeforeBreak = 1000
		fg := types.InitFontDef(nil, "sans-serif", 10, false, false, 0)
		fg.Anchor = types.FontDefAnchorEnum_right
		fentry := types.InitFontDef(nil, "serif", 7, false, false, 0)
		fentry.MaxLenBeforeBreak = 1000
		fevent := types.InitFontDef(nil, "sans-serif", 10, false, false, 0)
		defaultFill := "#f0f0f0"
		fill := types.FillDef{
			Color: &defaultFill,
		}

		f := DocGanttFormat{
			Font:      &font,
			GroupFont: &fg,
			EntryFont: &fentry,
			EventFont: &fevent,
			EntryFill: &fill,
		}
		d.Formats["default"] = f
	}
}

func (d *GanttDocument) initEvents(events []Event, startDate, endDate time.Time) {
	for _, event := range events {
		dg := NewDocGanttEvent()
		dg.Date = event.Date
		dg.Text = event.Text
		dg.Description = event.Description
		if event.EntryRefs != nil {
			for _, ref := range event.EntryRefs {
				de := DocEntryRef{
					GroupRef: ref.GroupRef,
					EntryRef: ref.EntryRef,
				}
				dg.EntryRefs = append(dg.EntryRefs, de)
			}
		}
		d.Events = append(d.Events, *dg)
	}
}

func (d *GanttDocument) initGroups(groups []Group, startDate, endDate time.Time) {
	for _, group := range groups {
		if group.Name == "" {
			// groups without name are ignored
			continue
		}
		var dg *DocGanttGroup
		if group.Start == nil && group.End == nil {
			// group is always present
			dg = d.CreateDocGanttGroup(&group, startDate, endDate)
		} else if group.Start == nil && group.End != nil && group.End.After(startDate) {
			// group ends after the start date
			dg = d.CreateDocGanttGroup(&group, startDate, endDate)
		} else if group.Start != nil && group.Start.Before(endDate) && group.End == nil {
			// group ends after the start date
			dg = d.CreateDocGanttGroup(&group, startDate, endDate)
		} else if group.Start.Before(endDate) || group.End.After(startDate) {
			// the group is in the range of the given times
			dg = d.CreateDocGanttGroup(&group, startDate, endDate)
		}
		if dg != nil {
			d.Groups = append(d.Groups, *dg)
		}
	}
}

func DocumentFromGantt(g *Gantt, startDate, endDate time.Time) *GanttDocument {
	doc := NewGanttDocument()
	doc.Title = g.Title
	doc.GlobalPadding = types.GlobalPadding
	doc.StartDate = &startDate
	doc.EndDate = &endDate
	doc.initFormats(g.Formats)
	doc.initGroups(g.Groups, startDate, endDate)
	doc.initEvents(g.Events, startDate, endDate)
	return doc
}

func (g *GanttDocument) DrawGantt(drawingImpl types.Drawing) error {
	//yOffset := g.GlobalPadding
	//DrawCalendar(*g.StartDate, *g.EndDate, drawingImpl, *g.GroupNameWidth, yOffset, g.Height)
	return nil
}
