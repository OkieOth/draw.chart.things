package boxes

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
	// TODO: Implement this
	return nil
}
