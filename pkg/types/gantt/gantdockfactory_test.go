package gantt_test

import (
	"testing"
	"time"

	"github.com/okieoth/draw.chart.things/pkg/types/gantt"
	"github.com/stretchr/testify/require"
)

func stringPointer(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func datePointer(s string) *time.Time {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return nil
	}
	return &t
}

func intPointer(i int) *int {
	if i == 0 {
		return nil
	}
	return &i
}

func TestCreateDocGanttGroup(t *testing.T) {
	var gg gantt.Group
	gg.Entries = []gantt.Entry{{
		Name:        "Test Entry",
		Start:       datePointer("2025-01-01"),
		End:         datePointer("2025-01-31"),
		Duration:    intPointer(30),
		Description: stringPointer("This is a test entry"),
		References: []gantt.EntryRef{
			{
				GroupRef: stringPointer("testGroup1"),
				EntryRef: stringPointer("testEntry1"),
			},
			{
				GroupRef: stringPointer("testGroup2"),
				EntryRef: stringPointer("testEntry2"),
			},
		},
	}}

	doc := gantt.NewGanttDocument()
	dgg := doc.CreateDocGanttGroup(&gg)
	require.Equal(t, "Test Entry", dgg.Entries[0].Name)
	require.Equal(t, datePointer("2025-01-01"), dgg.Entries[0].Start)
	require.Equal(t, datePointer("2025-01-31"), dgg.Entries[0].End)
	require.Equal(t, intPointer(30), dgg.Entries[0].Duration)
	require.Equal(t, stringPointer("This is a test entry"), dgg.Entries[0].Description)
	require.Len(t, dgg.Entries[0].References, 2)
	require.Equal(t, "testGroup1", *dgg.Entries[0].References[0].GroupRef)
	require.Equal(t, "testEntry1", *dgg.Entries[0].References[0].EntryRef)
	require.Equal(t, "testGroup2", *dgg.Entries[0].References[1].GroupRef)
	require.Equal(t, "testEntry2", *dgg.Entries[0].References[1].EntryRef)
}
