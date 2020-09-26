package entity

import (
	"github.com/milosgajdos/kraph/store"
	"gonum.org/v1/gonum/graph/encoding"
)

// Entity stores graph entity data
type Entity struct {
	id         string
	name       string
	metadata   store.Metadata
	attributes store.Attrs
}

// New creates new entity and returns it
func New(id, name string, opts ...store.Option) *Entity {
	o := store.NewOptions()
	for _, apply := range opts {
		apply(&o)
	}

	return &Entity{
		id:         id,
		name:       name,
		metadata:   o.Metadata,
		attributes: o.Attributes,
	}
}

// ID is entity unique ID
func (e Entity) ID() string {
	return e.id
}

// Name is entity name
func (e Entity) Name() string {
	return e.name
}

// Attributes returns entity attributes
func (e *Entity) Attrs() store.Attrs {
	return e.attributes
}

// Metadata returns entity metadata
func (e *Entity) Metadata() store.Metadata {
	return e.metadata
}

// Attributes returns entity attributes encoded as required by
// https://godoc.org/gonum.org/v1/gonum/graph/encoding/dot
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
