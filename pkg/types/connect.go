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

func (doc *BoxesDocument) getLeftNeighbor(startElem, currentElem, currentNeighbor *LayoutElement) *LayoutElement {
	if currentElem == nil {
		currentNeighbor = &doc.Boxes
	}

	return nil
}

func (doc *BoxesDocument) elementsAreHorizontalNeighbors(startElem, destElem *LayoutElement) bool {
	// leftNeighbor := doc.getLeftNeighbor(startElem, nil)
	// rightNeighbor := doc.getRightNeighbor(startElem, nil)
	// return leftNeighbor == destElem || rightNeighbor == destElem
	return false // TODO

}

func (doc *BoxesDocument) getConnectionParts(alreadyCollectedParts []ConnectionLine, startElem, destElem *LayoutElement) []ConnectionLine {
	if ((startElem.CenterY > destElem.Y) && (startElem.CenterY < destElem.Y+destElem.Height)) ||
		((destElem.CenterY > startElem.Y) && (destElem.CenterY < startElem.Y+startElem.Height)) {
		// the elements are on the same horizontal level
		if doc.elementsAreHorizontalNeighbors(startElem, destElem) {
			// TODO connect horizontally
			return alreadyCollectedParts
		}
	}
	if startElem.CenterY < destElem.CenterY {
		// connection from top to bottom
		if startElem.CenterX <= destElem.CenterX {
			// connection from left/top to right/down ...
			// connector starts at the bottom of the start element
			// TODO

		} else {
			// connection from right/top to left/down
			// TODO
		}
	} else {
		// connection from bottom to top
		if startElem.CenterX <= destElem.CenterX {
			// connection from left/down to right/top
			// TODO
		} else {
			// connection from right/down to left/top
			// TODO
		}
	}
	return alreadyCollectedParts
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
