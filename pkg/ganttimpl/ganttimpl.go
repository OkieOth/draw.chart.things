package ganttimpl

import (
	"fmt"
	"os"
	"time"

	"github.com/okieoth/draw.chart.things/pkg/svgdrawing"
	"github.com/okieoth/draw.chart.things/pkg/types"
	"github.com/okieoth/draw.chart.things/pkg/types/gantt"
)

func DrawGanttFromFile(inputFile, outputFile string, startDate, endDate time.Time) error {
	input, err := types.LoadInputFromFile[gantt.Gantt](inputFile)
	if err != nil {
		return err
	}
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
	initDocWidth(doc, startDate, endDate)
	drawing.Start(doc.Title, groupHeight+50, *doc.Width)
	w, _, _ := DrawCalendar(startDate, endDate, drawing, *doc.GroupNameWidth, 10, groupHeight)
	drawGroupLines(doc, drawing, 10+20, w, defaultFormat.GroupFont, textDimensionCalulator)
	// TODO - draw group caption
	// TODO - draw tasks
	// TODO - draw events
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

func drawGroupLines(doc *gantt.GanttDocument, drawing *svgdrawing.SvgDrawing, yOffset, width int, format *types.FontDef, c types.TextDimensionCalculator) {
	currentY := yOffset
	lineColor := "lightgrey"
	lineWidth := 2
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
			//drawing.DrawText(group.Name, x+5, y+height/2, width-10, &doc.Formats["default"].GroupFont)
		}
	}
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
			effectiveHeight := h + (2 * doc.GlobalPadding)
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
	lineWidth := 1
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
	dayWidth := 10
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
