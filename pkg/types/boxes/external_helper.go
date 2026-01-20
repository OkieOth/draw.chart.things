package boxes

import (
	"maps"
	"slices"
)

func (b *Boxes) MixinFormats(additional AdditionalFormats) {
	if len(additional.Formats) > 0 {
		maps.Copy(b.Formats, additional.Formats)
	}
	for _, i := range additional.Images {
		b.Images = append(b.Images, i)
	}
}

func hasConnection(connections []Connection, destId string) bool {
	return slices.ContainsFunc(connections, func(c Connection) bool {
		return c.DestId == destId
	})
}

func (b *Boxes) mixInConnectionsImplCont(cont []Layout, additional map[string]ConnectionCont) {
	for i := range len(cont) {
		b.mixInConnectionsImpl(&cont[i], additional)
	}
}

func (b *Boxes) mixInConnectionsImpl(l *Layout, additional map[string]ConnectionCont) {
	if l.Id != "" {
		if cc, ok := additional[l.Id]; ok {
			for _, c := range cc.Connections {
				if !hasConnection(l.Connections, c.DestId) {
					l.Connections = append(l.Connections, c)
				}
			}
		}
	}
	b.mixInConnectionsImplCont(l.Horizontal, additional)
	b.mixInConnectionsImplCont(l.Vertical, additional)
}

func (b *Boxes) MixinConnections(additional map[string]ConnectionCont) {
	b.mixInConnectionsImpl(&b.Boxes, additional)
}
