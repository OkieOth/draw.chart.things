package ganttimpl_test

import (
	"testing"
	"time"

	"github.com/okieoth/draw.chart.things/pkg/ganttimpl"
	"github.com/okieoth/draw.chart.things/pkg/svgdrawing"
	"github.com/okieoth/draw.chart.things/pkg/types"
	"github.com/stretchr/testify/require"
)

func TestInitialLayoutGantt(t *testing.T) {
	tests := []struct {
		name      string
		inputFile string
		startDate time.Time
		endDate   time.Time
	}{
		{
			name:      "simple",
			inputFile: "../../resources/examples_gantt/simple.yaml",
			startDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			endDate:   time.Date(2025, 3, 31, 0, 0, 0, 0, time.UTC),
		},
	}

	textDimensionCalulator := svgdrawing.NewSvgTextDimensionCalculator()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g, err := types.LoadInputFromFile[types.Gantt](tt.inputFile)

			doc, err := ganttimpl.InitialLayoutGantt(g, textDimensionCalulator, tt.startDate, tt.endDate)
			require.Nil(t, err)
			require.NotNil(t, doc)
		})
	}
}
