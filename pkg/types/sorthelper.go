package types

import (
	"sort"
)

func SortDescending(a []int) {
	sort.Sort(sort.Reverse(sort.IntSlice(a)))
}

func SortAscending(a []int) {
	sort.IntSlice(a).Sort()
}
