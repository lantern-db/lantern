package table

import (
	. "github.com/lantern-db/lantern/graph/model"
	"sort"
)

type EdgeTable struct {
	edges []Edge
}

func NewEdgeTableOf(edges []Edge) EdgeTable {
	t := EdgeTable{edges: edges}
	t.sort()
	return t
}

func NewEmptyEdgeTable() EdgeTable {
	return NewEdgeTableOf([]Edge{})
}

func (t *EdgeTable) sort() {
	sort.Slice(t.edges, func(i int, j int) bool { return t.edges[i].Expiration() > t.edges[j].Expiration() })
}

func (t *EdgeTable) Append(edge Edge) {
	t.edges = append(t.edges, edge)
	t.sort()
	t.flush()
}

func (t *EdgeTable) Weight() Weight {
	t.flush()
	w := Weight(0)
	for _, edge := range t.edges {
		w += edge.Weight()
	}
	return w
}

func (t *EdgeTable) flush() {
	for i, edge := range t.edges {
		if edge.Expiration().Dead() {
			t.edges = t.edges[:i]
		}
	}
}

func (t *EdgeTable) Len() int {
	t.flush()
	return len(t.edges)
}
