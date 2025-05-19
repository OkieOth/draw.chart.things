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
			newLine = newConnectionLine(newX, newY, newX, newY-RasterSize)
		} else {
			newLine = newConnectionLine(newX, newY, newX, newY+RasterSize)
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
			newLine = newConnectionLine(newX, newY, newX-RasterSize, newY)
		} else {
			newLine = newConnectionLine(newX, newY, newX+RasterSize, newY)
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

func (doc *BoxesDocument) addNewLineToVariant(useCurrentVariant bool, newLine *ConnectionLine, variants *[][]ConnectionLine, variantIndex, lineIndex int) (int, int) {
	var newVariantIndex, newLineIndex int
	if useCurrentVariant {
		// work with the current variant
		(*variants)[variantIndex] = append((*variants)[variantIndex], *newLine)
		newVariantIndex = variantIndex
		newLineIndex = len((*variants)[variantIndex]) - 1
	} else {
		newVariant := (*variants)[variantIndex][:]
		newVariant = append(newVariant, *newLine)
		*variants = append(*variants, newVariant)
		newVariantIndex = len(*variants) - 1
		newLineIndex = len(newVariant) - 1
	}
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

	if (startX == endX) && (startY == endY) {
		return
	}
	x, straightAhead, upwards, downwards, mostLeftX, err := doc.getNextJunctionLeft(startX, startY)
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

	if upwards {
		newLine := newConnectionLine(x, currentLine.EndY, x, currentLine.EndY-RasterSize)
		newVariantIndex, newLineIndex := doc.addNewLineToVariant((!straightAhead) && (!downwards), &newLine, variants, currentVariantIndex, currentLineIndex)
		doc.goToUp(x, newLine.EndY, endX, endY, newVariantIndex, newLineIndex, variants, startElem, destElem)
	}
	if downwards {
		newLine := newConnectionLine(x, currentLine.EndY, x, currentLine.EndY+RasterSize)
		newVariantIndex, newLineIndex := doc.addNewLineToVariant((!straightAhead) && (!downwards), &newLine, variants, currentVariantIndex, currentLineIndex)
		doc.goToDown(x, newLine.EndY, endX, endY, newVariantIndex, newLineIndex, variants, startElem, destElem)
	}
	if straightAhead {
		possibleEnd := x - 2*RasterSize
		if (mostLeftX == possibleEnd) && (endX == possibleEnd) && (endY == startY) {
			// reached the end
			changedCurrentLine := newConnectionLine(currentLine.StartX, currentLine.StartY, possibleEnd, currentLine.StartY)
			(*variants)[currentVariantIndex][currentLineIndex] = changedCurrentLine
			return
		}
		if mostLeftX > possibleEnd {
			// new downwards variant
			doc.goToLeft(x, startY, endX, endY, currentVariantIndex, currentLineIndex, variants, startElem, destElem)
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
	if startX > doc.Width {
		// moved too far right
		fmt.Println("DEBUG: goToRight - emptyVariant-1", startX, startY, endX, endY, currentVariantIndex)
		emptyVariant(currentVariantIndex, variants)
		return
	}

	if (startX == endX) && (startY == endY) {
		return
	}
	x, straightAhead, upwards, downwards, mostRightX, err := doc.getNextJunctionRight(startX, startY)
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

	if upwards {
		newLine := newConnectionLine(x, currentLine.EndY, x, currentLine.EndY-RasterSize)
		newVariantIndex, newLineIndex := doc.addNewLineToVariant((!straightAhead) && (!downwards), &newLine, variants, currentVariantIndex, currentLineIndex)
		doc.goToUp(x, newLine.EndY, endX, endY, newVariantIndex, newLineIndex, variants, startElem, destElem)
	}
	if downwards {
		newLine := newConnectionLine(x, currentLine.EndY, x, currentLine.EndY+RasterSize)
		newVariantIndex, newLineIndex := doc.addNewLineToVariant((!straightAhead) && (!downwards), &newLine, variants, currentVariantIndex, currentLineIndex)
		doc.goToDown(x, newLine.EndY, endX, endY, newVariantIndex, newLineIndex, variants, startElem, destElem)
	}
	if straightAhead {
		possibleEnd := x + 2*RasterSize
		if (mostRightX == possibleEnd) && (endX == possibleEnd) && (endY == startY) {
			// reached the end
			changedCurrentLine := newConnectionLine(currentLine.StartX, currentLine.StartY, possibleEnd, currentLine.StartY)
			(*variants)[currentVariantIndex][currentLineIndex] = changedCurrentLine
			return
		}
		if mostRightX > possibleEnd {
			// new downwards variant
			doc.goToRight(x, startY, endX, endY, currentVariantIndex, currentLineIndex, variants, startElem, destElem)
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
	if startY > doc.Height {
		// moved too far down
		fmt.Println("DEBUG: goToDown - emptyVariant-1", startX, startY, endX, endY, currentVariantIndex)
		emptyVariant(currentVariantIndex, variants)
		return
	}

	if (startX == endX) && (startY == endY) {
		return
	}
	_, maxY := minMax(startY, endY)
	y, junctionPointsToLeft, junctionPointsToRight, err := doc.checkRoadsToTheBottom(startX, startY, maxY)
	if err != nil {
		fmt.Println("DEBUG: goToDown - emptyVariant-2", startX, startY, endX, endY, currentVariantIndex)
		emptyVariant(currentVariantIndex, variants)
		return
	}
	if maxY == endY {
		// up to down line
		(*variants)[currentVariantIndex][currentLineIndex].EndY = y
	} else {
		(*variants)[currentVariantIndex][currentLineIndex].StartY = y
	}

	if startX < endX {
		// to the right
		doc.verticalNewVariants(startX, startY, endX, endY, currentVariantIndex, currentLineIndex,
			variants, startElem, destElem, false, false, junctionPointsToRight, false)
		return
	} else if startY > endY {
		// to the left
		doc.verticalNewVariants(startX, startY, endX, endY, currentVariantIndex, currentLineIndex,
			variants, startElem, destElem, true, false, junctionPointsToLeft, false)
		return
	} else {
		// same level: startY == endY
		if endY == y {
			// reached the end
			l := (*variants)[currentVariantIndex][currentLineIndex]
			changedCurrentLine := newConnectionLine(l.StartX, l.StartY, l.StartX, y)
			(*variants)[currentVariantIndex][currentLineIndex] = changedCurrentLine
			return
		} else if y < endY {
			// trying to the right
			doc.verticalNewVariants(startX, startY, endX, endY, currentVariantIndex, currentLineIndex,
				variants, startElem, destElem, false, true, []int{y}, false)
			// trying to the left
			doc.verticalNewVariants(startX, startY, endX, endY, currentVariantIndex, currentLineIndex,
				variants, startElem, destElem, true, false, []int{y}, false)
		} else {
			// error
			emptyVariant(currentVariantIndex, variants)
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
	if startY < 0 {
		// moved too far up
		fmt.Println("DEBUG: goToUp - emptyVariant-1", startX, startY, endX, endY, currentVariantIndex)
		emptyVariant(currentVariantIndex, variants)
		return
	}

	if (startX == endX) && (startY == endY) {
		return
	}

	minY, _ := minMax(startY, endY)
	y, junctionPointsToLeft, junctionPointsToRight, err := doc.checkRoadsToTheTop(startX, startY, minY)
	if err != nil {
		fmt.Println("DEBUG: goToUp - emptyVariant-2", startX, startY, endX, endY, currentVariantIndex)
		emptyVariant(currentVariantIndex, variants)
		return
	}
	if minY == startY {
		// up to down line
		(*variants)[currentVariantIndex][currentLineIndex].StartY = y
	} else {
		(*variants)[currentVariantIndex][currentLineIndex].EndY = y
	}

	if startX < endX {
		// to the right
		doc.verticalNewVariants(startX, startY, endX, endY, currentVariantIndex, currentLineIndex,
			variants, startElem, destElem, false, false, junctionPointsToRight, true)
		return
	} else if startY > endY {
		// to the left
		doc.verticalNewVariants(startX, startY, endX, endY, currentVariantIndex, currentLineIndex,
			variants, startElem, destElem, true, false, junctionPointsToLeft, true)
		return
	} else {
		// same level: startY == endY
		if endY == y {
			// reached the end
			l := (*variants)[currentVariantIndex][currentLineIndex]
			changedCurrentLine := newConnectionLine(l.StartX, l.StartY, l.StartX, y)
			(*variants)[currentVariantIndex][currentLineIndex] = changedCurrentLine
			return
		} else if y > endY {
			// trying to the right
			doc.verticalNewVariants(startX, startY, endX, endY, currentVariantIndex, currentLineIndex,
				variants, startElem, destElem, false, true, []int{y}, true)
			// trying to the left
			doc.verticalNewVariants(startX, startY, endX, endY, currentVariantIndex, currentLineIndex,
				variants, startElem, destElem, true, false, []int{y}, true)
		} else {
			// error
			emptyVariant(currentVariantIndex, variants)
		}
	}
}
