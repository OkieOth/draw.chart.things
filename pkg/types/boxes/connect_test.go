package boxes_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/okieoth/draw.chart.things/pkg/types/boxes"
)

func TestGetRelevantPoints(t *testing.T) {
	testData := []struct {
		v1             int
		v2             int
		other          int
		verticalBorder bool
		expectedLen    int
		expected1X     int
		expected1Y     int
		expected2X     int
		expected2Y     int
	}{
		{20, 50, 100, true, 6, 100, 24, 100, 44},
		{50, 20, 100, true, 6, 100, 24, 100, 44},
		{20, 50, 100, false, 6, 24, 100, 44, 100},
		{50, 20, 100, false, 6, 24, 100, 44, 100},
	}
	for _, d := range testData {
		points := boxes.GetRelevantPoints(d.v1, d.v2, d.other, d.verticalBorder)
		assert.Equal(t, d.expected1X, points[0].X)
		assert.Equal(t, d.expected1Y, points[0].Y)
		assert.False(t, points[0].HasCollision)
		last := d.expectedLen - 1
		assert.Equal(t, d.expected2X, points[last].X)
		assert.Equal(t, d.expected2Y, points[last].Y)
		assert.False(t, points[last].HasCollision)
	}
}
