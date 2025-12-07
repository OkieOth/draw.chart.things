package boxes

import (
	"fmt"

	"github.com/okieoth/draw.chart.things/pkg/types"
)

type PointToTest struct {
	X            int
	Y            int
	HasCollision bool
}

func getRelevantPoints(p1, p2, other int, verticalBorder bool) []PointToTest {
	// expects always to go from top to bottom or left to right
	start, end := p1, p2
	if p1 > p2 {
		start, end = p2, p1
	}
	endVal := end - types.RasterSize
	ret := make([]PointToTest, 0)
	for i := start + types.RasterSize; i <= endVal; i += types.RasterSize {
		var p PointToTest
		if verticalBorder {
			p.X = other
			p.Y = i
		} else {
			p.Y = other
			p.X = i
		}
		ret = append(ret, p)
	}
	return ret
}

type checkForCollFunc func(relevantPoints *[]PointToTest, currentElem, elemToHandle *LayoutElement)
type checkForCollFunc2 func(x, y int, currentElem, startElem, endElem *LayoutElement) bool

func (doc *BoxesDocument) checkForCollInContainer(cont *LayoutElemContainer, f checkForCollFunc, relevantPoints *[]PointToTest, elemToHandle *LayoutElement) {
	if cont != nil {
		l := len(cont.Elems)
		for i := 0; i < l; i++ {
			e := &cont.Elems[i]
			if e == elemToHandle {
				continue
			}
			f(relevantPoints, e, elemToHandle)
		}
	}
}

func (doc *BoxesDocument) checkForCollInContainer2(cont *LayoutElemContainer, f checkForCollFunc2, x, y int, startElem, endElem *LayoutElement) bool {
	if cont != nil {
		l := len(cont.Elems)
		for i := 0; i < l; i++ {
			e := &cont.Elems[i]
			if e == startElem || e == endElem {
				continue
			}
			if f(x, y, e, startElem, endElem) {
				return true
			}
		}
	}
	return false
}

// checks if a point is inside a box, returns true if so
func (doc *BoxesDocument) checkColl(x, y int, currentElem, startElem, endElem *LayoutElement) bool {
	if (currentElem != startElem) && (currentElem != endElem) &&
		(doc.ShouldHandle(currentElem)) &&
		((startElem == nil) || (!doc.isParent(currentElem, startElem))) &&
		((endElem == nil) || (!doc.isParent(currentElem, endElem))) {
		curMinX := currentElem.X
		curMaxX := currentElem.X + currentElem.Width
		curMinY := currentElem.Y
		curMaxY := currentElem.Y + currentElem.Height
		if (x <= curMaxX) && (x >= curMinX) && (y >= curMinY) && (y <= curMaxY) {
			return true
		}
	}
	if doc.checkForCollInContainer2(currentElem.Vertical, doc.checkColl, x, y, startElem, endElem) {
		return true
	}
	return doc.checkForCollInContainer2(currentElem.Horizontal, doc.checkColl, x, y, startElem, endElem)
}

func (doc *BoxesDocument) checkForLeftColl(relevantPoints *[]PointToTest, currentElem, elemToHandle *LayoutElement) {
	if (currentElem != elemToHandle) && doc.ShouldHandle(currentElem) && (!doc.isParent(currentElem, elemToHandle)) {
		curMinX := currentElem.X
		curMaxX := currentElem.X + currentElem.Width
		curMinY := currentElem.Y
		curMaxY := currentElem.Y + currentElem.Height
		for i := 0; i < len(*relevantPoints); i++ {
			x := (*relevantPoints)[i].X - (4 * types.RasterSize)
			y := (*relevantPoints)[i].Y
			if (x <= curMaxX) && (x > curMinX) && (y >= curMinY) && (y <= curMaxY) {
				(*relevantPoints)[i].HasCollision = true
			}
		}
	}
	doc.checkForCollInContainer(currentElem.Vertical, doc.checkForLeftColl, relevantPoints, elemToHandle)
	doc.checkForCollInContainer(currentElem.Horizontal, doc.checkForLeftColl, relevantPoints, elemToHandle)
}

func minMaxXY(currentElem *LayoutElement) (int, int, int, int) {
	curMinX := currentElem.X
	curMaxX := currentElem.X + currentElem.Width
	curMinY := currentElem.Y
	curMaxY := currentElem.Y + currentElem.Height
	return curMinX, curMaxX, curMinY, curMaxY
}

func (doc *BoxesDocument) checkForRightColl(relevantPoints *[]PointToTest, currentElem, elemToHandle *LayoutElement) {
	if (currentElem != elemToHandle) && doc.ShouldHandle(currentElem) && (!doc.isParent(currentElem, elemToHandle)) {
		curMinX, curMaxX, curMinY, curMaxY := minMaxXY(currentElem)
		for i := 0; i < len(*relevantPoints); i++ {
			x := (*relevantPoints)[i].X + (4 * types.RasterSize)
			y := (*relevantPoints)[i].Y
			if (x >= curMinX) && (x < curMaxX) && (y >= curMinY) && (y <= curMaxY) {
				(*relevantPoints)[i].HasCollision = true
			}
		}
	}
	doc.checkForCollInContainer(currentElem.Vertical, doc.checkForRightColl, relevantPoints, elemToHandle)
	doc.checkForCollInContainer(currentElem.Horizontal, doc.checkForRightColl, relevantPoints, elemToHandle)
}

func (doc *BoxesDocument) checkForAboveColl(relevantPoints *[]PointToTest, currentElem, elemToHandle *LayoutElement) {
	if (currentElem != elemToHandle) && doc.ShouldHandle(currentElem) && (!doc.isParent(currentElem, elemToHandle)) {
		curMinX, curMaxX, curMinY, curMaxY := minMaxXY(currentElem)
		for i := 0; i < len(*relevantPoints); i++ {
			x := (*relevantPoints)[i].X
			y := (*relevantPoints)[i].Y - (4 * types.RasterSize)
			if (y <= curMaxY) && (y > curMinY) && (x >= curMinX) && (x <= curMaxX) {
				(*relevantPoints)[i].HasCollision = true
			}
		}
	}
	doc.checkForCollInContainer(currentElem.Vertical, doc.checkForAboveColl, relevantPoints, elemToHandle)
	doc.checkForCollInContainer(currentElem.Horizontal, doc.checkForAboveColl, relevantPoints, elemToHandle)
}

func (doc *BoxesDocument) checkForBelowColl(relevantPoints *[]PointToTest, currentElem, elemToHandle *LayoutElement) {
	if (currentElem != elemToHandle) && doc.ShouldHandle(currentElem) && (!doc.isParent(currentElem, elemToHandle)) {
		curMinX, curMaxX, curMinY, curMaxY := minMaxXY(currentElem)
		for i := 0; i < len(*relevantPoints); i++ {
			x := (*relevantPoints)[i].X
			y := (*relevantPoints)[i].Y + (4 * types.RasterSize)
			if (y >= curMinY) && (y < curMaxY) && (x >= curMinX) && (x <= curMaxX) {
				(*relevantPoints)[i].HasCollision = true
			}
		}
	}
	doc.checkForCollInContainer(currentElem.Vertical, doc.checkForBelowColl, relevantPoints, elemToHandle)
	doc.checkForCollInContainer(currentElem.Horizontal, doc.checkForBelowColl, relevantPoints, elemToHandle)
}

func abs(v1 int) int {
	if v1 < 0 {
		return -v1
	}
	return v1
}

func valWithShortestDistanceToCenter(centerVal, currentVal, newVal int) int {
	if currentVal == -1 {
		return newVal
	}
	if abs(centerVal-newVal) < abs(centerVal-currentVal) {
		return newVal
	}
	return currentVal
}

// finds the best Y starting point for connections going from/to the right border
func (doc *BoxesDocument) initRightYToStart(elem *LayoutElement) {
	// Attention, this function works only with three raster points between boxes

	// the general algorithm is to take the relevant raster points from the border
	// and check the document elements if the point coordinate + 4 raster points in
	// the specific direction are free ... if so then the point is used. In case of
	// multiple matching points the one closest to the center of the element is used.
	relevantPoints := getRelevantPoints(elem.Y, elem.Y+elem.Height, elem.X+elem.Width, true)
	doc.checkForRightColl(&relevantPoints, &doc.Boxes, elem)
	found := false
	curVal := -1
	for _, p := range relevantPoints {
		if !p.HasCollision {
			if found {
				// check if the point is closer to the center of the element
				curVal = valWithShortestDistanceToCenter(elem.CenterY, curVal, p.Y)
			} else {
				found = true
				continue
			}
		}
		found = false
	}
	if curVal > -1 {
		elem.RightYToStart = &curVal
	} else {
		elem.RightYToStart = &elem.CenterY
	}
}

// finds the best Y starting point for connections going from/to the left border
func (doc *BoxesDocument) initLeftYToStart(elem *LayoutElement) {
	relevantPoints := getRelevantPoints(elem.Y, elem.Y+elem.Height, elem.X, true)
	doc.checkForLeftColl(&relevantPoints, &doc.Boxes, elem)
	found := false
	curVal := -1
	for _, p := range relevantPoints {
		if !p.HasCollision {
			if found {
				// check if the point is closer to the center of the element
				curVal = valWithShortestDistanceToCenter(elem.CenterY, curVal, p.Y)
			} else {
				found = true
				continue
			}
		}
		found = false
	}
	if curVal > -1 {
		elem.LeftYToStart = &curVal
	} else {
		elem.LeftYToStart = &elem.CenterY
	}
}

// finds the best X starting point for connections going from/to the top border
func (doc *BoxesDocument) initTopXToStart(elem *LayoutElement) {
	relevantPoints := getRelevantPoints(elem.X, elem.X+elem.Width, elem.Y, false)
	doc.checkForAboveColl(&relevantPoints, &doc.Boxes, elem)
	found := false
	curVal := -1
	for _, p := range relevantPoints {
		if !p.HasCollision {
			if found {
				// check if the point is closer to the center of the element
				curVal = valWithShortestDistanceToCenter(elem.CenterX, curVal, p.X)
			} else {
				found = true
				continue
			}
		}
		found = false
	}
	if curVal > -1 {
		elem.TopXToStart = &curVal
	} else {
		elem.TopXToStart = &elem.CenterX
	}
}

// finds the best X starting point for connections going from/to the bottom border
func (doc *BoxesDocument) initBottomXToStart(elem *LayoutElement) {
	relevantPoints := getRelevantPoints(elem.X, elem.X+elem.Width, elem.Y+elem.Height, false)
	doc.checkForBelowColl(&relevantPoints, &doc.Boxes, elem)
	found := false
	curVal := -1
	for _, p := range relevantPoints {
		if !p.HasCollision {
			if found {
				// check if the point is closer to the center of the element
				curVal = valWithShortestDistanceToCenter(elem.CenterX, curVal, p.X)
			} else {
				found = true
				continue
			}
		}
		found = false
	}
	if curVal > -1 {
		elem.BottomXToStart = &curVal
	} else {
		elem.BottomXToStart = &elem.CenterX
	}
}

func emptyVariant(indexToRemove int, variants *[][]ConnectionLine) {
	(*variants)[indexToRemove] = make([]ConnectionLine, 0)
}

func newConnectionLine(x1, y1, x2, y2 int) ConnectionLine {
	return ConnectionLine{
		StartX: x1,
		StartY: y1,
		EndX:   x2,
		EndY:   y2,
	}
}

// start from the right border - connect from right side to the left side
func (doc *BoxesDocument) connectFromRightBorderToLeftBorder(startElem, destElem *LayoutElement) [][]ConnectionLine {
	// connect from left to right
	if startElem.RightYToStart == nil {
		doc.initRightYToStart(startElem)
	}
	if destElem.LeftYToStart == nil {
		doc.initLeftYToStart(destElem)
	}
	startX := startElem.X + startElem.Width
	startY := *startElem.RightYToStart
	endX := destElem.X
	destY := destElem.LeftYToStart
	variant := []ConnectionLine{
		newConnectionLine(startX, startY, startX+types.RasterSize, startY)}

	ret, err := doc.goToRight(endX, *destY, variant, startElem, destElem)
	if err != nil {
		ret = make([][]ConnectionLine, 0)
	}
	return ret
}

// start from the right border - connect from left side to the right side
func (doc *BoxesDocument) connectFromLeftBorderToRightBorder(startElem, destElem *LayoutElement) [][]ConnectionLine {
	// connect from right to left
	if startElem.LeftYToStart == nil {
		doc.initLeftYToStart(startElem)
	}
	if destElem.LeftYToStart == nil {
		doc.initRightYToStart(destElem)
	}
	startX := startElem.X
	startY := startElem.LeftYToStart
	endX := destElem.X + destElem.Width
	destY := destElem.RightYToStart
	variant := []ConnectionLine{
		newConnectionLine(startX, *startY, startX-types.RasterSize, *startY)}
	ret, err := doc.goToLeft(endX, *destY, variant, startElem, destElem)
	if err != nil {
		ret = make([][]ConnectionLine, 0)
	}
	return ret
}

// start from the top - connect from top to top, left ro right ...
func (doc *BoxesDocument) connectFromTopBorderToRight(startElem, destElem *LayoutElement) [][]ConnectionLine {
	// horizontal connection to the right, starting at the top ... most likely a U from
	if startElem.TopXToStart == nil {
		doc.initTopXToStart(startElem)
	}
	if destElem.RightYToStart == nil {
		doc.initRightYToStart(destElem)
	}

	startX := startElem.TopXToStart
	startY := startElem.Y
	endX := destElem.X + destElem.Width
	destY := destElem.RightYToStart

	variant := []ConnectionLine{
		newConnectionLine(*startX, startY, *startX+types.RasterSize, startY)}
	ret, err := doc.goToRight(endX, *destY, variant, startElem, destElem)
	if err != nil {
		ret = make([][]ConnectionLine, 0)
	}
	return ret
}

// start from the bottom - connect from bottom to bottom, left ro right
func (doc *BoxesDocument) connectFromBottomBorderToRight(startElem, destElem *LayoutElement) [][]ConnectionLine {
	// horizontal connection to the right, starting at the bottom ... most likely a U from
	if startElem.BottomXToStart == nil {
		doc.initBottomXToStart(startElem)
	}
	if destElem.BottomXToStart == nil {
		doc.initBottomXToStart(destElem)
	}

	startX := startElem.BottomXToStart
	startY := startElem.Y + startElem.Height
	endX := destElem.BottomXToStart
	destY := destElem.Y + destElem.Height

	variant := []ConnectionLine{
		newConnectionLine(*startX, startY, *startX+types.RasterSize, startY)}
	ret, err := doc.goToRight(*endX, destY, variant, startElem, destElem)
	if err != nil {
		ret = make([][]ConnectionLine, 0)
	}
	return ret
}

// start from the top - connect from top to top, right to left
func (doc *BoxesDocument) connectFromTopBorderToLeft(startElem, destElem *LayoutElement) [][]ConnectionLine {
	// horizontal connection to the left, starting at the top ... most likely a U from
	if startElem.TopXToStart == nil {
		doc.initTopXToStart(startElem)
	}
	if destElem.TopXToStart == nil {
		doc.initTopXToStart(destElem)
	}

	startX := startElem.TopXToStart
	startY := startElem.Y
	endX := destElem.TopXToStart
	destY := destElem.Y

	variant := []ConnectionLine{
		newConnectionLine(*startX, startY, *startX-types.RasterSize, startY)}
	ret, err := doc.goToLeft(*endX, destY, variant, startElem, destElem)
	if err != nil {
		ret = make([][]ConnectionLine, 0)
	}
	return ret
}

// start from the bottom - connect from bottom to bottom, right to left
func (doc *BoxesDocument) connectFromBottomBorderLeft(startElem, destElem *LayoutElement) [][]ConnectionLine {
	// horizontal connection to the left, starting at the bottom ... most likely a U from
	if startElem.BottomXToStart == nil {
		doc.initBottomXToStart(startElem)
	}
	if destElem.BottomXToStart == nil {
		doc.initBottomXToStart(destElem)
	}

	startX := startElem.BottomXToStart
	startY := startElem.Y + startElem.Height
	endX := destElem.BottomXToStart
	destY := destElem.Y + destElem.Height

	variant := []ConnectionLine{
		newConnectionLine(*startX, startY, *startX-types.RasterSize, startY)}

	ret, err := doc.goToLeft(*endX, destY, variant, startElem, destElem)
	if err != nil {
		ret = make([][]ConnectionLine, 0)
	}
	return ret
}

// connect from with straight line: bottom to top
func (doc *BoxesDocument) connectFromBottomBorderToTopBorder(startElem, destElem *LayoutElement) [][]ConnectionLine {
	if startElem.BottomXToStart == nil {
		doc.initBottomXToStart(startElem)
	}
	if destElem.TopXToStart == nil {
		doc.initTopXToStart(destElem)
	}
	startX := startElem.BottomXToStart
	startY := startElem.Y + startElem.Height
	endX := destElem.TopXToStart
	destY := destElem.Y

	variant := []ConnectionLine{
		newConnectionLine(*startX, startY, *startX, startY+types.RasterSize)}
	ret, err := doc.goToDown(*endX, destY, variant, startElem, destElem)
	if err != nil {
		ret = make([][]ConnectionLine, 0)
	}
	return ret
}

// connect from left to left
func (doc *BoxesDocument) connectFromLeftBorderDown(startElem, destElem *LayoutElement) [][]ConnectionLine {
	if startElem.LeftYToStart == nil {
		doc.initLeftYToStart(startElem)
	}
	if destElem.LeftYToStart == nil {
		doc.initLeftYToStart(destElem)
	}

	startX := startElem.X
	startY := startElem.LeftYToStart
	endX := destElem.X
	destY := destElem.LeftYToStart

	variant := []ConnectionLine{
		newConnectionLine(startX, *startY, startX, *startY+types.RasterSize)}
	ret, err := doc.goToDown(endX, *destY, variant, startElem, destElem)
	if err != nil {
		ret = make([][]ConnectionLine, 0)
	}
	return ret
}

// connect from right to right
func (doc *BoxesDocument) connectFromRightBorderDown(startElem, destElem *LayoutElement) [][]ConnectionLine {
	if startElem.RightYToStart == nil {
		doc.initRightYToStart(startElem)
	}
	if destElem.RightYToStart == nil {
		doc.initRightYToStart(destElem)
	}

	startX := startElem.X + startElem.Width
	startY := startElem.RightYToStart
	endX := destElem.X + destElem.Width
	destY := destElem.RightYToStart

	variant := []ConnectionLine{
		newConnectionLine(startX, *startY, startX, *startY+types.RasterSize)}
	ret, err := doc.goToDown(endX, *destY, variant, startElem, destElem)
	if err != nil {
		ret = make([][]ConnectionLine, 0)
	}
	return ret
}

// connect from with straight line: top to bottom
func (doc *BoxesDocument) connectFromTopBorderToBottomBorder(startElem, destElem *LayoutElement) [][]ConnectionLine {
	if startElem.TopXToStart == nil {
		doc.initTopXToStart(startElem)
	}
	if destElem.BottomXToStart == nil {
		doc.initBottomXToStart(destElem)
	}
	startX := startElem.TopXToStart
	startY := startElem.Y
	endX := destElem.BottomXToStart
	destY := destElem.Y + destElem.Height

	variant := []ConnectionLine{
		newConnectionLine(*startX, startY, *startX, startY-types.RasterSize)}

	ret, err := doc.goToUp(*endX, destY, variant, startElem, destElem)
	if err != nil {
		ret = make([][]ConnectionLine, 0)
	}
	return ret
}

// connect from left to left
func (doc *BoxesDocument) connectFromLeftBorderUp(startElem, destElem *LayoutElement) [][]ConnectionLine {
	if startElem.LeftYToStart == nil {
		doc.initLeftYToStart(startElem)
	}
	if destElem.LeftYToStart == nil {
		doc.initLeftYToStart(destElem)
	}

	startX := startElem.X
	startY := startElem.LeftYToStart
	endX := destElem.X
	destY := destElem.LeftYToStart

	variant := []ConnectionLine{
		newConnectionLine(startX, *startY, startX, *startY-types.RasterSize)}
	ret, err := doc.goToUp(endX, *destY, variant, startElem, destElem)
	if err != nil {
		ret = make([][]ConnectionLine, 0)
	}
	return ret
}

// connect from right to right
func (doc *BoxesDocument) connectFromRightBorderUp(startElem, destElem *LayoutElement) [][]ConnectionLine {
	if startElem.RightYToStart == nil {
		doc.initRightYToStart(startElem)
	}
	if destElem.RightYToStart == nil {
		doc.initRightYToStart(destElem)
	}

	startX := startElem.X + startElem.Width
	startY := startElem.RightYToStart
	endX := destElem.X + destElem.Width
	destY := destElem.RightYToStart

	variant := []ConnectionLine{
		newConnectionLine(startX, *startY, startX, *startY-types.RasterSize)}
	ret, err := doc.goToUp(endX, *destY, variant, startElem, destElem)
	if err != nil {
		ret = make([][]ConnectionLine, 0)
	}
	return ret
}

// connect from bottom to left side
func (doc *BoxesDocument) connectFromBottomBorderToLeftBorder(startElem, destElem *LayoutElement) [][]ConnectionLine {
	if startElem.BottomXToStart == nil {
		doc.initBottomXToStart(startElem)
	}
	if destElem.LeftYToStart == nil {
		doc.initLeftYToStart(destElem)
	}
	startX := startElem.BottomXToStart
	startY := startElem.Y + startElem.Height
	endX := destElem.X
	destY := destElem.LeftYToStart

	variant := []ConnectionLine{
		newConnectionLine(*startX, startY, *startX, startY+types.RasterSize)}
	ret, err := doc.goToDown(endX, *destY, variant, startElem, destElem)
	if err != nil {
		ret = make([][]ConnectionLine, 0)
	}
	return ret
}

// connect from right to top side
func (doc *BoxesDocument) connectFromRightBorderToTopBorder(startElem, destElem *LayoutElement) [][]ConnectionLine {
	if startElem.RightYToStart == nil {
		doc.initRightYToStart(startElem)
	}
	if destElem.TopXToStart == nil {
		doc.initTopXToStart(destElem)
	}
	startX := startElem.X + startElem.Width
	startY := startElem.RightYToStart
	endX := destElem.TopXToStart
	destY := destElem.Y

	variant := []ConnectionLine{
		newConnectionLine(startX, *startY, startX+types.RasterSize, *startY)}
	ret, err := doc.goToRight(*endX, destY, variant, startElem, destElem)
	if err != nil {
		ret = make([][]ConnectionLine, 0)
	}
	return ret
}

// connect from bottom to right side
func (doc *BoxesDocument) connectFromBottomBorderToRightBorder(startElem, destElem *LayoutElement) [][]ConnectionLine {
	if startElem.BottomXToStart == nil {
		doc.initBottomXToStart(startElem)
	}
	if destElem.RightYToStart == nil {
		doc.initRightYToStart(destElem)
	}
	startX := startElem.BottomXToStart
	startY := startElem.Y + startElem.Height
	endX := destElem.X
	destY := destElem.RightYToStart

	variant := []ConnectionLine{
		newConnectionLine(*startX, startY, *startX, startY+types.RasterSize)}
	ret, err := doc.goToDown(endX, *destY, variant, startElem, destElem)
	if err != nil {
		ret = make([][]ConnectionLine, 0)
	}
	return ret
}

// connect from left to top side
func (doc *BoxesDocument) connectFromLeftBorderToTopBorder(startElem, destElem *LayoutElement) [][]ConnectionLine {
	if startElem.LeftYToStart == nil {
		doc.initLeftYToStart(startElem)
	}
	if destElem.TopXToStart == nil {
		doc.initTopXToStart(destElem)
	}
	startX := startElem.X
	startY := startElem.LeftYToStart
	endX := destElem.TopXToStart
	destY := destElem.Y

	variant := []ConnectionLine{
		newConnectionLine(startX, *startY, startX-types.RasterSize, *startY)}
	ret, err := doc.goToLeft(*endX, destY, variant, startElem, destElem)
	if err != nil {
		ret = make([][]ConnectionLine, 0)
	}
	return ret
}

// connect from top to left side
func (doc *BoxesDocument) connectFromTopBorderToLeftBorder(startElem, destElem *LayoutElement) [][]ConnectionLine {
	if startElem.TopXToStart == nil {
		doc.initTopXToStart(startElem)
	}
	if destElem.LeftYToStart == nil {
		doc.initLeftYToStart(destElem)
	}
	startX := startElem.TopXToStart
	startY := startElem.Y
	endX := destElem.X
	destY := destElem.LeftYToStart

	variant := []ConnectionLine{
		newConnectionLine(*startX, startY, *startX, startY-types.RasterSize)}
	ret, err := doc.goToUp(endX, *destY, variant, startElem, destElem)
	if err != nil {
		ret = make([][]ConnectionLine, 0)
	}
	return ret
}

// connect from right to bottom side
func (doc *BoxesDocument) connectFromRightBorderToBottomBorder(startElem, destElem *LayoutElement) [][]ConnectionLine {
	if startElem.RightYToStart == nil {
		doc.initRightYToStart(startElem)
	}
	if destElem.BottomXToStart == nil {
		doc.initBottomXToStart(destElem)
	}
	startX := startElem.X + startElem.Width
	startY := startElem.RightYToStart
	endX := destElem.BottomXToStart
	destY := destElem.Y + destElem.Height

	variant := []ConnectionLine{
		newConnectionLine(startX, *startY, startX+types.RasterSize, *startY)}
	ret, err := doc.goToRight(*endX, destY, variant, startElem, destElem)
	if err != nil {
		ret = make([][]ConnectionLine, 0)
	}
	return ret
}

// connect from top to right side
func (doc *BoxesDocument) connectFromTopBorderToRightBorder(startElem, destElem *LayoutElement) [][]ConnectionLine {
	if startElem.TopXToStart == nil {
		doc.initTopXToStart(startElem)
	}
	if destElem.RightYToStart == nil {
		doc.initRightYToStart(destElem)
	}
	startX := startElem.TopXToStart
	startY := startElem.Y
	endX := destElem.X
	destY := destElem.RightYToStart

	variant := []ConnectionLine{
		newConnectionLine(*startX, startY, *startX, startY-types.RasterSize)}
	ret, err := doc.goToUp(endX, *destY, variant, startElem, destElem)
	if err != nil {
		ret = make([][]ConnectionLine, 0)
	}
	return ret
}

// connect from left to bottom side
func (doc *BoxesDocument) connectFromLeftBorderToBottomBorder(startElem, destElem *LayoutElement) [][]ConnectionLine {
	if startElem.LeftYToStart == nil {
		doc.initLeftYToStart(startElem)
	}
	if destElem.BottomXToStart == nil {
		doc.initBottomXToStart(destElem)
	}
	startX := startElem.X
	startY := startElem.LeftYToStart
	endX := destElem.BottomXToStart
	destY := destElem.Y + destElem.Height

	variant := []ConnectionLine{
		newConnectionLine(startX, *startY, startX-types.RasterSize, *startY)}
	ret, err := doc.goToLeft(*endX, destY, variant, startElem, destElem)
	if err != nil {
		ret = make([][]ConnectionLine, 0)
	}
	return ret
}

func (doc *BoxesDocument) getConnectionVariants(startElem, destElem *LayoutElement) [][]ConnectionLine {
	connectionVariants := make([][]ConnectionLine, 0)
	if startElem.CenterY == destElem.CenterY {
		// horizontal connection
		if startElem.CenterX < destElem.CenterX {
			// connect from left to right
			// start from the right border - connect from right side to the left side
			v := doc.connectFromRightBorderToLeftBorder(startElem, destElem)
			connectionVariants = append(connectionVariants, v...)
			// start from the top - connect from top to top, left ro right
			v = doc.connectFromTopBorderToRight(startElem, destElem)
			connectionVariants = append(connectionVariants, v...)
			// start from the bottom - connect from bottom to bottom, left ro right
			v = doc.connectFromBottomBorderToRight(startElem, destElem)
			connectionVariants = append(connectionVariants, v...)
		} else {
			// connect from right to left
			// start from the right border - connect from left side to the right side
			v := doc.connectFromLeftBorderToRightBorder(startElem, destElem)
			connectionVariants = append(connectionVariants, v...)
			// start from the top - connect from top to top, right to left
			v = doc.connectFromTopBorderToLeft(startElem, destElem)
			connectionVariants = append(connectionVariants, v...)
			// start from the bottom - connect from bottom to bottom, right to left
			v = doc.connectFromBottomBorderLeft(startElem, destElem)
			connectionVariants = append(connectionVariants, v...)
		}
	} else if startElem.CenterX == destElem.CenterX {
		// vertical connection
		if startElem.CenterY < destElem.CenterY {
			// connect from with straight line: bottom to top
			v := doc.connectFromBottomBorderToTopBorder(startElem, destElem)
			connectionVariants = append(connectionVariants, v...)
			// connect from left to left
			v = doc.connectFromLeftBorderDown(startElem, destElem)
			connectionVariants = append(connectionVariants, v...)
			// connect from right to right
			v = doc.connectFromRightBorderDown(startElem, destElem)
			connectionVariants = append(connectionVariants, v...)
		} else {
			// connect from with straight line: top to bottom
			v := doc.connectFromTopBorderToBottomBorder(startElem, destElem)
			connectionVariants = append(connectionVariants, v...)
			// connect from left to left
			v = doc.connectFromLeftBorderUp(startElem, destElem)
			connectionVariants = append(connectionVariants, v...)
			// connect from right to right
			v = doc.connectFromRightBorderUp(startElem, destElem)
			connectionVariants = append(connectionVariants, v...)
		}
	} else if startElem.CenterY < destElem.CenterY {
		// top -> down: connect from bottom to top side
		v := doc.connectFromBottomBorderToTopBorder(startElem, destElem)
		connectionVariants = append(connectionVariants, v...)
		if startElem.CenterX < destElem.CenterX {
			// connect from bottom to left side
			v = doc.connectFromBottomBorderToLeftBorder(startElem, destElem)
			connectionVariants = append(connectionVariants, v...)
			// connect from right to left side
			v = doc.connectFromRightBorderToLeftBorder(startElem, destElem)
			connectionVariants = append(connectionVariants, v...)
			// connect from right to top side
			v := doc.connectFromRightBorderToTopBorder(startElem, destElem)
			connectionVariants = append(connectionVariants, v...)
		} else {
			// connect from bottom to right side
			v = doc.connectFromBottomBorderToRightBorder(startElem, destElem)
			connectionVariants = append(connectionVariants, v...)
			// connect from left to right side
			v = doc.connectFromLeftBorderToRightBorder(startElem, destElem)
			connectionVariants = append(connectionVariants, v...)
			// connect from left to top side
			v := doc.connectFromLeftBorderToTopBorder(startElem, destElem)
			connectionVariants = append(connectionVariants, v...)
		}
	} else {
		// down -> top: connect from top to bottom side
		v := doc.connectFromTopBorderToBottomBorder(startElem, destElem)
		connectionVariants = append(connectionVariants, v...)
		if startElem.CenterX < destElem.CenterX {
			// connect from top to left side
			v = doc.connectFromTopBorderToLeftBorder(startElem, destElem)
			connectionVariants = append(connectionVariants, v...)
			// connect from right to left side
			v = doc.connectFromRightBorderToLeftBorder(startElem, destElem)
			connectionVariants = append(connectionVariants, v...)
			// connect from right to bottom side
			v = doc.connectFromRightBorderToBottomBorder(startElem, destElem)
			connectionVariants = append(connectionVariants, v...)
		} else {
			// connect from top to right side
			v = doc.connectFromTopBorderToRightBorder(startElem, destElem)
			connectionVariants = append(connectionVariants, v...)
			// connect from left to right side
			v = doc.connectFromLeftBorderToRightBorder(startElem, destElem)
			connectionVariants = append(connectionVariants, v...)
			//connect from left to bottom side
			v = doc.connectFromLeftBorderToBottomBorder(startElem, destElem)
			connectionVariants = append(connectionVariants, v...)
		}
	}
	return connectionVariants
}

func isConnectedToDest(line *ConnectionLine, destElem *LayoutElement) bool {
	// either connected to top
	if line.EndX == *destElem.TopXToStart && line.EndY == destElem.Y {
		return true
	}
	// or to bottom
	if line.EndX == *destElem.BottomXToStart && line.EndY == destElem.Y+destElem.Height {
		return true
	}
	// or to left
	if line.EndX == destElem.X && line.EndY == *destElem.LeftYToStart {
		return true
	}
	// or to right
	if line.EndX == destElem.X+destElem.Width && line.EndY == *destElem.RightYToStart {
		return true
	}
	return false
}

func (doc *BoxesDocument) connectTwoElems(start, destElem *LayoutElement, lec *LayoutElemConnection) {
	variants := doc.getConnectionVariants(start, destElem)
	var connection []ConnectionLine
	for _, conn := range variants {
		l := len(conn)
		if l == 0 {
			continue
		}
		lastLine := conn[l-1]
		if l == 2 {
			fmt.Println("DEBUG: updated connection: ", conn)
		}
		if !isConnectedToDest(&lastLine, destElem) {
			continue
		}
		fmt.Println("DEBUG: updated connection: ", conn)
		if connection == nil || (len(conn) < len(connection)) {
			fmt.Println("DEBUG: updated connection: ", conn)
			connection = conn
		}
	}
	var ret ConnectionElem
	ret.DestArrow = &lec.DestArrow
	ret.SourceArrow = &lec.SourceArrow
	ret.Parts = connection
	doc.Connections = append(doc.Connections, ret)
}

func (doc *BoxesDocument) connectTwoElemsFull(start, destElem *LayoutElement, lec *LayoutElemConnection) {
	variants := doc.getConnectionVariants(start, destElem)
	for _, v := range variants {
		var ret ConnectionElem
		ret.DestArrow = &lec.DestArrow
		ret.SourceArrow = &lec.SourceArrow
		ret.Parts = v
		doc.Connections = append(doc.Connections, ret)
	}
}

func (doc *BoxesDocument) doConnect(elem *LayoutElement, full bool) {
	for _, conn := range elem.Connections {
		destElem, found := doc.findLayoutElementById(conn.DestId, &doc.Boxes)
		if !found {
			fmt.Println("Couldn't find destId: ", conn.DestId)
			continue
		}
		if destElem == elem {
			fmt.Println("Connection to self are not allowed and will be ignored: ", conn.DestId)
			continue
		}
		if full {
			// used to debug
			doc.connectTwoElemsFull(elem, destElem, &conn)
		} else {
			doc.connectTwoElems(elem, destElem, &conn)
		}
	}
}

func (doc *BoxesDocument) connectLayoutElem(le *LayoutElement, full bool) {
	doc.doConnect(le, full)
	doc.connectContainer(le.Vertical, full)
	doc.connectContainer(le.Horizontal, full)
}

func (doc *BoxesDocument) connectContainer(cont *LayoutElemContainer, full bool) {
	if cont != nil {
		for _, elem := range cont.Elems {
			doc.connectLayoutElem(&elem, full)
		}
	}
}

func (doc *BoxesDocument) _ConnectBoxes() {
	// TODO - Needs reimplementation
	doc.InitStartPositions()
	doc.InitRoads()
	// doc.connectLayoutElem(&doc.Boxes, false)
	// // doc.moveTooCloseVerticalConnectionLinesFromBorders()
	// // doc.moveTooCloseHorizontalConnectionLinesFromBorders()
	// // doc.moveTooCloseVerticalConnectionLines()
	// // doc.moveTooCloseHorizontalConnectionLines()
	// doc.truncateJoiningConnectionLines()
}

func (doc *BoxesDocument) _ConnectBoxesFull() {
	doc.InitStartPositions()
	doc.InitRoads()
	doc.connectLayoutElem(&doc.Boxes, true)
}

func (doc *BoxesDocument) InitStartPositions() {
	doc.initStartPositionsImpl(&doc.Boxes)
}

func (doc *BoxesDocument) initStartPositionsImpl(elem *LayoutElement) {
	if doc.ShouldHandle(elem) {
		doc.initBottomXToStart(elem)
		doc.initLeftYToStart(elem)
		doc.initRightYToStart(elem)
		doc.initTopXToStart(elem)
	}
	if elem.Vertical != nil {
		for i := 0; i < len(elem.Vertical.Elems); i++ {
			doc.initStartPositionsImpl(&elem.Vertical.Elems[i])
		}
	}
	if elem.Horizontal != nil {
		for i := 0; i < len(elem.Horizontal.Elems); i++ {
			doc.initStartPositionsImpl(&elem.Horizontal.Elems[i])
		}
	}
}
