package types_test

import (
	"reflect"
	"testing"

	"github.com/okieoth/draw.chart.things/pkg/types"
)

func TestInitLayoutElemArray(t *testing.T) {
	tests := []struct {
		name     string
		input    []types.Layout
		expected []types.LayoutElement
	}{
		{
			name:     "empty input",
			input:    []types.Layout{},
			expected: []types.LayoutElement{},
		},
		{
			name: "single element",
			input: []types.Layout{
				{
					Id:         "1",
					Caption:    "Caption1",
					Text1:      "Text1",
					Text2:      "Text2",
					Vertical:   []types.Layout{},
					Horizontal: []types.Layout{},
				},
			},
			expected: []types.LayoutElement{
				{
					Id:         "1",
					Caption:    "Caption1",
					Text1:      "Text1",
					Text2:      "Text2",
					Vertical:   []types.LayoutElement{},
					Horizontal: []types.LayoutElement{},
				},
			},
		},
		{
			name: "nested elements",
			input: []types.Layout{
				{
					Id:      "1",
					Caption: "Caption1",
					Text1:   "Text1",
					Text2:   "Text2",
					Vertical: []types.Layout{
						{
							Id:         "2",
							Caption:    "Caption2",
							Text1:      "Text3",
							Text2:      "Text4",
							Vertical:   []types.Layout{},
							Horizontal: []types.Layout{},
						},
					},
					Horizontal: []types.Layout{},
				},
			},
			expected: []types.LayoutElement{
				{
					Id:      "1",
					Caption: "Caption1",
					Text1:   "Text1",
					Text2:   "Text2",
					Vertical: []types.LayoutElement{
						{
							Id:         "2",
							Caption:    "Caption2",
							Text1:      "Text3",
							Text2:      "Text4",
							Vertical:   []types.LayoutElement{},
							Horizontal: []types.LayoutElement{},
						},
					},
					Horizontal: []types.LayoutElement{},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := types.ExpInitLayoutElemArray(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("InitLayoutElemArray() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestInitLayoutElement(t *testing.T) {
	tests := []struct {
		name     string
		input    types.Layout
		expected types.LayoutElement
	}{
		{
			name: "empty layout",
			input: types.Layout{
				Id:         "",
				Caption:    "",
				Text1:      "",
				Text2:      "",
				Vertical:   []types.Layout{},
				Horizontal: []types.Layout{},
			},
			expected: types.LayoutElement{
				Id:         "",
				Caption:    "",
				Text1:      "",
				Text2:      "",
				Vertical:   []types.LayoutElement{},
				Horizontal: []types.LayoutElement{},
			},
		},
		{
			name: "single element",
			input: types.Layout{
				Id:         "1",
				Caption:    "Caption1",
				Text1:      "Text1",
				Text2:      "Text2",
				Vertical:   []types.Layout{},
				Horizontal: []types.Layout{},
			},
			expected: types.LayoutElement{
				Id:         "1",
				Caption:    "Caption1",
				Text1:      "Text1",
				Text2:      "Text2",
				Vertical:   []types.LayoutElement{},
				Horizontal: []types.LayoutElement{},
			},
		},
		{
			name: "nested elements",
			input: types.Layout{
				Id:      "1",
				Caption: "Caption1",
				Text1:   "Text1",
				Text2:   "Text2",
				Vertical: []types.Layout{
					{
						Id:         "2",
						Caption:    "Caption2",
						Text1:      "Text3",
						Text2:      "Text4",
						Vertical:   []types.Layout{},
						Horizontal: []types.Layout{},
					},
				},
				Horizontal: []types.Layout{},
			},
			expected: types.LayoutElement{
				Id:      "1",
				Caption: "Caption1",
				Text1:   "Text1",
				Text2:   "Text2",
				Vertical: []types.LayoutElement{
					{
						Id:         "2",
						Caption:    "Caption2",
						Text1:      "Text3",
						Text2:      "Text4",
						Vertical:   []types.LayoutElement{},
						Horizontal: []types.LayoutElement{},
					},
				},
				Horizontal: []types.LayoutElement{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := types.ExpInitLayoutElement(&tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("InitLayoutElement() = %v, want %v", result, tt.expected)
			}
		})
	}
}
