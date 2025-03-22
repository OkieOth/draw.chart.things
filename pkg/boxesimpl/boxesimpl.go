package boxesimpl

import (
	"fmt"

	"github.com/okieoth/draw.chart.things/pkg/types"
)

func DrawBoxesFromFile(inputFile, outputFile string) error {

	layout, err := types.LoadInputFromFile[types.Boxes](inputFile)
	if err != nil {
		return err
	}
	fmt.Println(layout) // Dummy
	doc, err := InitialLayoutBoxes(layout)
	if err != nil {
		return err
	}
	fmt.Println(doc) // Dummy
	// TODO: Implement this

	// TODO: Draw the boxes
	return nil
}

func InitialLayoutBoxes(b *types.Boxes) (*types.BoxesDocument, error) {
	doc := types.DocumentFromBoxes(b)
	return doc, nil // TODO
}
