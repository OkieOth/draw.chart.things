package types

type TextDimensionCalculator interface {
	Dimensions(txt string, format *FontDef) (width, height int)
}

func (l *LayoutElement) incrementX(xOffset int) {
	l.X += xOffset
	if l.Vertical != nil {
		l.Vertical.X += xOffset
		for i := 0; i < len(l.Vertical.Elems); i++ {
			sub := &l.Vertical.Elems[i]
			sub.incrementX(xOffset)
		}
	}
	if l.Horizontal != nil {
		l.Horizontal.X += xOffset
		for i := 0; i < len(l.Horizontal.Elems); i++ {
			sub := &l.Horizontal.Elems[i]
			sub.incrementX(xOffset)
		}
	}
}

func (l *LayoutElement) incrementY(yOffset int) {
	l.Y += yOffset
	if l.Vertical != nil {
		l.Vertical.Y += yOffset
		for i := 0; i < len(l.Vertical.Elems); i++ {
			sub := &l.Vertical.Elems[i]
			sub.incrementY(yOffset)
		}
	}
	if l.Horizontal != nil {
		l.Horizontal.Y += yOffset
		for i := 0; i < len(l.Horizontal.Elems); i++ {
			sub := &l.Horizontal.Elems[i]
			sub.incrementY(yOffset)
		}
	}
}

func (l *LayoutElement) centerHorizontalElems() {
	if l.Horizontal != nil {
		cont := l.Horizontal
		contYOffset := (l.Height - cont.Height) / 2
		contXOffset := (l.Width - cont.Width) / 2
		cont.Y = l.Y + contYOffset
		cont.X = l.X + contXOffset
		for i := 0; i < len(cont.Elems); i++ {
			sub := &cont.Elems[i]
			offsetY := ((cont.Height - sub.Height) / 2)
			sub.incrementY(offsetY)
			sub.incrementX(contXOffset)
			//sub.centerHorizontalElems()
			sub.Center()
		}
	}
}

func (l *LayoutElement) centerVerticalElems() {
	if l.Vertical != nil {
		cont := l.Vertical
		contYOffset := (l.Y + l.Height - cont.Y - cont.Height) / 2
		contXOffset := (l.Width - cont.Width) / 2
		cont.Y = l.Y + contYOffset
		cont.X = l.X + contXOffset
		for i := 0; i < len(cont.Elems); i++ {
			sub := &cont.Elems[i]
			sub.incrementY(contYOffset)
			offsetX := contXOffset + ((cont.Width - sub.Width) / 2)
			sub.incrementX(offsetX)
			//sub.centerVerticalElems()
			sub.Center()
		}
	}
}

func (l *LayoutElement) AreOnTheSameVerticalLevel(otherElem *LayoutElement) bool {
	return ((l.CenterY > otherElem.Y) && (l.CenterY < otherElem.Y+otherElem.Height)) ||
		((otherElem.CenterY > l.Y) && (otherElem.CenterY < l.Y+l.Height))
}

type ConnDirection int

const (
	ConnDirectionLeft ConnDirection = iota
	ConnDirectionUp
	ConnDirectionRight
	ConnDirectionDown
)

func (l *LayoutElement) ConnectorStart(otherElem *LayoutElement) (int, int, ConnDirection) {
	if l.AreOnTheSameVerticalLevel(otherElem) {
		// the elements are on the same vertical level
		if l.CenterX < otherElem.CenterX {
			return l.X + l.Width, l.CenterY, ConnDirectionRight
		} else {
			return l.X, l.CenterY, ConnDirectionLeft
		}

	} else if l.CenterY < otherElem.CenterY {
		// connection from top to bottom
		return l.CenterX, l.Y + l.Height, ConnDirectionDown
	} else {
		// connection from bottom to top
		return l.CenterX, l.Y, ConnDirectionUp
	}
}

func (l *LayoutElement) fixVerticalCollisionDownWithCorner(connectionLines []ConnectionLine, index int, distanceToBorder int) []ConnectionLine {
	// I don't test the index, because if it crashes then the understanding of the algorithm is wrong
	conn := connectionLines[index]
	nextConn := connectionLines[index+1]
	// next Connection is a horizontal line
	leftX, _ := connToLeftRightX(nextConn)
	ret := make([]ConnectionLine, 3)
	ret[0] = ConnectionLine{
		StartX:   conn.StartX,
		StartY:   conn.StartY,
		EndX:     conn.StartX,
		EndY:     l.Y - distanceToBorder,
		MovedOut: true,
	}
	if leftX < conn.StartX {
		// next line is to the left
		ret[1] = ConnectionLine{
			StartX:   conn.StartX,
			StartY:   l.Y - distanceToBorder,
			EndX:     l.X - distanceToBorder,
			EndY:     l.Y - distanceToBorder,
			MovedOut: true,
		}
		ret[2] = ConnectionLine{
			StartX:   l.X - distanceToBorder,
			StartY:   l.Y - distanceToBorder,
			EndX:     l.X - distanceToBorder,
			EndY:     nextConn.StartY,
			MovedOut: true,
		}
		// adjust the next connection
		if nextConn.StartX > nextConn.EndX {
			// next line is to the left
			connectionLines[index+1].StartX = ret[2].EndX
		} else {
			connectionLines[index+1].EndX = ret[2].EndX
		}
	} else {
		// next line is to the right
		ret[1] = ConnectionLine{
			StartX:   conn.StartX,
			StartY:   l.Y - distanceToBorder,
			EndX:     l.X + l.Width + distanceToBorder,
			EndY:     l.Y - distanceToBorder,
			MovedOut: true,
		}
		ret[2] = ConnectionLine{
			StartX:   l.X - distanceToBorder,
			StartY:   l.Y - distanceToBorder,
			EndX:     l.X + l.Width + distanceToBorder,
			EndY:     nextConn.StartY,
			MovedOut: true,
		}
		// adjust the next connection
		if nextConn.StartX < nextConn.EndX {
			// next line is to the left
			connectionLines[index+1].StartX = ret[2].EndX
		} else {
			connectionLines[index+1].EndX = ret[2].EndX
		}
	}
	return ret
}

func (l *LayoutElement) fixVerticalCollisionDown(conn ConnectionLine, distanceToBorder int) []ConnectionLine {
	ret := make([]ConnectionLine, 5)
	ret[0] = ConnectionLine{
		StartX:   conn.StartX,
		StartY:   conn.StartY,
		EndX:     conn.StartX,
		EndY:     l.Y - distanceToBorder,
		MovedOut: true,
	}
	if conn.StartX < l.CenterX {
		// fix to the left
		ret[1] = ConnectionLine{
			StartX:   conn.StartX,
			StartY:   l.Y - distanceToBorder,
			EndX:     l.X - distanceToBorder,
			EndY:     l.Y - distanceToBorder,
			MovedOut: true,
		}
		ret[2] = ConnectionLine{
			StartX:   l.X - distanceToBorder,
			StartY:   l.Y - distanceToBorder,
			EndX:     l.X - distanceToBorder,
			EndY:     l.Y + l.Height + distanceToBorder,
			MovedOut: true,
		}
		ret[3] = ConnectionLine{
			StartX:   l.X - distanceToBorder,
			StartY:   l.Y + l.Height + distanceToBorder,
			EndX:     conn.StartX,
			EndY:     l.Y + l.Height + distanceToBorder,
			MovedOut: true,
		}

	} else {
		// fix to the right
		ret[1] = ConnectionLine{
			StartX:   conn.StartX,
			StartY:   l.Y - distanceToBorder,
			EndX:     l.X + l.Width + distanceToBorder,
			EndY:     l.Y - distanceToBorder,
			MovedOut: true,
		}
		ret[2] = ConnectionLine{
			StartX:   l.X + l.Width + distanceToBorder,
			StartY:   l.Y - distanceToBorder,
			EndX:     l.X + l.Width + distanceToBorder,
			EndY:     l.Y + l.Height + distanceToBorder,
			MovedOut: true,
		}
		ret[3] = ConnectionLine{
			StartX:   l.X + l.Width + distanceToBorder,
			StartY:   l.Y + l.Height + distanceToBorder,
			EndX:     conn.StartX,
			EndY:     l.Y + l.Height + distanceToBorder,
			MovedOut: true,
		}
	}
	// this part needs to be checked for further collisions
	ret[4] = ConnectionLine{
		StartX: conn.StartX,
		StartY: l.Y + l.Height + distanceToBorder,
		EndX:   conn.StartX,
		EndY:   conn.EndY,
	}
	return ret
}

func (l *LayoutElement) fixVerticalCollisionUpWithCorner(connectionLines []ConnectionLine, index int, distanceToBorder int) []ConnectionLine {
	// I don't test the index, because if it crashes then the understanding of the algorithm is wrong
	conn := connectionLines[index]
	nextConn := connectionLines[index+1]
	// next Connection is a horizontal line
	leftX, _ := connToLeftRightX(nextConn)
	ret := make([]ConnectionLine, 3)
	ret[0] = ConnectionLine{
		StartX:   conn.StartX,
		StartY:   conn.StartY,
		EndX:     conn.StartX,
		EndY:     l.Y + l.Height + distanceToBorder,
		MovedOut: true,
	}
	if leftX < conn.StartX {
		// next line is to the left
		ret[1] = ConnectionLine{
			StartX:   conn.StartX,
			StartY:   l.Y + l.Height + distanceToBorder,
			EndX:     l.X - distanceToBorder,
			EndY:     l.Y + l.Height + distanceToBorder,
			MovedOut: true,
		}
		ret[2] = ConnectionLine{
			StartX:   l.X - distanceToBorder,
			StartY:   l.Y + l.Height + distanceToBorder,
			EndX:     l.X - distanceToBorder,
			EndY:     nextConn.StartY,
			MovedOut: true,
		}
		// adjust the next connection
		if nextConn.StartX > nextConn.EndX {
			// next line is to the left
			connectionLines[index+1].StartX = ret[2].EndX
		} else {
			connectionLines[index+1].EndX = ret[2].EndX
		}
	} else {
		// next line is to the right
		ret[1] = ConnectionLine{
			StartX:   conn.StartX,
			StartY:   l.Y + l.Height + distanceToBorder,
			EndX:     l.X + l.Width + distanceToBorder,
			EndY:     l.Y + l.Height + distanceToBorder,
			MovedOut: true,
		}
		ret[2] = ConnectionLine{
			StartX:   l.X + l.Width + distanceToBorder,
			StartY:   l.Y + l.Height + distanceToBorder,
			EndX:     l.X + l.Width + distanceToBorder,
			EndY:     nextConn.StartY,
			MovedOut: true,
		}
		// adjust the next connection
		if nextConn.StartX < nextConn.EndX {
			// next line is to the left
			connectionLines[index+1].StartX = ret[2].EndX
		} else {
			connectionLines[index+1].EndX = ret[2].EndX
		}
	}
	return ret
}

func (l *LayoutElement) fixVerticalCollisionUp(conn ConnectionLine, distanceToBorder int) []ConnectionLine {
	ret := make([]ConnectionLine, 5)
	ret[0] = ConnectionLine{
		StartX:   conn.StartX,
		StartY:   conn.StartY,
		EndX:     conn.StartX,
		EndY:     l.Y + l.Height + distanceToBorder,
		MovedOut: true,
	}
	if conn.StartX < l.CenterX {
		// fix to the left
		ret[1] = ConnectionLine{
			StartX:   conn.StartX,
			StartY:   l.Y + l.Height + distanceToBorder,
			EndX:     l.X - distanceToBorder,
			EndY:     l.Y + l.Height + distanceToBorder,
			MovedOut: true,
		}
		ret[2] = ConnectionLine{
			StartX:   l.X - distanceToBorder,
			StartY:   l.Y + l.Height + distanceToBorder,
			EndX:     l.X - distanceToBorder,
			EndY:     l.Y - distanceToBorder,
			MovedOut: true,
		}
		ret[3] = ConnectionLine{
			StartX:   l.X - distanceToBorder,
			StartY:   l.Y - distanceToBorder,
			EndX:     conn.StartX,
			EndY:     l.Y - distanceToBorder,
			MovedOut: true,
		}
	} else {
		// fix to the right
		ret[1] = ConnectionLine{
			StartX:   conn.StartX,
			StartY:   l.Y + l.Height + distanceToBorder,
			EndX:     l.X + l.Width + distanceToBorder,
			EndY:     l.Y + l.Height + distanceToBorder,
			MovedOut: true,
		}
		ret[2] = ConnectionLine{
			StartX:   l.X + l.Width + distanceToBorder,
			StartY:   l.Y + l.Height + distanceToBorder,
			EndX:     l.X + l.Width + distanceToBorder,
			EndY:     l.Y - distanceToBorder,
			MovedOut: true,
		}
		ret[3] = ConnectionLine{
			StartX:   l.X + l.Width + distanceToBorder,
			StartY:   l.Y - distanceToBorder,
			EndX:     conn.StartX,
			EndY:     l.Y - distanceToBorder,
			MovedOut: true,
		}
	}
	// this part needs to be checked for further collisions
	ret[4] = ConnectionLine{
		StartX: conn.StartX,
		StartY: l.Y - distanceToBorder,
		EndX:   conn.StartX,
		EndY:   conn.EndY,
	}
	return ret
}

func (l *LayoutElement) fixHorizontalCollisionRightWithCorner(connectionLines []ConnectionLine, index int, distanceToBorder int) []ConnectionLine {
	// I don't test the index, because if it crashes then the understanding of the algorithm is wrong
	conn := connectionLines[index]
	nextConn := connectionLines[index+1]
	// next Connection is a vertical line
	topY, _ := connToUpperLowerY(nextConn)
	ret := make([]ConnectionLine, 3)
	ret[0] = ConnectionLine{
		StartX:   conn.StartX,
		StartY:   conn.StartY,
		EndX:     l.X - distanceToBorder,
		EndY:     conn.StartY,
		MovedOut: true,
	}
	if topY < conn.StartY {
		// next line is to the top
		ret[1] = ConnectionLine{
			StartX:   ret[0].EndX,
			StartY:   conn.StartY,
			EndX:     ret[0].EndX,
			EndY:     l.Y - distanceToBorder,
			MovedOut: true,
		}
		ret[2] = ConnectionLine{
			StartX:   ret[1].EndX,
			StartY:   ret[1].EndY,
			EndX:     nextConn.StartX,
			EndY:     ret[1].EndY,
			MovedOut: true,
		}
		// adjust the next connection
		if nextConn.StartY > nextConn.EndY {
			// next line is to the bottom
			connectionLines[index+1].StartY = ret[2].EndY
		} else {
			connectionLines[index+1].EndY = ret[2].EndY
		}
	} else {
		// next line is to the bottom
		ret[1] = ConnectionLine{
			StartX:   ret[0].EndX,
			StartY:   conn.StartY,
			EndX:     ret[0].EndX,
			EndY:     l.Y + l.Height + distanceToBorder,
			MovedOut: true,
		}
		ret[2] = ConnectionLine{
			StartX:   ret[1].EndX,
			StartY:   ret[1].EndY,
			EndX:     nextConn.StartX,
			EndY:     ret[1].EndY,
			MovedOut: true,
		}
		// adjust the next connection
		// adjust the next connection
		if nextConn.StartY < nextConn.EndY {
			// next line is to the bottom
			connectionLines[index+1].StartY = ret[2].EndY
		} else {
			connectionLines[index+1].EndY = ret[2].EndY
		}
	}
	return ret
}

func (l *LayoutElement) fixHorizontalCollisionRight(conn ConnectionLine, distanceToBorder int) []ConnectionLine {
	ret := make([]ConnectionLine, 5)
	ret[0] = ConnectionLine{
		StartX:   conn.StartX,
		StartY:   conn.StartY,
		EndX:     l.X - distanceToBorder,
		EndY:     conn.StartY,
		MovedOut: true,
	}
	if conn.StartY < l.CenterY {
		// fix to the top
		ret[1] = ConnectionLine{
			StartX:   l.X - distanceToBorder,
			StartY:   conn.StartY,
			EndX:     l.X - distanceToBorder,
			EndY:     l.Y - distanceToBorder,
			MovedOut: true,
		}
		ret[2] = ConnectionLine{
			StartX:   l.X - distanceToBorder,
			StartY:   l.Y - distanceToBorder,
			EndX:     l.X + l.Width + distanceToBorder,
			EndY:     l.Y - distanceToBorder,
			MovedOut: true,
		}
		ret[3] = ConnectionLine{
			StartX:   l.X + l.Width + distanceToBorder,
			StartY:   l.Y - distanceToBorder,
			EndX:     l.X + l.Width + distanceToBorder,
			EndY:     conn.StartY,
			MovedOut: true,
		}
	} else {
		// fix to the bottom
		ret[1] = ConnectionLine{
			StartX:   l.X - distanceToBorder,
			StartY:   conn.StartY,
			EndX:     l.X - distanceToBorder,
			EndY:     l.Y + l.Height + distanceToBorder,
			MovedOut: true,
		}
		ret[2] = ConnectionLine{
			StartX:   l.X - distanceToBorder,
			StartY:   l.Y + l.Height + distanceToBorder,
			EndX:     l.X + l.Width + distanceToBorder,
			EndY:     l.Y + l.Height + distanceToBorder,
			MovedOut: true,
		}
		ret[3] = ConnectionLine{
			StartX:   l.X + l.Width + distanceToBorder,
			StartY:   l.Y + l.Height + distanceToBorder,
			EndX:     l.X + l.Width + distanceToBorder,
			EndY:     conn.StartY,
			MovedOut: true,
		}
	}
	// this part needs to be checked for further collisions
	ret[4] = ConnectionLine{
		StartX: l.X + l.Width + distanceToBorder,
		StartY: conn.StartY,
		EndX:   conn.EndX,
		EndY:   conn.EndY,
	}
	return ret
}

func (l *LayoutElement) fixHorizontalCollisionLeftWithCorner(connectionLines []ConnectionLine, index int, distanceToBorder int) []ConnectionLine {
	// I don't test the index, because if it crashes then the understanding of the algorithm is wrong
	conn := connectionLines[index]
	nextConn := connectionLines[index+1]
	// next Connection is a vertical line
	topY, _ := connToUpperLowerY(nextConn)
	ret := make([]ConnectionLine, 3)
	ret[0] = ConnectionLine{
		StartX:   conn.StartX,
		StartY:   conn.StartY,
		EndX:     l.X + l.Width + distanceToBorder,
		EndY:     conn.StartY,
		MovedOut: true,
	}
	if topY < conn.StartY {
		// next line is to the top
		ret[1] = ConnectionLine{
			StartX:   ret[0].EndX,
			StartY:   conn.StartY,
			EndX:     ret[0].EndX,
			EndY:     l.Y - distanceToBorder,
			MovedOut: true,
		}
		ret[2] = ConnectionLine{
			StartX:   ret[1].EndX,
			StartY:   ret[1].EndY,
			EndX:     conn.EndX,
			EndY:     ret[1].EndY,
			MovedOut: true,
		}
		// adjust the next connection
		if nextConn.StartY > nextConn.EndY {
			// next line is to the bottom
			connectionLines[index+1].StartY = ret[2].EndY
		} else {
			connectionLines[index+1].EndY = ret[2].EndY
		}
	} else {
		// next line is to the bottom
		ret[1] = ConnectionLine{
			StartX:   ret[0].EndX,
			StartY:   conn.StartY,
			EndX:     ret[0].EndX,
			EndY:     l.Y + l.Height + distanceToBorder,
			MovedOut: true,
		}
		ret[2] = ConnectionLine{
			StartX:   ret[1].EndX,
			StartY:   ret[1].EndY,
			EndX:     conn.EndX,
			EndY:     ret[1].EndY,
			MovedOut: true,
		}
		// adjust the next connection
		if nextConn.StartY < nextConn.EndY {
			// next line is to the bottom
			connectionLines[index+1].StartY = ret[2].EndY
		} else {
			connectionLines[index+1].EndY = ret[2].EndY
		}
	}
	return ret
}

func (l *LayoutElement) fixHorizontalCollisionLeft(conn ConnectionLine, distanceToBorder int) []ConnectionLine {
	ret := make([]ConnectionLine, 5)
	ret[0] = ConnectionLine{
		StartX:   conn.StartX,
		StartY:   conn.StartY,
		EndX:     l.X + l.Width + distanceToBorder,
		EndY:     conn.StartY,
		MovedOut: true,
	}
	if conn.StartY < l.CenterY {
		// fix to the top
		ret[1] = ConnectionLine{
			StartX:   l.X + l.Width + distanceToBorder,
			StartY:   conn.StartY,
			EndX:     l.X + l.Width + distanceToBorder,
			EndY:     l.Y - distanceToBorder,
			MovedOut: true,
		}
		ret[2] = ConnectionLine{
			StartX:   l.X + l.Width + distanceToBorder,
			StartY:   l.Y - distanceToBorder,
			EndX:     l.X - distanceToBorder,
			EndY:     l.Y - distanceToBorder,
			MovedOut: true,
		}
		ret[3] = ConnectionLine{
			StartX:   l.X - distanceToBorder,
			StartY:   l.Y - distanceToBorder,
			EndX:     l.X - distanceToBorder,
			EndY:     conn.StartY,
			MovedOut: true,
		}
	} else {
		// fix to the bottom
		ret[1] = ConnectionLine{
			StartX:   l.X + l.Width + distanceToBorder,
			StartY:   conn.StartY,
			EndX:     l.X + l.Width + distanceToBorder,
			EndY:     l.Y + l.Height + distanceToBorder,
			MovedOut: true,
		}
		ret[2] = ConnectionLine{
			StartX:   l.X + l.Width + distanceToBorder,
			StartY:   l.Y + l.Height + distanceToBorder,
			EndX:     l.X - distanceToBorder,
			EndY:     l.Y + l.Height + distanceToBorder,
			MovedOut: true,
		}
		ret[3] = ConnectionLine{
			StartX:   l.X - distanceToBorder,
			StartY:   l.Y + l.Height + distanceToBorder,
			EndX:     l.X - distanceToBorder,
			EndY:     conn.StartY,
			MovedOut: true,
		}
	}
	// this part needs to be checked for further collisions
	ret[4] = ConnectionLine{
		StartX: l.X - distanceToBorder,
		StartY: conn.StartY,
		EndX:   conn.EndX,
		EndY:   conn.EndY,
	}
	return ret
}

func connToUpperLowerY(conn ConnectionLine) (int, int) {
	if conn.StartY < conn.EndY {
		return conn.StartY, conn.EndY
	}
	return conn.EndY, conn.StartY
}
func connToLeftRightX(conn ConnectionLine) (int, int) {
	if conn.StartX < conn.EndX {
		return conn.StartX, conn.EndX
	}
	return conn.EndX, conn.StartX
}

func (l *LayoutElement) ShouldBeDrawn() bool {
	return l.Caption != "" || l.Text1 != "" || l.Text2 != ""
}

func between(value, min, max int) bool {
	return value > min && value < max
}

func (l *LayoutElement) FixCollisionInCase(connectionLines []ConnectionLine, index int, distanceToBorder int) []ConnectionLine {
	if !l.ShouldBeDrawn() {
		return []ConnectionLine{connectionLines[index]}
	}
	conn := connectionLines[index]
	if conn.StartX == conn.EndX {
		// vertical line
		upperY, lowerY := connToUpperLowerY(conn)
		if (l.X < conn.StartX) && ((l.X + l.Width) > conn.StartX) &&
			(((l.Y > upperY) && ((l.Y + l.Height) < lowerY)) || between(lowerY, l.Y, l.Y+l.Height) || between(upperY, l.Y, l.Y+l.Height)) {
			// has collision
			if conn.StartY < conn.EndY {
				// going down
				if lowerY > (l.Y + l.Height) {
					// line going full through the box
					return l.fixVerticalCollisionDown(conn, distanceToBorder)
				} else {
					// line going only partially through the box
					return l.fixVerticalCollisionDownWithCorner(connectionLines, index, distanceToBorder)
				}
			} else {
				// going up
				if upperY < l.Y {
					// line going full through the box
					return l.fixVerticalCollisionUp(conn, distanceToBorder)
				} else {
					// line going only partially through the box
					return l.fixVerticalCollisionUpWithCorner(connectionLines, index, distanceToBorder)
				}
			}
		}
	} else {
		// horizontal line
		leftX, rightX := connToLeftRightX(conn)
		if (l.Y < conn.StartY) && ((l.Y + l.Height) > conn.StartY) &&
			(((l.X > leftX) && ((l.X + l.Width) < rightX)) || between(leftX, l.X, l.X+l.Width) || between(rightX, l.X, l.X+l.Width)) {
			// has collision
			if conn.StartX < conn.EndX {
				// going right
				if rightX > (l.X + l.Width) {
					// line going full through the box
					return l.fixHorizontalCollisionRight(conn, distanceToBorder)
				} else {
					// line going only partially through the box
					return l.fixHorizontalCollisionRightWithCorner(connectionLines, index, distanceToBorder)
				}
			} else {
				// going left
				if leftX < l.X {
					// line going full through the box
					return l.fixHorizontalCollisionLeft(conn, distanceToBorder)
				} else {
					// line going only partially through the box
					return l.fixHorizontalCollisionLeftWithCorner(connectionLines, index, distanceToBorder)
				}
			}
		}
	}
	return []ConnectionLine{conn}
}

func (l *LayoutElement) initVertical(c TextDimensionCalculator, yInnerOffset, defaultPadding, defaultBoxMargin int) {
	if l.Vertical != nil && len(l.Vertical.Elems) > 0 {
		curX := l.X
		l.Vertical.X = curX
		curY := l.Y + yInnerOffset
		l.Vertical.Y = curY
		var h, w int
		var hasChilds bool
		for i := 0; i < len(l.Vertical.Elems); i++ {
			sub := &l.Vertical.Elems[i]
			if (sub.Horizontal != nil && len(sub.Horizontal.Elems) > 0) || (sub.Vertical != nil && len(sub.Vertical.Elems) > 0) {
				hasChilds = true
			}
			if h > 0 {
				h += defaultBoxMargin
			}
			sub.X = curX
			sub.Y = curY
			sub.InitDimensions(c, defaultPadding, defaultBoxMargin)
			curY += (sub.Height + defaultBoxMargin)
			h += sub.Height
			if sub.Width > w {
				w = sub.Width
			}
			if sub.Width > l.Vertical.Width {
				l.Vertical.Width = sub.Width
			}
		}
		if !hasChilds {
			for i := 0; i < len(l.Vertical.Elems); i++ {
				sub := &l.Vertical.Elems[i]
				sub.Width = w
			}
		}

		l.Vertical.Height = h + defaultPadding
		l.Height += l.Vertical.Height
		if w > l.Width {
			l.Width = w
		}
	}
}

func (l *LayoutElement) initHorizontal(c TextDimensionCalculator, yInnerOffset, defaultPadding, defaultBoxMargin int) {
	if l.Horizontal != nil && len(l.Horizontal.Elems) > 0 {
		curX := l.X
		l.Horizontal.X = curX
		curY := l.Y + yInnerOffset
		l.Horizontal.Y = curY
		var h, w int
		var hasChilds bool
		for i := 0; i < len(l.Horizontal.Elems); i++ {
			sub := &l.Horizontal.Elems[i]
			if (sub.Horizontal != nil && len(sub.Horizontal.Elems) > 0) || (sub.Vertical != nil && len(sub.Vertical.Elems) > 0) {
				hasChilds = true
			}
			if w > 0 {
				w += defaultBoxMargin
			}
			sub.X = curX
			sub.Y = curY
			sub.InitDimensions(c, defaultPadding, defaultBoxMargin)
			curX += (sub.Width + defaultBoxMargin)
			w += sub.Width
			if sub.Height > h {
				h = sub.Height
			}
			if sub.Height > l.Horizontal.Height {
				l.Horizontal.Height = sub.Height
			}
		}
		if !hasChilds {
			for i := 0; i < len(l.Horizontal.Elems); i++ {
				sub := &l.Horizontal.Elems[i]
				sub.Height = h
			}
		}

		l.Horizontal.Height = h + defaultPadding
		l.Height += l.Horizontal.Height
		l.Horizontal.Width = w

		if l.Width < w {
			l.Width = w
		}
	}
}

func (l *LayoutElement) InitDimensions(c TextDimensionCalculator, defaultPadding, defaultBoxMargin int) {
	var cW, cH, t1W, t1H, t2W, t2H int
	//var yCaptionOffset, yText1Offset, yText2Offset, yInnerOffset int
	var yInnerOffset int
	if l.Caption != "" || l.Text1 != "" || l.Text2 != "" {
		l.Height = (2 * l.Format.Padding)
		if l.Caption != "" {
			cW, cH = c.Dimensions(l.Caption, &l.Format.FontCaption)
			l.Height += cH + l.Format.FontCaption.SpaceTop + l.Format.FontCaption.SpaceBottom
		}
		if l.Text1 != "" {
			p := l.Format.Padding
			if l.Format.FontText1.SpaceTop > 0 {
				p = l.Format.FontText1.SpaceTop
			}
			t1W, t1H = c.Dimensions(l.Text1, &l.Format.FontText1)
			l.Height += t1H + p + l.Format.FontText1.SpaceBottom
		}
		if l.Text2 != "" {
			p := l.Format.Padding
			if l.Format.FontText2.SpaceTop > 0 {
				p = l.Format.FontText2.SpaceTop
			}
			t2W, t2H = c.Dimensions(l.Text2, &l.Format.FontText2)
			l.Height += t2H + p + l.Format.FontText2.SpaceBottom
		}
		//yInnerOffset = l.Format.Padding + max(yCaptionOffset, max(yText1Offset, yText2Offset))
		//yInnerOffset = l.Format.Padding + l.Height
		yInnerOffset = l.Height
		l.Width = max(cW, max(t1W, t2W)) + (2 * l.Format.Padding)
	}
	l.initVertical(c, yInnerOffset, defaultPadding, defaultBoxMargin)
	l.initHorizontal(c, yInnerOffset, defaultPadding, defaultBoxMargin)
}

func (l *LayoutElement) Center() {
	l.CenterX = l.X + (l.Width / 2)
	l.CenterY = l.Y + (l.Height / 2)
	l.centerHorizontalElems()
	l.centerVerticalElems()
}
