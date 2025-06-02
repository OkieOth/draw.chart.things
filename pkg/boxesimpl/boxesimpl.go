package boxesimpl

import (
	"fmt"

	"github.com/okieoth/draw.chart.things/pkg/svgdrawing"
	"github.com/okieoth/draw.chart.things/pkg/types"
	"github.com/okieoth/draw.chart.things/pkg/types/boxes"
)

func DrawBoxesFromFile(inputFile, outputFile string) error {

	layout, err := types.LoadInputFromFile[boxes.Boxes](inputFile)
	if err != nil {
		return err
	}
	textDimensionCalulator := svgdrawing.NewSvgTextDimensionCalculator()

	doc, err := InitialLayoutBoxes(layout, textDimensionCalulator)
	if err != nil {
		return err
	}
	fmt.Println(doc) // Dummy
	// TODO: Implement this

	// TODO: Draw the boxes
	return nil
}

func InitialLayoutBoxes(b *boxes.Boxes, c types.TextDimensionCalculator) (*boxes.BoxesDocument, error) {
	doc := DocumentFromBoxes(b)
	doc.Boxes.X = doc.GlobalPadding
	doc.Boxes.Y = doc.GlobalPadding
	doc.Boxes.InitDimensions(c, doc.GlobalPadding, doc.MinBoxMargin)
	doc.Width = doc.Boxes.Width + doc.GlobalPadding*2
	doc.Height = doc.Boxes.Height + doc.GlobalPadding*2
	if doc.Title != "" {
		defaultFormat := doc.Formats["default"] // risky doesn't check if default exists
		w, h := c.Dimensions(doc.Title, &defaultFormat.FontCaption)
		doc.Height += h + (2 * doc.GlobalPadding)
		if w > doc.Width {
			doc.Width = w + (2 * doc.GlobalPadding)
		}
	}
	doc.Boxes.Center()
	doc.Height = doc.AdjustDocHeight(&doc.Boxes, 0) + doc.GlobalPadding
	return doc, nil
}
