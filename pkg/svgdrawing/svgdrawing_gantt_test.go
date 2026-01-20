package svgdrawing_test

import (
	"os"
	"time"

	"testing"

	"github.com/okieoth/draw.chart.things/pkg/ganttimpl"
	"github.com/okieoth/draw.chart.things/pkg/svgdrawing"
	"github.com/stretchr/testify/require"
)

func TestSimpleCalendar(t *testing.T) {
	startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 3, 31, 0, 0, 0, 0, time.UTC)
	outputFile := "../../temp/SimpleCalendar.svg"
	output, err := os.Create(outputFile)
	require.Nil(t, err)
	drawing := svgdrawing.NewDrawing(output)
	drawing.Start("Test Calendar", 2000, 2000)
	ganttimpl.DrawCalendar(startDate, endDate, drawing, 10, 10, 500)
	drawing.Done()
	output.Close()

}
