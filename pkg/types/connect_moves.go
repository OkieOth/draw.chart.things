package types

import "fmt"

func (doc *BoxesDocument) horizontalNewVariants(
	startX, startY, endX, endY int,
	currentVariantIndex, currentLineIndex int,
	variants [][]ConnectionLine,
	startElem, destElem *LayoutElement, up, alwaysNewVariant bool, junctionPoints []int, toLeft bool) [][]ConnectionLine {
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
		newVariant := &variants[currentVariantIndex]
		changedCurrentLine := newConnectionLine((*newVariant)[currentLineIndex].StartX, startY, p, startY)
		newX := changedCurrentLine.EndX
		newY := changedCurrentLine.EndY
		if i < lastElem || alwaysNewVariant {
			tmp := variants[currentVariantIndex][:]
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
			variants = append(variants, *newVariant)
			newVariantIndex = len(variants) - 1
		}
		if up {
			variants = doc.goToUp(newX, newY, endX, endY, newVariantIndex, newLineIndex, variants, startElem, destElem)
		} else {
			variants = doc.goToDown(newX, newY, endX, endY, newVariantIndex, newLineIndex, variants, startElem, destElem)
		}
	}
	return variants
}

func (doc *BoxesDocument) verticalNewVariants(
	startX, startY, endX, endY int,
	currentVariantIndex, currentLineIndex int,
	variants [][]ConnectionLine,
	startElem, destElem *LayoutElement, left, alwaysNewVariant bool, junctionPoints []int, upwards bool) [][]ConnectionLine {
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

		newVariant := &variants[currentVariantIndex]
		changedCurrentLine := newConnectionLine(startX, (*newVariant)[currentLineIndex].StartY, startX, p)
		newX := changedCurrentLine.EndX
		newY := changedCurrentLine.EndY
		if i < lastElem || alwaysNewVariant {
			tmp := variants[currentVariantIndex][:]
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
			variants = append(variants, *newVariant)
			newVariantIndex = len(variants) - 1
		}
		if left {
			variants = doc.goToLeft(newX, newY, endX, endY, newVariantIndex, newLineIndex, variants, startElem, destElem)
		} else {
			variants = doc.goToRight(newX, newY, endX, endY, newVariantIndex, newLineIndex, variants, startElem, destElem)
		}
	}
	return variants
}

func (doc *BoxesDocument) goToLeft(
	startX, startY, endX, endY int,
	currentVariantIndex, currentLineIndex int,
	variants [][]ConnectionLine,
	startElem, destElem *LayoutElement) [][]ConnectionLine {
	fmt.Println("DEBUG: goToLeft", startX, startY, endX, endY, currentVariantIndex, currentLineIndex)
	if len(variants[currentVariantIndex]) == 0 && currentLineIndex > 0 {
		return variants
	}
	if startX < endX {
		// moved too far left
		fmt.Println("DEBUG: goToLeft - emptyVariant-1", startX, startY, endX, endY, currentVariantIndex)
		return emptyVariant(currentVariantIndex, variants)
	}

	if (startX == endX) && (startY == endY) {
		return variants
	}

	x, junctionPointsUpwards, junctionPointsDownwards, err := doc.checkRoadsToTheLeft(startX, startY)
	if err != nil {
		fmt.Println("DEBUG: goToLeft - emptyVariant-2", startX, startY, endX, endY, currentVariantIndex)
		return emptyVariant(currentVariantIndex, variants)
	}
	if startY < endY {
		// downwards
		return doc.horizontalNewVariants(startX, startY, endX, endY, currentVariantIndex, currentLineIndex,
			variants, startElem, destElem, false, false, junctionPointsDownwards, true)
	} else if startY > endY {
		// upwards
		return doc.horizontalNewVariants(startX, startY, endX, endY, currentVariantIndex, currentLineIndex,
			variants, startElem, destElem, true, false, junctionPointsUpwards, true)
	} else {
		// same level: startY == endY
		if endX == x {
			// reached the end
			l := variants[currentVariantIndex][currentLineIndex]
			changedCurrentLine := newConnectionLine(l.StartX, l.StartY, x, l.StartY)
			variants[currentVariantIndex][currentLineIndex] = changedCurrentLine
			return variants
		} else if x > endX {
			// trying downwards
			variants := doc.horizontalNewVariants(startX, startY, endX, endY, currentVariantIndex, currentLineIndex,
				variants, startElem, destElem, false, true, []int{x}, true)
			return doc.horizontalNewVariants(startX, startY, endX, endY, currentVariantIndex, currentLineIndex,
				variants, startElem, destElem, true, false, []int{x}, true)
		} else {
			// error
			return emptyVariant(currentVariantIndex, variants)
		}
	}
}

func (doc *BoxesDocument) goToRight(
	startX, startY, endX, endY int,
	currentVariantIndex, currentLineIndex int,
	variants [][]ConnectionLine,
	startElem, destElem *LayoutElement) [][]ConnectionLine {
	fmt.Println("DEBUG: goToRight", startX, startY, endX, endY, currentVariantIndex, currentLineIndex)
	if len(variants[currentVariantIndex]) == 0 && currentLineIndex > 0 {
		return variants
	}
	if startX > endX {
		// moved too far right
		fmt.Println("DEBUG: goToRight - emptyVariant-1", startX, startY, endX, endY, currentVariantIndex)
		return emptyVariant(currentVariantIndex, variants)
	}

	if (startX == endX) && (startY == endY) {
		return variants
	}

	x, junctionPointsUpwards, junctionPointsDownwards, err := doc.checkRoadsToTheRight(startX, startY)
	if err != nil {
		fmt.Println("DEBUG: goToRight - emptyVariant-2", startX, startY, endX, endY, currentVariantIndex)
		return emptyVariant(currentVariantIndex, variants)
	}
	if startY < endY {
		// downwards
		return doc.horizontalNewVariants(startX, startY, endX, endY, currentVariantIndex, currentLineIndex,
			variants, startElem, destElem, false, false, junctionPointsDownwards, false)
	} else if startY > endY {
		// upwards
		return doc.horizontalNewVariants(startX, startY, endX, endY, currentVariantIndex, currentLineIndex,
			variants, startElem, destElem, true, false, junctionPointsUpwards, false)
	} else {
		// same level: startY == endY
		if endX == x {
			// reached the end
			l := variants[currentVariantIndex][currentLineIndex]
			changedCurrentLine := newConnectionLine(l.StartX, l.StartY, x, l.StartY)
			variants[currentVariantIndex][currentLineIndex] = changedCurrentLine
			return variants
		} else if x < endX {
			// trying downwards
			variants = doc.horizontalNewVariants(startX, startY, endX, endY, currentVariantIndex, currentLineIndex,
				variants, startElem, destElem, false, true, []int{x}, false)
			// trying upwards
			return doc.horizontalNewVariants(startX, startY, endX, endY, currentVariantIndex, currentLineIndex,
				variants, startElem, destElem, true, false, []int{x}, false)
		} else {
			// error
			return emptyVariant(currentVariantIndex, variants)
		}
	}
}

func (doc *BoxesDocument) goToDown(
	startX, startY, endX, endY int,
	currentVariantIndex, currentLineIndex int,
	variants [][]ConnectionLine,
	startElem, destElem *LayoutElement) [][]ConnectionLine {
	fmt.Println("DEBUG: goToDown", startX, startY, endX, endY, currentVariantIndex, currentLineIndex)
	if len(variants[currentVariantIndex]) == 0 && currentLineIndex > 0 {
		return variants
	}
	if startY > endY {
		// moved too far down
		fmt.Println("DEBUG: goToDown - emptyVariant-1", startX, startY, endX, endY, currentVariantIndex)
		return emptyVariant(currentVariantIndex, variants)
	}

	if (startX == endX) && (startY == endY) {
		return variants
	}

	y, junctionPointsToLeft, junctionPointsToRight, err := doc.checkRoadsToTheBottom(startX, startY)
	if err != nil {
		fmt.Println("DEBUG: goToDown - emptyVariant-2", startX, startY, endX, endY, currentVariantIndex)
		return emptyVariant(currentVariantIndex, variants)
	}
	if startX < endX {
		// to the right
		return doc.verticalNewVariants(startX, startY, endX, endY, currentVariantIndex, currentLineIndex,
			variants, startElem, destElem, false, false, junctionPointsToRight, false)
	} else if startY > endY {
		// to the left
		return doc.verticalNewVariants(startX, startY, endX, endY, currentVariantIndex, currentLineIndex,
			variants, startElem, destElem, true, false, junctionPointsToLeft, false)
	} else {
		// same level: startY == endY
		if endY == y {
			// reached the end
			l := variants[currentVariantIndex][currentLineIndex]
			changedCurrentLine := newConnectionLine(l.StartX, l.StartY, l.StartX, y)
			variants[currentVariantIndex][currentLineIndex] = changedCurrentLine
			return variants
		} else if y < endY {
			// trying to the right
			variants = doc.verticalNewVariants(startX, startY, endX, endY, currentVariantIndex, currentLineIndex,
				variants, startElem, destElem, false, true, []int{y}, false)
			// trying to the left
			return doc.verticalNewVariants(startX, startY, endX, endY, currentVariantIndex, currentLineIndex,
				variants, startElem, destElem, true, false, []int{y}, false)
		} else {
			// error
			return emptyVariant(currentVariantIndex, variants)
		}
	}
}

func (doc *BoxesDocument) goToUp(
	startX, startY, endX, endY int,
	currentVariantIndex, currentLineIndex int,
	variants [][]ConnectionLine,
	startElem, destElem *LayoutElement) [][]ConnectionLine {
	fmt.Println("DEBUG: goToUp", startX, startY, endX, endY, currentVariantIndex, currentLineIndex)
	if len(variants[currentVariantIndex]) == 0 && currentLineIndex > 0 {
		return variants
	}
	if startY < endY {
		// moved too far up
		fmt.Println("DEBUG: goToUp - emptyVariant-1", startX, startY, endX, endY, currentVariantIndex)
		return emptyVariant(currentVariantIndex, variants)
	}

	if (startX == endX) && (startY == endY) {
		return variants
	}

	y, junctionPointsToLeft, junctionPointsToRight, err := doc.checkRoadsToTheTop(startX, startY)
	if err != nil {
		fmt.Println("DEBUG: goToUp - emptyVariant-2", startX, startY, endX, endY, currentVariantIndex)
		return emptyVariant(currentVariantIndex, variants)
	}
	if startX < endX {
		// to the right
		return doc.verticalNewVariants(startX, startY, endX, endY, currentVariantIndex, currentLineIndex,
			variants, startElem, destElem, false, false, junctionPointsToRight, true)
	} else if startY > endY {
		// to the left
		return doc.verticalNewVariants(startX, startY, endX, endY, currentVariantIndex, currentLineIndex,
			variants, startElem, destElem, true, false, junctionPointsToLeft, true)
	} else {
		// same level: startY == endY
		if endY == y {
			// reached the end
			l := variants[currentVariantIndex][currentLineIndex]
			changedCurrentLine := newConnectionLine(l.StartX, l.StartY, l.StartX, y)
			variants[currentVariantIndex][currentLineIndex] = changedCurrentLine
			return variants
		} else if y > endY {
			// trying to the right
			variants = doc.verticalNewVariants(startX, startY, endX, endY, currentVariantIndex, currentLineIndex,
				variants, startElem, destElem, false, true, []int{y}, true)
			// trying to the left
			return doc.verticalNewVariants(startX, startY, endX, endY, currentVariantIndex, currentLineIndex,
				variants, startElem, destElem, true, false, []int{y}, true)
		} else {
			// error
			return emptyVariant(currentVariantIndex, variants)
		}
	}
}
