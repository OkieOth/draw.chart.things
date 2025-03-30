package boxesimpl

import (
	"fmt"

	"github.com/okieoth/draw.chart.things/pkg/svgdrawing"
	"github.com/okieoth/draw.chart.things/pkg/types"
)

func DrawBoxesFromFile(inputFile, outputFile string) error {

	layout, err := types.LoadInputFromFile[types.Boxes](inputFile)
	if err != nil {
		return err
	}
	fmt.Println(layout) // Dummy
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

func InitialLayoutBoxes(b *types.Boxes, c types.TextDimensionCalculator) (*types.BoxesDocument, error) {
	doc := types.DocumentFromBoxes(b)
	doc.Boxes.X = doc.GlobalPadding
	doc.Boxes.Y = doc.GlobalPadding
	doc.Boxes.InitDimensions(c, doc.GlobalPadding, doc.MinBoxMargin)
	doc.Width = doc.Boxes.Width + doc.GlobalPadding*2
	doc.Height = doc.Boxes.Height + doc.GlobalPadding*2
	if doc.Title != "" {
		defaultFormat := doc.Formats["default"] // risky doesn't check if default exists
		w, h := c.Text1Dimensions(doc.Title, &defaultFormat.FontCaption)
		doc.Height += h + (2 * doc.GlobalPadding)
		if w > doc.Width {
			doc.Width = w + (2 * doc.GlobalPadding)
		}
	}
	doc.Boxes.Center()
	return doc, nil
}
