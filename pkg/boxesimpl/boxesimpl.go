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
	doc.Boxes.X = types.GlobalPadding
	doc.Boxes.Y = types.GlobalPadding
	doc.Boxes.InitDimensions(c)
	doc.Width = doc.Boxes.Width + types.GlobalPadding*2
	doc.Height = doc.Boxes.Height + types.GlobalPadding*2
	if doc.Title != "" {
		w, h := c.Text1Dimensions(doc.Title)
		doc.Height += h + (2 * types.GlobalPadding)
		if w > doc.Width {
			doc.Width = w + (2 * types.GlobalPadding)
		}
	}
	return doc, nil
}
