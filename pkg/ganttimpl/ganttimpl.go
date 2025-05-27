package ganttimpl

import (
	"fmt"

	types "github.com/okieoth/draw.chart.things/pkg/gantttypes"
	"github.com/okieoth/draw.chart.things/pkg/svgdrawing"
)

func DrawBoxesFromFile(inputFile, outputFile string) error {
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
