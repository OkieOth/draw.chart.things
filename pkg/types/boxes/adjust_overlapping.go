package boxes

import "github.com/okieoth/draw.chart.things/pkg/types"

// This nasty function tries to resolve collisions of connections lines.
// The approach to solve the issue is, to move overlapping parts and strech
// the elements and the picture in cases.
func (doc *BoxesDocument) adjustForOverlappingConnections() {
	doc.adjustForOverlappingHorizontalLines()
	doc.adjustForOverlappingVerticalLines()
}

func (doc *BoxesDocument) addYAdjustment(layoutId string, line *ConnectionLine, endLineAdjustments *map[string]EndLineAdjustments) {
	alreadyStored, ok := (*endLineAdjustments)[layoutId]
	if !ok {
		alreadyStored = NewEndLineAdjustments()
	}
	line.EndY += types.LineDist
	if alreadyStored.maxLen < line.EndY {
		alreadyStored.maxLen = line.EndY
	}
	alreadyStored.maxLen += types.LineDist
	alreadyStored.lines = append(alreadyStored.lines, line)
	(*endLineAdjustments)[layoutId] = alreadyStored
}

func (doc *BoxesDocument) addXAdjustment(layoutId string, line *ConnectionLine, endLineAdjustments *map[string]EndLineAdjustments) {
	alreadyStored, ok := (*endLineAdjustments)[layoutId]
	if !ok {
		alreadyStored = NewEndLineAdjustments()
	}
	line.EndX += types.LineDist
	if alreadyStored.maxLen < line.EndX {
		alreadyStored.maxLen = line.EndX
	}
	alreadyStored.lines = append(alreadyStored.lines, line)
	(*endLineAdjustments)[layoutId] = alreadyStored
}

func (doc *BoxesDocument) fixVerticalStartAndEndOfHorizontalLine(xStart, xEnd, y, connectionIndex int, endLineAdjustments *map[string]EndLineAdjustments) {
	count := len(doc.VerticalLines)
	startFound := false
	endFound := false
	for i := range count {
		line := &doc.VerticalLines[i]
		if line.ConnectionIndex != connectionIndex {
			continue
		}
		if line.StartX == xStart || line.StartX == xEnd {
			// the vertical line cuts the triggering line
			if line.StartY == y {
				// vertical line goes down from the triggering line ... shorten the v-line
				line.StartY += types.LineDist
				line.EndY += types.LineDist
				if line.DestLayoutId != nil {
					doc.addYAdjustment(*line.DestLayoutId, line, endLineAdjustments)
				}
				if line.SrcLayoutId != nil {
					doc.addYAdjustment(*line.SrcLayoutId, line, endLineAdjustments)
				}

				startFound = true
				if startFound && endFound {
					return
				}
				continue
			}
			if line.EndY == y {
				// vertical line goes down to the triggering line ... increase the v-line
				line.EndY += types.LineDist
				//line.StartY += types.LineDist
				if line.DestLayoutId != nil {
					doc.addYAdjustment(*line.DestLayoutId, line, endLineAdjustments)
				}
				if line.SrcLayoutId != nil {
					doc.addYAdjustment(*line.SrcLayoutId, line, endLineAdjustments)
				}
				endFound = true
				if startFound && endFound {
					return
				}
				continue
			}
		}
	}
}

func (doc *BoxesDocument) fixHorizontalStartAndEndOfVerticalLine(offset, yStart, yEnd, x, connectionIndex int, endLineAdjustments *map[string]EndLineAdjustments) {
	count := len(doc.HorizontalLines)
	startFound := false
	endFound := false
	for i := range count {
		line := &doc.HorizontalLines[i]
		if line.ConnectionIndex != connectionIndex {
			continue
		}
		if line.StartY == yStart || line.StartY == yEnd {
			// the horizontal line cuts the triggering line
			if line.StartX == x {
				// horizontal line goes right from the triggering line
				line.StartX += offset
				line.EndX += offset
				startFound = true
				if line.DestLayoutId != nil {
					doc.addXAdjustment(*line.DestLayoutId, line, endLineAdjustments)
				}
				// seems not to be needed
				// if line.SrcLayoutId != nil {
				// 	doc.addXAdjustment(*line.SrcLayoutId, line, endLineAdjustments)
				// }
				if startFound && endFound {
					return
				}
				continue
			}
			if line.EndX == x {
				// horizontal line goes right to the triggering line
				line.EndX += offset
				if line.StartX >= x {
					line.StartX += offset
				}
				endFound = true
				if line.DestLayoutId != nil {
					doc.addXAdjustment(*line.DestLayoutId, line, endLineAdjustments)
				}
				if line.SrcLayoutId != nil {
					doc.addXAdjustment(*line.SrcLayoutId, line, endLineAdjustments)
				}
				if startFound && endFound {
					return
				}
				continue
			}
		}
	}
}

func (doc *BoxesDocument) adjustForOverlappingHorizontalLines() {
	xOffset := 0
	horizontalCount := len(doc.HorizontalLines)
	endLineAdjustments := make(map[string]EndLineAdjustments, 0)
	for i := range horizontalCount {
		if i == 0 {
			continue
		}
		last := &doc.HorizontalLines[i-1]
		current := &doc.HorizontalLines[i]
		if current.StartY == last.StartY &&
			overlapsHorizontally(current.StartX, current.EndX, last.StartX, last.EndX) {
			// adjust in case start and end lines ... but it's important before moving the h-line
			doc.fixVerticalStartAndEndOfHorizontalLine(current.StartX, current.EndX, current.StartY, current.ConnectionIndex, &endLineAdjustments)
			current.StartY += types.LineDist
			current.EndY = current.StartY
			for j := i + 1; j < horizontalCount; j++ {
				cur2 := &doc.HorizontalLines[j]
				// adjust in case start and end lines ... but it's important before moving the h-line
				doc.fixVerticalStartAndEndOfHorizontalLine(cur2.StartX, cur2.EndX, cur2.StartY, cur2.ConnectionIndex, &endLineAdjustments)
				cur2.StartY += types.LineDist
				cur2.EndY = cur2.StartY
			}
			doc.MoveBoxesHorizontal(current.StartY, types.LineDist)
			xOffset += types.LineDist
		}
	}
	// fix the endlines that have a different lenght ... as reason of the movements
	for _, v := range endLineAdjustments {
		for _, l := range v.lines {
			l.EndY = v.maxLen
		}
	}
}

func (doc *BoxesDocument) adjustForOverlappingVerticalLines() {
	yOffset := 0
	verticalCount := len(doc.VerticalLines)
	endLineAdjustments := make(map[string]EndLineAdjustments, 0)
	for i := range verticalCount {
		if i == 0 {
			continue
		}
		last := &doc.VerticalLines[i-1]
		current := &doc.VerticalLines[i]
		if current.StartX == last.StartX &&
			overlapsVertically(current.StartY, current.EndY, last.StartY, last.EndY) {
			doc.fixHorizontalStartAndEndOfVerticalLine(types.LineDist, current.StartY, current.EndY, current.StartX, current.ConnectionIndex, &endLineAdjustments)
			current.StartX += types.LineDist
			current.EndX = current.StartX
			for j := i + 1; j < verticalCount; j++ {
				cur2 := &doc.VerticalLines[j]
				doc.fixHorizontalStartAndEndOfVerticalLine(types.LineDist, cur2.StartY, cur2.EndY, cur2.StartX, cur2.ConnectionIndex, &endLineAdjustments)
				cur2.StartX += types.LineDist
				cur2.EndX = cur2.StartX
				// TODO adjust in case start and end lines
			}
			doc.MoveBoxesVertical(current.StartX, types.LineDist)
			yOffset += types.LineDist
		}
	}
	// fix the endlines that have a different lenght ... as reason of the movements
	for _, v := range endLineAdjustments {
		for _, l := range v.lines {
			l.EndX = v.maxLen
		}
	}
}

type EndLineAdjustments struct {
	lines  []*ConnectionLine
	maxLen int
}

func NewEndLineAdjustments() EndLineAdjustments {
	return EndLineAdjustments{
		lines:  make([]*ConnectionLine, 0),
		maxLen: 0,
	}
}

func overlapsVertically(yTop1, yBottom1, yTop2, yBottom2 int) bool {
	return ((yTop1 >= yTop2 && yTop1 <= yBottom2) || // yTop1 is in range of v-line 2
		(yBottom1 >= yTop2 && yBottom1 <= yBottom2)) || // ...yBottom1 is in range of v-line2
		((yTop2 >= yTop1 && yTop2 <= yBottom1) ||
			(yBottom2 >= yTop1 && yBottom2 <= yBottom1)) ||
		(yTop1 == yTop2) || (yTop1 == yBottom2) ||
		(yBottom1 == yTop2) || (yBottom1 == yBottom2)
}

func overlapsHorizontally(xLeft1, xRight1, xLeft2, xRight2 int) bool {
	return ((xLeft1 >= xLeft2 && xLeft1 <= xRight2) || // yTop1 is in range of v-line 2
		(xRight1 >= xLeft2 && xRight1 <= xRight2)) || // ...yBottom1 is in range of v-line2
		((xLeft2 >= xLeft1 && xLeft2 <= xRight1) ||
			(xRight2 >= xLeft1 && xRight2 <= xRight1)) ||
		(xLeft1 == xLeft2) || (xLeft1 == xRight2) ||
		(xRight1 == xLeft2) || (xRight1 == xRight2)
}
