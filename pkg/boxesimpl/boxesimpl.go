package boxesimpl

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/okieoth/draw.chart.things/pkg/svgdrawing"
	"github.com/okieoth/draw.chart.things/pkg/types"
	"github.com/okieoth/draw.chart.things/pkg/types/boxes"
)

func LoadBoxesFromFile(inputFile string) (*boxes.Boxes, error) {
	layout, err := types.LoadInputFromFile[boxes.Boxes](inputFile)
	if err != nil {
		return nil, fmt.Errorf("error while loading boxes from file: %w", err)
	}
	err = replaceAlternativePaths(&layout.Boxes, inputFile)
	if err != nil {
		return nil, fmt.Errorf("error while replacing relative paths: %w", err)
	}
	return layout, nil
}

func DrawBoxesFromFile(inputFile, outputFile string) error {
	layout, err := LoadBoxesFromFile(inputFile)
	if err != nil {
		return fmt.Errorf("error while loading input: %v", err)
	}
	textDimensionCalulator := svgdrawing.NewSvgTextDimensionCalculator()

	doc, err := InitialLayoutBoxes(layout, textDimensionCalulator)
	if err != nil {
		return err
	}

	// FIXME, TODO: this doesn't terminate!!!!
	//doc.ConnectBoxes()
	output, err := os.Create(outputFile)
	svgdrawing := svgdrawing.NewDrawing(output)
	svgdrawing.Start(doc.Title, doc.Height, doc.Width)
	svgdrawing.DrawRaster(doc.Width, doc.Height, types.RasterSize)
	doc.DrawBoxes(svgdrawing)
	doc.DrawConnections(svgdrawing)
	svgdrawing.Done()
	output.Close()
	return nil
}

type UIReturn struct {
	SVG      string
	ErrorMsg string
}

func isRelatedToId(b boxes.Layout, filter []string) bool {
	if slices.Contains(filter, b.Id) {
		return true
	}
	for _, e := range b.Horizontal {
		if isRelatedToId(e, filter) {
			return true
		}
	}
	for _, e := range b.Vertical {
		if isRelatedToId(e, filter) {
			return true
		}
	}
	return false
}

func truncBoxes(b boxes.Layout, currentDepth, maxDepth int, filter []string) boxes.Layout {
	if (currentDepth >= maxDepth) && (!isRelatedToId(b, filter)) {
		b.Horizontal = make([]boxes.Layout, 0)
		b.Vertical = make([]boxes.Layout, 0)
	} else {
		for i := range len(b.Horizontal) {
			b.Horizontal[i] = truncBoxes(b.Horizontal[i], currentDepth+1, maxDepth, filter)
		}
		for i := range len(b.Vertical) {
			b.Vertical[i] = truncBoxes(b.Vertical[i], currentDepth+1, maxDepth, filter)
		}
	}
	return b
}

func filterBoxes(layout boxes.Boxes, defaultDepth int, filter []string) boxes.Boxes {
	filteredBoxes := boxes.CopyBoxes(&layout)
	filteredBoxes.Boxes = truncBoxes(layout.Boxes, 0, defaultDepth, filter)
	return *filteredBoxes
}

func DrawBoxesFiltered(layout boxes.Boxes, defaultDepth int, filter []string) UIReturn {
	textDimensionCalulator := svgdrawing.NewSvgTextDimensionCalculator()
	filteredLayout := filterBoxes(layout, defaultDepth, filter)
	doc, err := InitialLayoutBoxes(&filteredLayout, textDimensionCalulator)
	if err != nil {
		return UIReturn{ErrorMsg: fmt.Sprintf("error while initialy layout: %v", err)}
	}

	// FIXME, TODO: this doesn't terminate!!!!
	//doc.ConnectBoxes()
	var svgBuilder strings.Builder
	svgdrawing := svgdrawing.NewDrawing(&svgBuilder)
	svgdrawing.Start(doc.Title, doc.Height, doc.Width)
	svgdrawing.DrawRaster(doc.Width, doc.Height, types.RasterSize)
	doc.DrawBoxes(svgdrawing)
	doc.DrawConnections(svgdrawing)
	svgdrawing.Done()
	return UIReturn{SVG: svgBuilder.String()}
}

func initContainerFromAlternativePath(extPathToLoad, origInputFile string) ([]boxes.Layout, error) {
	var pathToLoad string
	if path.IsAbs(extPathToLoad) {
		pathToLoad = extPathToLoad
	} else {
		// relative path, so we need to resolve it against the original input file
		origDir := filepath.Dir(origInputFile)
		pathToLoad = filepath.Join(origDir, extPathToLoad)
	}
	layoutArray, err := types.LoadInputFromFile[[]boxes.Layout](pathToLoad)
	if err != nil {
		return make([]boxes.Layout, 0), err
	}
	return *layoutArray, nil
}

func replaceAlternativePaths(l *boxes.Layout, inputFile string) error {
	if l == nil {
		return nil
	}
	if l.ExtHorizontal != nil {
		if cont, err := initContainerFromAlternativePath(*l.ExtHorizontal, inputFile); err != nil {
			return fmt.Errorf("error while initialize horizontal container from relative path: %w", err)
		} else {
			l.ExtHorizontal = nil // clear the alternative path
			l.Horizontal = cont
		}
	}
	if l.ExtVertical != nil {
		if cont, err := initContainerFromAlternativePath(*l.ExtVertical, inputFile); err != nil {
			return fmt.Errorf("error while initialize vertical container from relative path: %w", err)
		} else {
			l.ExtVertical = nil // clear the alternative path
			l.Vertical = cont
		}
	}
	if err := replaceAlternativePathsForContainer(l.Horizontal, inputFile); err != nil {
		return fmt.Errorf("error while replacing alternative paths for horizontal container: %w", err)
	}
	if err := replaceAlternativePathsForContainer(l.Vertical, inputFile); err != nil {
		return fmt.Errorf("error while replacing alternative paths for vertical container: %w", err)
	}
	return nil
}

func replaceAlternativePathsForContainer(cont []boxes.Layout, inputFile string) error {
	if cont == nil {
		return nil
	}
	for i, _ := range cont {
		if err := replaceAlternativePaths(&cont[i], inputFile); err != nil {
			return fmt.Errorf("error while replacing alternative paths for container: %w", err)
		}
	}
	return nil
}

func InitialLayoutBoxes(b *boxes.Boxes, c types.TextDimensionCalculator) (*boxes.BoxesDocument, error) {
	doc := DocumentFromBoxes(b)
	doc.Boxes.X = doc.GlobalPadding
	doc.Boxes.Y = doc.GlobalPadding
	doc.Boxes.InitDimensions(c)
	doc.Width = doc.Boxes.Width + doc.GlobalPadding*2
	doc.Height = doc.Boxes.Height + doc.GlobalPadding*2
	if doc.Title != "" {
		defaultFormat := doc.Formats["default"] // risky doesn't check if default exists
		w, h := c.Dimensions(doc.Title, &defaultFormat.FontCaption)
		doc.Height += h + (2 * doc.GlobalPadding)
		if w > doc.Width {
			doc.Width = w + (2 * doc.GlobalPadding)
		}
	}
	doc.Boxes.Center()
	doc.Height = doc.AdjustDocHeight(&doc.Boxes, 0) + (2 * doc.GlobalPadding)
	return doc, nil
}

func LayoutBoxesWithFilter(b *boxes.Boxes, c types.TextDimensionCalculator, defaultDepth int) (*boxes.BoxesDocument, error) {
	doc := DocumentFromBoxes(b)
	doc.Boxes.X = doc.GlobalPadding
	doc.Boxes.Y = doc.GlobalPadding
	doc.Boxes.InitDimensions(c)
	doc.Width = doc.Boxes.Width + doc.GlobalPadding*2
	doc.Height = doc.Boxes.Height + doc.GlobalPadding*2
	if doc.Title != "" {
		defaultFormat := doc.Formats["default"] // risky doesn't check if default exists
		w, h := c.Dimensions(doc.Title, &defaultFormat.FontCaption)
		doc.Height += h + (2 * doc.GlobalPadding)
		if w > doc.Width {
			doc.Width = w + (2 * doc.GlobalPadding)
		}
	}
	doc.Boxes.Center()
	doc.Height = doc.AdjustDocHeight(&doc.Boxes, 0) + (2 * doc.GlobalPadding)
	return doc, nil
}
