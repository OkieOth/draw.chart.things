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
	for _, h := range doc.HorizontalRoads {
		// horizontal line
		newH := true
		minHX, maxHX := minMax(h.StartX, h.EndX)
		for j, v := range doc.VerticalRoads {
			// verticalLine
			minVY, maxVY := minMax(v.StartY, v.EndY)
			if (minVY <= h.StartY) && (maxVY >= h.StartY) && // the vertical road covers the y-range of the horizontal line
				(minHX <= v.StartX) && (v.StartX <= maxHX) {
				newNode := CreateConnectionNode(v.StartX, h.StartY)
				if j > 0 && !newH {
					newEdge := ConnectionEdge{
						X: v.StartX,
						Y: h.StartY,
					}
					doc.ConnectionNodes[len(doc.ConnectionNodes)-1].Edges = append(doc.ConnectionNodes[len(doc.ConnectionNodes)-1].Edges, newEdge)
				}
				doc.ConnectionNodes = append(doc.ConnectionNodes, *newNode)
				newH = false
			}
		}
	}

	// adding vertical edges
	for _, v := range doc.VerticalRoads {
		// vertical line
		newV := true
		minVY, maxVY := minMax(v.StartY, v.EndY)
		var lastY int
		for j, h := range doc.HorizontalRoads {
			// verticalLine
			minHX, maxHX := minMax(h.StartX, h.EndX)
			if (minHX <= v.StartX) && (maxHX >= v.StartX) && // the vertical road covers the y-range of the horizontal line
				(minVY <= h.StartY) && (h.StartY <= maxVY) {
				nodeX, nodeY := v.StartX, h.StartY
				if j > 0 && !newV {
					if nodeIndex := doc.getMatchingNodeIndexes(nodeX, nodeY); nodeIndex > 0 {
						newEdge := ConnectionEdge{
							X: v.StartX,
							Y: lastY,
						}
						doc.ConnectionNodes[nodeIndex].Edges = append(doc.ConnectionNodes[nodeIndex].Edges, newEdge)
					}
				}
				lastY = h.StartY
				newV = false
			}
		}
	}

}

func (doc *BoxesDocument) getMatchingNodeIndexes(x, y int) int {
	for i, n := range doc.ConnectionNodes {
		if n.X == x && n.Y == y {
			return i
		}
	}
	return -1
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
			n := CreateConnectionNode(*elem.BottomXToStart, elem.Y)
			n.Edges = append(n.Edges, ConnectionEdge{
				X: *elem.BottomXToStart,
				Y: elem.Y - (2 * types.RasterSize),
			})
			doc.ConnectionNodes = append(doc.ConnectionNodes, *n)

			// add bottomX
			n = CreateConnectionNode(*elem.BottomXToStart, elem.Y+elem.Height)
			n.Edges = append(n.Edges, ConnectionEdge{
				X: *elem.BottomXToStart,
				Y: elem.Y + elem.Height + (2 * types.RasterSize),
			})
			doc.ConnectionNodes = append(doc.ConnectionNodes, *n)
			// add leftY
			n = CreateConnectionNode(elem.X, *elem.LeftYToStart)
			n.Edges = append(n.Edges, ConnectionEdge{
				X: elem.X - (2 * types.RasterSize),
				Y: *elem.LeftYToStart,
			})
			doc.ConnectionNodes = append(doc.ConnectionNodes, *n)
			// add rightY
			n = CreateConnectionNode(elem.X+elem.Width, *elem.RightYToStart)
			n.Edges = append(n.Edges, ConnectionEdge{
				X: elem.X + elem.Width + (2 * types.RasterSize),
				Y: *elem.LeftYToStart,
			})
			doc.ConnectionNodes = append(doc.ConnectionNodes, *n)
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
