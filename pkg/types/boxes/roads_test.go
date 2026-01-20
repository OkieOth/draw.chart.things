package boxes_test

import (
	"testing"

	"github.com/okieoth/draw.chart.things/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestSortAscending(t *testing.T) {
	testData := []struct {
		input    []int
		expected []int
	}{
		{[]int{3, 1, 2}, []int{1, 2, 3}},
		{[]int{5, 4, 3, 2, 1}, []int{1, 2, 3, 4, 5}},
		{[]int{}, []int{}},
	}

	for _, d := range testData {
		types.SortAscending(d.input)
		assert.Equal(t, d.expected, d.input)
	}
}

func TestSortDescending(t *testing.T) {
	testData := []struct {
		input    []int
		expected []int
	}{
		{[]int{1, 2, 3}, []int{3, 2, 1}},
		{[]int{5, 4, 3, 2, 1}, []int{5, 4, 3, 2, 1}},
		{[]int{}, []int{}},
	}

	for _, d := range testData {
		types.SortDescending(d.input)
		assert.Equal(t, d.expected, d.input)
	}
}
