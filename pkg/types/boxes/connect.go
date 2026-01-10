package boxes

import (
	"slices"

	"github.com/okieoth/draw.chart.things/pkg/types"
)

type PointToTest struct {
	X            int
	Y            int
	HasCollision bool
}

type checkForCollFunc2 func(x, y int, currentElem, startElem, endElem *LayoutElement) CollisionType

func (doc *BoxesDocument) checkForCollInContainer2(cont *LayoutElemContainer, f checkForCollFunc2, x, y int, startElem, endElem *LayoutElement) CollisionType {
	if cont != nil {
		l := len(cont.Elems)
		for i := range l {
			e := &cont.Elems[i]
			if e == startElem || e == endElem {
				continue
			}
			if ct := f(x, y, e, startElem, endElem); ct != CollisionType_NoCollision {
				return ct
			}
		}
	}
	return CollisionType_NoCollision
}

type CollisionType int

const (
	CollisionType_NoCollision CollisionType = iota
	CollisionType_WithElem
	CollisionType_WithSurroundings
)

// checks if a point is inside a box, returns true if so
func (doc *BoxesDocument) checkColl(x, y int, currentElem, startElem, endElem *LayoutElement) CollisionType {
	if (currentElem != startElem) && (currentElem != endElem) &&
		(doc.ShouldHandle(currentElem)) {
		currentElemIsParentToStart := false
		if startElem != nil {
			currentElemIsParentToStart = doc.isParent(currentElem, startElem)
		}
		currentElemIsParentToEnd := false
		if endElem != nil {
			currentElemIsParentToEnd = doc.isParent(currentElem, endElem)
		}
		if !currentElemIsParentToStart && !currentElemIsParentToEnd {
			curMinX := currentElem.X
			curMaxX := currentElem.X + currentElem.Width
			curMinY := currentElem.Y
			curMaxY := currentElem.Y + currentElem.Height
			if (x <= curMaxX) && (x >= curMinX) &&
				(y >= curMinY) && (y <= curMaxY) {
				return CollisionType_WithElem
			}
			if (x <= (curMaxX + types.RasterSize)) && (x >= (curMinX - types.RasterSize)) &&
				(y >= (curMinY - types.RasterSize)) && (y <= (curMaxY + types.RasterSize)) {
				return CollisionType_WithSurroundings
			}
		}

		if currentElemIsParentToStart {
			// check if there is a collision with the text
			if currentElem.WidthTextBox != nil {
				if (x <= (*currentElem.XTextBox + *currentElem.WidthTextBox)) && (x >= *currentElem.XTextBox) &&
					(y >= *currentElem.YTextBox) && (y <= (*currentElem.YTextBox + *currentElem.HeightTextBox)) {
					return CollisionType_WithElem
				}
			}
		}
	}
	if ct := doc.checkForCollInContainer2(currentElem.Vertical, doc.checkColl, x, y, startElem, endElem); ct != CollisionType_NoCollision {
		return ct
	}
	return doc.checkForCollInContainer2(currentElem.Horizontal, doc.checkColl, x, y, startElem, endElem)
}

func newConnectionLine(x1, y1, x2, y2 int) ConnectionLine {
	return ConnectionLine{
		StartX: x1,
		StartY: y1,
		EndX:   x2,
		EndY:   y2,
	}
}

func (doc *BoxesDocument) ConnectBoxes() {
	doc.InitStartPositions()
	doc.InitRoads()
	doc.Roads2ConnectionNodes()
	// TODO - Needs reimplementation
}

func (doc *BoxesDocument) Roads2ConnectionNodes() {
	// TODO
	for _, h := range doc.HorizontalRoads {
		// horizontal line
		minHX, maxHX := minMax(h.StartX, h.EndX)
		for _, v := range doc.VerticalRoads {
			// verticalLine
			minVY, maxVY := minMax(v.StartY, v.EndY)
			if (minVY <= h.StartY) && (maxVY >= h.StartY) && // the vertical road covers the y-range of the horizontal line
				(minHX <= v.StartX) && (v.StartX <= maxHX) {
				doc.ConnectionNodes = append(doc.ConnectionNodes, *CreateConnectionNode(v.StartX, h.StartY))
			}
		}
	}
}

func (doc *BoxesDocument) InitStartPositions() {
	doc.initStartPositionsImpl(&doc.Boxes)
}

func CreateConnectionNode(x, y int) *ConnectionNode {
	node := NewConnectionNode()
	node.X = x
	node.Y = y
	return node
}

func (doc *BoxesDocument) initStartPositionsImpl(elem *LayoutElement) {
	if doc.ShouldHandle(elem) {
		if slices.Contains(doc.ConnectedElems, elem.Id) {
			elem.BottomXToStart = &elem.CenterX
			elem.TopXToStart = &elem.CenterX
			elem.LeftYToStart = &elem.CenterY
			elem.RightYToStart = &elem.CenterY
			// add topX
			doc.ConnectionNodes = append(doc.ConnectionNodes, *CreateConnectionNode(*elem.BottomXToStart, elem.Y))
			// add bottomX
			doc.ConnectionNodes = append(doc.ConnectionNodes, *CreateConnectionNode(*elem.BottomXToStart, elem.Y+elem.Height))
			// add leftY
			doc.ConnectionNodes = append(doc.ConnectionNodes, *CreateConnectionNode(elem.X, *elem.LeftYToStart))
			// add rightY
			doc.ConnectionNodes = append(doc.ConnectionNodes, *CreateConnectionNode(elem.X+elem.Width, *elem.RightYToStart))
		}
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
