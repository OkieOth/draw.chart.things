package types

import (
	"fmt"
)

func (doc *BoxesDocument) findLayoutElementById(id string, startElem *LayoutElement) (*LayoutElement, bool) {
	if startElem.Id == id {
		return startElem, true
	}
	for _, elem := range startElem.Vertical.Elems {
		if elem.Id == id {
			return &elem, true
		}
		if e, found := doc.findLayoutElementById(id, &elem); found {
			return e, true
		}
	}
	for _, elem := range startElem.Horizontal.Elems {
		if elem.Id == id {
			return &elem, true
		}
		if e, found := doc.findLayoutElementById(id, &elem); found {
			return e, true
		}
	}
	return nil, false
}

func (doc *BoxesDocument) checkForCollisionLeft(startX, startY, endX, endY int, startElem *LayoutElement) (int, int, *ConnDirection, int, int) {
	closest := doc.elementClosestLeft(&doc.Boxes, nil, startElem)
	if (closest == nil) || (closest.X <= endX) {
		return endX, endY, nil, endX, endY
	} else {
		var nextDirection ConnDirection
		var x, y int
		x = closest.X + (doc.MinBoxMargin / 2)
		if closest.CenterY > startY {
			nextDirection = ConnDirectionUp
			y = closest.Y - (doc.MinBoxMargin / 2)
		} else {
			nextDirection = ConnDirectionDown
			y = closest.Y + closest.Height + (doc.MinBoxMargin / 2)
		}
		return x, endY, &nextDirection, x, y
	}
}

func (doc *BoxesDocument) checkForCollisionRight(startX, startY, endX, endY int, startElem *LayoutElement) (int, int, *ConnDirection, int, int) {
	closest := doc.elementClosestRight(&doc.Boxes, nil, startElem)
	if (closest == nil) || (closest.X >= endX) {
		return endX, endY, nil, endX, endY
	} else {
		var nextDirection ConnDirection
		var x, y int
		x = closest.X - (doc.MinBoxMargin / 2)
		if closest.CenterY > startY {
			nextDirection = ConnDirectionUp
			y = closest.Y - (doc.MinBoxMargin / 2)
		} else {
			nextDirection = ConnDirectionDown
			y = closest.Y + closest.Height + (doc.MinBoxMargin / 2)
		}
		return x, endY, &nextDirection, x, y
	}
}

func (doc *BoxesDocument) checkForCollisionUp(startX, startY, endX, endY int, startElem *LayoutElement) (int, int, *ConnDirection, int, int) {
	closest := doc.elementClosestTop(&doc.Boxes, nil, startElem)
	if (closest == nil) || (closest.Y <= endY) {
		return endX, endY, nil, endX, endY
	} else {
		var nextDirection ConnDirection
		var x, y int
		y = closest.Y + (doc.MinBoxMargin / 2)
		if closest.CenterX > startX {
			nextDirection = ConnDirectionLeft
			x = closest.X + (doc.MinBoxMargin / 2)
		} else {
			nextDirection = ConnDirectionRight
			x = closest.X + closest.Width + (doc.MinBoxMargin / 2)
		}
		return endX, y, &nextDirection, x, y
	}
}

func (doc *BoxesDocument) checkForCollisionDown(startX, startY, endX, endY int, startElem *LayoutElement) (int, int, *ConnDirection, int, int) {
	closest := doc.elementClosestBottom(&doc.Boxes, nil, startElem)
	if (closest == nil) || (closest.Y >= endY) {
		return endX, endY, nil, endX, endY
	} else {
		var nextDirection ConnDirection
		var x, y int
		y = closest.Y - (doc.MinBoxMargin / 2)
		if closest.CenterX > startX {
			nextDirection = ConnDirectionLeft
			x = closest.X + (doc.MinBoxMargin / 2)
		} else {
			nextDirection = ConnDirectionRight
			x = closest.X + closest.Width + (doc.MinBoxMargin / 2)
		}
		return endX, y, &nextDirection, x, y
	}
}

func (doc *BoxesDocument) isParentInContainer(container *LayoutElemContainer, possibleParent, elemToCheckFor *LayoutElement) bool {
	if container != nil {
		for _, subElem := range container.Elems {
			if subElem.Id == elemToCheckFor.Id {
				return true
			}
			if doc.isParent(&subElem, elemToCheckFor) {
				return true
			}
		}
	}
	return false
}

func (doc *BoxesDocument) isParent(possibleParent, elemToCheckFor *LayoutElement) bool {
	if doc.isParentInContainer(possibleParent.Vertical, possibleParent, elemToCheckFor) {
		return true
	}
	return doc.isParentInContainer(possibleParent.Horizontal, possibleParent, elemToCheckFor)
}

func (doc *BoxesDocument) elementClosestLeft(elem *LayoutElement, soFarClosest, originalStartElem *LayoutElement) *LayoutElement {
	if soFarClosest == nil {
		soFarClosest = elem
	} else {
		// TODO only if elem is not the originalStartElem and
		// elem is not parent that contains originalStartElem
		if (elem.Id != originalStartElem.Id) && (!doc.isParent(elem, soFarClosest)) &&
			(originalStartElem.CenterY > elem.Y) && (originalStartElem.CenterY < elem.Y+elem.Height) {
			if (elem.X > soFarClosest.X) && (elem.X < originalStartElem.X) {
				soFarClosest = elem
			}
		}
	}
	if elem.Vertical != nil {
		for _, subElem := range elem.Vertical.Elems {
			soFarClosest = doc.elementClosestLeft(&subElem, soFarClosest, originalStartElem)
		}
	}
	if elem.Horizontal != nil {
		for _, subElem := range elem.Horizontal.Elems {
			soFarClosest = doc.elementClosestLeft(&subElem, soFarClosest, originalStartElem)
		}
	}
	return soFarClosest
}

func (doc *BoxesDocument) elementClosestRight(elem *LayoutElement, soFarClosest, originalStartElem *LayoutElement) *LayoutElement {
	if soFarClosest == nil {
		soFarClosest = elem
	} else {
		// TODO only if elem is not the originalStartElem and
		// elem is not parent that contains originalStartElem
		if (elem.Id != originalStartElem.Id) && (!doc.isParent(elem, soFarClosest)) &&
			(originalStartElem.CenterY > elem.Y) && (originalStartElem.CenterY < elem.Y+elem.Height) {
			if (elem.X < soFarClosest.X) && (elem.X > originalStartElem.X) {
				soFarClosest = elem
			}
		}
	}
	if elem.Vertical != nil {
		for _, subElem := range elem.Vertical.Elems {
			soFarClosest = doc.elementClosestRight(&subElem, soFarClosest, originalStartElem)
		}
	}
	if elem.Horizontal != nil {
		for _, subElem := range elem.Horizontal.Elems {
			soFarClosest = doc.elementClosestRight(&subElem, soFarClosest, originalStartElem)
		}
	}
	return soFarClosest
}

func (doc *BoxesDocument) elementClosestTop(elem *LayoutElement, soFarClosest, originalStartElem *LayoutElement) *LayoutElement {
	if soFarClosest == nil {
		soFarClosest = elem
	} else {
		// TODO only if elem is not the originalStartElem and
		// elem is not parent that contains originalStartElem
		if (elem.Id != originalStartElem.Id) && (!doc.isParent(elem, soFarClosest)) &&
			(originalStartElem.CenterX > elem.X) && (originalStartElem.CenterX < elem.X+elem.Width) {
			if (elem.Y > soFarClosest.Y) && (elem.Y < originalStartElem.Y) {
				soFarClosest = elem
			}
		}
	}
	if elem.Vertical != nil {
		for _, subElem := range elem.Vertical.Elems {
			soFarClosest = doc.elementClosestTop(&subElem, soFarClosest, originalStartElem)
		}
	}
	if elem.Horizontal != nil {
		for _, subElem := range elem.Horizontal.Elems {
			soFarClosest = doc.elementClosestTop(&subElem, soFarClosest, originalStartElem)
		}
	}
	return soFarClosest
}

func (doc *BoxesDocument) elementClosestBottom(elem *LayoutElement, soFarClosest, originalStartElem *LayoutElement) *LayoutElement {
	if soFarClosest == nil {
		soFarClosest = elem
	} else {
		// TODO only if elem is not the originalStartElem and
		// elem is not parent that contains originalStartElem
		if (elem.Id != originalStartElem.Id) && (!doc.isParent(elem, soFarClosest)) &&
			(originalStartElem.CenterX > elem.X) && (originalStartElem.CenterX < elem.X+elem.Width) {
			if (elem.Y < soFarClosest.Y) && (elem.Y > originalStartElem.Y) {
				soFarClosest = elem
			}
		}
	}
	if elem.Vertical != nil {
		for _, subElem := range elem.Vertical.Elems {
			soFarClosest = doc.elementClosestBottom(&subElem, soFarClosest, originalStartElem)
		}
	}
	if elem.Horizontal != nil {
		for _, subElem := range elem.Horizontal.Elems {
			soFarClosest = doc.elementClosestBottom(&subElem, soFarClosest, originalStartElem)
		}
	}
	return soFarClosest
}

func (doc *BoxesDocument) findNextConnectionParts(alreadyCollectedParts []ConnectionLine, startX, startY, endX, endY int, nextDirection *ConnDirection) []ConnectionLine {
	return alreadyCollectedParts // TODO
}

func (doc *BoxesDocument) handleDirection(direction ConnDirection, startX, startY, endX, endY int, startElem *LayoutElement) (int, int, *ConnDirection, int, int) {
	var nextX, nextY, nextX2, nextY2 int
	var nextDirection *ConnDirection
	switch direction {
	case ConnDirectionLeft:
		if startY == endY {
			// straight line
			nextX = endX
		} else {
			nextX = (startX - endX) / 2
		}
		nextX, nextY, nextDirection, nextX2, nextY2 = doc.checkForCollisionLeft(startX, startY, nextX, startY, startElem)
	case ConnDirectionRight:
		if startY == endY {
			// straight line
			nextX = endX
		} else {
			nextX = (endX - startX) / 2
		}
		nextX, nextY, nextDirection, nextX2, nextY2 = doc.checkForCollisionRight(startX, startY, nextX, startY, startElem)
	case ConnDirectionUp:
		if startX == endX {
			// straight line
			nextY = endY
		} else {
			nextY = (startY - endY) / 2
		}
		nextX, nextY, nextDirection, nextX2, nextY2 = doc.checkForCollisionUp(startX, startY, startX, nextY, startElem)
	case ConnDirectionDown:
		if startX == endX {
			// straight line
			nextY = endY
		} else {
			nextY = (endY - startY) / 2
		}
		nextX, nextY, nextDirection, nextX2, nextY2 = doc.checkForCollisionDown(startX, startY, startX, nextY, startElem)
	}
	return nextX, nextY, nextDirection, nextX2, nextY2
}

// connects two points with a straight line
func (doc *BoxesDocument) initialLineConnection(startX, startY, endX, endY int) []ConnectionLine {
	connection := make([]ConnectionLine, 0)
	connection = append(connection, ConnectionLine{
		StartX: startX,
		StartY: startY,
		EndX:   endX,
		EndY:   endY,
	})
	return connection
}

// connects two points with a U-shaped line
func (doc *BoxesDocument) initialUConnection(startX, startY, endX, endY, startLen int) []ConnectionLine {
	connection := make([]ConnectionLine, 0)
	y2 := startY + startLen
	connection = append(connection, ConnectionLine{
		StartX: startX,
		StartY: startY,
		EndX:   startX,
		EndY:   y2,
	})
	connection = append(connection, ConnectionLine{
		StartX: startX,
		StartY: y2,
		EndX:   endX,
		EndY:   y2,
	})
	connection = append(connection, ConnectionLine{
		StartX: endX,
		StartY: y2,
		EndX:   endX,
		EndY:   endY,
	})
	return connection
}

// connects two points with a L-shaped line, starting in a horizontal direction
func (doc *BoxesDocument) initialHLConnection(startX, startY, endX, endY int) []ConnectionLine {
	connection := make([]ConnectionLine, 0)
	connection = append(connection, ConnectionLine{
		StartX: startX,
		StartY: startY,
		EndX:   endX,
		EndY:   startY,
	})
	connection = append(connection, ConnectionLine{
		StartX: endX,
		StartY: startY,
		EndX:   endX,
		EndY:   endY,
	})
	return connection
}

// connects two points with a L-shaped line, starting in a vertical direction
func (doc *BoxesDocument) initialVLConnection(startX, startY, endX, endY int) []ConnectionLine {
	connection := make([]ConnectionLine, 0)
	connection = append(connection, ConnectionLine{
		StartX: startX,
		StartY: startY,
		EndX:   startX,
		EndY:   endY,
	})
	connection = append(connection, ConnectionLine{
		StartX: startX,
		StartY: endY,
		EndX:   endX,
		EndY:   endY,
	})
	return connection
}

// connects two points with a vertical "S"-shaped line
func (doc *BoxesDocument) initialVSConnection(startX, startY, endX, endY int) []ConnectionLine {
	connection := make([]ConnectionLine, 0)
	y2 := startY + (endY-startY)/2
	connection = append(connection, ConnectionLine{
		StartX: startX,
		StartY: startY,
		EndX:   startX,
		EndY:   y2,
	})
	connection = append(connection, ConnectionLine{
		StartX: startX,
		StartY: y2,
		EndX:   endX,
		EndY:   y2,
	})
	connection = append(connection, ConnectionLine{
		StartX: endX,
		StartY: y2,
		EndX:   endX,
		EndY:   endY,
	})
	return connection
}

// connects two points with a horizontal "S"-shaped line
func (doc *BoxesDocument) initialHSConnection(startX, startY, endX, endY int) []ConnectionLine {
	connection := make([]ConnectionLine, 0)
	x2 := startX + (endX-startX)/2
	connection = append(connection, ConnectionLine{
		StartX: startX,
		StartY: startY,
		EndX:   x2,
		EndY:   startY,
	})
	connection = append(connection, ConnectionLine{
		StartX: x2,
		StartY: startY,
		EndX:   x2,
		EndY:   endY,
	})
	connection = append(connection, ConnectionLine{
		StartX: x2,
		StartY: endY,
		EndX:   endX,
		EndY:   endY,
	})
	return connection
}

func (doc *BoxesDocument) solveCollisions(connection []ConnectionLine) []ConnectionLine {
	// TODO
	return connection
}

func (doc *BoxesDocument) lineConnection(startX, startY, endX, endY int) []ConnectionLine {
	connection := doc.initialLineConnection(startX, startY, endX, endY)
	return doc.solveCollisions(connection)
}

func (doc *BoxesDocument) uConnection(startX, startY, endX, endY, startLen int) []ConnectionLine {
	connection := doc.initialUConnection(startX, startY, endX, endY, startLen)
	return doc.solveCollisions(connection)
}

func (doc *BoxesDocument) vlConnection(startX, startY, endX, endY int) []ConnectionLine {
	connection := doc.initialVLConnection(startX, startY, endX, endY)
	return doc.solveCollisions(connection)
}

func (doc *BoxesDocument) hlConnection(startX, startY, endX, endY int) []ConnectionLine {
	connection := doc.initialHLConnection(startX, startY, endX, endY)
	return doc.solveCollisions(connection)
}

func (doc *BoxesDocument) hsConnection(startX, startY, endX, endY int) []ConnectionLine {
	connection := doc.initialHSConnection(startX, startY, endX, endY)
	return doc.solveCollisions(connection)
}

func (doc *BoxesDocument) vsConnection(startX, startY, endX, endY int) []ConnectionLine {
	connection := doc.initialVSConnection(startX, startY, endX, endY)
	return doc.solveCollisions(connection)
}

func (doc *BoxesDocument) getConnectionParts(startElem, destElem *LayoutElement) []ConnectionLine {
	connectionVariants := make([][]ConnectionLine, 0)
	if startElem.AreOnTheSameVerticalLevel(destElem) {
		// 1. connect from right side to the left side
		connectionVariants = append(connectionVariants, doc.lineConnection(startElem.X+startElem.Width, startElem.CenterY, destElem.X, destElem.CenterY))
		// 2. connect from top to top
		connectionVariants = append(connectionVariants, doc.uConnection(startElem.CenterX, startElem.Y, destElem.CenterX, destElem.Y, doc.MinBoxMargin/-2))
		// 3. connect from bottom to bottom
		connectionVariants = append(connectionVariants, doc.uConnection(startElem.CenterX, startElem.Y+startElem.Height, destElem.CenterX, destElem.Y+destElem.Height, doc.MinBoxMargin/2))
	} else if startElem.CenterY < destElem.CenterY {
		// 1. connect from bottom to top side
		connectionVariants = append(connectionVariants, doc.vsConnection(startElem.CenterX, startElem.Y+startElem.Height, destElem.CenterX, destElem.Y))
		// 2. connect from left to top side
		connectionVariants = append(connectionVariants, doc.hlConnection(startElem.CenterX, startElem.CenterY, destElem.CenterX, destElem.Y))
		// 3. connect from right to top side
		connectionVariants = append(connectionVariants, doc.hlConnection(startElem.X+startElem.Width, startElem.CenterY, destElem.CenterX, destElem.Y))
		if startElem.CenterX < destElem.CenterX {
			// 4. connect from bottom to left side
			connectionVariants = append(connectionVariants, doc.vlConnection(startElem.CenterX, startElem.Y+startElem.Height, destElem.X, destElem.CenterY))
			// 5. connect from right to left side
			connectionVariants = append(connectionVariants, doc.hsConnection(startElem.X+startElem.Width, startElem.CenterY, destElem.X, destElem.CenterY))
		} else {
			// 6. connect from bottom to right side
			connectionVariants = append(connectionVariants, doc.vlConnection(startElem.CenterX, startElem.Y+startElem.Height, destElem.X+destElem.Width, destElem.CenterY))
			// 7. connect from left to right side
			connectionVariants = append(connectionVariants, doc.hsConnection(startElem.X, startElem.CenterY, destElem.X+destElem.Width, destElem.CenterY))
		}
	} else {
		// 1. connect from top to bottom side
		connectionVariants = append(connectionVariants, doc.vsConnection(startElem.CenterX, startElem.Y, destElem.CenterX, destElem.Y+destElem.Height))
		// 2. connect from left to bottom side
		connectionVariants = append(connectionVariants, doc.hlConnection(startElem.X, startElem.CenterY, destElem.CenterX, destElem.Y+destElem.Height))
		// 3. connect from right to bottom side
		connectionVariants = append(connectionVariants, doc.hlConnection(startElem.X+startElem.Width, startElem.CenterY, destElem.CenterX, destElem.Y+destElem.Height))
		if startElem.CenterX < destElem.CenterX {
			// 4. connect from top to left side
			connectionVariants = append(connectionVariants, doc.vlConnection(startElem.CenterX, startElem.Y, destElem.X, destElem.CenterY))
			// 5. connect from right to left side
			connectionVariants = append(connectionVariants, doc.hsConnection(startElem.X+startElem.Width, startElem.CenterY, destElem.X, destElem.CenterY))
		} else {
			// 6. connect from top to right side
			connectionVariants = append(connectionVariants, doc.vlConnection(startElem.CenterX, startElem.Y, destElem.X+destElem.Width, destElem.CenterY))
			// 7. connect from left to right side
			connectionVariants = append(connectionVariants, doc.hsConnection(startElem.X, startElem.CenterY, destElem.X+destElem.Width, destElem.CenterY))
		}
	}
	var connection []ConnectionLine
	for _, conn := range connectionVariants {
		if connection == nil || len(conn) < len(connection) {
			connection = conn
		}
	}
	return connection
}

func (doc *BoxesDocument) connectTwoElems(start, destElem *LayoutElement, lec *LayoutElemConnection) ConnectionElem {
	var ret ConnectionElem
	ret.DestArrow = &lec.DestArrow
	ret.SourceArrow = &lec.SourceArrow
	ret.Parts = doc.getConnectionParts(start, destElem)
	return ret
}

func (doc *BoxesDocument) doConnect(elem *LayoutElement) {
	for _, conn := range elem.Connections {
		destElem, found := doc.findLayoutElementById(conn.DestId, elem)
		if !found {
			fmt.Println("Couldn't find destId: ", conn.DestId)
			continue
		}
		if destElem == elem {
			fmt.Println("Connection to self are not allowed and will be ignored: ", conn.DestId)
			continue
		}
		connectionElem := doc.connectTwoElems(elem, destElem, &conn)
		doc.Connections = append(doc.Connections, connectionElem)
	}
}

func (doc *BoxesDocument) connectLayoutElem(le *LayoutElement) {
	doc.doConnect(le)
	doc.connectContainer(le.Vertical)
	doc.connectContainer(le.Horizontal)
}

func (doc *BoxesDocument) connectContainer(cont *LayoutElemContainer) {
	if cont != nil {
		for _, elem := range cont.Elems {
			doc.connectLayoutElem(&elem)
		}
	}
}

func (doc *BoxesDocument) ConnectBoxes() {
	doc.connectLayoutElem(&doc.Boxes)
}
