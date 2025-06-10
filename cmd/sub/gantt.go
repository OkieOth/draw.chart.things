package sub

import (
	"fmt"
	"time"

	"github.com/okieoth/draw.chart.things/pkg/ganttimpl"

	"github.com/spf13/cobra"
)

var startDate string
var endDate string
var groups []string
var title string

var GanttCmd = &cobra.Command{
	Use:   "gantt",
	Short: "Draws a gantt diagram from a YAML file",
	Long:  `Draws a gantt diagram with possible extensions from a given YAML file.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("TODO: I will draw boxes :D")
		start, err := time.Parse("2006-01-02", startDate)
		if err != nil {
			fmt.Println("Invalid start date format. Please use YYYY-MM-DD.")
			return
		}
		end, err := time.Parse("2006-01-02", endDate)
		if err != nil {
			fmt.Println("Invalid end date format. Please use YYYY-MM-DD.")
			return
		}
		err = ganttimpl.DrawGanttFromFile(From, Output, start, end, groups, title)
		if err != nil {
			fmt.Println("Error while drawing gantt: ", err)
		}
	},
}

func init() {
	initDefaultFlags(GanttCmd)
	GanttCmd.Flags().StringVarP(&startDate, "start", "s", "", "Start date for the gantt diagram (format: YYYY-MM-DD)")
	GanttCmd.Flags().StringVarP(&endDate, "end", "e", "", "End date for the gantt diagram (format: YYYY-MM-DD)")
	GanttCmd.Flags().StringVarP(&title, "title", "t", "", "Title for the diagram to create")
	GanttCmd.Flags().StringSliceVarP(&groups, "group", "g", make([]string, 0), "file paths to group files that will be merged with the input file")
	if err := GanttCmd.MarkFlagRequired("start"); err != nil {
		fmt.Println("Error marking start date flag as required:", err)
	}
	if err := GanttCmd.MarkFlagRequired("end"); err != nil {
		fmt.Println("Error marking end date flag as required:", err)
	}
}
