package entity

import "github.com/milosgajdos/kraph/store"

// Entity stores graph entity data
type Entity struct {
	metadata   store.Metadata
	attributes store.Attributes
}

// New creates new entity and returns it
func New(opts ...store.Option) store.Entity {
	o := store.NewOptions()
	for _, apply := range opts {
		apply(&o)
	}

	return &Entity{
		metadata:   o.Metadata,
		attributes: o.Attributes,
	}
}

// Properties returns entity attributes
func (e *Entity) Properties() store.Attributes {
	return e.attributes
}

// Metadata returns metadata
func (e *Entity) Metadata() store.Metadata {
	return e.metadata
}
