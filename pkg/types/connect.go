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

func (doc *BoxesDocument) checkForCollisionLeft(startX, startY, endX, endY int, startElem *LayoutElement) (int, int, *ConnDirection) {
	closest := doc.elementClosestLeft(&doc.Boxes, nil, startElem)
	if (closest == nil) || (closest.X <= endX) {
		return endX, endY, nil
	} else {
		var nextDirection ConnDirection
		if (closest.Y + (closest.Height / 2)) > startY {
			nextDirection = ConnDirectionUp
		} else {
			nextDirection = ConnDirectionDown
		}
		return closest.X + (doc.MinBoxMargin / 2), endY, &nextDirection
	}
}

func (doc *BoxesDocument) checkForCollisionRight(startX, startY, endX, endY int, startElem *LayoutElement) (int, int, *ConnDirection) {
	closest := doc.elementClosestRight(&doc.Boxes, nil, startElem)
	if (closest == nil) || (closest.X >= endX) {
		return endX, endY, nil
	} else {
		var nextDirection ConnDirection
		if (closest.Y + (closest.Height / 2)) > startY {
			nextDirection = ConnDirectionUp
		} else {
			nextDirection = ConnDirectionDown
		}
		return closest.X - (doc.MinBoxMargin / 2), endY, &nextDirection
	}
}

func (doc *BoxesDocument) checkForCollisionUp(startX, startY, endX, endY int, startElem *LayoutElement) (int, int, *ConnDirection) {
	closest := doc.elementClosestTop(&doc.Boxes, nil, startElem)
	if (closest == nil) || (closest.Y <= endY) {
		return endX, endY, nil
	} else {
		var nextDirection ConnDirection
		if (closest.X + (closest.Width / 2)) > startX {
			nextDirection = ConnDirectionLeft
		} else {
			nextDirection = ConnDirectionRight
		}
		return endX, closest.Y + (doc.MinBoxMargin / 2), &nextDirection
	}
}

func (doc *BoxesDocument) checkForCollisionDown(startX, startY, endX, endY int, startElem *LayoutElement) (int, int, *ConnDirection) {
	closest := doc.elementClosestBottom(&doc.Boxes, nil, startElem)
	if (closest == nil) || (closest.Y >= endY) {
		return endX, endY, nil
	} else {
		var nextDirection ConnDirection
		if (closest.X + (closest.Width / 2)) > startX {
			nextDirection = ConnDirectionLeft
		} else {
			nextDirection = ConnDirectionRight
		}
		return endX, closest.Y - (doc.MinBoxMargin / 2), &nextDirection
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

func (doc *BoxesDocument) getConnectionParts(alreadyCollectedParts []ConnectionLine, startElem, destElem *LayoutElement) []ConnectionLine {
	startX, startY, startDirection := startElem.ConnectorStart(destElem)
	endX, endY, _ := destElem.ConnectorStart(startElem)
	var nextX, nextY int
	var nextDirection *ConnDirection
	switch startDirection {
	case ConnDirectionLeft:
		if startY == endY {
			// straight line
			nextX = endX
		} else {
			nextX = (startX - endX) / 2
		}
		nextX, nextY, nextDirection = doc.checkForCollisionLeft(startX, startY, nextX, startY, startElem)
	case ConnDirectionRight:
		if startY == endY {
			// straight line
			nextX = endX
		} else {
			nextX = (endX - startX) / 2
		}
		nextX, nextY, nextDirection = doc.checkForCollisionRight(startX, startY, nextX, startY, startElem)
	case ConnDirectionUp:
		if startX == endX {
			// straight line
			nextY = endY
		} else {
			nextY = (startY - endY) / 2
		}
		nextX, nextY, nextDirection = doc.checkForCollisionUp(startX, startY, startX, nextY, startElem)
	case ConnDirectionDown:
		if startX == endX {
			// straight line
			nextY = endY
		} else {
			nextY = (endY - startY) / 2
		}
		nextX, nextY, nextDirection = doc.checkForCollisionDown(startX, startY, startX, nextY, startElem)
	}
	newAlreadyCollectedParts := append(alreadyCollectedParts, ConnectionLine{
		StartX: startX,
		StartY: startY,
		EndX:   endX,
		EndY:   endY,
	})
	if nextX == endX && nextY == endY {
		return newAlreadyCollectedParts
	}
	return doc.findNextConnectionParts(newAlreadyCollectedParts, nextX, nextY, endX, endY, nextDirection)
}

func (doc *BoxesDocument) connectTwoElems(start, destElem *LayoutElement, lec *LayoutElemConnection) ConnectionElem {
	var ret ConnectionElem
	ret.DestArrow = &lec.DestArrow
	ret.SourceArrow = &lec.SourceArrow
	ret.Parts = doc.getConnectionParts(make([]ConnectionLine, 0), start, destElem)
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
