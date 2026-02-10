package boxes

import "github.com/okieoth/draw.chart.things/pkg/types"

func (doc *BoxesDocument) IncludeOverlays(c types.TextDimensionCalculator) error {
	if len(doc.Overlays) > 0 {
		for overlayIndex := range doc.Overlays {
			overlay := doc.Overlays[overlayIndex]
			for layoutKey, _ := range overlay.Layouts {
				overlayEntry := overlay.Layouts[layoutKey]
				doc.initOverlayForLayout(&doc.Boxes, &overlay, layoutKey, &overlayEntry)
			}
		}
	}
	return nil
}

func (doc *BoxesDocument) initOverlayForLayout(l *LayoutElement, overlay *DocOverlay, entryKey string, overlayEntry *OverlayEntry) bool {
	if l.Id != "" && (l.Id == entryKey || l.Caption == entryKey) {
		// found the related layout element
		doc.initOverlayXY(l, overlay, overlayEntry)
		doc.initOverlayRadius(overlay, overlayEntry)
		doc.initOverlayFormat(overlay, overlayEntry)
		return true
	}
	if doc.initOverlaysFromLayoutCont(l.Horizontal, overlay, entryKey, overlayEntry) {
		return true
	}
	return doc.initOverlaysFromLayoutCont(l.Vertical, overlay, entryKey, overlayEntry)
}

func (doc *BoxesDocument) initOverlayXY(l *LayoutElement, overlay *DocOverlay, overlayEntry *OverlayEntry) {
	overlayEntry.X = l.CenterX
	overlayEntry.Y = l.CenterY
	if overlay.CenterXOffset != 0 {
		diff := float64((l.CenterX - l.X)) * overlay.CenterXOffset
		overlayEntry.X += int(diff)
	}
	if overlay.CenterYOffset != 0 {
		diff := float64((l.CenterY - l.Y)) * overlay.CenterYOffset
		overlayEntry.Y += int(diff)
	}
}

func (doc *BoxesDocument) initOverlayRadius(overlay *DocOverlay, overlayEntry *OverlayEntry) {
	if overlay.RefValue == 0.0 {
		overlayEntry.Radius = overlayEntry.Value
	} else {
		overlayEntry.Radius = (overlayEntry.Value / overlay.RefValue) * overlay.RadiusDefs.Max
	}
	if overlay.RadiusDefs.Min > 0 && overlayEntry.Radius < overlay.RadiusDefs.Min {
		overlayEntry.Radius = overlay.RadiusDefs.Min
	}
	if overlay.RadiusDefs.Max > 0 && overlayEntry.Radius > overlay.RadiusDefs.Max {
		overlayEntry.Radius = overlay.RadiusDefs.Max
	}
}

func (doc *BoxesDocument) getDefaultOverlayFormat() BoxFormat {
	return BoxFormat{
		Line: types.InitLineDef2("blue", 0.6),
		Fill: types.InitFillDef2("blue", 0.4),
	}
}

func (doc *BoxesDocument) initOverlayFormat(overlay *DocOverlay, overlayEntry *OverlayEntry) {
	var format BoxFormat
	if overlay.Formats.Default != "" {
		f, ok := doc.Formats[overlay.Formats.Default]
		if !ok {
			f = doc.getDefaultOverlayFormat()
		} else {
			f = doc.getDefaultOverlayFormat()
		}
		format = f
	} else {
		format = doc.getDefaultOverlayFormat()
	}

	if len(overlay.Formats.Gradations) > 0 {
		var lastFormat, format2Use string
		for i := range overlay.Formats.Gradations {
			g := overlay.Formats.Gradations[i]
			if g.Limit > overlayEntry.Value {
				format2Use = lastFormat
				break
			}
			lastFormat = g.Format
		}
		if format2Use != "" {
			if f, ok := doc.Formats[format2Use]; ok {
				format = f
			}
		}
	}
	overlayEntry.Format = format
}

func (doc *BoxesDocument) initOverlaysFromLayoutCont(cont *LayoutElemContainer, overlay *DocOverlay, entryKey string, overlayEntry *OverlayEntry) bool {
	if cont != nil {
		for i := range cont.Elems {
			if doc.initOverlayForLayout(&cont.Elems[i], overlay, entryKey, overlayEntry) {
				return true
			}
		}
	}
	return false
}
