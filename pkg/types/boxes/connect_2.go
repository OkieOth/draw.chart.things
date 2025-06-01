package boxes

import (
	"fmt"
)

func (doc *BoxesDocument) findLayoutElementById(id string, startElem *LayoutElement) (*LayoutElement, bool) {
	if startElem.Id == id {
		return startElem, true
	}
	if startElem.Vertical != nil {
		for _, elem := range startElem.Vertical.Elems {
			if elem.Id == id {
				return &elem, true
			}
			if e, found := doc.findLayoutElementById(id, &elem); found {
				return e, true
			}
		}
	}
	if startElem.Horizontal != nil {
		for _, elem := range startElem.Horizontal.Elems {
			if elem.Id == id {
				return &elem, true
			}
			if e, found := doc.findLayoutElementById(id, &elem); found {
				return e, true
			}
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
func (doc *BoxesDocument) initialHUConnection(startX, startY, endX, endY, startLen int) []ConnectionLine {
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

func (doc *BoxesDocument) initialVUConnection(startX, startY, endX, endY, x2 int) []ConnectionLine {
	connection := make([]ConnectionLine, 0)

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

// connects two points with a L-shaped line, starting in a horizontal direction
func (doc *BoxesDocument) initialHLConnection(startX, startY, endX, endY int) ([]ConnectionLine, error) {
	connection := make([]ConnectionLine, 0)

	upperY, lowerY := minMax(startY, endY)

	if (lowerY - upperY) < doc.MinBoxMargin/2 {
		return nil, fmt.Errorf("HL connection would be too small")
	}

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
	return connection, nil
}

func minMax(a, b int) (int, int) {
	if a < b {
		return a, b
	}
	return b, a
}

// connects two points with a L-shaped line, starting in a vertical direction
func (doc *BoxesDocument) initialVLConnection(startX, startY, endX, endY int) ([]ConnectionLine, error) {
	connection := make([]ConnectionLine, 0)

	upperY, lowerY := minMax(startY, endY)

	if (lowerY - upperY) < doc.MinBoxMargin/2 {
		return nil, fmt.Errorf("HS connection would be too small")
	}

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
	return connection, nil
}

// connects two points with a vertical "S"-shaped line
func (doc *BoxesDocument) initialVSConnection(startX, startY, endX, endY int) ([]ConnectionLine, error) {
	connection := make([]ConnectionLine, 0)
	y2 := startY + (endY-startY)/2
	if y2 < doc.MinBoxMargin/2 {
		return nil, fmt.Errorf("VS connection would be too small")
	}
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
	return connection, nil
}

// connects two points with a horizontal "S"-shaped line
func (doc *BoxesDocument) initialHSConnection(startX, startY, endX, endY int) ([]ConnectionLine, error) {
	connection := make([]ConnectionLine, 0)
	x2 := startX + (endX-startX)/2
	if x2 < doc.MinBoxMargin/2 {
		return nil, fmt.Errorf("HS connection would be too small")
	}
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
	return connection, nil
}

func needsInverseOrder(conn ConnectionLine) bool {
	return conn.StartX > conn.EndX || conn.StartY > conn.EndY
}

func (doc *BoxesDocument) checkAndSolveCollisionImpl(layoutElement *LayoutElement, connectionLines []ConnectionLine, index int, startElem, destElem *LayoutElement) ([]ConnectionLine, error) {
	ret := make([]ConnectionLine, 0)
	isParentOfStartElem := doc.isParent(layoutElement, startElem)
	isParentOfDestElem := doc.isParent(layoutElement, destElem)
	if (!isParentOfStartElem) && (!isParentOfDestElem) {
		// current element is not a parent of start or dest element
		beforeFixing := len(connectionLines)
		ret := layoutElement.FixCollisionInCase(connectionLines, index, doc.MinBoxMargin/2)
		if len(ret) > beforeFixing {
			return ret, nil
		}
	} else {
		if layoutElement.Caption != "" || layoutElement.Text1 != "" || layoutElement.Text2 != "" {
			// current element is a parent of start or dest element
			// and has a caption or text ... if one of the lines goes vertical through the parent, then is has a collision
			for _, l := range connectionLines {
				if l.StartX == l.EndX {
					// vertical line ...
					leftX, rightX := connToLeftRightX(l)
					topY, bottomY := connToUpperLowerY(l)
					if ((leftX >= layoutElement.X) && (rightX <= layoutElement.X+layoutElement.Width)) &&
						(((topY <= layoutElement.Y) && (bottomY > layoutElement.Y)) || ((topY > layoutElement.Y) && (topY < layoutElement.Y+layoutElement.Height))) {
						return nil, fmt.Errorf("collision in caption or text: %s", layoutElement.Caption)
					}
				}
			}
		}
	}
	connToCheck := index
	connectionLinesToCheck := connectionLines
	inverseOrder := needsInverseOrder(connectionLines[index])
	if layoutElement.Vertical != nil {
		lv := len(layoutElement.Vertical.Elems)
		for i := 0; i < lv; i++ {
			index := i
			if inverseOrder {
				index = lv - 1 - i
			}
			subElem := layoutElement.Vertical.Elems[index]
			beforeFixing := len(connectionLinesToCheck)
			fixedConnection, err := doc.checkAndSolveCollisionImpl(&subElem, connectionLinesToCheck, connToCheck, startElem, destElem)
			if err != nil {
				return nil, err
			}
			if len(fixedConnection) > beforeFixing {
				// connection was changed
				l := len(fixedConnection)
				connToCheck = 0
				connectionLinesToCheck = fixedConnection[l-1 : l]
				ret = append(ret, fixedConnection[0:l-1]...)
				break
			}
		}
	}
	if layoutElement.Horizontal != nil {
		lh := len(layoutElement.Horizontal.Elems)
		for i := 0; i < lh; i++ {
			index := i
			if inverseOrder {
				index = lh - 1 - i
			}
			subElem := layoutElement.Horizontal.Elems[index]
			beforeFixing := len(connectionLinesToCheck)
			fixedConnection, err := doc.checkAndSolveCollisionImpl(&subElem, connectionLinesToCheck, connToCheck, startElem, destElem)
			if err != nil {
				return nil, err
			}
			if len(fixedConnection) > beforeFixing {
				// connection was changed
				l := len(fixedConnection)
				connToCheck = 0
				connectionLinesToCheck = fixedConnection[l-1 : l]
				ret = append(ret, fixedConnection[0:l-1]...)
				break
			}
		}
	}
	if len(ret) > 0 {
		ret = append(ret, connectionLinesToCheck...)
	} else {
		ret = append(ret, connectionLines[index])
	}
	return ret, nil
}

func (doc *BoxesDocument) checkAndSolveCollision(connectionLines []ConnectionLine, index int, startElem, destElem *LayoutElement) ([]ConnectionLine, error) {
	return doc.checkAndSolveCollisionImpl(&doc.Boxes, connectionLines, index, startElem, destElem)
}

func isStepForAndBack(previous, current ConnectionLine) bool {
	xdiff1 := current.StartX - previous.StartX
	if xdiff1 < 0 {
		xdiff1 = -xdiff1
	}
	ydiff1 := current.StartY - previous.StartY
	if ydiff1 < 0 {
		ydiff1 = -ydiff1
	}
	xdiff2 := current.EndX - previous.EndX
	if xdiff2 < 0 {
		xdiff2 = -xdiff2
	}
	ydiff2 := current.EndY - previous.EndY
	if ydiff2 < 0 {
		ydiff2 = -ydiff2
	}
	xdiff := xdiff1 + xdiff2
	ydiff := ydiff1 + ydiff2

	// two steps draw a horizontal line
	if xdiff <= 2 || ydiff <= 2 {
		return true
	}
	return false
}

func (doc *BoxesDocument) solveCollisions(connectionLines []ConnectionLine, startElem, destElem *LayoutElement) ([]ConnectionLine, error) {
	//return connectionLines, nil // TODO, DEBUG
	ret := make([]ConnectionLine, 0)
	for i := range len(connectionLines) {
		if connectionLines[i].MovedOut {
			ret = append(ret, connectionLines[i])
		} else {
			r, err := doc.checkAndSolveCollision(connectionLines, i, startElem, destElem)
			if err != nil {
				return nil, err
			}
			ret = append(ret, r...)
		}
	}
	// check if the connection is valid
	for i := 0; i < len(ret); i++ {
		if i > 0 {
			current := ret[i]
			previous := ret[i-1]
			if (current.StartX != previous.EndX) || (current.StartY != previous.EndY) {
				return nil, fmt.Errorf("invalid connection line: %v", current)
			}
			if isStepForAndBack(previous, current) {
				return nil, fmt.Errorf("invalid connection line: %v", current)
			}
		}
	}
	return ret, nil
}

func (doc *BoxesDocument) lineConnection(startX, startY, endX, endY int, startElem, destElem *LayoutElement) ([]ConnectionLine, error) {
	connection := doc.initialLineConnection(startX, startY, endX, endY)
	return doc.solveCollisions(connection, startElem, destElem)
}

func (doc *BoxesDocument) huConnection(startX, startY, endX, endY, startLen int, startElem, destElem *LayoutElement) ([]ConnectionLine, error) {
	connection := doc.initialHUConnection(startX, startY, endX, endY, startLen)
	return doc.solveCollisions(connection, startElem, destElem)
}

func (doc *BoxesDocument) vuConnection(startX, startY, endX, endY, startLen int, startElem, destElem *LayoutElement) ([]ConnectionLine, error) {
	connection := doc.initialVUConnection(startX, startY, endX, endY, startLen)
	return doc.solveCollisions(connection, startElem, destElem)
}

func (doc *BoxesDocument) vlConnection(startX, startY, endX, endY int, startElem, destElem *LayoutElement) ([]ConnectionLine, error) {
	connection, err := doc.initialVLConnection(startX, startY, endX, endY)
	if err != nil {
		return nil, err
	}
	return doc.solveCollisions(connection, startElem, destElem)
}

func (doc *BoxesDocument) hlConnection(startX, startY, endX, endY int, startElem, destElem *LayoutElement) ([]ConnectionLine, error) {
	connection, err := doc.initialHLConnection(startX, startY, endX, endY)
	if err != nil {
		return nil, err
	}
	return doc.solveCollisions(connection, startElem, destElem)
}

func (doc *BoxesDocument) hsConnection(startX, startY, endX, endY int, startElem, destElem *LayoutElement) ([]ConnectionLine, error) {
	connection, err := doc.initialHSConnection(startX, startY, endX, endY)
	if err != nil {
		return nil, err
	}
	return doc.solveCollisions(connection, startElem, destElem)
}

func (doc *BoxesDocument) vsConnection(startX, startY, endX, endY int, startElem, destElem *LayoutElement) ([]ConnectionLine, error) {
	connection, err := doc.initialVSConnection(startX, startY, endX, endY)
	if err != nil {
		return nil, err
	}
	return doc.solveCollisions(connection, startElem, destElem)
}
