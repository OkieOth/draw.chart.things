package boxesimpl_test

import (
	"os"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"testing"

	"github.com/okieoth/draw.chart.things/pkg/boxesimpl"
	"github.com/okieoth/draw.chart.things/pkg/svgdrawing"
	"github.com/okieoth/draw.chart.things/pkg/types"
	"github.com/okieoth/draw.chart.things/pkg/types/boxes"
)

type DummyDimensionCalculator struct {
	width  int
	height int
}

func (d *DummyDimensionCalculator) Dimensions(txt string, format *types.FontDef) (width, height int) {
	return d.width, d.height
}

func (d *DummyDimensionCalculator) SplitTxt(txt string, format *types.FontDef) (width, height int, lines []types.TextAndDimensions) {
	return 0, 0, make([]types.TextAndDimensions, 0)
}

func NewDummyDimensionCalculator(width, height int) *DummyDimensionCalculator {
	return &DummyDimensionCalculator{
		width:  width,
		height: height,
	}
}

func TestDrawBoxesFromFile(t *testing.T) {
	tests := []struct {
		inputFile string
	}{
		// {
		// 	inputFile:  "../../resources/examples_boxes/simple_box.yaml",
		// },
		{
			inputFile: "../../resources/examples_boxes/simple_diamond.yaml",
		}}

	textDimensionCalulator := svgdrawing.NewSvgTextDimensionCalculator()

	for _, test := range tests {
		b, err := types.LoadInputFromFile[boxes.Boxes](test.inputFile)
		require.Nil(t, err)
		doc, err := boxesimpl.InitialLayoutBoxes(b, textDimensionCalulator)
		require.Nil(t, err)
		require.NotNil(t, doc)
	}

}

// func tested types.InitDimensions
func TestInitDimensions(t *testing.T) {
	tests := []struct {
		layout         boxes.Layout
		expectedHeight int
		expectedWidth  int
	}{
		{
			layout: boxes.Layout{
				Caption: "test1",
			},
			expectedHeight: 80,
			expectedWidth:  150,
		},
		{
			// extends "test1" with an additional text1
			layout: boxes.Layout{
				Caption: "test2",
				Text1:   "test2-text1",
			},
			expectedHeight: 150,
			expectedWidth:  150,
		},
		{
			// extends "test2" with an additional text2
			layout: boxes.Layout{
				Caption: "test3",
				Text1:   "test3-text1",
				Text2:   "test3-text2",
			},
			expectedHeight: 200,
			expectedWidth:  150,
		},
		{
			// basic vertical layout test
			layout: boxes.Layout{
				Caption: "test4",
				Text1:   "test4-text1",
				Text2:   "test4-text2",
				Vertical: []boxes.Layout{
					{
						Caption: "test4-V1",
					},
					{
						Caption: "test4-V2",
					},
					{
						Caption: "test3-V3",
					},
				},
			},
			expectedHeight: 500,
			expectedWidth:  150,
		},
		{
			// basic horizontal layout test
			layout: boxes.Layout{
				Caption: "test5",
				Text1:   "test5-text1",
				Text2:   "test5-text2",
				Horizontal: []boxes.Layout{
					{
						Caption: "test5-V1",
					},
					{
						Caption: "test5-V2",
					},
					{
						Caption: "test5-V3",
					},
				},
			},
			expectedHeight: 300,
			expectedWidth:  530,
		},
	}

	dc := NewDummyDimensionCalculator(100, 50)
	emptyFormats := map[string]boxes.BoxFormat{}
	for _, test := range tests {
		dummy := make([]string, 0)
		le := boxesimpl.ExpInitLayoutElement(&test.layout, emptyFormats, &dummy)
		le.InitDimensions(dc)
		assert.Equal(t, test.expectedHeight, le.Height)
		assert.Equal(t, test.expectedWidth, le.Width)
	}
}

func TestLoadBoxesFromFile(t *testing.T) {
	tests := []struct {
		testFile        string
		shouldLoad      bool
		fileToCompareTo string
	}{
		{
			testFile:        "../../resources/examples_boxes/horizontal_nested_diamond_ext.yaml",
			shouldLoad:      true,
			fileToCompareTo: "../../resources/examples_boxes/horizontal_nested_diamond.yaml",
		},
		{
			testFile:        "../../resources/examples_boxes/horizontal_nested_diamond_ext2.yaml",
			shouldLoad:      true,
			fileToCompareTo: "../../resources/examples_boxes/horizontal_nested_diamond.yaml",
		},
		{
			testFile:        "../../resources/examples_boxes/horizontal_nested_diamond_ext2_fail.yaml",
			shouldLoad:      false,
			fileToCompareTo: "",
		},
	}
	for _, test := range tests {
		b, err := boxesimpl.LoadBoxesFromFile(test.testFile)
		if test.shouldLoad {
			require.Nil(t, err)
			require.NotNil(t, b)
			// Compare the loaded file with the original file
			b2, err := types.LoadInputFromFile[boxes.Boxes](test.fileToCompareTo)
			require.Nil(t, err)
			assert.Equal(t, b, b2, "The loaded boxes should match the expected boxes from the file")
		} else {
			require.NotNil(t, err)
		}
	}
}

func TestDrawBoxesForUi(t *testing.T) {
	tests := []struct {
		inputFile   string
		outputFile  string
		depth       int
		expanded    []string
		blacklisted []string
	}{
		{
			inputFile:   "../../resources/examples_boxes/complex_complex.yaml",
			outputFile:  "../../temp/complex_complex_filtered_1.svg",
			depth:       1,
			expanded:    []string{},
			blacklisted: []string{},
		},
		{
			inputFile:   "../../resources/examples_boxes/complex_complex.yaml",
			outputFile:  "../../temp/complex_complex_filtered_2.svg",
			depth:       2,
			expanded:    []string{},
			blacklisted: []string{},
		},
	}
	for i, test := range tests {
		b, err := types.LoadInputFromFile[boxes.Boxes](test.inputFile)
		require.Nil(t, err, "error while loading input file for test", i)
		svgReturn := boxesimpl.DrawBoxesFiltered(*b, test.depth, test.expanded, test.blacklisted, true)

		require.Equal(t, "", svgReturn.ErrorMsg, "error generating SVG output for test", i)

		err = os.WriteFile(test.outputFile, []byte(svgReturn.SVG), 0600)
		require.Nil(t, err, "error while writing output file for test", i)
		require.FileExists(t, test.outputFile, "can't find created output file", test.outputFile)
	}
}

func TestFilterBoxes(t *testing.T) {
	tests := []struct {
		inputFile   string
		checkFunc   func(b *boxes.Boxes)
		depth       int
		expanded    []string
		blacklisted []string
	}{
		// {
		// 	inputFile: "../../resources/examples_boxes/complex_complex.yaml",
		// 	checkFunc: func(b *boxes.Boxes) {
		// 		for _, e := range b.Boxes.Horizontal {
		// 			require.Equal(t, 0, len(e.Horizontal), "got unexpected horizontal childs (1-1)")
		// 			require.Equal(t, 0, len(e.Vertical), "got unexpected vertical childs (1-1)")
		// 		}
		// 		for _, e := range b.Boxes.Vertical {
		// 			require.Equal(t, 0, len(e.Horizontal), "got unexpected horizontal childs (1-2)")
		// 			require.Equal(t, 0, len(e.Vertical), "got unexpected vertical childs (1-2)")
		// 		}
		// 	},
		// 	depth:       1,
		// 	expanded:    []string{},
		// 	blacklisted: []string{},
		// },
		// {
		// 	inputFile: "../../resources/examples_boxes/complex_complex.yaml",
		// 	checkFunc: func(b *boxes.Boxes) {
		// 		found := false
		// 		for _, e := range b.Boxes.Horizontal {
		// 			if len(e.Horizontal) > 0 {
		// 				found = true
		// 			}
		// 			if len(e.Vertical) > 0 {
		// 				found = true
		// 			}
		// 		}
		// 		for _, e := range b.Boxes.Vertical {
		// 			if len(e.Horizontal) > 0 {
		// 				found = true
		// 			}
		// 			if len(e.Vertical) > 0 {
		// 				found = true
		// 			}
		// 		}
		// 		require.True(t, found, "didn't find second level")
		// 	},
		// 	depth:       2,
		// 	expanded:    []string{},
		// 	blacklisted: []string{},
		// },
		{
			inputFile: "../../resources/examples_boxes/complex_complex.yaml",
			checkFunc: func(b *boxes.Boxes) {
				found := false
				for _, e := range b.Boxes.Horizontal {
					if len(e.Horizontal) > 0 {
						found = true
					}
					if len(e.Vertical) > 0 {
						found = true
					}
				}
				for _, e := range b.Boxes.Vertical {
					if len(e.Horizontal) > 0 {
						found = true
					}
					if len(e.Vertical) > 0 {
						found = true
					}
				}
				require.True(t, found, "didn't find second level")
			},
			depth:       20,
			expanded:    []string{},
			blacklisted: []string{"r2_2", "r4_1"},
		},
	}
	for i, test := range tests {
		b, err := types.LoadInputFromFile[boxes.Boxes](test.inputFile)
		require.Nil(t, err, "error while loading input file for test", i)
		filtered := boxesimpl.FilterBoxes(*b, test.depth, test.expanded, test.blacklisted)
		test.checkFunc(&filtered)
	}
}
