package types_test

import (
	"testing"

	"github.com/okieoth/draw.chart.things/pkg/types"
)

func TestInitLayoutElement(t *testing.T) {
	tests := []struct {
		name         string
		layout       types.Layout
		inputFormats map[string]types.Format
		expected     types.LayoutElement
	}{
		{
			name: "Test with empty layout and formats",
			layout: types.Layout{
				Id: "test1",
			},
			inputFormats: map[string]types.Format{},
			expected: types.LayoutElement{
				Id: "test1",
				Format: types.BoxFormat{
					Padding: 5,
					FontCaption: types.FontDef{
						Size:       10,
						Type:       &types.ExpFontDefTypeEnum_normal,
						Weight:     &types.ExpFontDefWeightEnum_normal,
						LineHeight: 1.5,
						Color:      "black",
						Aligned:    &types.ExpFontDefAlignedEnum_left,
					},
					FontText1: types.FontDef{
						Size:       10,
						Type:       &types.ExpFontDefTypeEnum_normal,
						Weight:     &types.ExpFontDefWeightEnum_normal,
						LineHeight: 1.5,
						Color:      "black",
						Aligned:    &types.ExpFontDefAlignedEnum_left,
					},
					FontText2: types.FontDef{
						Size:       10,
						Type:       &types.ExpFontDefTypeEnum_normal,
						Weight:     &types.ExpFontDefWeightEnum_normal,
						LineHeight: 1.5,
						Color:      "black",
						Aligned:    &types.ExpFontDefAlignedEnum_left,
					},
				},
			},
		},
		{
			name: "Test with layout and formats",
			layout: types.Layout{
				Id:      "test2",
				Caption: "Test Caption",
				Text1:   "Test Text1",
				Text2:   "Test Text2",
				Tags:    []string{"tag1"},
			},
			inputFormats: map[string]types.Format{
				"tag1": {
					FontCaption: &types.FontDef{
						Size:       12,
						Type:       &types.ExpFontDefTypeEnum_normal,
						Weight:     &types.ExpFontDefWeightEnum_bold,
						LineHeight: 2.0,
						Color:      "red",
						Aligned:    &types.ExpFontDefAlignedEnum_center,
					},
				},
			},
			expected: types.LayoutElement{
				Id:      "test2",
				Caption: "Test Caption",
				Text1:   "Test Text1",
				Text2:   "Test Text2",
				Format: types.BoxFormat{
					Padding: 5,
					FontCaption: types.FontDef{
						Size:       12,
						Type:       &types.ExpFontDefTypeEnum_normal,
						Weight:     &types.ExpFontDefWeightEnum_bold,
						LineHeight: 2.0,
						Color:      "red",
						Aligned:    &types.ExpFontDefAlignedEnum_center,
					},
					FontText1: types.FontDef{
						Size:       10,
						Type:       &types.ExpFontDefTypeEnum_normal,
						Weight:     &types.ExpFontDefWeightEnum_normal,
						LineHeight: 1.5,
						Color:      "black",
						Aligned:    &types.ExpFontDefAlignedEnum_left,
					},
					FontText2: types.FontDef{
						Size:       10,
						Type:       &types.ExpFontDefTypeEnum_normal,
						Weight:     &types.ExpFontDefWeightEnum_normal,
						LineHeight: 1.5,
						Color:      "black",
						Aligned:    &types.ExpFontDefAlignedEnum_left,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := types.ExpInitLayoutElement(&tt.layout, tt.inputFormats)
			if result.Id != tt.expected.Id {
				t.Errorf("expected Id %v, got %v", tt.expected.Id, result.Id)
			}
			if result.Caption != tt.expected.Caption {
				t.Errorf("expected Caption %v, got %v", tt.expected.Caption, result.Caption)
			}
			if result.Text1 != tt.expected.Text1 {
				t.Errorf("expected Text1 %v, got %v", tt.expected.Text1, result.Text1)
			}
			if result.Text2 != tt.expected.Text2 {
				t.Errorf("expected Text2 %v, got %v", tt.expected.Text2, result.Text2)
			}
			if result.Format.Padding != tt.expected.Format.Padding {
				t.Errorf("expected Padding %v, got %v", tt.expected.Format.Padding, result.Format.Padding)
			}
			if result.Format.FontCaption.Size != tt.expected.Format.FontCaption.Size {
				t.Errorf("expected FontCaption Size %v, got %v", tt.expected.Format.FontCaption.Size, result.Format.FontCaption.Size)
			}
			if *result.Format.FontCaption.Type != *tt.expected.Format.FontCaption.Type {
				t.Errorf("expected FontCaption Type %v, got %v", tt.expected.Format.FontCaption.Type, result.Format.FontCaption.Type)
			}
			if *result.Format.FontCaption.Weight != *tt.expected.Format.FontCaption.Weight {
				t.Errorf("expected FontCaption Weight %v, got %v", tt.expected.Format.FontCaption.Weight, result.Format.FontCaption.Weight)
			}
			if result.Format.FontCaption.LineHeight != tt.expected.Format.FontCaption.LineHeight {
				t.Errorf("expected FontCaption LineHeight %v, got %v", tt.expected.Format.FontCaption.LineHeight, result.Format.FontCaption.LineHeight)
			}
			if result.Format.FontCaption.Color != tt.expected.Format.FontCaption.Color {
				t.Errorf("expected FontCaption Color %v, got %v", tt.expected.Format.FontCaption.Color, result.Format.FontCaption.Color)
			}
			if *result.Format.FontCaption.Aligned != *tt.expected.Format.FontCaption.Aligned {
				t.Errorf("expected FontCaption Aligned %v, got %v", tt.expected.Format.FontCaption.Aligned, result.Format.FontCaption.Aligned)
			}
		})
	}
}
