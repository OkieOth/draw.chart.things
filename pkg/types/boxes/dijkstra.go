package boxes

import (
	"container/heap"
	"math"
)

type Graph map[string]ConnectionNode

// Priority queue item
type Item struct {
	node     string
	distance int
}

type PriorityQueue []Item

func (pq PriorityQueue) Len() int            { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool  { return pq[i].distance < pq[j].distance }
func (pq PriorityQueue) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x interface{}) { *pq = append(*pq, x.(Item)) }
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[:n-1]
	return item
}

func (doc *BoxesDocument) getBBox(startPointNodes []ConnectionNode) (minX, maxX, minY, maxY int) {
	minX = math.MaxInt
	minY = math.MaxInt
	for _, n := range startPointNodes {
		if n.X < minX {
			minX = n.X
		}
		if n.X > maxX {
			maxX = n.X
		}
		if n.Y < minY {
			minY = n.Y
		}
		if n.Y > maxY {
			maxY = n.Y
		}
	}
	return minX, maxX, minY, maxY
}

func (doc *BoxesDocument) createSourceCenterNode(nodeId string, startPointNodes []ConnectionNode) *ConnectionNode {
	minX, maxX, minY, maxY := doc.getBBox(startPointNodes)
	nX := minX + ((maxX - minX) / 2)
	nY := minY + ((maxY - minY) / 2)
	n := CreateConnectionNode(nX, nY)
	n.NodeId = &nodeId
	for _, startPointNode := range startPointNodes {
		newEdge := CreateConnectionEdge(startPointNode.X, startPointNode.Y, 1)
		newEdge.DestNodeId = startPointNode.NodeId
		n.Edges = append(n.Edges, newEdge)
	}
	return n
}

func (doc *BoxesDocument) createDestCenterNode(nodeId string, destPointNodes []ConnectionNode, graph *map[string]ConnectionNode) *ConnectionNode {
	minX, maxX, minY, maxY := doc.getBBox(destPointNodes)
	nX := minX + ((maxX - minX) / 2)
	nY := minY + ((maxY - minY) / 2)
	n := CreateConnectionNode(nX, nY)
	n.NodeId = &nodeId
	for _, destPointNode := range destPointNodes {
		newEdge := CreateConnectionEdge(nX, nY, 1)
		newEdge.DestNodeId = n.NodeId
		destPointNode.Edges = append(n.Edges, newEdge)
		(*graph)[*destPointNode.NodeId] = destPointNode
	}
	return n
}

func (doc *BoxesDocument) createGraph(sourceId, destId string) Graph {
	ret := make(map[string]ConnectionNode, 0)
	sourceNodes := make([]ConnectionNode, 0)
	destNodes := make([]ConnectionNode, 0)
	for _, n := range doc.ConnectionNodes {
		ret[*n.NodeId] = n
		if n.BoxId != nil {
			switch *n.BoxId {
			case sourceId:
				sourceNodes = append(sourceNodes, n)
			case destId:
				destNodes = append(destNodes, n)
			}
		}
	}
	// add center nodes to src and dest, to guarantee that the shortest path ...
	// ... always considers all connection possibilities
	ret[sourceId] = *doc.createSourceCenterNode(sourceId, sourceNodes)
	ret[destId] = *doc.createDestCenterNode(destId, destNodes, &ret)
	return ret
}

func (doc *BoxesDocument) DijkstraPath(source, target string) ([]ConnectionNode, int, bool) {
	dist := make(map[string]int)
	path := make([]ConnectionNode, 0)

	graph := doc.createGraph(source, target)

	for node := range graph {
		dist[node] = math.MaxInt
	}
	dist[source] = 0

	pq := &PriorityQueue{}
	heap.Init(pq)
	heap.Push(pq, Item{node: source, distance: 0})

	for pq.Len() > 0 {
		current := heap.Pop(pq).(Item)
		u := current.node

		if current.distance > dist[u] {
			continue
		}

		// Early exit if target reached
		if u == target {
			break
		}

		for _, edge := range graph[u].Edges {
			v := edge.DestNodeId
			if v == nil {
				continue
			}
			alt := dist[u] + edge.Weight
			if alt < dist[*v] {
				dist[*v] = alt
				path = append(path, graph[u])
				destNode := graph[*v]
				heap.Push(pq, Item{node: *destNode.NodeId, distance: alt})
			}
		}
	}

	// No path found
	if _, ok := dist[target]; !ok || dist[target] == math.MaxInt {
		return nil, 0, false
	}

	return path, dist[target], true
}
