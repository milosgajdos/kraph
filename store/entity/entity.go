package entity

import (
	"github.com/milosgajdos/kraph/store"
	"gonum.org/v1/gonum/graph/encoding"
)

// Entity stores graph entity data
type Entity struct {
	metadata   store.Metadata
	attributes store.Attrs
}

// New creates new entity and returns it
func New(opts ...store.Option) store.Entity {
	o := store.NewOptions()
	for _, apply := range opts {
		apply(&o)
	}

	return &Entity{
		metadata:   o.Metadata,
		attributes: o.EntAttrs,
	}
}

// Attributes returns entity attributes
func (e *Entity) Attrs() store.Attrs {
	return e.attributes
}

// Metadata returns metadata
func (e *Entity) Metadata() store.Metadata {
	return e.metadata
}

// Attributes returns all attributes encoded as required by https://godoc.org/gonum.org/v1/gonum/graph/encoding/dot
func (e Entity) Attributes() []encoding.Attribute {
	keys := e.attributes.Keys()

	attrs := make([]encoding.Attribute, len(keys))

	i := 0
	for _, k := range keys {
		attrs[i] = encoding.Attribute{
			Key:   k,
			Value: e.attributes.Get(k),
		}
		i++
	}

	return attrs
}
