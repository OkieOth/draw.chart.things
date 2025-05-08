package types

func isInDistance(v1, v2, distance int) (bool, bool, int) {
	if v1 > v2 {
		dist := v1 - v2
		return (dist <= distance), true, dist
	}
	dist := v2 - v1
	return (dist <= distance), false, dist
}

func (doc *BoxesDocument) moveAllObjectsLowerAndEqualDown(le *LayoutElement, y, dist int) {
	// TODO
}
func (doc *BoxesDocument) moveAllObjectsLowerDown(le *LayoutElement, y, dist int) {
	// TODO
}
func (doc *BoxesDocument) moveAllObjectsMoreRightAndEqualRight(le *LayoutElement, x, dist int) {
	// TODO
}
func (doc *BoxesDocument) moveAllObjectsMoreRightRight(le *LayoutElement, x, dist int) {
	// TODO
}

func (doc *BoxesDocument) moveAllConnectorsLowerAndEqualDown(y, dist int) {
	for i := 0; i < len(doc.Connections); i++ {
		for j := 0; j < len(doc.Connections[i].Parts); j++ {
			part := doc.Connections[i].Parts[j]
			if part.StartY != part.EndY {
				// vertical line
				if part.StartY < part.EndY {
					// line moves down
					if part.EndY > y {
						if part.StartY >= y {
							// line starts after or at the y position
							doc.Connections[i].Parts[j].StartY += dist
						}
						doc.Connections[i].Parts[j].EndY += dist
					} else if part.EndY == y {
						// line ends at the y position
						doc.Connections[i].Parts[j].EndY += dist
					}
				} else {
					// line moves up
					if part.StartY > y {
						if part.EndY >= y {
							// line starts before or at the y position
							doc.Connections[i].Parts[j].EndY += dist
						}
						doc.Connections[i].Parts[j].StartY += dist
					} else if part.StartY == y {
						// line starts at the y position
						doc.Connections[i].Parts[j].StartY += dist
					}
				}
			} else {
				// horizontal line
				if part.StartY >= y {
					doc.Connections[i].Parts[j].StartY += dist
					doc.Connections[i].Parts[j].EndY += dist
				}
			}
		}
	}
}
func (doc *BoxesDocument) moveAllConnectorsLowerDown(y, dist int) {
	for i := 0; i < len(doc.Connections); i++ {
		for j := 0; j < len(doc.Connections[i].Parts); j++ {
			part := doc.Connections[i].Parts[j]
			if part.StartY != part.EndY {
				// vertical line
				if part.StartY < part.EndY {
					// line moves down
					if part.EndY > y {
						if part.StartY > y {
							// line starts after the y position
							doc.Connections[i].Parts[j].StartY += dist
						}
						doc.Connections[i].Parts[j].EndY += dist
					}
				} else {
					// line moves up
					if part.StartY > y {
						if part.EndY > y {
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
			// TODO
		}
	}
}
func (doc *BoxesDocument) moveAllConnectorsMoreRightAndEqualRight(x, dist int) {
	for i := 0; i < len(doc.Connections); i++ {
		for j := 0; j < len(doc.Connections[i].Parts); j++ {
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
						if part.StartX > x {
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
						if part.EndX > x {
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
		for j := 0; j < len(doc.Connections[i].Parts); j++ {
			part := doc.Connections[i].Parts[j]
			if part.StartY != part.EndY {
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
						if part.StartX > x {
							// line starts after the x position
							doc.Connections[i].Parts[j].StartX += dist
						}
						doc.Connections[i].Parts[j].EndX += dist
					}
				} else {
					// line moves to the left
					if part.StartX > x {
						if part.EndX > x {
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
	if isIn, isMoreDown, dist := isInDistance(y, le.Y, doc.MinConnectorMargin); isIn {
		if isMoreDown {
			// the connector is to close to the lower border of the box ... move the connector down
			// ... and all objects equal or lower to y will be moved down, too
			doc.moveAllObjectsLowerAndEqualDown(&doc.Boxes, y, dist)
			doc.moveAllConnectorsLowerAndEqualDown(y, dist)
		} else {
			// the connector is to close to the upper border of the box ... move all objects lower than y, down
			doc.moveAllObjectsLowerDown(&doc.Boxes, y, dist)
			doc.moveAllConnectorsLowerDown(y, dist)
		}
	}
}

func (doc *BoxesDocument) checkLayoutElemForHorizontalCollision(le *LayoutElement, connIndex, partIndex, x int) {
	if isIn, isMoreRight, dist := isInDistance(x, le.X, doc.MinConnectorMargin); isIn {
		if isMoreRight {
			// the connector is to close to the right border of the box ... move the connector to the right
			// ... and all objects equal or more to the right will be moved to the right, too
			doc.moveAllObjectsMoreRightAndEqualRight(&doc.Boxes, x, dist)
			doc.moveAllConnectorsMoreRightAndEqualRight(x, dist)
		} else {
			// the connector is to close to the left border of the box ... move all objects more to the right of x
			// to the right
			doc.moveAllObjectsMoreRightRight(&doc.Boxes, x, dist)
			doc.moveAllConnectorsMoreRightRight(x, dist)
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
