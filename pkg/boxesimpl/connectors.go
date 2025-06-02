package boxesimpl

import (
	"fmt"

	"github.com/okieoth/draw.chart.things/pkg/types/boxes"
)

type Point struct {
	X int
	Y int
}

type Line struct {
	Start Point
	End   Point
}

func initPossibleVerticalConnectors(doc *boxes.BoxesDocument) []Line {
	return nil // TODO
}

func initPossibleHorizontalConnectors(doc *boxes.BoxesDocument) []Line {
	return nil // TODO
}

func ConnectBoxex(doc *boxes.BoxesDocument) error {
	verticalConnectors := initPossibleVerticalConnectors(doc)
	horizontalConnectors := initPossibleHorizontalConnectors(doc)
	fmt.Println("verticalConnectors", verticalConnectors)     // Dummy
	fmt.Println("horizontalConnectors", horizontalConnectors) // Dummy
	return nil                                                // TODO
}
