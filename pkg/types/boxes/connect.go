package boxes

import (
	"fmt"
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

func (doc *BoxesDocument) connectContImpl(layoutCont []LayoutElement) {
	for i := range len(layoutCont) {
		doc.connectImpl(&layoutCont[i])
	}
}

func (doc *BoxesDocument) createAConnectionPath(path []ConnectionNode) {
	if len(path) < 2 {
		return
	}
	pathToDraw := path[:len(path)-1]
	connElem := NewConnectionElem()
	connElem.From = path[0].NodeId
	connElem.To = path[0].NodeId
	var lastX, lastY int
	for i, p := range pathToDraw {
		if i > 0 {
			line := ConnectionLine{
				StartX: lastX,
				StartY: lastY,
				EndX:   p.X,
				EndY:   p.Y,
			}
			connElem.Parts = append(connElem.Parts, line)
		}
		lastX = p.X
		lastY = p.Y
	}
	doc.Connections = append(doc.Connections, *connElem)
}

func (doc *BoxesDocument) connectImpl(layout *LayoutElement) {
	if len(layout.Connections) > 0 {
		for _, c := range layout.Connections {
			srcId, destId := layout.Id, c.DestId
			if path, dist, ok := doc.DijkstraPath(srcId, destId); ok {
				fmt.Printf("Found path: src=%s, dest=%s, dist=%d\n", srcId, destId, dist)
				doc.createAConnectionPath(path)
			} else {
				fmt.Printf("Couldn't calculate path: src=%s, dest=%s\n", srcId, destId)
			}
		}
	}
	if layout.Horizontal != nil {
		doc.connectContImpl(layout.Horizontal.Elems)
	}
	if layout.Vertical != nil {
		doc.connectContImpl(layout.Vertical.Elems)
	}
}

func (doc *BoxesDocument) ConnectBoxes() {
	doc.InitStartPositions()
	doc.InitRoads()
	doc.Roads2ConnectionNodes()
	doc.connectImpl(&doc.Boxes)
}

func (doc *BoxesDocument) horizontalRoads2ConnectionNodes() {
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
				newNodeIndex := len(doc.ConnectionNodes)
				id := fmt.Sprintf("__n_%d", newNodeIndex)
				newNode.NodeId = &id
				if j > 0 && !newH {
					weight := v.StartX - doc.ConnectionNodes[len(doc.ConnectionNodes)-1].X
					newEdge := CreateConnectionEdge(v.StartX, h.StartY, weight)
					newEdge.DestNodeId = &id
					newEdge.DestNodeIndex = &newNodeIndex
					doc.ConnectionNodes[len(doc.ConnectionNodes)-1].Edges = append(doc.ConnectionNodes[len(doc.ConnectionNodes)-1].Edges, newEdge)
					newEdge2 := CreateConnectionEdge(doc.ConnectionNodes[len(doc.ConnectionNodes)-1].X, doc.ConnectionNodes[len(doc.ConnectionNodes)-1].Y, weight)
					newEdge2.DestNodeId = doc.ConnectionNodes[len(doc.ConnectionNodes)-1].NodeId
					lastNodeIndex := newNodeIndex - 1
					newEdge2.DestNodeIndex = &lastNodeIndex
					newNode.Edges = append(newNode.Edges, newEdge2)
				}
				doc.ConnectionNodes = append(doc.ConnectionNodes, *newNode)
				newH = false
			}
			if v.StartX > maxHX {
				// the other vertical lines are more right then the current horizontal line
				break
			}
		}
	}
}

func (doc *BoxesDocument) verticalEdges() {
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
					nodeIndex := doc.getMatchingNodeIndexes(nodeX, nodeY)
					lastNodeIndex := doc.getMatchingNodeIndexes(nodeX, lastY)
					if nodeIndex > -1 && lastNodeIndex > -1 {
						weight := nodeY - lastY
						newEdgeUp := CreateConnectionEdge(v.StartX, lastY, weight)
						newEdgeUp.DestNodeId = doc.ConnectionNodes[lastNodeIndex].NodeId
						newEdgeUp.DestNodeIndex = &lastNodeIndex
						doc.ConnectionNodes[nodeIndex].Edges = append(doc.ConnectionNodes[nodeIndex].Edges, newEdgeUp)

						newEdgeDown := CreateConnectionEdge(nodeX, nodeY, weight)
						newEdgeDown.DestNodeId = doc.ConnectionNodes[nodeIndex].NodeId
						newEdgeDown.DestNodeIndex = &nodeIndex
						doc.ConnectionNodes[lastNodeIndex].Edges = append(doc.ConnectionNodes[lastNodeIndex].Edges, newEdgeDown)
					}
				}
				j++
				lastY = h.StartY
				newV = false
			}
			if h.StartY > maxVY {
				// the other horizontal lines are more below the length of the current vertical line
				break
			}
		}
	}
}

func (doc *BoxesDocument) initEdgesForBoxConnections() {
	for i := range len(doc.ConnectionNodes) {
		if doc.ConnectionNodes[i].BoxId != nil {
			// node on a potential start point of a connected box found
			// ... should currently have only one outbound edge!!!
			nodeIndex := doc.getMatchingNodeIndexes(doc.ConnectionNodes[i].Edges[0].X, doc.ConnectionNodes[i].Edges[0].Y)
			if nodeIndex > -1 {
				weight := doc.ConnectionNodes[i].Edges[0].Weight
				doc.ConnectionNodes[i].Edges[0].DestNodeId = doc.ConnectionNodes[nodeIndex].NodeId
				doc.ConnectionNodes[i].Edges[0].DestNodeIndex = &nodeIndex
				newEdge := CreateConnectionEdge(doc.ConnectionNodes[i].X, doc.ConnectionNodes[i].Y, weight)
				newEdge.DestNodeId = doc.ConnectionNodes[i].NodeId
				newEdge.DestNodeIndex = &i
				doc.ConnectionNodes[nodeIndex].Edges = append(doc.ConnectionNodes[nodeIndex].Edges, newEdge)
			}
		}
	}
}

func (doc *BoxesDocument) Roads2ConnectionNodes() {
	doc.horizontalRoads2ConnectionNodes()
	doc.verticalEdges()
	doc.initEdgesForBoxConnections()
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

func CreateConnectionEdge(x, y, weight int) ConnectionEdge {
	return ConnectionEdge{
		X:      x,
		Y:      y,
		Weight: weight,
	}
}

func CreateConnectionEdge2(x, y int, id string, index int) ConnectionEdge {
	return ConnectionEdge{
		X:          x,
		Y:          y,
		DestNodeId: &id,
	}
}

func (doc *BoxesDocument) newConnectionNodeFromStartPos(boxId string, x, y int, edges []ConnectionEdge) {
	n := CreateConnectionNode(x, y)
	index := len(doc.ConnectionNodes)
	nodeId := fmt.Sprintf("__n_%d", index)
	n.NodeId = &nodeId
	n.BoxId = &boxId
	for _, e := range edges {
		n.Edges = append(n.Edges, e)
	}
	doc.ConnectionNodes = append(doc.ConnectionNodes, *n)
}

func (doc *BoxesDocument) initStartPositionsImpl(elem *LayoutElement) {
	if doc.ShouldHandle(elem) {
		if slices.Contains(doc.ConnectedElems, elem.Id) {
			elem.BottomXToStart = &elem.CenterX
			elem.TopXToStart = &elem.CenterX
			elem.LeftYToStart = &elem.CenterY
			elem.RightYToStart = &elem.CenterY
			// add topX
			weight := 2 * types.RasterSize
			doc.newConnectionNodeFromStartPos(elem.Id, *elem.TopXToStart, elem.Y,
				[]ConnectionEdge{
					CreateConnectionEdge(*elem.TopXToStart, elem.Y-weight, weight),
				})
			// add bottomX
			doc.newConnectionNodeFromStartPos(elem.Id, *elem.BottomXToStart, elem.Y+elem.Height,
				[]ConnectionEdge{
					CreateConnectionEdge(*elem.BottomXToStart, elem.Y+elem.Height+weight, weight),
				})
			// add leftY
			doc.newConnectionNodeFromStartPos(elem.Id, elem.X, *elem.LeftYToStart,
				[]ConnectionEdge{
					CreateConnectionEdge(elem.X-weight, *elem.LeftYToStart, weight),
				})
			// add rightY
			doc.newConnectionNodeFromStartPos(elem.Id, elem.X+elem.Width, *elem.RightYToStart,
				[]ConnectionEdge{
					CreateConnectionEdge(elem.X+elem.Width+weight, *elem.LeftYToStart, weight),
				})
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
