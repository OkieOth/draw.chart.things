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

// used in filter situations in cases where no ID are provided
var globalId int

func LoadBoxesFromFile(inputFile string) (*boxes.Boxes, error) {
	layout, err := types.LoadInputFromFile[boxes.Boxes](inputFile)
	if err != nil {
		return nil, fmt.Errorf("error while loading boxes from file: %w", err)
	}
	err = replaceAlternativePaths(&layout.Boxes, inputFile)
	if err != nil {
		return nil, fmt.Errorf("error while replacing relative paths: %w", err)
	}
	if layout != nil && layout.Version != nil && layout.Title != "" {
		layout.Title += fmt.Sprintf(" [%s]", *layout.Version)
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

	doc.ConnectBoxes()
	doc.AdjustDocHeightTLegend(textDimensionCalulator)
	doc.IncludeComments(textDimensionCalulator)
	output, err := os.Create(outputFile)
	svgdrawing := svgdrawing.NewDrawing(output)
	svgdrawing.Start(doc.Title, doc.Height, doc.Width)
	svgdrawing.InitImages(doc.Images)
	svgdrawing.DrawRaster(doc.Width, doc.Height, types.RasterSize)
	doc.DrawBoxes(svgdrawing)
	doc.DrawConnections(svgdrawing)
	doc.DrawTitle(svgdrawing, textDimensionCalulator)
	doc.DrawLegend(svgdrawing, textDimensionCalulator)
	doc.DrawComments(svgdrawing, textDimensionCalulator)
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

func collectTruncatedFromCont(b []boxes.Layout, newId string, truncated *map[string]TruncatedInfo, hideComments bool) {
	for _, l := range b {
		collectTruncated(l, newId, truncated, hideComments)
	}
}

func collectTruncated(b boxes.Layout, newId string, truncated *map[string]TruncatedInfo, hideComments bool) {
	if b.Id != "" {
		if hideComments && b.Comment != nil {
			b.Comment = nil
		}
		(*truncated)[b.Id] = TruncatedInfo{
			truncated: b,
			newId:     newId,
		}
	}
	collectTruncatedFromCont(b.Horizontal, newId, truncated, hideComments)
	collectTruncatedFromCont(b.Vertical, newId, truncated, hideComments)
}

type TruncatedInfo struct {
	truncated boxes.Layout
	newId     string
}

func addTruncated(src map[string]TruncatedInfo, dest *map[string]TruncatedInfo) {
	for k, v := range src {
		(*dest)[k] = v
	}
}

func getNewId() string {
	globalId++
	return fmt.Sprintf("xxx_%d", globalId)
}

func removeCommentsInCase(b boxes.Layout, hideComments bool) boxes.Layout {
	if hideComments {
		if b.Comment != nil {
			b.Comment = nil
		}
		for i := range b.Connections {
			c := b.Connections[i]
			if c.Comment != nil {
				c.Comment = nil
			}
		}
	}
	return b
}

func truncBoxes(b boxes.Layout, currentDepth, maxDepth int, expanded, blacklisted []string, hideComments bool) (boxes.Layout, map[string]TruncatedInfo) {
	truncatedBoxes := make(map[string]TruncatedInfo, 0)
	b = removeCommentsInCase(b, hideComments)
	if (currentDepth >= maxDepth) && (!b.Expand) && (!isRelatedToId(b, expanded)) {
		// possible removed connections to this object
		if b.Id == "" {
			b.Id = getNewId()
		}
		collectTruncatedFromCont(b.Horizontal, b.Id, &truncatedBoxes, hideComments)
		collectTruncatedFromCont(b.Vertical, b.Id, &truncatedBoxes, hideComments)
		b.Horizontal = make([]boxes.Layout, 0)
		b.Vertical = make([]boxes.Layout, 0)
	} else {
		if len(b.Horizontal) > 0 {
			cont := make([]boxes.Layout, 0)
			for i := range len(b.Horizontal) {
				id := b.Horizontal[i].Id
				if blacklisted != nil && id != "" && slices.Contains(blacklisted, id) {
					// possible removed connections to this object
					if b.Id == "" {
						b.Id = getNewId()
					}
					collectTruncated(b.Horizontal[i], b.Id, &truncatedBoxes, hideComments)
					continue
				}
				l, trunc := truncBoxes(b.Horizontal[i], currentDepth+1, maxDepth, expanded, blacklisted, hideComments)
				addTruncated(trunc, &truncatedBoxes)
				cont = append(cont, l)
			}
			b.Horizontal = cont
		}
		if len(b.Vertical) > 0 {
			cont := make([]boxes.Layout, 0)
			for i := range len(b.Vertical) {
				id := b.Vertical[i].Id
				if blacklisted != nil && id != "" && slices.Contains(blacklisted, id) {
					// possible removed connections to this object
					if b.Id == "" {
						b.Id = getNewId()
					}
					collectTruncated(b.Vertical[i], b.Id, &truncatedBoxes, hideComments)
					continue
				}
				l, trunc := truncBoxes(b.Vertical[i], currentDepth+1, maxDepth, expanded, blacklisted, hideComments)
				addTruncated(trunc, &truncatedBoxes)
				cont = append(cont, l)
			}
			b.Vertical = cont
		}
	}
	return b, truncatedBoxes
}

func connectionExistsByDestId(connections []boxes.Connection, destId string) bool {
	return slices.ContainsFunc(connections, func(c boxes.Connection) bool {
		return c.DestId == destId
	})
}

func copyTruncatedComments(layout *boxes.Layout, truncatedObjects map[string]TruncatedInfo) {
	// find comments from truncated objects to the annotate current object
	found := false
	for _, v := range truncatedObjects {
		if v.newId == layout.Id {
			if v.truncated.Comment != nil && (v.truncated.Comment.Text != "" || v.truncated.Comment.Label != nil) {
				found = true
				break
			}
		}
	}
	if found {
		layout.HiddenComments = true
	}
}

func copyTruncatedConnections(layout *boxes.Layout, truncatedObjects map[string]TruncatedInfo, hideComments bool) {
	// copy all needed connections from truncated objects to the current object
	for _, v := range truncatedObjects {
		if v.newId == layout.Id {
			// copy all truncated connections
			for _, c := range v.truncated.Connections {
				if c.DestId == layout.Id {
					continue
				}
				// destId is also replaced
				destIdToUse := c.DestId
				if obj, ok := truncatedObjects[c.DestId]; ok {
					destIdToUse = obj.newId
				}
				if !connectionExistsByDestId(layout.Connections, destIdToUse) {
					if destIdToUse != layout.Id {
						c.DestId = destIdToUse
						if hideComments && c.Comment != nil {
							c.Comment = nil
						}
						(*layout).Connections = append(layout.Connections, c)
					}
				}
			}
		}
	}
}

func adjustDestIdInRespectOfTruncated(layout *boxes.Layout, truncatedObjects map[string]TruncatedInfo) {
	normalizedConnections := make([]boxes.Connection, 0)
	for _, c := range layout.Connections {
		if trunc, found := truncatedObjects[c.DestId]; found {
			c.DestId = trunc.newId
		}
		if !connectionExistsByDestId(normalizedConnections, c.DestId) {
			normalizedConnections = append(normalizedConnections, c)
		}
	}
	layout.Connections = normalizedConnections
}

func adjustTruncated(layout *boxes.Layout, truncatedObjects map[string]TruncatedInfo, hideComments bool) {
	copyTruncatedConnections(layout, truncatedObjects, hideComments)
	copyTruncatedComments(layout, truncatedObjects)
	adjustDestIdInRespectOfTruncated(layout, truncatedObjects)
	if len(layout.Horizontal) > 0 {
		adjustTruncatedForCont(&layout.Horizontal, truncatedObjects, hideComments)
	}
	if len(layout.Vertical) > 0 {
		adjustTruncatedForCont(&layout.Vertical, truncatedObjects, hideComments)
	}
}

func adjustTruncatedForCont(cont *[]boxes.Layout, truncatedObjects map[string]TruncatedInfo, hideComments bool) {
	if cont != nil {
		for i := range len(*cont) {
			adjustTruncated(&(*cont)[i], truncatedObjects, hideComments)
		}
	}
}

func FilterBoxes(layout boxes.Boxes, defaultDepth int, expanded, blacklisted []string) boxes.Boxes {
	return FilterBoxesComments(layout, defaultDepth, expanded, blacklisted, false)
}

func FilterBoxesComments(layout boxes.Boxes, defaultDepth int, expanded, blacklisted []string, hideComments bool) boxes.Boxes {
	filteredBoxes := boxes.CopyBoxes(&layout)
	b, truncatedObjects := truncBoxes(layout.Boxes, 0, defaultDepth, expanded, blacklisted, hideComments)
	adjustTruncated(&b, truncatedObjects, hideComments)
	filteredBoxes.Boxes = b
	return *filteredBoxes
}

func DrawBoxesFilteredExt(
	layout boxes.Boxes,
	mixins []boxes.BoxesFileMixings,
	defaultDepth int,
	expanded, blacklisted []string,
	debug bool) UIReturn {
	for _, m := range mixins {
		layout.MixinThings(m)
	}
	return DrawBoxesFiltered(layout, defaultDepth, expanded, blacklisted, debug)
}

func DrawBoxesFiltered(layout boxes.Boxes, defaultDepth int, expanded, blacklisted []string, debug bool) UIReturn {
	return DrawBoxesFilteredComments(layout, defaultDepth, expanded, blacklisted, false, debug)
}

func DrawBoxesFilteredComments(layout boxes.Boxes, defaultDepth int, expanded, blacklisted []string, hideComments, debug bool) UIReturn {
	textDimensionCalulator := svgdrawing.NewSvgTextDimensionCalculator()
	filteredLayout := FilterBoxesComments(layout, defaultDepth, expanded, blacklisted, hideComments)
	doc, err := InitialLayoutBoxes(&filteredLayout, textDimensionCalulator)
	if err != nil {
		return UIReturn{ErrorMsg: fmt.Sprintf("error while initialy layout: %v", err)}
	}
	doc.ConnectBoxes()
	doc.AdjustDocHeightTLegend(textDimensionCalulator)
	doc.IncludeComments(textDimensionCalulator)
	var svgBuilder strings.Builder
	svgdrawing := svgdrawing.NewDrawing(&svgBuilder)
	svgdrawing.Start(doc.Title, doc.Height, doc.Width)
	svgdrawing.InitImages(doc.Images)
	if debug {
		svgdrawing.DrawRaster(doc.Width, doc.Height, types.RasterSize)
	}
	doc.DrawBoxes(svgdrawing)
	if debug {
		doc.InitStartPositions() // not needed to be called separately
		doc.InitRoads()          // not needed to be called separately
		doc.DrawRoads(svgdrawing)
		doc.DrawStartPositions(svgdrawing)
		doc.DrawConnectionNodes(svgdrawing)
		doc.Boxes.DrawTextBoxes(svgdrawing)
	}
	doc.DrawConnections(svgdrawing)
	doc.DrawTitle(svgdrawing, textDimensionCalulator)
	doc.DrawLegend(svgdrawing, textDimensionCalulator)
	doc.DrawComments(svgdrawing, textDimensionCalulator)
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
	doc, err := DocumentFromBoxes(b)
	if err != nil {
		return nil, err
	}

	doc.Boxes.X = doc.GlobalPadding
	doc.Boxes.Y = doc.GlobalPadding

	doc.Boxes.InitDimensions(c)
	doc.Width = doc.Boxes.Width + doc.GlobalPadding*2
	doc.Height = doc.Boxes.Height + doc.GlobalPadding*2
	doc.Boxes.Center()
	doc.Height = doc.AdjustDocHeight(&doc.Boxes, 0) + (2 * doc.GlobalPadding)

	if doc.Title != "" {
		format2Use := doc.GetTitleFormat()
		w, h := c.Dimensions(doc.Title, &format2Use)
		doc.Height += h + (2 * doc.GlobalPadding)
		if w > doc.Width {
			doc.Width = w + (2 * doc.GlobalPadding)
		}
	}

	return doc, nil
}

func LayoutBoxesWithFilter(b *boxes.Boxes, c types.TextDimensionCalculator, defaultDepth int) (*boxes.BoxesDocument, error) {
	doc, err := DocumentFromBoxes(b)
	if err != nil {
		return nil, err
	}
	doc.Boxes.X = doc.GlobalPadding
	doc.Boxes.Y = doc.GlobalPadding
	doc.Boxes.InitDimensions(c)
	doc.Width = doc.Boxes.Width + doc.GlobalPadding*2
	doc.Height = doc.Boxes.Height + doc.GlobalPadding*2
	doc.Boxes.Center()
	doc.Height = doc.AdjustDocHeight(&doc.Boxes, 0) + (2 * doc.GlobalPadding)

	if doc.Title != "" {
		format2Use := doc.GetTitleFormat()
		w, h := c.Dimensions(doc.Title, &format2Use)
		doc.Height += h + (2 * doc.GlobalPadding)
		if w > doc.Width {
			doc.Width = w + (2 * doc.GlobalPadding)
		}
	}
	return doc, nil
}
