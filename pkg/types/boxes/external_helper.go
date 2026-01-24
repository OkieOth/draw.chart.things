package boxes

import (
	"maps"
	"slices"

	"github.com/okieoth/draw.chart.things/pkg/types"
)

func (b *Boxes) MixinFormats(additional AdditionalFormats) {
	if len(additional.Formats) > 0 {
		if b.Formats == nil {
			b.Formats = make(map[string]Format)
		}
		maps.Copy(b.Formats, additional.Formats)
	}
	if len(additional.Images) > 0 {
		if b.Images == nil {
			b.Images = make(map[string]types.ImageDef)
		}
		maps.Copy(b.Images, additional.Images)
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

func (b *Boxes) MixinConnections(additional AdditionalConnections) {
	if len(additional.Formats) > 0 {
		maps.Copy(b.Formats, additional.Formats)
	}
	b.mixInConnectionsImpl(&b.Boxes, additional.Connections)
}
