package sub_test

import (
	"testing"

	"github.com/okieoth/draw.chart.things/cmd/sub"
	"github.com/stretchr/testify/require"
)

func TestTruncBoxes(t *testing.T) {
	tests := []struct {
		inputFile      string
		outputFile     string
		depth          int
		expandedIds    []string
		blacklistedIds []string
	}{
		{
			inputFile:      "../../ui/data/boxes_random.yaml",
			outputFile:     "../../temp/boxes_random_truncated.yaml",
			depth:          2,
			expandedIds:    []string{"id_1_1"},
			blacklistedIds: []string{},
		},
	}
	for i, test := range tests {
		err := sub.TruncBoxes(test.inputFile, test.outputFile, test.depth, test.expandedIds, test.blacklistedIds)
		require.Nil(t, err, "error in TruncBoxes call, test:", i)
		require.FileExists(t, test.outputFile)
	}
}
