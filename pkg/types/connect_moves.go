package types

func (doc *BoxesDocument) horizontalNewVariants(
	startY, endX, endY int,
	currentVariantIndex, currentLineIndex int,
	variants [][]ConnectionLine,
	startElem, destElem *LayoutElement, up, alwaysNewVariant bool, junctionPoints []int) [][]ConnectionLine {
	for i, p := range junctionPoints {
		newVariant := variants[currentVariantIndex]
		changedCurrentLine := newConnectionLine(newVariant[currentLineIndex].StartX, startY, p, startY)
		newX := changedCurrentLine.EndX
		newY := changedCurrentLine.EndY
		if i > 0 || alwaysNewVariant {
			newVariant = variants[currentVariantIndex][:]
		}
		newVariant[currentLineIndex] = changedCurrentLine
		newLine := newConnectionLine(newX, newY, newX, newY)
		newVariant = append(newVariant, newLine)
		newVariantIndex := currentVariantIndex
		newLineIndex := len(newVariant) - 1
		if i > 0 || alwaysNewVariant {
			variants = append(variants, newVariant)
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
	startX, endX, endY int,
	currentVariantIndex, currentLineIndex int,
	variants [][]ConnectionLine,
	startElem, destElem *LayoutElement, left, alwaysNewVariant bool, junctionPoints []int) [][]ConnectionLine {
	for i, p := range junctionPoints {
		newVariant := variants[currentVariantIndex]
		changedCurrentLine := newConnectionLine(startX, newVariant[currentLineIndex].StartY, startX, p)
		newX := changedCurrentLine.EndX
		newY := changedCurrentLine.EndY
		if i > 0 || alwaysNewVariant {
			newVariant = variants[currentVariantIndex][:]
		}
		newVariant[currentLineIndex] = changedCurrentLine
		newLine := newConnectionLine(newX, newY, newX, newY)
		newVariant = append(newVariant, newLine)
		newVariantIndex := currentVariantIndex
		newLineIndex := len(newVariant) - 1
		if i > 0 || alwaysNewVariant {
			variants = append(variants, newVariant)
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
	if len(variants[currentVariantIndex]) == 0 && currentLineIndex > 0 {
		return variants
	}
	if startX < endX {
		// moved too far left
		return emptyVariant(currentVariantIndex, variants)
	}

	if (startX == endX) && (startY == endY) {
		return variants
	}

	x, junctionPointsUpwards, junctionPointsDownwards, err := doc.checkRoadsToTheLeft(startX, startY)
	if err != nil {
		return emptyVariant(currentVariantIndex, variants)
	}
	if startY < endY {
		// downwards
		return doc.horizontalNewVariants(startY, endX, endY, currentVariantIndex, currentLineIndex,
			variants, startElem, destElem, false, false, junctionPointsDownwards)
	} else if startY > endY {
		// upwards
		return doc.horizontalNewVariants(startY, endX, endY, currentVariantIndex, currentLineIndex,
			variants, startElem, destElem, true, false, junctionPointsUpwards)
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
			variants = doc.horizontalNewVariants(startY, endX, endY, currentVariantIndex, currentLineIndex,
				variants, startElem, destElem, false, true, []int{x})
			// trying upwards
			return doc.horizontalNewVariants(startY, endX, endY, currentVariantIndex, currentLineIndex,
				variants, startElem, destElem, true, false, []int{x})
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
	if len(variants[currentVariantIndex]) == 0 && currentLineIndex > 0 {
		return variants
	}
	if startX > endX {
		// moved too far right
		return emptyVariant(currentVariantIndex, variants)
	}

	if (startX == endX) && (startY == endY) {
		return variants
	}

	x, junctionPointsUpwards, junctionPointsDownwards, err := doc.checkRoadsToTheRight(startX, startY)
	if err != nil {
		return emptyVariant(currentVariantIndex, variants)
	}
	if startY < endY {
		// downwards
		return doc.horizontalNewVariants(startY, endX, endY, currentVariantIndex, currentLineIndex,
			variants, startElem, destElem, false, false, junctionPointsDownwards)
	} else if startY > endY {
		// upwards
		return doc.horizontalNewVariants(startY, endX, endY, currentVariantIndex, currentLineIndex,
			variants, startElem, destElem, true, false, junctionPointsUpwards)
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
			variants = doc.horizontalNewVariants(startY, endX, endY, currentVariantIndex, currentLineIndex,
				variants, startElem, destElem, false, true, []int{x})
			// trying upwards
			return doc.horizontalNewVariants(startY, endX, endY, currentVariantIndex, currentLineIndex,
				variants, startElem, destElem, true, false, []int{x})
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
	if len(variants[currentVariantIndex]) == 0 && currentLineIndex > 0 {
		return variants
	}
	if startY > endY {
		// moved too far down
		return emptyVariant(currentVariantIndex, variants)
	}

	if (startX == endX) && (startY == endY) {
		return variants
	}

	y, junctionPointsToLeft, junctionPointsToRight, err := doc.checkRoadsToTheBottom(startX, startY)
	if err != nil {
		return emptyVariant(currentVariantIndex, variants)
	}
	if startX < endX {
		// to the right
		return doc.verticalNewVariants(startX, endX, endY, currentVariantIndex, currentLineIndex,
			variants, startElem, destElem, false, false, junctionPointsToRight)
	} else if startY > endY {
		// to the left
		return doc.verticalNewVariants(startX, endX, endY, currentVariantIndex, currentLineIndex,
			variants, startElem, destElem, true, false, junctionPointsToLeft)
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
			variants = doc.verticalNewVariants(startX, endX, endY, currentVariantIndex, currentLineIndex,
				variants, startElem, destElem, false, true, []int{y})
			// trying to the left
			return doc.verticalNewVariants(startX, endX, endY, currentVariantIndex, currentLineIndex,
				variants, startElem, destElem, true, false, []int{y})
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
	if len(variants[currentVariantIndex]) == 0 && currentLineIndex > 0 {
		return variants
	}
	if startY < endY {
		// moved too far up
		return emptyVariant(currentVariantIndex, variants)
	}

	if (startX == endX) && (startY == endY) {
		return variants
	}

	y, junctionPointsToLeft, junctionPointsToRight, err := doc.checkRoadsToTheTop(startX, startY)
	if err != nil {
		return emptyVariant(currentVariantIndex, variants)
	}
	if startX < endX {
		// to the right
		return doc.verticalNewVariants(startX, endX, endY, currentVariantIndex, currentLineIndex,
			variants, startElem, destElem, false, false, junctionPointsToRight)
	} else if startY > endY {
		// to the left
		return doc.verticalNewVariants(startX, endX, endY, currentVariantIndex, currentLineIndex,
			variants, startElem, destElem, true, false, junctionPointsToLeft)
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
			variants = doc.verticalNewVariants(startX, endX, endY, currentVariantIndex, currentLineIndex,
				variants, startElem, destElem, false, true, []int{y})
			// trying to the left
			return doc.verticalNewVariants(startX, endX, endY, currentVariantIndex, currentLineIndex,
				variants, startElem, destElem, true, false, []int{y})
		} else {
			// error
			return emptyVariant(currentVariantIndex, variants)
		}
	}
}
