package ganttimpl_test

import (
	"testing"
	"time"

	"github.com/okieoth/draw.chart.things/pkg/ganttimpl"
	"github.com/okieoth/draw.chart.things/pkg/svgdrawing"
	"github.com/okieoth/draw.chart.things/pkg/types"
	"github.com/okieoth/draw.chart.things/pkg/types/gantt"
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
			g, err := types.LoadInputFromFile[gantt.Gantt](tt.inputFile)

			doc, err := ganttimpl.InitialLayoutGantt(g, textDimensionCalulator, tt.startDate, tt.endDate)
			require.Nil(t, err)
			require.NotNil(t, doc)
		})
	}
}

func TestDrawGanttFromFile(t *testing.T) {
	startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 3, 31, 0, 0, 0, 0, time.UTC)
	outputFile := "../../temp/simple_gantt.svg"
	inputFile := "../../resources/examples_gantt/simple.yaml"
	err := ganttimpl.DrawGanttFromFile(inputFile, outputFile, startDate, endDate, nil, "", "")
	require.Nil(t, err)
}

func TestLoadAdditionalEventsFileAndMerge(t *testing.T) {
	inputFile := "../../resources/examples_gantt/small.yaml"
	eventFile := "../../resources/examples_gantt/small_events.yaml"
	input, err := types.LoadInputFromFile[gantt.Gantt](inputFile)
	require.Nil(t, err)
	require.Equal(t, 1, len(input.Events))
	require.Equal(t, 2, len(input.Groups))
	input, err = ganttimpl.LoadAdditionalEventsFileAndMerge(input, eventFile)
	require.Nil(t, err)
	require.Equal(t, 4, len(input.Events))
}

func TestLoadAdditionalGroupFilesAndMerge(t *testing.T) {
	inputFile := "../../resources/examples_gantt/small.yaml"
	groupFiles := []string{"../../resources/examples_gantt/small_group1.yaml",
		"../../resources/examples_gantt/small_group1.yaml"}
	input, err := types.LoadInputFromFile[gantt.Gantt](inputFile)
	require.Nil(t, err)
	require.Equal(t, 1, len(input.Events))
	require.Equal(t, 2, len(input.Groups))
	input, err = ganttimpl.LoadAdditionalGroupFilesAndMerge(input, groupFiles)
	require.Nil(t, err)
	require.Equal(t, 1, len(input.Events))
	require.Equal(t, 4, len(input.Groups))
}

func TestTrimInputToDates(t *testing.T) {
	testData := []struct {
		inputFile         string
		startDate         string
		endDate           string
		groupEntriesCount []int
		eventCount        int
	}{
		{
			inputFile:         "../../resources/examples_gantt/simple.yaml",
			startDate:         "2020-01-01",
			endDate:           "2030-01-01",
			groupEntriesCount: []int{2, 4, 1, 3},
			eventCount:        5,
		},
		{
			inputFile:         "../../resources/examples_gantt/simple.yaml",
			startDate:         "2025-01-01",
			endDate:           "2025-03-31",
			groupEntriesCount: []int{2, 4, 1, 3},
			eventCount:        5,
		},
		{
			inputFile:         "../../resources/examples_gantt/simple.yaml",
			startDate:         "2025-02-01",
			endDate:           "2025-03-01",
			groupEntriesCount: []int{2, 3, 1, 2},
			eventCount:        2,
		},
	}
	for _, td := range testData {
		g, err := types.LoadInputFromFile[gantt.Gantt](td.inputFile)
		require.Nil(t, err)
		var start time.Time
		var end time.Time
		if td.startDate != "" {
			start, err = time.Parse("2006-01-02", td.startDate)
			require.Nil(t, err)
		}
		if td.endDate != "" {
			end, err = time.Parse("2006-01-02", td.endDate)
			require.Nil(t, err)
		}
		g = ganttimpl.TrimInputToDates(g, start, end)
		require.NotNil(t, g)
		require.Equal(t, len(td.groupEntriesCount), len(g.Groups))
		for i, count := range td.groupEntriesCount {
			require.Equal(t, count, len(g.Groups[i].Entries), "Group %d should have %d entries, but has %d", i, count, len(g.Groups[i].Entries))
		}
		require.Equal(t, td.eventCount, len(g.Events), "Expected %d events, but got %d", td.eventCount, len(g.Events))
	}
}

func TestIsDateToBeConsidered(t *testing.T) {
	testData := []struct {
		minDate   string
		maxDate   string
		startDate string
		endDate   string
		expected  bool
	}{
		{
			minDate:   "",
			maxDate:   "",
			startDate: "2020-01-01",
			endDate:   "2030-01-01",
			expected:  true,
		},
		{
			minDate:   "2020-01-01",
			maxDate:   "",
			startDate: "2020-01-01",
			endDate:   "2030-01-01",
			expected:  true,
		},
		{
			minDate:   "2019-11-30",
			maxDate:   "2019-12-31",
			startDate: "2020-01-01",
			endDate:   "2030-01-01",
			expected:  false,
		},
		{
			minDate:   "2030-01-02",
			maxDate:   "",
			startDate: "2020-01-01",
			endDate:   "2030-01-01",
			expected:  false,
		},
		{
			minDate:   "2022-01-02",
			maxDate:   "2023-01-02",
			startDate: "2020-01-01",
			endDate:   "2030-01-01",
			expected:  true,
		},
	}
	for _, td := range testData {
		var minDate, maxDate *time.Time
		if td.minDate != "" {
			parsedMinDate, err := time.Parse("2006-01-02", td.minDate)
			require.Nil(t, err)
			minDate = &parsedMinDate
		}
		if td.maxDate != "" {
			parsedMaxDate, err := time.Parse("2006-01-02", td.maxDate)
			require.Nil(t, err)
			maxDate = &parsedMaxDate
		}
		startDate, err := time.Parse("2006-01-02", td.startDate)
		require.Nil(t, err)
		endDate, err := time.Parse("2006-01-02", td.endDate)
		require.Nil(t, err)
		ret := ganttimpl.IsDateToBeConsidered(minDate, maxDate, startDate, endDate)
		require.Equal(t, td.expected, ret, "Expected %v but got %v for minDate: %v, maxDate: %v, startDate: %v, endDate: %v")
	}
}
