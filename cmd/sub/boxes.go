package sub

import (
	"fmt"

	"github.com/okieoth/draw.chart.things/pkg/boxes"

	"github.com/spf13/cobra"
)

var BoxesCmd = &cobra.Command{
	Use:   "boxes",
	Short: "Draws boxes based on given layout",
	Long:  `Draws boxes and their connections and layouts them.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("TODO: I will draw boxes :D")
		err := boxes.DrawBoxesFromFile(From, Output)
		if err != nil {
			fmt.Println("Error while drawing boxes: ", err)
		}
	},
}

func init() {
	initDefaultFlags(BoxesCmd)
}
