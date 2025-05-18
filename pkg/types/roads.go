package types

import (
	"fmt"
	"slices"
	"sort"
)

// traces the document and finds all reasonable "roads" to connect
func (doc *BoxesDocument) InitRoads() {
	doc.initRoadsImpl(&doc.Boxes)
}

func (doc *BoxesDocument) addRoad(line ConnectionLine, roads *[]ConnectionLine) {
	for i := 0; i < len(*roads); i++ {
		l := (*roads)[i]
		if (l.StartX == line.StartX) &&
			(l.EndX == line.EndX) &&
			(l.StartY == line.StartY) &&
			(l.EndY == line.EndY) {
			// already exists
			return
		}
		if (l.StartX == line.EndX) &&
			(l.EndX == line.StartX) &&
			(l.StartY == line.EndY) &&
			(l.EndY == line.StartY) {
			// already exists
			return
		}
	}
	*roads = append(*roads, line)
}

func (doc *BoxesDocument) roadUp(line *ConnectionLine, startElem *LayoutElement) {
	if line.EndY == 0 {
		return
	}
	if doc.checkColl(line.EndX, line.EndY, &doc.Boxes, startElem, nil) {
		// has collision
		line.EndY += 2 * RasterSize
		return
	} else {
		line.EndY -= RasterSize
		doc.roadUp(line, startElem)
	}
}

func (doc *BoxesDocument) roadDown(line *ConnectionLine, startElem *LayoutElement) {
	if line.EndY == doc.Height {
		return
	}
	if doc.checkColl(line.EndX, line.EndY, &doc.Boxes, startElem, nil) {
		// has collision
		line.EndY -= 2 * RasterSize
		return
	} else {
		line.EndY += RasterSize
		doc.roadDown(line, startElem)
	}
}

func (doc *BoxesDocument) roadLeft(line *ConnectionLine, startElem *LayoutElement) {
	if line.EndX == 0 {
		return
	}
	if doc.checkColl(line.EndX, line.EndY, &doc.Boxes, startElem, nil) {
		// has collision
		line.EndX += 2 * RasterSize
		return
	} else {
		line.EndX -= RasterSize
		doc.roadLeft(line, startElem)
	}
}

func (doc *BoxesDocument) roadRight(line *ConnectionLine, startElem *LayoutElement) {
	if line.EndX >= doc.Width {
		return
	}
	if doc.checkColl(line.EndX, line.EndY, &doc.Boxes, startElem, nil) {
		// has collision
		line.EndX -= 2 * RasterSize
		return
	} else {
		line.EndX += RasterSize
		doc.roadRight(line, startElem)
	}
}

func (doc *BoxesDocument) initRoadsImpl(elem *LayoutElement) {
	if doc.ShouldHandle(elem) {
		stepSize := 2 * RasterSize
		// draw line from the top x start, till the first collision
		upRoad := newConnectionLine(*elem.TopXToStart, elem.Y, *elem.TopXToStart, elem.Y+stepSize)
		doc.roadUp(&upRoad, elem)
		doc.addRoad(upRoad, &doc.VerticalRoads)

		// draw line from the bottom x start, till the first collision
		downRoad := newConnectionLine(*elem.BottomXToStart, elem.Y+elem.Height,
			*elem.BottomXToStart, elem.Y+elem.Height+stepSize)
		doc.roadDown(&downRoad, elem)
		doc.addRoad(downRoad, &doc.VerticalRoads)

		// draw line from the left y start, till the first collision
		leftRoad := newConnectionLine(elem.X, *elem.LeftYToStart, elem.X-stepSize, *elem.LeftYToStart)
		doc.roadLeft(&leftRoad, elem)
		doc.addRoad(leftRoad, &doc.HorizontalRoads)

		// draw line from the right y start, till the first collision
		rightRoad := newConnectionLine(elem.X+elem.Width, *elem.RightYToStart, elem.X+elem.Width+stepSize, *elem.RightYToStart)
		doc.roadRight(&rightRoad, elem)
		doc.addRoad(rightRoad, &doc.HorizontalRoads)

		// draw the line parallel to the left border, till the first collision, up and down
		upRoad = newConnectionLine(elem.X-stepSize, elem.CenterY, elem.X-stepSize, elem.Y-stepSize)
		doc.roadUp(&upRoad, elem)
		downRoad = newConnectionLine(elem.X-stepSize, elem.CenterY, elem.X-stepSize, elem.Y+elem.Height+stepSize)
		doc.roadDown(&downRoad, elem)
		l := newConnectionLine(upRoad.EndX, upRoad.EndY, downRoad.EndX, downRoad.EndY)
		doc.addRoad(l, &doc.HorizontalRoads)

		// draw the line parallel to the right border, till the first collision, up and down
		upRoad = newConnectionLine(elem.X+elem.Width+stepSize, elem.CenterY, elem.X+elem.Width+stepSize, elem.Y-stepSize)
		doc.roadUp(&upRoad, elem)
		downRoad = newConnectionLine(elem.X+elem.Width+stepSize, elem.CenterY, elem.X+elem.Width+stepSize, elem.Y+elem.Height+stepSize)
		doc.roadDown(&downRoad, elem)
		l = newConnectionLine(upRoad.EndX, upRoad.EndY, downRoad.EndX, downRoad.EndY)
		doc.addRoad(l, &doc.HorizontalRoads)

		// draw the line parallel to the top border, till the first collision, left and right
		leftRoad = newConnectionLine(elem.CenterX, elem.Y-stepSize, elem.X-stepSize, elem.Y-stepSize)
		doc.roadLeft(&leftRoad, elem)
		rightRoad = newConnectionLine(elem.CenterX, elem.Y-stepSize, elem.X+elem.Width+stepSize, elem.Y-stepSize)
		doc.roadRight(&rightRoad, elem)
		l = newConnectionLine(leftRoad.EndX, leftRoad.EndY, rightRoad.EndX, rightRoad.EndY)
		doc.addRoad(l, &doc.HorizontalRoads)

		// draw the line parallel to the bottom border, till the first collision, left and right
		leftRoad = newConnectionLine(elem.CenterX, elem.Y+elem.Height+stepSize, elem.X-stepSize, elem.Y+elem.Height+stepSize)
		doc.roadLeft(&leftRoad, elem)
		rightRoad = newConnectionLine(elem.CenterX, elem.Y+elem.Height+stepSize, elem.X+elem.Width+stepSize, elem.Y+elem.Height+stepSize)
		doc.roadRight(&rightRoad, elem)
		l = newConnectionLine(leftRoad.EndX, leftRoad.EndY, rightRoad.EndX, rightRoad.EndY)
		doc.addRoad(l, &doc.HorizontalRoads)
	}
	if elem.Vertical != nil {
		for i := 0; i < len(elem.Vertical.Elems); i++ {
			doc.initRoadsImpl(&elem.Vertical.Elems[i])
		}
	}
	if elem.Horizontal != nil {
		for i := 0; i < len(elem.Horizontal.Elems); i++ {
			doc.initRoadsImpl(&elem.Horizontal.Elems[i])
		}
	}
}

// function goes over the horizontal roads, finds a related road and
// returns as first value the most left value of the road. The second value
// is the list of x-coordinates where the road has leftwards junctions up.
// The third value is the list of x-coordinates where the road has
// leftwards junctions down.
func (doc *BoxesDocument) checkRoadsToTheLeft(startX, startY int) (int, []int, []int, error) {
	for _, h := range doc.HorizontalRoads {
		leftX, rightX := minMax(h.StartX, h.EndX)
		leftX2, rightX2 := minMax(leftX, startX)
		if (h.StartY == startY) && (leftX <= startX) && (rightX >= startX) {
			// found the road
			upJunctions := make([]int, 0)
			downJunctions := make([]int, 0)
			retX := leftX
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
			SortDescending(upJunctions)
			SortAscending(downJunctions)
			return retX, upJunctions, downJunctions, nil
		}
	}
	return 0, nil, nil, fmt.Errorf("no road found to move left for %d,%d", startX, startY)
}

func SortDescending(a []int) {
	sort.Sort(sort.Reverse(sort.IntSlice(a)))
}

func SortAscending(a []int) {
	sort.IntSlice(a).Sort()
}

// function goes over the horizontal roads, finds a related road and
// returns as first value the most right value of the road. The second value
// is the list of x-coordinates where the road has rightwards junctions up.
// The third value is the list of x-coordinates where the road has
// rightwards junctions down.
func (doc *BoxesDocument) checkRoadsToTheRight(startX, startY int) (int, []int, []int, error) {
	for _, h := range doc.HorizontalRoads {
		leftX, rightX := minMax(h.StartX, h.EndX)
		leftX2, rightX2 := minMax(rightX, startX)
		if (h.StartY == startY) && (leftX <= startX) && (rightX >= startX) {
			// found the road
			upJunctions := make([]int, 0)
			downJunctions := make([]int, 0)
			retX := rightX
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
			SortDescending(upJunctions)
			SortAscending(downJunctions)
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
func (doc *BoxesDocument) checkRoadsToTheTop(startX, startY int) (int, []int, []int, error) {
	for _, v := range doc.VerticalRoads {
		topY, bottomY := minMax(v.StartY, v.EndY)
		topY2, bottomY2 := minMax(topY, startY)

		if (v.StartX == startX) && (topY <= startY) && (bottomY >= startY) {
			// found the road
			leftJunctions := make([]int, 0)
			rightJunctions := make([]int, 0)
			retY := topY
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
			SortDescending(leftJunctions)
			SortAscending(rightJunctions)
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
func (doc *BoxesDocument) checkRoadsToTheBottom(startX, startY int) (int, []int, []int, error) {
	for _, v := range doc.VerticalRoads {
		topY, bottomY := minMax(v.StartY, v.EndY)
		topY2, bottomY2 := minMax(bottomY, startY)
		if (v.StartX == startX) && (topY <= startY) && (bottomY >= startY) {
			// found the road
			leftJunctions := make([]int, 0)
			rightJunctions := make([]int, 0)
			retY := bottomY
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
			SortDescending(leftJunctions)
			SortAscending(rightJunctions)
			return retY, leftJunctions, rightJunctions, nil
		}
	}
	return 0, nil, nil, fmt.Errorf("no road found to move downwards for %d,%d", startX, startY)
}
