package gantt

import (
	"time"

	"github.com/okieoth/draw.chart.things/pkg/types"
)

func CreateDocGanttGroup(gg *Group, startDate, endDate time.Time) *DocGanttGroup {
	g := NewDocGanttGroup()
	g.Name = gg.Name
	g.Start = gg.Start
	g.End = gg.End
	if gg.Entries != nil {
		// TODO
	}
	return g
}

func DocumentFromGantt(g *Gantt, startDate, endDate time.Time) *GanttDocument {
	doc := NewGanttDocument()
	doc.Title = g.Title
	doc.GlobalPadding = types.GlobalPadding
	for _, group := range g.Groups {
		if group.Name == "" {
			// groups without name are ignored
			continue
		}
		var dg *DocGanttGroup
		if group.Start == nil && group.End == nil {
			// group is always present
			CreateDocGanttGroup(&group, startDate, endDate)
		} else if group.Start == nil && group.End != nil && group.End.After(startDate) {
			// group ends after the start date
			CreateDocGanttGroup(&group, startDate, endDate)
		} else if group.Start != nil && group.Start.Before(endDate) && group.End == nil {
			// group ends after the start date
			CreateDocGanttGroup(&group, startDate, endDate)
		} else if group.Start.Before(endDate) || group.End.After(startDate) {
			// the group is in the range of the given times
			CreateDocGanttGroup(&group, startDate, endDate)
		}
		if dg != nil {
			doc.Groups = append(doc.Groups, *dg)
		}
	}
	return doc
}

func (g *GanttDocument) DrawGantt(drawingImpl types.Drawing) error {
	//yOffset := g.GlobalPadding
	//DrawCalendar(*g.StartDate, *g.EndDate, drawingImpl, *g.GroupNameWidth, yOffset, g.Height)
	return nil
}
