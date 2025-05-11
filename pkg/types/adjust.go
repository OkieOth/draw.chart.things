package types

func isInDistance(v1, v2, distance int) (bool, bool, int) {
	if v1 > v2 {
		dist := v1 - v2
		return (dist <= distance), true, dist
	}
	dist := v2 - v1
	return (dist <= distance), false, dist
}

func (doc *BoxesDocument) moveAllObjectsLowerDown(le *LayoutElement, y, dist int, allowEqual bool) {
	if le != &doc.Boxes {
		if ((le.Y) > y) || (allowEqual && (le.Y) == y) {
			// the object is lower than y
			le.Y += dist
			le.CenterY += dist
			if le.Horizontal != nil {
				le.Horizontal.Y += dist
			}
			if le.Vertical != nil {
				le.Vertical.Y += dist
			}
			doc.Height += dist
		}
	}
	if le.Horizontal != nil {
		for i := 0; i < len(le.Horizontal.Elems); i++ {
			doc.moveAllObjectsLowerDown(&le.Horizontal.Elems[i], y, dist, allowEqual)
		}
	}
	if le.Vertical != nil {
		for i := 0; i < len(le.Vertical.Elems); i++ {
			doc.moveAllObjectsLowerDown(&le.Vertical.Elems[i], y, dist, allowEqual)
		}
	}
}

func (doc *BoxesDocument) moveAllObjectsMoreRightRight(le *LayoutElement, x, dist int, allowEqual bool) {
	if le != &doc.Boxes {
		if (le.X > x) || (allowEqual && le.X == x) {
			le.X += dist
			le.CenterX += dist
			if le.Horizontal != nil {
				le.Horizontal.X += dist
			}
			if le.Vertical != nil {
				le.Vertical.X += dist
			}
			doc.Width += dist
		}
	}
	if le.Horizontal != nil {
		for i := 0; i < len(le.Horizontal.Elems); i++ {
			doc.moveAllObjectsMoreRightRight(&le.Horizontal.Elems[i], x, dist, allowEqual)
		}
	}
	if le.Vertical != nil {
		for i := 0; i < len(le.Vertical.Elems); i++ {
			doc.moveAllObjectsMoreRightRight(&le.Vertical.Elems[i], x, dist, allowEqual)
		}
	}
}

func (doc *BoxesDocument) isConnectedToUpperOrLowerBoxBorder(le *LayoutElement, line *ConnectionLine) bool {
	// check upper border
	if le.Id == "r6_2" {
		// Debug
		le.Id = le.Id
	}
	if (le.Y == line.StartY || le.Y == line.EndY) && (line.StartX > le.X && line.StartX < le.X+le.Width) {
		return true
	}
	if ((le.Y+le.Height) == line.StartY || (le.Y+le.Height) == line.EndY) && (line.StartX > le.X && line.StartX < le.X+le.Width) {
		return true
	}
	if le.Vertical != nil {
		for i := 0; i < len(le.Vertical.Elems); i++ {
			if doc.isConnectedToUpperOrLowerBoxBorder(&le.Vertical.Elems[i], line) {
				return true
			}
		}
	}
	if le.Horizontal != nil {
		for i := 0; i < len(le.Horizontal.Elems); i++ {
			if doc.isConnectedToUpperOrLowerBoxBorder(&le.Horizontal.Elems[i], line) {
				return true
			}
		}
	}
	return false
}

func (doc *BoxesDocument) partIsntConnectedToUpperOrLowerBoxBorder(line *ConnectionLine, linePos, connPartCount int) bool {
	if linePos != 0 && linePos != connPartCount-1 {
		// the line is not the first or last part of the connection
		return true
	}
	if line.StartX != line.EndX {
		// the line is not vertical ... should be an error
		return true
	}
	return !doc.isConnectedToUpperOrLowerBoxBorder(&doc.Boxes, line)
}

func (doc *BoxesDocument) moveAllConnectorsLowerDown(y, dist int) {
	for i := 0; i < len(doc.Connections); i++ {
		lp := len(doc.Connections[i].Parts)
		for j := 0; j < lp; j++ {
			part := doc.Connections[i].Parts[j]
			if part.StartY != part.EndY {
				// vertical line
				if part.StartY < part.EndY {
					// line moves down
					if part.EndY > y {
						// ... need to check that the line isn't a connector start that belongs to a unmoved box
						if (part.StartY > y) && doc.partIsntConnectedToUpperOrLowerBoxBorder(&part, j, lp) {
							// line starts after the y position
							doc.Connections[i].Parts[j].StartY += dist
						}
						doc.Connections[i].Parts[j].EndY += dist
					}
				} else {
					// line moves up
					if part.StartY > y {
						// ... need to check that the line isn't a connector start that belongs to a unmoved box
						if (part.EndY > y) && doc.partIsntConnectedToUpperOrLowerBoxBorder(&part, j, lp) {
							// line starts before or at the y position
							doc.Connections[i].Parts[j].EndY += dist
						}
						doc.Connections[i].Parts[j].StartY += dist
					}
				}
			} else {
				// horizontal line
				if part.StartY > y {
					doc.Connections[i].Parts[j].StartY += dist
					doc.Connections[i].Parts[j].EndY += dist
				}
			}
		}
	}
}
func (doc *BoxesDocument) moveAllConnectorsMoreRightAndEqualRight(x, dist int) {
	for i := 0; i < len(doc.Connections); i++ {
		lp := len(doc.Connections[i].Parts)
		for j := 0; j < lp; j++ {
			part := doc.Connections[i].Parts[j]
			if part.StartY != part.EndY {
				// vertical line
				if part.StartX >= x {
					doc.Connections[i].Parts[j].StartX += dist
					doc.Connections[i].Parts[j].EndX += dist
				}
			} else {
				// horizontal line
				if part.StartX < part.EndX {
					// line moves to the right
					if part.EndX >= x {
						if (part.StartX > x) && (j > 0) && (j < lp-1) {
							// line starts after the x position
							doc.Connections[i].Parts[j].StartX += dist
						}
						doc.Connections[i].Parts[j].EndX += dist
					} else if part.EndX == x {
						// line ends at the x position
						doc.Connections[i].Parts[j].EndX += dist
					}
				} else {
					// line moves to the left
					if part.StartX >= x {
						if (part.EndX > x) && (j > 0) && (j < lp-1) {
							// line starts before or at the y position
							doc.Connections[i].Parts[j].EndX += dist
						}
						doc.Connections[i].Parts[j].StartX += dist
					} else if part.StartX == x {
						// line starts at the x position
						doc.Connections[i].Parts[j].StartX += dist
					}
				}
			}
		}
	}
}
func (doc *BoxesDocument) moveAllConnectorsMoreRightRight(x, dist int) {
	for i := 0; i < len(doc.Connections); i++ {
		lp := len(doc.Connections[i].Parts)
		for j := 0; j < lp; j++ {
			part := doc.Connections[i].Parts[j]
			if part.StartX == part.EndX {
				// vertical line
				if part.StartX > x {
					doc.Connections[i].Parts[j].StartX += dist
					doc.Connections[i].Parts[j].EndX += dist
				}
			} else {
				// horizontal line
				if part.StartX < part.EndX {
					// line moves to the right
					if part.EndX > x {
						if (part.StartX > x) && (j > 0) && (j < lp-1) {
							// line starts after the x position
							doc.Connections[i].Parts[j].StartX += dist
						}
						doc.Connections[i].Parts[j].EndX += dist
					}
				} else {
					// line moves to the left
					if part.StartX > x {
						if (part.EndX > x) && (j > 0) && (j < lp-1) {
							// line starts before or at the y position
							doc.Connections[i].Parts[j].EndX += dist
						}
						doc.Connections[i].Parts[j].StartX += dist
					}
				}
			}
		}
	}
}

func (doc *BoxesDocument) checkLayoutElemForVerticalCollision(le *LayoutElement, connIndex, partIndex, y int) {
	if le != &doc.Boxes && (le.Id != "" || le.Caption != "" || le.Text1 != "" || le.Text2 != "") {
		part := doc.Connections[connIndex].Parts[partIndex]
		if !le.IsInXRange(part.StartX, part.EndX) {
			return
		}
		if isIn, isMoreDown, dist := isInDistance(y, le.Y+le.Height, doc.MinConnectorMargin); isIn {
			if isMoreDown {
				// the connector is to close to the lower border of the box ... move the connector down
				// ... and all objects equal or lower to y will be moved down, too
				doc.moveAllObjectsLowerDown(&doc.Boxes, y, dist, true)
				doc.moveAllConnectorsLowerDown(y, dist)
			} else {
				// the connector is to close to the upper border of the box ... move all objects lower than y, down
				doc.moveAllObjectsLowerDown(&doc.Boxes, y, dist, false)
				doc.moveAllConnectorsLowerDown(y, dist)
			}
			return
		}
		if isIn, isMoreDown, dist := isInDistance(y, le.Y, doc.MinConnectorMargin); isIn {
			if isMoreDown {
				// the connector is to close to the upper border of the box ... move the connector down
				// ... and all objects equal or lower to y will be moved down, too
				doc.moveAllObjectsLowerDown(&doc.Boxes, y, dist, false)
				doc.moveAllConnectorsLowerDown(y, dist)
			} else {
				// the connector is to close to the upper border of the box ... move all objects lower than y, down
				doc.moveAllObjectsLowerDown(&doc.Boxes, y, dist, false)
				doc.moveAllConnectorsLowerDown(y, dist)
			}
			return
		}
	}
	if le.Horizontal != nil {
		for i := 0; i < len(le.Horizontal.Elems); i++ {
			doc.checkLayoutElemForVerticalCollision(&le.Horizontal.Elems[i], connIndex, partIndex, y)
		}
	}
	if le.Vertical != nil {
		for i := 0; i < len(le.Vertical.Elems); i++ {
			doc.checkLayoutElemForVerticalCollision(&le.Vertical.Elems[i], connIndex, partIndex, y)
		}
	}
}

func (doc *BoxesDocument) checkLayoutElemForHorizontalCollision(le *LayoutElement, connIndex, partIndex, x int) {
	if le != &doc.Boxes && (le.Id != "" || le.Caption != "" || le.Text1 != "" || le.Text2 != "") {
		// ignores the first box of the document
		part := doc.Connections[connIndex].Parts[partIndex]
		if !le.IsInYRange(part.StartY, part.EndY) {
			return
		}
		if isIn, isMoreRight, dist := isInDistance(x, le.X+le.Width, doc.MinConnectorMargin); isIn {
			if isMoreRight {
				// the connector is to close to the right border of the box ... move the connector to the right
				// ... and all objects equal or more to the right will be moved to the right, too
				doc.moveAllObjectsMoreRightRight(&doc.Boxes, x, dist, true)
				doc.moveAllConnectorsMoreRightRight(x, dist)
			} else {
				// the connector is to close to the left border of the box ... move all objects more to the right of x
				// to the right
				doc.moveAllObjectsMoreRightRight(&doc.Boxes, x, dist, false)
				doc.moveAllConnectorsMoreRightRight(x, dist)
			}
			return
		}
		if isIn, isMoreRight, dist := isInDistance(x, le.X, doc.MinConnectorMargin); isIn {
			if isMoreRight {
				// the connector is to close to the right border of the box ... move the connector to the right
				// ... and all objects equal or more to the right will be moved to the right, too
				doc.moveAllObjectsMoreRightRight(&doc.Boxes, x, dist, true)
				doc.moveAllConnectorsMoreRightRight(x, dist)
			} else {
				// the connector is to close to the left border of the box ... move all objects more to the right of x
				// to the right
				doc.moveAllObjectsMoreRightRight(&doc.Boxes, x, dist, false)
				doc.moveAllConnectorsMoreRightRight(x, dist)
			}
			return
		}
	}
	if le.Horizontal != nil {
		for i := 0; i < len(le.Horizontal.Elems); i++ {
			doc.checkLayoutElemForHorizontalCollision(&le.Horizontal.Elems[i], connIndex, partIndex, x)
		}
	}
	if le.Vertical != nil {
		for i := 0; i < len(le.Vertical.Elems); i++ {
			doc.checkLayoutElemForHorizontalCollision(&le.Vertical.Elems[i], connIndex, partIndex, x)
		}
	}
}

// moves all vertical connection lines that are too close to the side borders of the boxes.
// In case of a collision all impacted objects are moved to the right.
func (doc *BoxesDocument) moveTooCloseVerticalConnectionLinesFromBorders() {
	for i := 0; i < len(doc.Connections); i++ {
		for j := 0; j < len(doc.Connections[i].Parts); j++ {
			part := doc.Connections[i].Parts[j]
			if part.StartX != part.EndX {
				continue // is no vertical line
			}
			doc.checkLayoutElemForHorizontalCollision(&doc.Boxes, i, j, part.StartX)
			// in case that the left border triggered a change, we have to check the right border, too
			if part.StartX != doc.Connections[i].Parts[j].StartX || part.StartY != doc.Connections[i].Parts[j].StartY {
				doc.checkLayoutElemForHorizontalCollision(&doc.Boxes, i, j, doc.Connections[i].Parts[j].StartX)
			}
		}
	}
}

// moves all horizontal connection lines that are too close to the top or bottom borders of the boxes.
// In case of a collision all impacted objects are moved down.
func (doc *BoxesDocument) moveTooCloseHorizontalConnectionLinesFromBorders() {
	for i := 0; i < len(doc.Connections); i++ {
		conn := doc.Connections[i]
		for j := 0; j < len(conn.Parts); j++ {
			part := conn.Parts[j]
			if part.StartY != part.EndY {
				continue // is no horizontal line
			}
			doc.checkLayoutElemForVerticalCollision(&doc.Boxes, i, j, part.StartY)
			if part.StartX != doc.Connections[i].Parts[j].StartX || part.StartY != doc.Connections[i].Parts[j].StartY {
				doc.checkLayoutElemForVerticalCollision(&doc.Boxes, i, j, doc.Connections[i].Parts[j].StartY)
			}
		}
	}
}

// moves all vertical connection lines that are too close to the other vertical connection lines.
// In case of a collision all impacted objects are moved to the left.
func (doc *BoxesDocument) moveTooCloseVerticalConnectionLines() {
	for i := 0; i < len(doc.Connections); i++ {
		conn := doc.Connections[i]
		for j := 0; j < len(conn.Parts); j++ {
			part := conn.Parts[j]
			if part.StartX != part.EndX {
				continue // is no vertical line
			}
			// TODO
		}
	}
}

// moves all horizontal connection lines that are too close to the other horizontal connection lines.
// In case of a collision all impacted objects are moved down.
func (doc *BoxesDocument) moveTooCloseHorizontalConnectionLines() {
	for i := 0; i < len(doc.Connections); i++ {
		conn := doc.Connections[i]
		for j := 0; j < len(conn.Parts); j++ {
			part := conn.Parts[j]
			if part.StartY != part.EndY {
				continue // is no horizontal line
			}
			// TODO
		}
	}
}

// removes overlapping connection line ends - and marks them to be drawn as circles
func (doc *BoxesDocument) truncateJoiningConnectionLines() {
	// TODO
}
