package sub

import (
	"fmt"

	"github.com/okieoth/draw.chart.things/pkg/ganttimpl"

	"github.com/spf13/cobra"
)

var GanttCmd = &cobra.Command{
	Use:   "gantt",
	Short: "Draws a gantt diagram from a YAML file",
	Long:  `Draws a gantt diagram with possible extensions from a given YAML file.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("TODO: I will draw boxes :D")
		err := ganttimpl.DrawGanttFromFile(From, Output)
		if err != nil {
			fmt.Println("Error while drawing gantt: ", err)
		}
	},
}

func init() {
	initDefaultFlags(GanttCmd)
}
