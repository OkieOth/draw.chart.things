package types

import "fmt"

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
	doc.VerticalRoads = append(doc.VerticalRoads, line)
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
	if startElem.Id == "r5_1" {
		fmt.Println("Debug: r5_1")
	}
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
	for _, r := range doc.HorizontalRoads {
		leftX, rightX := minMax(r.StartX, r.EndX)
		if (r.StartY == startY) && (leftX <= startX) && (rightX >= startX) {
			// found the road
			// TODO
		}
	}
	return 0, nil, nil, fmt.Errorf("no road found to move left for %d,%d", startX, startY)
}

// function goes over the horizontal roads, finds a related road and
// returns as first value the most right value of the road. The second value
// is the list of x-coordinates where the road has rightwards junctions up.
// The third value is the list of x-coordinates where the road has
// rightwards junctions down.
func (doc *BoxesDocument) checkRoadsToTheRight(startX, startY int) (int, []int, []int, error) {
	for _, r := range doc.HorizontalRoads {
		leftX, rightX := minMax(r.StartX, r.EndX)
		if (r.StartY == startY) && (leftX <= startX) && (rightX >= startX) {
			// found the road
			// TODO
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
	for _, r := range doc.VerticalRoads {
		topY, bottomY := minMax(r.StartY, r.EndY)
		if (r.StartX == startX) && (topY <= startY) && (bottomY >= startY) {
			// found the road
			// TODO
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
	for _, r := range doc.VerticalRoads {
		topY, bottomY := minMax(r.StartY, r.EndY)
		if (r.StartX == startX) && (topY <= startY) && (bottomY >= startY) {
			// found the road
			// TODO
		}
	}
	return 0, nil, nil, fmt.Errorf("no road found to move downwards for %d,%d", startX, startY)
}
