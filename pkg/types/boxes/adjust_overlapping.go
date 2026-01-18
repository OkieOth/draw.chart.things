package boxes

import "github.com/okieoth/draw.chart.things/pkg/types"

// This nasty function tries to resolve collisions of connections lines.
// The approach to solve the issue is, to move overlapping parts and strech
// the elements and the picture in cases.
func (doc *BoxesDocument) adjustForOverlappingConnections() {
	doc.adjustForOverlappingHorizontalLines()
	doc.adjustForOverlappingVerticalLines()
}

func (doc *BoxesDocument) adjustForOverlappingHorizontalLines() {
	xOffset := 0
	horizontalCount := len(doc.HorizontalLines)
	for i := range horizontalCount {
		if i == 0 {
			continue
		}
		last := &doc.HorizontalLines[i-1]
		current := &doc.HorizontalLines[i]
		if current.StartY == last.StartY &&
			((current.StartX >= last.StartX && current.StartX <= last.EndX) ||
				(current.EndX >= last.StartX && current.EndX <= last.EndX)) {
			current.StartY += types.LineDist
			current.EndY = current.StartY
			// TODO adjust in case start and end lines
			for j := i + 1; j < horizontalCount; j++ {
				cur2 := &doc.HorizontalLines[j]
				cur2.StartY += types.LineDist
				cur2.EndY = cur2.StartY
				// TODO adjust in case start and end lines
			}
			xOffset += types.LineDist
		}
	}
}

func (doc *BoxesDocument) adjustForOverlappingVerticalLines() {
	yOffset := 0
	verticalCount := len(doc.VerticalLines)
	for i := range verticalCount {
		if i == 0 {
			continue
		}
		last := &doc.VerticalLines[i-1]
		current := &doc.VerticalLines[i]
		if current.StartX == last.StartX &&
			((current.StartY >= last.StartY && current.StartY <= last.EndY) ||
				(current.EndY >= last.StartY && current.EndY <= last.EndY)) {
			current.StartX += types.LineDist
			current.EndX = current.StartX
			// TODO adjust in case start and end lines
			for j := i + 1; j < verticalCount; j++ {
				cur2 := &doc.VerticalLines[j]
				cur2.StartX += types.LineDist
				cur2.EndX = cur2.StartX
				// TODO adjust in case start and end lines
			}
			yOffset += types.LineDist
		}
	}
}
