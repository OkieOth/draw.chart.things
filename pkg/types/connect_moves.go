package types

import "fmt"

func (doc *BoxesDocument) horizontalNewVariants(
	startX, startY, endX, endY int,
	currentVariantIndex, currentLineIndex int,
	variants *[][]ConnectionLine,
	startElem, destElem *LayoutElement, up, alwaysNewVariant bool, junctionPoints []int, toLeft bool) {
	lastElem := len(junctionPoints) - 1
	for i, p := range junctionPoints {
		if p == startX {
			continue
		}
		if toLeft {
			if p < endX {
				continue
			}
		} else {
			if p > endX {
				continue
			}
		}
		newVariant := &(*variants)[currentVariantIndex]
		changedCurrentLine := newConnectionLine((*newVariant)[currentLineIndex].StartX, startY, p, startY)
		newX := changedCurrentLine.EndX
		newY := changedCurrentLine.EndY
		if i < lastElem || alwaysNewVariant {
			tmp := (*variants)[currentVariantIndex][:]
			newVariant = &tmp
		}
		(*newVariant)[currentLineIndex] = changedCurrentLine
		var newLine ConnectionLine
		if up {
			newLine = newConnectionLine(newX, newY, newX, newY-(2*RasterSize))
		} else {
			newLine = newConnectionLine(newX, newY, newX, newY+(2*RasterSize))
		}
		*newVariant = append(*newVariant, newLine)
		newVariantIndex := currentVariantIndex
		newLineIndex := len(*newVariant) - 1
		if i < lastElem || alwaysNewVariant {
			*variants = append(*variants, *newVariant)
			newVariantIndex = len(*variants) - 1
		}
		if up {
			doc.goToUp(newX, newY, endX, endY, newVariantIndex, newLineIndex, variants, startElem, destElem)
		} else {
			doc.goToDown(newX, newY, endX, endY, newVariantIndex, newLineIndex, variants, startElem, destElem)
		}
	}
}

func (doc *BoxesDocument) verticalNewVariants(
	startX, startY, endX, endY int,
	currentVariantIndex, currentLineIndex int,
	variants *[][]ConnectionLine,
	startElem, destElem *LayoutElement, left, alwaysNewVariant bool, junctionPoints []int, upwards bool) {
	lastElem := len(junctionPoints) - 1
	for i, p := range junctionPoints {
		if p == startY {
			continue
		}
		if upwards {
			if p < endY {
				continue
			}
		} else {
			if p > endY {
				continue
			}
		}

		newVariant := &(*variants)[currentVariantIndex]
		changedCurrentLine := newConnectionLine(startX, (*newVariant)[currentLineIndex].StartY, startX, p)
		newX := changedCurrentLine.EndX
		newY := changedCurrentLine.EndY
		if i < lastElem || alwaysNewVariant {
			tmp := (*variants)[currentVariantIndex][:]
			newVariant = &tmp
		}
		(*newVariant)[currentLineIndex] = changedCurrentLine
		var newLine ConnectionLine
		if left {
			newLine = newConnectionLine(newX, newY, newX-(2*RasterSize), newY)
		} else {
			newLine = newConnectionLine(newX, newY, newX+(2*RasterSize), newY)
		}
		*newVariant = append(*newVariant, newLine)
		newVariantIndex := currentVariantIndex
		newLineIndex := len(*newVariant) - 1
		if i < lastElem || alwaysNewVariant {
			*variants = append((*variants), *newVariant)
			newVariantIndex = len((*variants)) - 1
		}
		if left {
			doc.goToLeft(newX, newY, endX, endY, newVariantIndex, newLineIndex, variants, startElem, destElem)
		} else {
			doc.goToRight(newX, newY, endX, endY, newVariantIndex, newLineIndex, variants, startElem, destElem)
		}
	}
}

func (doc *BoxesDocument) addNewLineToVariant(newLine *ConnectionLine, baseVariant *[]ConnectionLine, variants *[][]ConnectionLine) (int, int) {
	var newVariantIndex, newLineIndex int
	newVariant := (*baseVariant)[:]
	newVariant = append(newVariant, *newLine)
	*variants = append(*variants, newVariant)
	newVariantIndex = len(*variants) - 1
	newLineIndex = len(newVariant) - 1
	return newVariantIndex, newLineIndex
}

func (doc *BoxesDocument) goToLeft(
	startX, startY, endX, endY int,
	currentVariantIndex, currentLineIndex int,
	variants *[][]ConnectionLine,
	startElem, destElem *LayoutElement) {
	fmt.Println("DEBUG: goToLeft", startX, startY, endX, endY, currentVariantIndex, currentLineIndex)
	if len((*variants)[currentVariantIndex]) == 0 && currentLineIndex > 0 {
		return
	}
	if startX < 0 {
		// moved too far left
		fmt.Println("DEBUG: goToLeft - emptyVariant-1", startX, startY, endX, endY, currentVariantIndex)
		emptyVariant(currentVariantIndex, variants)
		return
	}
	if startX < endX {
		// wrong direction
		fmt.Println("DEBUG: goToLeft - emptyVariant-3", startX, startY, endX, endY, currentVariantIndex)
		emptyVariant(currentVariantIndex, variants)
		return
	}

	if (startX == endX) && (startY == endY) {
		return
	}
	x, upwards, downwards, straightAhead, _, err := doc.getNextJunctionLeft(startX, startY)
	if err != nil {
		fmt.Println("DEBUG: goToLeft - emptyVariant-2", startX, startY, endX, endY, currentVariantIndex)
		emptyVariant(currentVariantIndex, variants)
		return
	}
	currentLine := (*variants)[currentVariantIndex][currentLineIndex]
	if currentLine.StartX == currentLine.EndX {
		// vertical line ... create a new line to x
		newLine := newConnectionLine(currentLine.EndX, currentLine.EndY, x, currentLine.EndY)
		(*variants)[currentVariantIndex] = append((*variants)[currentVariantIndex], newLine)
		currentLineIndex = len((*variants)[currentVariantIndex]) - 1
	} else {
		// horizontal line ... extend the current line to x
		(*variants)[currentVariantIndex][currentLineIndex].EndX = x
	}
	if ((*variants)[currentVariantIndex][currentLineIndex].EndX == endX) && ((*variants)[currentVariantIndex][currentLineIndex].EndY == endY) {
		// reached the end
		return
	}
	possibleEnd := x - 2*RasterSize
	if (endX == possibleEnd) && (endY == startY) {
		// reached the end
		changedCurrentLine := newConnectionLine(currentLine.StartX, currentLine.StartY, possibleEnd, currentLine.StartY)
		(*variants)[currentVariantIndex][currentLineIndex] = changedCurrentLine
		return
	}
	baseVariant := (*variants)[currentVariantIndex][:]
	if straightAhead {
		// not already visited this junction ... no circle
		doc.goToLeft(x, startY, endX, endY, currentVariantIndex, currentLineIndex, variants, startElem, destElem)
		if len((*variants)[currentVariantIndex]) == 0 && currentLineIndex > 0 {
			return
		}
		currentLine := (*variants)[currentVariantIndex][currentLineIndex]
		if currentLine.StartX == currentLine.EndX {
			return
		}
	}
	if upwards && currentLine.EndY > destElem.Y {
		newLineEndX, newLineEndY := x, currentLine.EndY-(2*RasterSize)
		newLine := newConnectionLine(x, currentLine.EndY, newLineEndX, newLineEndY)
		if !alreadyVisitedJunction(&newLine, currentVariantIndex, variants) {
			newVariantIndex, newLineIndex := doc.addNewLineToVariant(&newLine, &baseVariant, variants)
			doc.goToUp(x, newLine.EndY, endX, endY, newVariantIndex, newLineIndex, variants, startElem, destElem)
		}
	}
	if downwards && currentLine.EndY < destElem.Y {
		newLineEndX, newLineEndY := x, currentLine.EndY+(2*RasterSize)
		newLine := newConnectionLine(x, currentLine.EndY, newLineEndX, newLineEndY)
		if !alreadyVisitedJunction(&newLine, currentVariantIndex, variants) {
			newVariantIndex, newLineIndex := doc.addNewLineToVariant(&newLine, &baseVariant, variants)
			doc.goToDown(x, newLine.EndY, endX, endY, newVariantIndex, newLineIndex, variants, startElem, destElem)
		}
	}
}

func (doc *BoxesDocument) goToRight(
	startX, startY, endX, endY int,
	currentVariantIndex, currentLineIndex int,
	variants *[][]ConnectionLine,
	startElem, destElem *LayoutElement) {
	fmt.Println("DEBUG: goToRight", startX, startY, endX, endY, currentVariantIndex, currentLineIndex)
	if len((*variants)[currentVariantIndex]) == 0 && currentLineIndex > 0 {
		return
	}
	if startX > endX {
		// wrong direction
		fmt.Println("DEBUG: goToRight - emptyVariant-3", startX, startY, endX, endY, currentVariantIndex)
		emptyVariant(currentVariantIndex, variants)
		return
	}

	if startX > doc.Width {
		// moved too far right
		fmt.Println("DEBUG: goToRight - emptyVariant-1", startX, startY, endX, endY, currentVariantIndex)
		emptyVariant(currentVariantIndex, variants)
		return
	}

	if (startX == endX) && (startY == endY) {
		return
	}
	x, upwards, downwards, straightAhead, _, err := doc.getNextJunctionRight(startX, startY)
	if err != nil {
		fmt.Println("DEBUG: goToRight - emptyVariant-2", startX, startY, endX, endY, currentVariantIndex)
		emptyVariant(currentVariantIndex, variants)
		return
	}
	currentLine := (*variants)[currentVariantIndex][currentLineIndex]
	if currentLine.StartX == currentLine.EndX {
		// vertical line ... create a new line to x
		newLine := newConnectionLine(currentLine.EndX, currentLine.EndY, x, currentLine.EndY)
		(*variants)[currentVariantIndex] = append((*variants)[currentVariantIndex], newLine)
		currentLineIndex = len((*variants)[currentVariantIndex]) - 1
	} else {
		// horizontal line ... extend the current line to x
		(*variants)[currentVariantIndex][currentLineIndex].EndX = x
	}
	if ((*variants)[currentVariantIndex][currentLineIndex].EndX == endX) && ((*variants)[currentVariantIndex][currentLineIndex].EndY == endY) {
		// reached the end
		return
	}
	possibleEnd := x + 2*RasterSize
	if (endX == possibleEnd) && (endY == startY) {
		// reached the end
		changedCurrentLine := newConnectionLine(currentLine.StartX, currentLine.StartY, possibleEnd, currentLine.StartY)
		(*variants)[currentVariantIndex][currentLineIndex] = changedCurrentLine
		return
	}
	baseVariant := (*variants)[currentVariantIndex][:]
	if straightAhead {
		doc.goToRight(x, startY, endX, endY, currentVariantIndex, currentLineIndex, variants, startElem, destElem)
		if len((*variants)[currentVariantIndex]) == 0 && currentLineIndex > 0 {
			return
		}
		currentLine := (*variants)[currentVariantIndex][currentLineIndex]
		if currentLine.StartX == currentLine.EndX {
			return
		}
	}
	if upwards && currentLine.EndY > destElem.Y {
		newLineEndX, newLineEndY := x, currentLine.EndY-(2*RasterSize)
		// not already visited this junction ... no circle
		newLine := newConnectionLine(x, currentLine.EndY, newLineEndX, newLineEndY)
		if !alreadyVisitedJunction(&newLine, currentVariantIndex, variants) {
			newVariantIndex, newLineIndex := doc.addNewLineToVariant(&newLine, &baseVariant, variants)
			doc.goToUp(x, newLine.EndY, endX, endY, newVariantIndex, newLineIndex, variants, startElem, destElem)
		}
	}
	if downwards && currentLine.EndY < destElem.Y {
		newLineEndX, newLineEndY := x, currentLine.EndY+(2*RasterSize)
		newLine := newConnectionLine(x, currentLine.EndY, newLineEndX, newLineEndY)
		if !alreadyVisitedJunction(&newLine, currentVariantIndex, variants) {
			newVariantIndex, newLineIndex := doc.addNewLineToVariant(&newLine, &baseVariant, variants)
			doc.goToDown(x, newLine.EndY, endX, endY, newVariantIndex, newLineIndex, variants, startElem, destElem)
		}
	}
}

func (doc *BoxesDocument) goToDown(
	startX, startY, endX, endY int,
	currentVariantIndex, currentLineIndex int,
	variants *[][]ConnectionLine,
	startElem, destElem *LayoutElement) {
	fmt.Println("DEBUG: goToDown", startX, startY, endX, endY, currentVariantIndex, currentLineIndex)
	if len((*variants)[currentVariantIndex]) == 0 && currentLineIndex > 0 {
		return
	}
	if startY > endY {
		// wrong direction
		fmt.Println("DEBUG: goToUp - emptyVariant-3", startX, startY, endX, endY, currentVariantIndex)
		emptyVariant(currentVariantIndex, variants)
		return
	}
	if startY > doc.Height {
		// moved too far down
		fmt.Println("DEBUG: goToDown - emptyVariant-1", startX, startY, endX, endY, currentVariantIndex)
		emptyVariant(currentVariantIndex, variants)
		return
	}

	if (startX == endX) && (startY == endY) {
		return
	}

	y, leftwards, rightwards, straightAhead, _, err := doc.getNextJunctionDown(startX, startY)
	if err != nil {
		fmt.Println("DEBUG: goToDown - emptyVariant-2", startX, startY, endX, endY, currentVariantIndex)
		emptyVariant(currentVariantIndex, variants)
		return
	}
	currentLine := (*variants)[currentVariantIndex][currentLineIndex]
	if currentLine.StartX != currentLine.EndX {
		// horizontal line ... create a new line to y
		newLine := newConnectionLine(currentLine.EndX, currentLine.EndY, currentLine.EndX, y)
		(*variants)[currentVariantIndex] = append((*variants)[currentVariantIndex], newLine)
		currentLineIndex = len((*variants)[currentVariantIndex]) - 1
	} else {
		// vertical line ... extend the current line to y
		(*variants)[currentVariantIndex][currentLineIndex].EndY = y
	}
	if ((*variants)[currentVariantIndex][currentLineIndex].EndX == endX) && ((*variants)[currentVariantIndex][currentLineIndex].EndY == endY) {
		// reached the end
		return
	}
	possibleEnd := y + 2*RasterSize
	if (endY == possibleEnd) && (endX == startX) {
		// reached the end
		changedCurrentLine := newConnectionLine(currentLine.StartX, currentLine.StartY, currentLine.StartX, possibleEnd)
		(*variants)[currentVariantIndex][currentLineIndex] = changedCurrentLine
		return
	}
	baseVariant := (*variants)[currentVariantIndex][:]
	if straightAhead {
		doc.goToDown(startX, y, endX, endY, currentVariantIndex, currentLineIndex, variants, startElem, destElem)
		if len((*variants)[currentVariantIndex]) == 0 && currentLineIndex > 0 {
			return
		}
		currentLine := (*variants)[currentVariantIndex][currentLineIndex]
		if currentLine.StartX == currentLine.EndX {
			return
		}
	}
	if leftwards && currentLine.EndX > destElem.X {
		newLineEndX, newLineEndY := currentLine.EndX-(2*RasterSize), y
		newLine := newConnectionLine(currentLine.EndX, y, newLineEndX, newLineEndY)
		if !alreadyVisitedJunction(&newLine, currentVariantIndex, variants) {
			newVariantIndex, newLineIndex := doc.addNewLineToVariant(&newLine, &baseVariant, variants)
			doc.goToLeft(newLine.EndX, y, endX, endY, newVariantIndex, newLineIndex, variants, startElem, destElem)
		}
	}
	if rightwards && currentLine.EndX > destElem.X {
		newLineEndX, newLineEndY := currentLine.EndX+(2*RasterSize), y
		newLine := newConnectionLine(currentLine.EndX, y, newLineEndX, newLineEndY)
		if !alreadyVisitedJunction(&newLine, currentVariantIndex, variants) {
			newVariantIndex, newLineIndex := doc.addNewLineToVariant(&newLine, &baseVariant, variants)
			doc.goToRight(newLine.EndX, y, endX, endY, newVariantIndex, newLineIndex, variants, startElem, destElem)
		}
	}
}

func (doc *BoxesDocument) goToUp(
	startX, startY, endX, endY int,
	currentVariantIndex, currentLineIndex int,
	variants *[][]ConnectionLine,
	startElem, destElem *LayoutElement) {
	if endY == 192 && startY == 192 {
		fmt.Println("DEBUG: goToUp - emptyVariant-0000", startX, startY, endX, endY, currentVariantIndex)
	}
	fmt.Println("DEBUG: goToUp", startX, startY, endX, endY, currentVariantIndex, currentLineIndex)
	if len((*variants)[currentVariantIndex]) == 0 && currentLineIndex > 0 {
		return
	}
	if startY < endY {
		// wrong direction
		fmt.Println("DEBUG: goToUp - emptyVariant-3", startX, startY, endX, endY, currentVariantIndex)
		emptyVariant(currentVariantIndex, variants)
		return
	}

	if startY < 0 {
		// moved too far up
		fmt.Println("DEBUG: goToUp - emptyVariant-1", startX, startY, endX, endY, currentVariantIndex)
		emptyVariant(currentVariantIndex, variants)
		return
	}

	if (startX == endX) && (startY == endY) {
		return
	}

	y, leftwards, rightwards, straightAhead, _, err := doc.getNextJunctionUp(startX, startY)
	if err != nil {
		fmt.Println("DEBUG: goToUp - emptyVariant-2", startX, startY, endX, endY, currentVariantIndex)
		emptyVariant(currentVariantIndex, variants)
		return
	}
	currentLine := (*variants)[currentVariantIndex][currentLineIndex]
	if currentLine.StartX != currentLine.EndX {
		// horizontal line ... create a new line to y
		newLine := newConnectionLine(currentLine.EndX, currentLine.EndY, currentLine.EndX, y)
		(*variants)[currentVariantIndex] = append((*variants)[currentVariantIndex], newLine)
		currentLineIndex = len((*variants)[currentVariantIndex]) - 1
	} else {
		// vertical line ... extend the current line to y
		(*variants)[currentVariantIndex][currentLineIndex].EndY = y
	}
	if ((*variants)[currentVariantIndex][currentLineIndex].EndX == endX) && ((*variants)[currentVariantIndex][currentLineIndex].EndY == endY) {
		// reached the end
		return
	}
	possibleEnd := y - 2*RasterSize
	if (endY == possibleEnd) && (endX == startX) {
		// reached the end
		changedCurrentLine := newConnectionLine(currentLine.StartX, currentLine.StartY, currentLine.StartX, possibleEnd)
		(*variants)[currentVariantIndex][currentLineIndex] = changedCurrentLine
		return
	}
	baseVariant := (*variants)[currentVariantIndex][:]
	if straightAhead {
		doc.goToUp(startX, y, endX, endY, currentVariantIndex, currentLineIndex, variants, startElem, destElem)
		if straightAhead {
			doc.goToDown(startX, y, endX, endY, currentVariantIndex, currentLineIndex, variants, startElem, destElem)
			if len((*variants)[currentVariantIndex]) == 0 && currentLineIndex > 0 {
				return
			}
			currentLine := (*variants)[currentVariantIndex][currentLineIndex]
			if currentLine.StartX == currentLine.EndX {
				return
			}
		}
	}
	if leftwards && currentLine.EndX > destElem.X {
		newLineEndX, newLineEndY := currentLine.EndX-(2*RasterSize), y
		newLine := newConnectionLine(currentLine.EndX, y, newLineEndX, newLineEndY)
		if !alreadyVisitedJunction(&newLine, currentVariantIndex, variants) {
			newVariantIndex, newLineIndex := doc.addNewLineToVariant(&newLine, &baseVariant, variants)
			doc.goToLeft(newLine.EndX, y, endX, endY, newVariantIndex, newLineIndex, variants, startElem, destElem)
		}
	}
	if rightwards && currentLine.EndX < destElem.X {
		newLineEndX, newLineEndY := currentLine.EndX+(2*RasterSize), y
		newLine := newConnectionLine(currentLine.EndX, y, newLineEndX, newLineEndY)
		if !alreadyVisitedJunction(&newLine, currentVariantIndex, variants) {
			newVariantIndex, newLineIndex := doc.addNewLineToVariant(&newLine, &baseVariant, variants)
			doc.goToRight(newLine.EndX, y, endX, endY, newVariantIndex, newLineIndex, variants, startElem, destElem)
		}
	}
}

func alreadyVisitedJunction(newLine *ConnectionLine, currentVariantIndex int, variants *[][]ConnectionLine) bool {
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
