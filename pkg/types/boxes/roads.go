package boxes

import (
	"fmt"
	"slices"
	"sort"

	"github.com/okieoth/draw.chart.things/pkg/types"
)

// traces the document and finds all reasonable "roads" to connect
func (doc *BoxesDocument) InitRoads() {
	doc.initRoadsImpl(&doc.Boxes, DefRoadType_All)
	// sort roads
	sort.Slice(doc.VerticalRoads, func(i, j int) bool {
		return doc.VerticalRoads[i].StartX < doc.VerticalRoads[j].StartX
	})
	sort.Slice(doc.HorizontalRoads, func(i, j int) bool {
		return doc.HorizontalRoads[i].StartY < doc.HorizontalRoads[j].StartY
	})
}

func (doc *BoxesDocument) addVerticalRoad(line ConnectionLine) {
	l := line
	if line.StartX > line.EndX || line.StartY > line.EndY {
		l = newConnectionLine(line.EndX, line.EndY, line.StartX, line.StartY)
	}
	for i := range len(doc.VerticalRoads) {
		cl := &doc.VerticalRoads[i]
		if (cl.StartX == l.StartX) &&
			(cl.EndX == l.EndX) &&
			(cl.StartY == l.StartY) &&
			(cl.EndY == l.EndY) {
			// already exists
			return
		}
		if cl.StartX == l.StartX {
			if (l.StartY < cl.StartY) && (l.EndY >= cl.StartY) {
				// extend the existing line to the top
				cl.StartY = l.StartY
				return
			} else if (l.EndY > cl.EndY) && (l.StartY <= cl.EndY) {
				// extend the existing line to the bottom
				cl.EndY = l.EndY
				return
			}
		}
	}
	doc.VerticalRoads = append(doc.VerticalRoads, l)
}

func (doc *BoxesDocument) addHorizontalRoad(line ConnectionLine) {
	l := line
	if line.StartX > line.EndX || line.StartY > line.EndY {
		l = newConnectionLine(line.EndX, line.EndY, line.StartX, line.StartY)
	}
	for i := range len(doc.HorizontalRoads) {
		cl := &doc.HorizontalRoads[i]
		if (cl.StartX == l.StartX) &&
			(cl.EndX == l.EndX) &&
			(cl.StartY == l.StartY) &&
			(cl.EndY == l.EndY) {
			// already exists
			return
		}
		if cl.StartY == l.StartY {
			if (l.StartX < cl.StartX) && (l.EndX >= cl.StartX) {
				// extend the existing line to the left
				cl.StartX = l.StartX
				return
			} else if (l.EndX > cl.EndX) && (l.StartX <= cl.EndX) {
				// extend the existing line to the right
				cl.EndX = l.EndX
				return
			}
		}
	}
	doc.HorizontalRoads = append(doc.HorizontalRoads, l)
}

func (doc *BoxesDocument) pointHasCollision(x, y int, startElem *LayoutElement, isForHorizontalLine bool) bool {
	switch doc.checkColl(x, y, &doc.Boxes, startElem, nil, isForHorizontalLine) {
	case CollisionType_WithElem:
		return true
	case CollisionType_WithSurroundings:
		return true
	default:
		return false
	}
}

func (doc *BoxesDocument) roadUp(line *ConnectionLine, startElem *LayoutElement) {
	if line.EndY <= 0 {
		line.EndY = 0
		return
	}
	switch doc.checkColl(line.EndX, line.EndY, &doc.Boxes, startElem, nil, false) {
	case CollisionType_WithElem:
		line.EndY += 2 * types.RasterSize
		return
	case CollisionType_WithSurroundings:
		line.EndY += types.RasterSize
		return
	default:
		line.EndY -= types.RasterSize
		doc.roadUp(line, startElem)
	}
}

func (doc *BoxesDocument) roadDown(line *ConnectionLine, startElem *LayoutElement) {
	if line.EndY >= doc.Height {
		line.EndY = doc.Height
		return
	}
	switch doc.checkColl(line.EndX, line.EndY, &doc.Boxes, startElem, nil, false) {
	case CollisionType_WithElem:
		line.EndY -= 2 * types.RasterSize
		return
	case CollisionType_WithSurroundings:
		line.EndY -= types.RasterSize
		return
	default:
		line.EndY += types.RasterSize
		doc.roadDown(line, startElem)
	}
}

func (doc *BoxesDocument) roadLeft(line *ConnectionLine, startElem *LayoutElement) {
	if line.EndX <= 0 {
		line.EndX = 0
		return
	}
	switch doc.checkColl(line.EndX, line.EndY, &doc.Boxes, startElem, nil, true) {
	case CollisionType_WithElem:
		line.EndX += 2 * types.RasterSize
		return
	case CollisionType_WithSurroundings:
		line.EndX += types.RasterSize // bug???
		return
	default:
		line.EndX -= types.RasterSize
		doc.roadLeft(line, startElem)
	}
}

func (doc *BoxesDocument) roadRight(line *ConnectionLine, startElem *LayoutElement) {
	if line.EndX >= doc.Width {
		line.EndX = doc.Width
		return
	}
	switch doc.checkColl(line.EndX, line.EndY, &doc.Boxes, startElem, nil, true) {
	case CollisionType_WithElem:
		line.EndX -= 2 * types.RasterSize
		return
	case CollisionType_WithSurroundings:
		line.EndX -= types.RasterSize
		return
	default:
		line.EndX += types.RasterSize
		doc.roadRight(line, startElem)
	}
}

func (doc *BoxesDocument) elemHasParentWithTextImpl(elemToCheck, currentElem *LayoutElement, parentTxt bool) bool {
	if currentElem.Id != "" && currentElem.Id == elemToCheck.Id {
		// reached end of recursion
		return parentTxt
	}

	var textOverlaps bool
	if currentElem.XTextBox != nil && currentElem.WidthTextBox != nil {
		textOverlaps = elemToCheck.CenterX >= *currentElem.XTextBox && elemToCheck.CenterX <= (*currentElem.XTextBox+*currentElem.WidthTextBox)
	}

	myParentTxt := textOverlaps || parentTxt
	if doc.elemHasParentWithTextCont(elemToCheck, myParentTxt, currentElem.Horizontal) {
		return true
	}
	return doc.elemHasParentWithTextCont(elemToCheck, myParentTxt, currentElem.Vertical)
}

func (doc *BoxesDocument) elemHasParentWithTextCont(elem *LayoutElement, parentTxt bool, cont *LayoutElemContainer) bool {
	if cont == nil {
		return false
	}
	for _, e := range cont.Elems {
		if doc.elemHasParentWithTextImpl(elem, &e, parentTxt) {
			return true
		}
	}
	return false
}

func (doc *BoxesDocument) elemHasParentWithText(elem *LayoutElement) bool {
	return doc.elemHasParentWithTextImpl(elem, &doc.Boxes, false)
}

type DefRoadType int

const (
	DefRoadType_All DefRoadType = iota
	DefRoadType_Vertical
	DefRoadType_Horizontal
	DefRoadType_None
)

func (doc *BoxesDocument) initRoadsImpl(elem *LayoutElement, defRoadType DefRoadType) {
	if elem.Caption == "" && elem.Text1 == "" && elem.Text2 == "" {
		defRoadType = DefRoadType_None
	}
	stepSize := 2 * types.RasterSize
	if elem.TopXToStart != nil { // elem has connections by itself
		// draw line from the top x start, till the first collision
		if !doc.pointHasCollision(*elem.TopXToStart, elem.Y+stepSize, elem, false) {
			upRoad := newConnectionLine(*elem.TopXToStart, elem.Y, *elem.TopXToStart, elem.Y+stepSize)
			doc.roadUp(&upRoad, elem)
			doc.addVerticalRoad(upRoad)
		}
		// draw line from the bottom x start, till the first collision
		if !doc.pointHasCollision(*elem.BottomXToStart, elem.Y+elem.Height+stepSize, elem, false) {
			downRoad := newConnectionLine(*elem.BottomXToStart, elem.Y+elem.Height,
				*elem.BottomXToStart, elem.Y+elem.Height+stepSize)
			doc.roadDown(&downRoad, elem)
			doc.addVerticalRoad(downRoad)
		}

		// draw line from the left y start, till the first collision
		if !doc.pointHasCollision(elem.X-stepSize, *elem.LeftYToStart, elem, true) {
			leftRoad := newConnectionLine(elem.X, *elem.LeftYToStart, elem.X-stepSize, *elem.LeftYToStart)
			doc.roadLeft(&leftRoad, elem)
			doc.addHorizontalRoad(leftRoad)
		}
		// draw line from the right y start, till the first collision
		if !doc.pointHasCollision(elem.X+elem.Width+stepSize, *elem.RightYToStart, elem, true) {
			rightRoad := newConnectionLine(elem.X+elem.Width, *elem.RightYToStart, elem.X+elem.Width+stepSize, *elem.RightYToStart)
			doc.roadRight(&rightRoad, elem)
			doc.addHorizontalRoad(rightRoad)
		}
	}
	if elem.TopXToStart != nil {
		defRoadType = DefRoadType_All
	}
	if defRoadType == DefRoadType_All {
		// draw the line parallel to the right border, till the first collision, up and down
		if !doc.pointHasCollision(elem.X+elem.Width+stepSize, elem.CenterY, elem, false) {
			upRoad := newConnectionLine(elem.X+elem.Width+stepSize, elem.CenterY, elem.X+elem.Width+stepSize, elem.Y-stepSize)
			doc.roadUp(&upRoad, elem)
			downRoad := newConnectionLine(elem.X+elem.Width+stepSize, elem.CenterY, elem.X+elem.Width+stepSize, elem.Y+elem.Height+stepSize)
			doc.roadDown(&downRoad, elem)
			l := newConnectionLine(upRoad.EndX, upRoad.EndY, downRoad.EndX, downRoad.EndY)
			doc.addVerticalRoad(l)
		}
	}
	if defRoadType == DefRoadType_Horizontal || defRoadType == DefRoadType_All {
		// draw the line parallel to the left border, till the first collision, up and down
		if !doc.pointHasCollision(elem.X-stepSize, elem.CenterY, elem, false) {
			upRoad := newConnectionLine(elem.X-stepSize, elem.CenterY, elem.X-stepSize, elem.Y-stepSize)
			doc.roadUp(&upRoad, elem)
			downRoad := newConnectionLine(elem.X-stepSize, elem.CenterY, elem.X-stepSize, elem.Y+elem.Height+stepSize)
			doc.roadDown(&downRoad, elem)
			l := newConnectionLine(upRoad.EndX, upRoad.EndY, downRoad.EndX, downRoad.EndY)
			doc.addVerticalRoad(l)
		}
	}
	if defRoadType == DefRoadType_All {
		// draw the line parallel to the top border, till the first collision, left and right
		if !doc.pointHasCollision(elem.CenterX, elem.Y-stepSize, elem, true) {
			leftRoad := newConnectionLine(elem.CenterX, elem.Y-stepSize, elem.X-stepSize, elem.Y-stepSize)
			doc.roadLeft(&leftRoad, elem)
			rightRoad := newConnectionLine(elem.CenterX, elem.Y-stepSize, elem.X+elem.Width+stepSize, elem.Y-stepSize)
			doc.roadRight(&rightRoad, elem)
			l := newConnectionLine(leftRoad.EndX, leftRoad.EndY, rightRoad.EndX, rightRoad.EndY)
			doc.addHorizontalRoad(l)
		}
	}
	if defRoadType == DefRoadType_Vertical || defRoadType == DefRoadType_All {
		// draw the line parallel to the bottom border, till the first collision, left and right
		if !doc.pointHasCollision(elem.CenterX, elem.Y+elem.Height+stepSize, elem, true) {
			leftRoad := newConnectionLine(elem.CenterX, elem.Y+elem.Height+stepSize, elem.X-stepSize, elem.Y+elem.Height+stepSize)
			doc.roadLeft(&leftRoad, elem)
			rightRoad := newConnectionLine(elem.CenterX, elem.Y+elem.Height+stepSize, elem.X+elem.Width+stepSize, elem.Y+elem.Height+stepSize)
			doc.roadRight(&rightRoad, elem)
			l := newConnectionLine(leftRoad.EndX, leftRoad.EndY, rightRoad.EndX, rightRoad.EndY)
			doc.addHorizontalRoad(l)
		}
	}
	if elem.Vertical != nil {
		for i := range len(elem.Vertical.Elems) {
			defRoadType := DefRoadType_Vertical
			if i == 0 {
				defRoadType = DefRoadType_None
			}
			doc.initRoadsImpl(&elem.Vertical.Elems[i], defRoadType)
		}
	}
	if elem.Horizontal != nil {
		for i := range len(elem.Horizontal.Elems) {
			defRoadType := DefRoadType_Horizontal
			if i == 0 {
				defRoadType = DefRoadType_None
			}
			doc.initRoadsImpl(&elem.Horizontal.Elems[i], defRoadType)
		}
	}
}

// function goes over the horizontal roads, finds a related road and
// returns as first value the most left value of the road. The second value
// is the list of x-coordinates where the road has leftwards junctions up.
// The third value is the list of x-coordinates where the road has
// leftwards junctions down.
func (doc *BoxesDocument) remove_checkRoadsToTheLeft(startX, startY, minX int) (int, []int, []int, error) {
	for _, h := range doc.HorizontalRoads {
		leftX, rightX := minMax(h.StartX, h.EndX)
		leftX2, rightX2 := minMax(leftX, minX)
		if (h.StartY == startY) && (leftX <= startX) && (rightX >= startX) {
			// found the road
			upJunctions := make([]int, 0)
			downJunctions := make([]int, 0)
			retX := rightX2
			for _, v := range doc.VerticalRoads {
				if (v.StartX >= leftX2) && (v.StartX <= rightX2) {
					minY, maxY := minMax(v.StartY, v.EndY)
					if startY >= minY && startY <= maxY {
						if minY < startY {
							// junction upwards
							if !slices.Contains(upJunctions, v.StartX) {
								upJunctions = append(upJunctions, v.StartX)
							}
						}
						if maxY > startY {
							// junction downwards
							if !slices.Contains(downJunctions, v.StartX) {
								downJunctions = append(downJunctions, v.StartX)
							}
						}
					}
				}
			}
			// sort the junctions, so they are going to the left
			types.SortDescending(upJunctions)
			types.SortAscending(downJunctions)
			return retX, upJunctions, downJunctions, nil
		}
	}
	return 0, nil, nil, fmt.Errorf("no road found to move left for %d,%d", startX, startY)
}

// function goes over the horizontal roads, finds a related road and
// returns as first value the most right value of the road. The second value
// is the list of x-coordinates where the road has rightwards junctions up.
// The third value is the list of x-coordinates where the road has
// rightwards junctions down.
func (doc *BoxesDocument) remove_checkRoadsToTheRight(startX, startY, maxX int) (int, []int, []int, error) {
	for _, h := range doc.HorizontalRoads {
		leftX, rightX := minMax(h.StartX, h.EndX)
		leftX2, rightX2 := minMax(rightX, maxX)
		if (h.StartY == startY) && (leftX <= startX) && (rightX >= startX) {
			// found the road
			upJunctions := make([]int, 0)
			downJunctions := make([]int, 0)
			retX := leftX2
			for _, v := range doc.VerticalRoads {
				if (v.StartX > leftX2) && (v.StartX <= rightX2) {
					minY, maxY := minMax(v.StartY, v.EndY)
					if startY >= minY && startY <= maxY {
						if minY < startY {
							// junction upwards
							if !slices.Contains(upJunctions, v.StartX) {
								upJunctions = append(upJunctions, v.StartX)
							}
						}
						if maxY > startY {
							// junction downwards
							if !slices.Contains(downJunctions, v.StartX) {
								downJunctions = append(downJunctions, v.StartX)
							}
						}
					}
				}
			}
			// sort the junctions, so they are going to the left
			types.SortDescending(upJunctions)
			types.SortAscending(downJunctions)
			return retX, upJunctions, downJunctions, nil
		}
	}
	return 0, nil, nil, fmt.Errorf("no road found to move right for %d,%d", startX, startY)
}

// function goes over the vertical roads, finds a related road and
// returns as first value the most top value (lowest Y) of the road. The second value
// is the list of y-coordinates where the road has upwards junctions to the left.
// The third value is the list of y-coordinates where the road has
// upwards junctions to the right.
func (doc *BoxesDocument) remove_checkRoadsToTheTop(startX, startY, minY int) (int, []int, []int, error) {
	for _, v := range doc.VerticalRoads {
		topY, bottomY := minMax(v.StartY, v.EndY)
		topY2, bottomY2 := minMax(topY, minY)

		if (v.StartX == startX) && (topY <= startY) && (bottomY >= startY) {
			// found the road
			leftJunctions := make([]int, 0)
			rightJunctions := make([]int, 0)
			retY := bottomY2
			for _, h := range doc.HorizontalRoads {
				if (h.StartY >= topY2) && (h.StartY <= bottomY2) {
					minX, maxX := minMax(h.StartX, h.EndX)
					if startX >= minX && startX <= maxX {
						if minX < startX {
							// junction to the left
							if !slices.Contains(leftJunctions, h.StartY) {
								leftJunctions = append(leftJunctions, h.StartY)
							}
						}
						if maxX > startX {
							// junction to the right
							if !slices.Contains(rightJunctions, h.StartY) {
								rightJunctions = append(rightJunctions, h.StartY)
							}
						}
					}
				}
			}
			// sort the junctions, so they are going to the left
			types.SortDescending(leftJunctions)
			types.SortAscending(rightJunctions)
			return retY, leftJunctions, rightJunctions, nil
		}
	}
	return 0, nil, nil, fmt.Errorf("no road found to move upwards for %d,%d", startX, startY)
}

// function goes over the vertical roads, finds a related road and
// returns as first value the most down value (biggest Y) of the road. The second value
// is the list of y-coordinates where the road has downwards junctions to the left.
// The third value is the list of y-coordinates where the road has
// downwards junctions to the right.
func (doc *BoxesDocument) remove_checkRoadsToTheBottom(startX, startY, maxY int) (int, []int, []int, error) {
	for _, v := range doc.VerticalRoads {
		topY, bottomY := minMax(v.StartY, v.EndY)
		topY2, bottomY2 := minMax(bottomY, maxY)
		if (v.StartX == startX) && (topY <= startY) && (bottomY >= startY) {
			// found the road
			leftJunctions := make([]int, 0)
			rightJunctions := make([]int, 0)
			retY := topY2
			for _, h := range doc.HorizontalRoads {
				if (h.StartY >= topY2) && (h.StartY <= bottomY2) {
					minX, maxX := minMax(h.StartX, h.EndX)
					if startX >= minX && startX <= maxX {
						if minX < startX {
							// junction to the left
							leftJunctions = append(leftJunctions, h.StartY)
						}
						if maxX > startX {
							// junction to the right
							rightJunctions = append(rightJunctions, h.StartY)
						}
					}
				}
			}
			// sort the junctions, so they are going to the left
			types.SortDescending(leftJunctions)
			types.SortAscending(rightJunctions)
			return retY, leftJunctions, rightJunctions, nil
		}
	}
	return 0, nil, nil, fmt.Errorf("no road found to move downwards for %d,%d", startX, startY)
}

// Search and return the next junction to the left of the given coordinates.
// The function returns the x-coordinate of the junction, a boolean indicating if you reach from
// this junction other junctions upwards, downwards or in the same direction. It also returns the most
// left x-coordinate of the road. If no junction is found, it returns an error.
// example call: x, straightAhead, upwards, downwards, mostLeftX, err := doc.getNextJunctionLeft(startX, startY)
func (doc *BoxesDocument) getNextJunctionLeft(startX, startY int) (int, bool, bool, bool, int, error) {
	for _, h := range doc.HorizontalRoads {
		leftX, rightX := minMax(h.StartX, h.EndX)
		if (h.StartY == startY) && (leftX <= startX) && (rightX >= startX) {
			// found the road
			// search for next junction to the left
			var nextRoad *ConnectionLine
			var upwards, downwards, straightAhead bool
			for _, v := range doc.VerticalRoads {
				if (v.StartX >= leftX) && (v.StartX <= rightX) && (v.StartX < startX) {
					if (nextRoad == nil) || (v.StartX > nextRoad.StartX) {
						minY, maxY := minMax(v.StartY, v.EndY)
						if startY >= minY && startY <= maxY {
							nextRoad = &v
							straightAhead = leftX < nextRoad.StartX
							upwards = minY < startY
							downwards = maxY > startY
						}
					}
				}
			}
			if nextRoad != nil {
				// DEBUG
				// DEBUG
				straightAhead = true
				upwards = true
				downwards = true
				return nextRoad.StartX, upwards, downwards, straightAhead, leftX, nil
			} else {
				return leftX, false, false, false, leftX, nil
			}
		}
	}
	return 0, false, false, false, 0, fmt.Errorf("no horizontal road found to move left from %d,%d", startX, startY)
}

// Search and return the next junction to the right of the given coordinates.
// The function returns the x-coordinate of the junction, a boolean indicating if you reach from
// this junction other junctions upwards, downwards or in the same direction. It also returns the most
// right x-coordinate of the road. If no junction is found, it returns an error.
func (doc *BoxesDocument) getNextJunctionRight(startX, startY int) (int, bool, bool, bool, int, error) {
	for _, h := range doc.HorizontalRoads {
		leftX, rightX := minMax(h.StartX, h.EndX)
		if (h.StartY == startY) && (leftX <= startX) && (rightX > startX) {
			// found the road
			// search for next junction to the left
			var nextRoad *ConnectionLine
			var upwards, downwards, straightAhead bool
			for _, v := range doc.VerticalRoads {
				if (v.StartX >= leftX) && (v.StartX <= rightX) && (v.StartX >= startX) {
					if (nextRoad == nil) || (v.StartX < nextRoad.StartX) {
						minY, maxY := minMax(v.StartY, v.EndY)
						if startY >= minY && startY <= maxY {
							nextRoad = &v
							straightAhead = rightX > nextRoad.StartX
							upwards = minY < startY
							downwards = maxY > startY
						}
					}
				}
			}
			if nextRoad != nil {
				// DEBUG
				straightAhead = true
				upwards = true
				downwards = true
				return nextRoad.StartX, upwards, downwards, straightAhead, rightX, nil
			} else {
				return rightX, false, false, false, rightX, nil
			}
		}
	}
	return 0, false, false, false, 0, fmt.Errorf("no horizontal road found to move right from %d,%d", startX, startY)
}

// Search and return the next junction to the top of the given coordinates.
// The function returns the y-coordinate of the junction, a boolean indicating if you reach from
// this junction other junctions leftwards, rightwards or in the same direction. It also returns the most
// top y-coordinate (lower thant startY) of the road. If no junction is found, it returns an error.
func (doc *BoxesDocument) getNextJunctionUp(startX, startY int) (int, bool, bool, bool, int, error) {
	for _, v := range doc.VerticalRoads {
		topY, bottomY := minMax(v.StartY, v.EndY)

		if (v.StartX == startX) && (topY <= startY) && (bottomY >= startY) {
			// found the road
			// search for next junction up
			var nextRoad *ConnectionLine
			var leftwards, rightwards, straightAhead bool
			for _, h := range doc.HorizontalRoads {
				if (h.StartY < bottomY) && (h.StartY >= topY) && (h.StartY < startY) {
					if (nextRoad == nil) || (h.StartY > nextRoad.StartY) {
						minX, maxX := minMax(h.StartX, h.EndX)
						if startX >= minX && startX <= maxX {
							nextRoad = &h
							straightAhead = topY < nextRoad.StartY // what if the road is in the inverse direction?
							leftwards = minX < startX
							rightwards = maxX > startX
						}
					}
				}
			}
			if nextRoad != nil {
				// DEBUG
				straightAhead = true
				leftwards = true
				rightwards = true
				return nextRoad.StartY, leftwards, rightwards, straightAhead, topY, nil
			} else {
				return topY, false, false, false, topY, nil
			}
		}
	}
	return 0, false, false, false, 0, fmt.Errorf("no vertical road found to move up from %d,%d", startX, startY)
}

// Search and return the next junction to the bottom of the given coordinates.
// The function returns the y-coordinate of the junction, a boolean indicating if you reach from
// this junction other junctions leftwards, rightwards or in the same direction. It also returns the most
// bottom y-coordinate of the road. If no junction is found, it returns an error.
func (doc *BoxesDocument) getNextJunctionDown(startX, startY int) (int, bool, bool, bool, int, error) {
	for _, v := range doc.VerticalRoads {
		topY, bottomY := minMax(v.StartY, v.EndY)

		if (v.StartX == startX) && (topY <= startY) && (bottomY >= startY) {
			// found the road
			// search for next junction down
			var nextRoad *ConnectionLine
			var leftwards, rightwards, straightAhead bool
			for _, h := range doc.HorizontalRoads {
				if (h.StartY <= bottomY) && (h.StartY >= topY) && (h.StartY > startY) {
					if (nextRoad == nil) || (h.StartY < nextRoad.StartY) {
						minX, maxX := minMax(h.StartX, h.EndX)
						if startX >= minX && startX <= maxX {
							nextRoad = &h
							straightAhead = bottomY > nextRoad.StartY // what if the road is in the inverse direction?
							leftwards = minX < startX
							rightwards = maxX > startX
						}
					}
				}
			}
			if nextRoad != nil {
				// DEBUG
				straightAhead = true
				leftwards = true
				rightwards = true
				return nextRoad.StartY, leftwards, rightwards, straightAhead, bottomY, nil
			} else {
				return bottomY, false, false, false, bottomY, nil
			}
		}
	}
	return 0, false, false, false, 0, fmt.Errorf("no vertical road found to move down from %d,%d", startX, startY)
}
