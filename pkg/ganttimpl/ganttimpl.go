package ganttimpl

import (
	"fmt"
	"time"

	"github.com/okieoth/draw.chart.things/pkg/svgdrawing"
	"github.com/okieoth/draw.chart.things/pkg/types"
)

func DrawGanttFromFile(inputFile, outputFile string) error {
	input, err := types.LoadInputFromFile[types.Gantt](inputFile)
	if err != nil {
		return err
	}
	textDimensionCalulator := svgdrawing.NewSvgTextDimensionCalculator()

	doc, err := InitialLayoutGantt(input, textDimensionCalulator)
	if err != nil {
		return err
	}
	fmt.Println(doc) // Dummy
	// TODO: Implement this

	// TODO: Draw the boxes
	return nil
}

func InitialLayoutGantt(b *types.Gantt, c types.TextDimensionCalculator) (*types.GanttDocument, error) {
	doc := types.DocumentFromGantt(b)
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

func DrawCalendar(startDate, endDate time.Time, drawing *svgdrawing.Drawing, xOffset, yOffset, length int) error {
	if startDate.After(endDate) {
		return fmt.Errorf("end date must be after start date")
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
	if width > 10 {
		monthStr := endDate.Format("01-2006")
		drawing.DrawText(monthStr, lastMonthX, yStart+monthStartOffset-20, width, &monthFormat)
	}
	return nil
}
