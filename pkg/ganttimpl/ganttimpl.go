package ganttimpl

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/okieoth/draw.chart.things/pkg/svgdrawing"
	"github.com/okieoth/draw.chart.things/pkg/types"
	"github.com/okieoth/draw.chart.things/pkg/types/gantt"
)

func DrawGanttFromFile(inputFile, outputFile string, startDate, endDate time.Time, additionalGroupFiles []string, additionalEventsFile string, title string) error {
	input, err := types.LoadInputFromFile[gantt.Gantt](inputFile)
	if err != nil {
		return err
	}

	if title != "" {
		input.Title = title
	}

	if additionalGroupFiles != nil && len(additionalGroupFiles) > 0 {
		input, err = loadAdditionalGroupFilesAndMerge(input, additionalGroupFiles)
		if err != nil {
			return fmt.Errorf("failed to load additional group files: %w", err)
		}
	}

	if additionalEventsFile != "" {
		input, err = loadAdditionalEventsFileAndMerge(input, additionalEventsFile)
		if err != nil {
			return fmt.Errorf("failed to load additional events file: %w", err)
		}
	}

	input = trimInputToDates(input, startDate, endDate)

	textDimensionCalulator := svgdrawing.NewSvgTextDimensionCalculator()

	doc, err := InitialLayoutGantt(input, textDimensionCalulator, startDate, endDate)
	if err != nil {
		return err
	}
	output, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer output.Close()
	drawing := svgdrawing.NewDrawing(output)

	defaultFormat := doc.Formats["default"]
	groupHeight := calcGroupHeight(doc, textDimensionCalulator, defaultFormat.GroupFont)
	eventTxtHeight := calcEventTextHeight(doc, textDimensionCalulator, defaultFormat.EventFont)
	yOffset := eventTxtHeight
	if doc.Title != "" {
		_, h := textDimensionCalulator.Dimensions(doc.Title, defaultFormat.Font)
		yOffset += h + (2 * doc.GlobalPadding)
	}

	initDocWidth(doc, startDate, endDate)
	drawing.Start(doc.Title, groupHeight+yOffset+groupHeight+eventTxtHeight, *doc.Width)
	drawing.DrawText(doc.Title, 0, doc.GlobalPadding, *doc.Width, defaultFormat.Font)
	w, calendarHeight, _ := DrawCalendar(startDate, endDate, drawing, *doc.GroupNameWidth, yOffset, groupHeight)
	yOffset += 20
	drawGroupLines(doc, drawing, yOffset, w, defaultFormat.GroupFont, textDimensionCalulator)
	drawGroupEntries(doc, drawing, yOffset+2, defaultFormat.EntryFont, textDimensionCalulator)
	drawEvents(doc, drawing, yOffset, calendarHeight, eventTxtHeight, defaultFormat.EventFont, textDimensionCalulator)

	drawing.Done()
	output.Close()

	return nil
}

func initDocWidth(doc *gantt.GanttDocument, startDate, endDate time.Time) {
	if startDate.After(endDate) {
		return
	}
	// Calculate the number of days between startDate and endDate
	days := int(endDate.Sub(startDate).Hours() / 24)
	calendarWidth := days * 10 // assuming each day is 10 units wide
	w := calendarWidth + doc.GlobalPadding*2 + *doc.GroupNameWidth
	doc.Width = &w
}

func findEntry(doc *gantt.GanttDocument, ref *gantt.DocEntryRef) *gantt.DocGanttEntry {
	if ref == nil || ref.GroupRef == nil || ref.EntryRef == nil {
		return nil
	}
	for _, group := range doc.Groups {
		lowerGroupName := strings.ToLower(group.Name)
		lowerRefGroup := strings.ToLower(*ref.GroupRef)
		if lowerGroupName == lowerRefGroup {
			for _, entry := range group.Entries {
				lowerEntryName := strings.ToLower(entry.Name)
				lowerRefEntry := strings.ToLower(*ref.EntryRef)
				if lowerEntryName == lowerRefEntry {
					return &entry
				}
			}
		}
	}
	return nil
}

func drawEvents(doc *gantt.GanttDocument, drawing *svgdrawing.SvgDrawing, yOffset, calendarHeight, eventTxtHeight int, format *types.FontDef, c types.TextDimensionCalculator) {
	startDate := doc.StartDate
	endDate := doc.EndDate.Add(time.Hour * 24) // extend end date by one day to include the last day

	for _, event := range doc.Events {
		// TODO
		if event.Date.Equal(*startDate) || event.Date.Equal(endDate) || (event.Date.After(*startDate) && event.Date.Before(endDate)) {
			// Calculate x position based on event date
			xOffset := int(event.Date.Sub(*startDate).Hours()/24)*types.GlobalDayWidth + (types.GlobalDayWidth / 2) // center the event on the day
			// // Draw the event text
			// drawing.DrawText(event.Text, doc.GlobalPadding+xOffset, yOffset, *doc.GroupNameWidth, format)
			// // Draw a line for the event
			lineColor := "red"
			lineWidth := 1.0
			lineFormat := types.LineDef{
				Color: &lineColor,
				Width: &lineWidth,
			}
			drawing.DrawLine(*doc.GroupNameWidth+xOffset, yOffset-20, *doc.GroupNameWidth+xOffset, calendarHeight, lineFormat)
			// Draw small circles for the references
			for _, ref := range event.EntryRefs {
				foundEntry := findEntry(doc, &ref)
				if foundEntry != nil {
					circleColor := "red"
					circleRadius := 2
					drawing.DrawSolidCircle(*doc.GroupNameWidth+xOffset, foundEntry.Y+yOffset+types.GlobalGanttEntryHeight-2, circleRadius, circleColor)
				}
				// entry := doc.Groups[ref.GroupRef].Entries[ref.EntryRef]
			}
			// Print top text
			drawing.DrawVerticalText(event.Text, *doc.GroupNameWidth+xOffset-(format.Size/2)-2, yOffset-25, eventTxtHeight, format)
			// Print bottom text
			format.Anchor = types.FontDefAnchorEnum_right
			drawing.DrawVerticalText(event.Text, *doc.GroupNameWidth+xOffset-(format.Size/2)-2, calendarHeight+3, eventTxtHeight, format)
			format.Anchor = types.FontDefAnchorEnum_left
		}
	}
}

func drawGroupEntries(doc *gantt.GanttDocument, drawing *svgdrawing.SvgDrawing, yOffset int, format *types.FontDef, c types.TextDimensionCalculator) {
	lineWidth := 0.5
	borderOpacity := 1.0
	borderColor := "black"
	startDate := doc.StartDate
	endDate := doc.EndDate.Add(time.Hour * 24) // extend end date by one day to include the last day
	for _, group := range doc.Groups {
		for _, entry := range group.Entries {
			f := types.LineDef{
				Color:   &borderColor,
				Width:   &lineWidth,
				Opacity: &borderOpacity,
			}
			if entry.Format.EntryFill != nil && entry.Format.EntryFill.Color != nil {
				f.Color = entry.Format.EntryFill.Color
			}
			// TODO - calc start x and width based on entry start and end dates
			entryStart := startDate
			entryEnd := endDate
			if entry.Start != nil && entry.Start.After(*doc.StartDate) {
				entryStart = entry.Start
			}
			if entry.End != nil && entry.End.Before(*doc.EndDate) {
				entryEnd = *entry.End
			}
			daysDiff := int(entryEnd.Sub(*entryStart).Hours() / 24)
			width := daysDiff * types.GlobalDayWidth // assuming each day is 10 units wide
			xOffset := 0
			if entryStart != startDate {
				xOffset = int(entryStart.Sub(*doc.StartDate).Hours()/24) * types.GlobalDayWidth
			}
			textFormat := format
			if entry.Format != nil && entry.Format.EntryFont != nil {
				textFormat = entry.Format.EntryFont
			}

			opacity := 0.9
			if entry.Format != nil && entry.Format.EntryFill != nil && entry.Format.EntryFill.Opacity != nil {
				opacity = *entry.Format.EntryFill.Opacity
			}

			fill := types.FillDef{
				Color:   entry.Format.EntryFill.Color,
				Opacity: &opacity}

			rf := types.RectWithTextFormat{
				Padding:     0,
				FontCaption: *textFormat,
				Border:      &f,
				Fill:        &fill,
			}
			drawing.DrawRectWithText("", entry.Name, "", "", *doc.GroupNameWidth+xOffset, entry.Y+yOffset, width, types.GlobalGanttEntryHeight, rf)
		}
	}
}

func drawGroupLines(doc *gantt.GanttDocument, drawing *svgdrawing.SvgDrawing, yOffset, width int, format *types.FontDef, c types.TextDimensionCalculator) {
	currentY := yOffset
	lineColor := "lightgrey"
	lineWidth := 2.0
	lineFormat := types.LineDef{
		Color: &lineColor, // dark grey
		Width: &lineWidth,
	}
	lastY := currentY
	drawing.DrawLine(doc.GlobalPadding, currentY, *doc.GroupNameWidth+width, currentY, lineFormat)
	for _, group := range doc.Groups {
		if group.Name != "" {
			currentY += group.Height
			drawing.DrawLine(doc.GlobalPadding, currentY, *doc.GroupNameWidth+width, currentY, lineFormat)
			_, _, lines := c.SplitTxt(group.Name, format)
			y := lastY + 5
			for _, l := range lines {
				drawing.DrawText(l.Text, doc.GlobalPadding, y, *doc.GroupNameWidth, format)
				y += l.Height
			}
			lastY = currentY
		}
	}
}

func calcEventTextHeight(doc *gantt.GanttDocument, c types.TextDimensionCalculator, format *types.FontDef) int {
	if len(doc.Events) == 0 {
		return 0
	}
	height := 0
	for _, event := range doc.Events {
		if event.Text != "" {
			w, _ := c.Dimensions(event.Text, format)
			if w > height {
				height = w
			}
		}
	}
	return height
}

func calcGroupHeight(doc *gantt.GanttDocument, c types.TextDimensionCalculator, format *types.FontDef) int {
	if len(doc.Groups) == 0 {
		return 0
	}
	height := 0
	maxWidth := 0
	for i, group := range doc.Groups {
		if group.Name != "" {
			w, h := c.Dimensions(group.Name, format)
			curTaskY := height + 1
			entriesHeight := 1
			for j, _ := range group.Entries {
				doc.Groups[i].Entries[j].Y = curTaskY
				curTaskY += types.GlobalGanttEntryHeight + 2
				entriesHeight += types.GlobalGanttEntryHeight + 2
			}
			entriesHeight += 4 // add padding for aesthetics
			effectiveHeight := h + (doc.GlobalPadding)
			if entriesHeight > effectiveHeight {
				effectiveHeight = entriesHeight
			}
			height += effectiveHeight
			doc.Groups[i].Height = effectiveHeight
			if w > maxWidth {
				maxWidth = w
			}
		}
	}
	maxWidth += doc.GlobalPadding * 2
	doc.GroupNameWidth = &maxWidth
	return height // add padding for aesthetics
}

func InitialLayoutGantt(b *gantt.Gantt, c types.TextDimensionCalculator, startDate, endDate time.Time) (*gantt.GanttDocument, error) {
	doc := gantt.DocumentFromGantt(b, startDate, endDate)
	// doc.Boxes.X = doc.GlobalPadding
	// doc.Boxes.Y = doc.GlobalPadding
	// doc.Boxes.InitDimensions(c, doc.GlobalPadding, doc.MinBoxMargin)
	// doc.Width = doc.Boxes.Width + doc.GlobalPadding*2
	// doc.Height = doc.Boxes.Height + doc.GlobalPadding*2
	// if doc.Title != "" {
	// 	defaultFormat := doc.Formats["default"] // risky doesn't check if default exists
	// 	w, h := c.Dimensions(doc.Title, &defaultFormat.FontCaption)
	// 	doc.Height += h + (2 * doc.GlobalPadding)
	// 	if w > doc.Width {
	// 		doc.Width = w + (2 * doc.GlobalPadding)
	// 	}
	// }
	// doc.Boxes.Center()
	// doc.Height = doc.AdjustDocHeight(&doc.Boxes, 0) + doc.GlobalPadding
	return doc, nil
}

func DrawCalendar(startDate, endDate time.Time, drawing *svgdrawing.SvgDrawing, xOffset, yOffset, length int) (int, int, error) {
	if startDate.After(endDate) {
		return 0, 0, fmt.Errorf("end date must be after start date")
	}
	yStart := yOffset
	currentX := xOffset
	lineStyle := types.LineDefStyleEnum_solid
	lineWidth := 1.0
	lineColor := "lightgrey"
	weekendColor := "#f0f0f0" // light grey for weekends
	lineFormat := types.LineDef{
		Width: &lineWidth,
		Style: &lineStyle,
		Color: &lineColor,
	}
	weekendBoxFormat := types.LineDef{
		Width: &lineWidth,
		Style: &lineStyle,
		Color: &weekendColor,
	}
	dayFormat := types.FontDef{
		Size:              7,
		Color:             "grey",
		Font:              "sans-serif",
		Anchor:            types.FontDefAnchorEnum_left,
		MaxLenBeforeBreak: 20,
	}
	monthFormat := types.FontDef{
		Size:              10,
		Color:             "grey",
		Font:              "sans-serif",
		MaxLenBeforeBreak: 100,
	}
	dayWidth := types.GlobalDayWidth
	monthStartOffset := 20
	lastMonthX := currentX

	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		if d.Weekday() == time.Sunday || d.Weekday() == time.Saturday {
			drawing.DrawSolidRect(currentX, yStart+monthStartOffset, dayWidth, length, weekendBoxFormat)
		}
		if d.Day() == 1 {
			// new month
			width := currentX - lastMonthX
			if width > 10 {
				lastMonthDay := d.AddDate(0, 0, -1)
				monthStr := lastMonthDay.Format("01-2006")
				drawing.DrawText(monthStr, lastMonthX, yStart+monthStartOffset-20, width, &monthFormat)
				drawing.DrawText(monthStr, lastMonthX, yStart+monthStartOffset+length+10, width, &monthFormat)
			}
			drawing.DrawLine(currentX, yStart, currentX, yStart+length+(2*monthStartOffset), lineFormat)
			lastMonthX = currentX
		} else {
			drawing.DrawLine(currentX, yStart+monthStartOffset, currentX, yStart+monthStartOffset+length, lineFormat)
		}
		dayStr := fmt.Sprintf("%02d", d.Day())
		drawing.DrawText(dayStr, currentX+1, yStart+monthStartOffset-9, dayWidth, &dayFormat)
		drawing.DrawText(dayStr, currentX+1, yStart+monthStartOffset+length+1, dayWidth, &dayFormat)
		currentX += dayWidth
	}
	drawing.DrawLine(currentX, yStart+monthStartOffset, currentX, yStart+monthStartOffset+length, lineFormat)
	width := currentX - lastMonthX
	height := yStart + monthStartOffset + length + 10
	if width > 10 {
		monthStr := endDate.Format("01-2006")
		drawing.DrawText(monthStr, lastMonthX, yStart+monthStartOffset-20, width, &monthFormat)
		drawing.DrawText(monthStr, lastMonthX, height, width, &monthFormat)
		height += (monthFormat.Size + types.GlobalMinBoxMargin) // add space for month label
	}
	return currentX - xOffset, height, nil
}

func loadAdditionalEventsFileAndMerge(input *gantt.Gantt, eventFile string) (*gantt.Gantt, error) {
	e, err := types.LoadInputFromFile[[]gantt.Event](eventFile)
	if err != nil {
		return input, err
	}
	input.Events = append(input.Events, *e...)
	return input, nil
}

func loadAdditionalGroupFilesAndMerge(input *gantt.Gantt, groupFiles []string) (*gantt.Gantt, error) {
	for _, groupFile := range groupFiles {
		g, err := types.LoadInputFromFile[gantt.Group](groupFile)
		if err != nil {
			return input, err
		}
		input.Groups = append(input.Groups, *g)
	}
	return input, nil
}

func trimInputToDates(input *gantt.Gantt, startDate, endDate time.Time) *gantt.Gantt {
	groups := make([]gantt.Group, 0)
	for _, group := range input.Groups {
		if isDateToBeConsidered(group.Start, group.End, startDate, endDate) {
			g := gantt.NewGroup()
			g.End = group.End
			g.Start = group.Start
			g.Format = group.Format
			g.Name = group.Name
			g.Entries = filterEntriesByDates(group.Entries, startDate, endDate)
			groups = append(groups, *g)
		}
	}
	input.Groups = groups
	events := make([]gantt.Event, 0)
	for _, event := range input.Events {
		if event.Date.Equal(startDate) || event.Date.Equal(endDate) || (event.Date.After(startDate) && event.Date.Before(endDate)) {
			// Only keep events that are within the date range
			events = append(events, event)
		}
	}
	input.Events = events
	return input
}

func filterEntriesByDates(entries []gantt.Entry, startDate time.Time, endDate time.Time) []gantt.Entry {
	ret := make([]gantt.Entry, 0)
	for _, entry := range entries {
		if isDateToBeConsidered(entry.Start, entry.End, startDate, endDate) {
			var e gantt.Entry
			e.End = entry.End
			e.Start = entry.Start
			e.Format = entry.Format
			e.Name = entry.Name
			ret = append(ret, e)
		}
	}
	return ret
}

func isDateToBeConsidered(minDate, maxDate *time.Time, startDate, endDate time.Time) bool {
	if minDate != nil && minDate.After(endDate) {
		return false
	}
	if maxDate != nil && maxDate.Before(startDate) {
		return false
	}
	return true
}
