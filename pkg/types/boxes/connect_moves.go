package boxes

import (
	"fmt"
	"slices"

	"github.com/okieoth/draw.chart.things/pkg/types"
)

type CircleTrace struct {
	StartX int
	StartY int
	EndX   int
	EndY   int
}

// func connectionExistsByDestId(connections []boxes.Connection, destId string) bool {
// 	return slices.ContainsFunc(connections, func(c boxes.Connection) bool {
// 		return c.DestId == destId
// 	})
// }

func foundXBasedCircle(circleTrace *[]CircleTrace, start, end int) bool {
	if ret := slices.ContainsFunc(*circleTrace, func(ct CircleTrace) bool {
		return ct.StartX == start && ct.EndX == end
	}); ret {
		return true
	}
	addXToCircleTrace(circleTrace, start, end)
	return false
}

func foundYBasedCircle(circleTrace *[]CircleTrace, start, end int) bool {
	if ret := slices.ContainsFunc(*circleTrace, func(ct CircleTrace) bool {
		return ct.StartY == start && ct.EndY == end
	}); ret {
		return true
	}
	addYToCircleTrace(circleTrace, start, end)
	return false
}

func addXToCircleTrace(circleTrace *[]CircleTrace, start, end int) {
	*circleTrace = append(*circleTrace, CircleTrace{
		StartX: start,
		EndX:   end,
	})
}

func addYToCircleTrace(circleTrace *[]CircleTrace, start, end int) {
	*circleTrace = append(*circleTrace, CircleTrace{
		StartY: start,
		EndY:   end,
	})
}

func (doc *BoxesDocument) goToLeft(
	endX, endY int,
	variant []ConnectionLine,
	startElem, destElem *LayoutElement, circleTrace *[]CircleTrace) ([][]ConnectionLine, error) {
	startX := variant[len(variant)-1].EndX
	startY := variant[len(variant)-1].EndY
	if foundXBasedCircle(circleTrace, startX, endX) {
		return nil, fmt.Errorf("goToLeft - reached circle %d %d %d %d", startX, startY, endX, endY)
	}
	if startX < 0 {
		// moved too far left
		return nil, fmt.Errorf("goToLeft - too far left %d %d %d %d", startX, startY, endX, endY)
	}

	if (startX == endX) && (startY == endY) {
		ret := make([][]ConnectionLine, 1)
		ret[0] = variant
		return ret, nil
	}
	x, upwards, downwards, straightAhead, _, err := doc.getNextJunctionLeft(startX, startY)
	if err != nil {
		return nil, err
	}
	currentLineIndex := len(variant) - 1
	currentLine := variant[currentLineIndex]
	if currentLine.StartX == currentLine.EndX {
		// vertical line ... create a new line to x
		newLine := newConnectionLine(currentLine.EndX, currentLine.EndY, x, currentLine.EndY)
		variant = append(variant, newLine)
		currentLineIndex = len(variant) - 1
	} else {
		// horizontal line ... extend the current line to x
		variant[currentLineIndex].EndX = x
	}
	if (variant[currentLineIndex].EndX == endX) && (variant[currentLineIndex].EndY == endY) {
		ret := make([][]ConnectionLine, 1)
		ret[0] = variant
		return ret, nil
	}
	possibleEnd := x - 2*types.RasterSize
	if (endX == possibleEnd) && (endY == startY) {
		// reached the end
		changedCurrentLine := newConnectionLine(currentLine.StartX, currentLine.StartY, possibleEnd, currentLine.StartY)
		variant[currentLineIndex] = changedCurrentLine
		ret := make([][]ConnectionLine, 1)
		ret[0] = variant
		return ret, nil
	}
	ret := make([][]ConnectionLine, 1)
	ret[0] = variant
	if straightAhead {
		// going straight ahead
		newVariant := make([]ConnectionLine, len(variant))
		copy(newVariant, variant)
		newVariant[currentLineIndex].EndX = x - types.RasterSize
		newVariants, err := doc.goToLeft(endX, endY, newVariant, startElem, destElem, circleTrace)
		if err == nil && len(newVariants) > 0 {
			ret = append(ret, newVariants...)
		}
	}
	if upwards {
		newLineEndX, newLineEndY := x, currentLine.EndY-(2*types.RasterSize)
		// not already visited this junction ... no circle
		newLine := newConnectionLine(x, currentLine.EndY, newLineEndX, newLineEndY)
		if !alreadyVisitedJunction(&newLine, variant) {
			newVariant := make([]ConnectionLine, len(variant))
			copy(newVariant, variant)
			newVariant = append(newVariant, newLine)
			newVariants, err := doc.goToUp(endX, endY, newVariant, startElem, destElem, circleTrace)
			if err == nil && len(newVariants) > 0 {
				ret = append(ret, newVariants...)
			}
		}
	}
	if downwards {
		newLineEndX, newLineEndY := x, currentLine.EndY+(2*types.RasterSize)
		newLine := newConnectionLine(x, currentLine.EndY, newLineEndX, newLineEndY)
		if !alreadyVisitedJunction(&newLine, variant) {
			newVariant := make([]ConnectionLine, len(variant))
			copy(newVariant, variant)
			newVariant = append(newVariant, newLine)
			newVariants, err := doc.goToDown(endX, endY, newVariant, startElem, destElem, circleTrace)
			if err == nil && len(newVariants) > 0 {
				ret = append(ret, newVariants...)
			}
		}
	}
	if len(ret) == 0 {
		return nil, fmt.Errorf("goToLeft - no new variants %d %d %d %d", startX, startY, endX, endY)
	}
	return ret, nil
}

func (doc *BoxesDocument) goToRight(
	endX, endY int,
	variant []ConnectionLine,
	startElem, destElem *LayoutElement, circleTrace *[]CircleTrace) ([][]ConnectionLine, error) {
	startX := variant[len(variant)-1].EndX
	startY := variant[len(variant)-1].EndY

	if foundXBasedCircle(circleTrace, startX, endX) {
		return nil, fmt.Errorf("goToRight - reached circle %d %d %d %d", startX, startY, endX, endY)
	}

	if startX >= doc.Width {
		// moved too far right
		return nil, fmt.Errorf("goToRight - too far right %d %d %d %d", startX, startY, endX, endY)
	}

	if (startX == endX) && (startY == endY) {
		ret := make([][]ConnectionLine, 1)
		ret[0] = variant
		return ret, nil
	}
	x, upwards, downwards, straightAhead, _, err := doc.getNextJunctionRight(startX, startY)
	if err != nil {
		return nil, err
	}
	currentLineIndex := len(variant) - 1
	currentLine := variant[currentLineIndex]
	if currentLine.StartX == currentLine.EndX {
		// vertical line ... create a new line to x
		newLine := newConnectionLine(currentLine.EndX, currentLine.EndY, x, currentLine.EndY)
		variant = append(variant, newLine)
		currentLineIndex = len(variant) - 1
	} else {
		// horizontal line ... extend the current line to x
		variant[currentLineIndex].EndX = x
	}
	if (variant[currentLineIndex].EndX == endX) && (variant[currentLineIndex].EndY == endY) {
		ret := make([][]ConnectionLine, 1)
		ret[0] = variant
		return ret, nil
	}
	possibleEnd := x + 2*types.RasterSize
	if (endX == possibleEnd) && (endY == startY) {
		// reached the end
		changedCurrentLine := newConnectionLine(currentLine.StartX, currentLine.StartY, possibleEnd, currentLine.StartY)
		variant[currentLineIndex] = changedCurrentLine
		ret := make([][]ConnectionLine, 1)
		ret[0] = variant
		return ret, nil
	}
	ret := make([][]ConnectionLine, 1)
	ret[0] = variant
	if straightAhead {
		// going straight ahead
		newVariant := make([]ConnectionLine, len(variant))
		copy(newVariant, variant)
		newVariant[currentLineIndex].EndX = x + types.RasterSize
		newVariants, err := doc.goToRight(endX, endY, newVariant, startElem, destElem, circleTrace)
		if err == nil && len(newVariants) > 0 {
			ret = append(ret, newVariants...)
		}
	}
	if upwards {
		newLineEndX, newLineEndY := x, currentLine.EndY-(2*types.RasterSize)
		// not already visited this junction ... no circle
		newLine := newConnectionLine(x, currentLine.EndY, newLineEndX, newLineEndY)
		if !alreadyVisitedJunction(&newLine, variant) {
			newVariant := make([]ConnectionLine, len(variant))
			copy(newVariant, variant)
			newVariant = append(newVariant, newLine)
			newVariants, err := doc.goToUp(endX, endY, newVariant, startElem, destElem, circleTrace)
			if err == nil && len(newVariants) > 0 {
				ret = append(ret, newVariants...)
			}
		}
	}
	if downwards {
		newLineEndX, newLineEndY := x, currentLine.EndY+(2*types.RasterSize)
		newLine := newConnectionLine(x, currentLine.EndY, newLineEndX, newLineEndY)
		if !alreadyVisitedJunction(&newLine, variant) {
			newVariant := make([]ConnectionLine, len(variant))
			copy(newVariant, variant)
			newVariant = append(newVariant, newLine)
			newVariants, err := doc.goToDown(endX, endY, newVariant, startElem, destElem, circleTrace)
			if err == nil && len(newVariants) > 0 {
				ret = append(ret, newVariants...)
			}
		}
	}
	if len(ret) == 0 {
		return nil, fmt.Errorf("goToRight - no new variants %d %d %d %d", startX, startY, endX, endY)
	}
	return ret, nil
}

func (doc *BoxesDocument) goToDown(
	endX, endY int,
	variant []ConnectionLine,
	startElem, destElem *LayoutElement, circleTrace *[]CircleTrace) ([][]ConnectionLine, error) {
	startX := variant[len(variant)-1].EndX
	startY := variant[len(variant)-1].EndY

	if foundYBasedCircle(circleTrace, startY, endY) {
		return nil, fmt.Errorf("goToDown - reached circle %d %d %d %d", startX, startY, endX, endY)
	}

	if startY >= doc.Height {
		// moved too far down
		return nil, fmt.Errorf("goToDown - too far down %d %d %d %d", startX, startY, endX, endY)
	}

	if (startX == endX) && (startY == endY) {
		ret := make([][]ConnectionLine, 1)
		ret[0] = variant
		return ret, nil
	}

	y, leftwards, rightwards, straightAhead, _, err := doc.getNextJunctionDown(startX, startY)
	if err != nil {
		return nil, err
	}
	currentLineIndex := len(variant) - 1
	currentLine := variant[currentLineIndex]
	if currentLine.StartX != currentLine.EndX {
		// horizontal line ... create a new line to y
		newLine := newConnectionLine(currentLine.EndX, currentLine.EndY, currentLine.EndX, y)
		variant = append(variant, newLine)
		currentLineIndex = len(variant) - 1
	} else {
		// vertical line ... extend the current line to y
		variant[currentLineIndex].EndY = y
	}
	if (variant[currentLineIndex].EndX == endX) && (variant[currentLineIndex].EndY == endY) {
		ret := make([][]ConnectionLine, 1)
		ret[0] = variant
		return ret, nil
	}
	possibleEnd := y + 2*types.RasterSize
	if (endY == possibleEnd) && (endX == startX) {
		// reached the end
		changedCurrentLine := newConnectionLine(currentLine.StartX, currentLine.StartY, currentLine.StartX, possibleEnd)
		variant[currentLineIndex] = changedCurrentLine
		ret := make([][]ConnectionLine, 1)
		ret[0] = variant
		return ret, nil
	}
	ret := make([][]ConnectionLine, 1)
	ret[0] = variant
	if straightAhead {
		// going straight ahead
		newVariant := make([]ConnectionLine, len(variant))
		copy(newVariant, variant)
		newVariant[currentLineIndex].EndY = y + types.RasterSize
		newVariants, err := doc.goToDown(endX, endY, newVariant, startElem, destElem, circleTrace)
		if err == nil && len(newVariants) > 0 {
			ret = append(ret, newVariants...)
		}
	}
	if leftwards {
		newLineEndX, newLineEndY := currentLine.EndX-(2*types.RasterSize), y
		newLine := newConnectionLine(currentLine.EndX, y, newLineEndX, newLineEndY)
		if !alreadyVisitedJunction(&newLine, variant) {
			newVariant := make([]ConnectionLine, len(variant))
			copy(newVariant, variant)
			newVariant = append(newVariant, newLine)
			newVariants, err := doc.goToLeft(endX, endY, newVariant, startElem, destElem, circleTrace)
			if err == nil && len(newVariants) > 0 {
				ret = append(ret, newVariants...)
			}
		}
	}
	if rightwards {
		newLineEndX, newLineEndY := currentLine.EndX+(2*types.RasterSize), y
		newLine := newConnectionLine(currentLine.EndX, y, newLineEndX, newLineEndY)
		if !alreadyVisitedJunction(&newLine, variant) {
			newVariant := make([]ConnectionLine, len(variant))
			copy(newVariant, variant)
			newVariant = append(newVariant, newLine)
			newVariants, err := doc.goToRight(endX, endY, newVariant, startElem, destElem, circleTrace)
			if err == nil && len(newVariants) > 0 {
				ret = append(ret, newVariants...)
			}
		}
	}
	if len(ret) == 0 {
		return nil, fmt.Errorf("goToDown - no new variants %d %d %d %d", startX, startY, endX, endY)
	}
	return ret, nil
}

func (doc *BoxesDocument) goToUp(
	endX, endY int,
	variant []ConnectionLine,
	startElem, destElem *LayoutElement, circleTrace *[]CircleTrace) ([][]ConnectionLine, error) {
	startX := variant[len(variant)-1].EndX
	startY := variant[len(variant)-1].EndY

	if foundYBasedCircle(circleTrace, startY, endY) {
		return nil, fmt.Errorf("goToUp - reached circle %d %d %d %d", startX, startY, endX, endY)
	}

	if startY < 0 {
		// moved too far up
		return nil, fmt.Errorf("goToUp - too far up %d %d %d %d", startX, startY, endX, endY)
	}

	if (startX == endX) && (startY == endY) {
		ret := make([][]ConnectionLine, 1)
		ret[0] = variant
		return ret, nil
	}

	y, leftwards, rightwards, straightAhead, _, err := doc.getNextJunctionUp(startX, startY)
	if err != nil {
		return nil, err
	}
	currentLineIndex := len(variant) - 1
	currentLine := variant[currentLineIndex]
	if currentLine.StartX != currentLine.EndX {
		// horizontal line ... create a new line to y
		newLine := newConnectionLine(currentLine.EndX, currentLine.EndY, currentLine.EndX, y)
		variant = append(variant, newLine)
		currentLineIndex = len(variant) - 1
	} else {
		// vertical line ... extend the current line to y
		variant[currentLineIndex].EndY = y
	}
	if (variant[currentLineIndex].EndX == endX) && (variant[currentLineIndex].EndY == endY) {
		ret := make([][]ConnectionLine, 1)
		ret[0] = variant
		return ret, nil
	}
	possibleEnd := y - 2*types.RasterSize
	if (endY == possibleEnd) && (endX == startX) {
		// reached the end
		changedCurrentLine := newConnectionLine(currentLine.StartX, currentLine.StartY, currentLine.StartX, possibleEnd)
		variant[currentLineIndex] = changedCurrentLine
		ret := make([][]ConnectionLine, 1)
		ret[0] = variant
		return ret, nil
	}

	ret := make([][]ConnectionLine, 1)
	ret[0] = variant
	if straightAhead {
		// going straight ahead
		newVariant := make([]ConnectionLine, len(variant))
		copy(newVariant, variant)
		newVariant[currentLineIndex].EndY = y - types.RasterSize
		newVariants, err := doc.goToUp(endX, endY, newVariant, startElem, destElem, circleTrace)
		if err == nil && len(newVariants) > 0 {
			ret = append(ret, newVariants...)
		}
	}
	if leftwards {
		newLineEndX, newLineEndY := currentLine.EndX-(2*types.RasterSize), y
		newLine := newConnectionLine(currentLine.EndX, y, newLineEndX, newLineEndY)
		if !alreadyVisitedJunction(&newLine, variant) {
			newVariant := make([]ConnectionLine, len(variant))
			copy(newVariant, variant)
			newVariant = append(newVariant, newLine)
			newVariants, err := doc.goToLeft(endX, endY, newVariant, startElem, destElem, circleTrace)
			if err == nil && len(newVariants) > 0 {
				ret = append(ret, newVariants...)
			}
		}
	}
	if rightwards {
		newLineEndX, newLineEndY := currentLine.EndX+(2*types.RasterSize), y
		newLine := newConnectionLine(currentLine.EndX, y, newLineEndX, newLineEndY)
		if !alreadyVisitedJunction(&newLine, variant) {
			newVariant := make([]ConnectionLine, len(variant))
			copy(newVariant, variant)
			newVariant = append(newVariant, newLine)
			newVariants, err := doc.goToRight(endX, endY, newVariant, startElem, destElem, circleTrace)
			if err == nil && len(newVariants) > 0 {
				ret = append(ret, newVariants...)
			}
		}
	}
	if len(ret) == 0 {
		return nil, fmt.Errorf("goToUp - no new variants %d %d %d %d", startX, startY, endX, endY)
	}
	return ret, nil
}

func alreadyVisitedJunction(newLine *ConnectionLine, variant []ConnectionLine) bool {
	for _, l := range variant {
		if l.StartX == newLine.StartX &&
			l.StartY == newLine.StartY &&
			l.EndX == newLine.EndX &&
			l.EndY == newLine.EndY {
			return true
		}
	}
	return false
}

func alreadyVisitedJunction2(newLine *ConnectionLine, currentVariantIndex int, variants *[][]ConnectionLine) bool {
	for _, l := range (*variants)[currentVariantIndex] {
		if l.StartX == newLine.StartX &&
			l.StartY == newLine.StartY &&
			l.EndX == newLine.EndX &&
			l.EndY == newLine.EndY {
			return true
		}
	}
	return false
}
