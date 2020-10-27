package memory

import (
	"github.com/milosgajdos/kraph/pkg/attrs"
	"github.com/milosgajdos/kraph/pkg/store/entity"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding"
)

// Line implements graph.WeightedLine
type Line struct {
	*entity.Edge
	id    int64
	from  *Node
	to    *Node
	dotid string
}

// NewLine creates a new line and returns it
func NewLine(id int64, uid, dotid string, from, to *Node, opts ...entity.Option) *Line {
	edgeOpts := entity.NewOptions()
	for _, apply := range opts {
		apply(&edgeOpts)
	}

	edge := entity.NewEdge(uid, from.Node, to.Node, opts...)

	return &Line{
		Edge:  edge,
		id:    id,
		from:  from,
		to:    to,
		dotid: dotid,
	}
}

// ID is line ID
func (l *Line) ID() int64 {
	return l.id
}

// From returns the from node of the first non-nil edge, or nil.
func (l *Line) From() graph.Node {
	return l.from
}

// To returns the to node of the first non-nil edge, or nil.
func (l *Line) To() graph.Node {
	return l.to
}

// Weight returns edge weight
func (l Line) Weight() float64 {
	return l.Options().Weight
}

// ReversedLine returns a new line with end points of the pair swapped
func (l *Line) ReversedLine() graph.Line {
	opts := []entity.Option{
		entity.Attrs(l.Options().Attrs),
		entity.Metadata(l.Options().Metadata),
		entity.Weight(l.Options().Weight),
	}

	edge := entity.NewEdge(l.UID(), l.to.Node, l.from.Node, opts...)

	return &Line{
		Edge:  edge,
		id:    l.id,
		from:  l.to,
		to:    l.from,
		dotid: l.dotid,
	}
}

// DOTID returns the edge's DOT ID.
func (l *Line) DOTID() string {
	return l.dotid
}

// SetDOTID sets the edge's DOT ID.
func (l *Line) SetDOTID(id string) {
	l.Edge.Attrs().Set("dotid", id)
	l.dotid = id
}

// Attributes implements store.DOTAttrs
func (l *Line) Attributes() []encoding.Attribute {
	if a, ok := l.Attrs().(attrs.DOT); ok {
		return a.Attributes()
	}

	attrs := make([]encoding.Attribute, len(l.Attrs().Keys()))

	i := 0
	for _, k := range l.Attrs().Keys() {
		attrs[i] = encoding.Attribute{
			Key:   k,
			Value: l.Attrs().Get(k),
		}
		i++
	}

	return attrs
}
