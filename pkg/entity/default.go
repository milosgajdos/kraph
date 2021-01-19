package entity

import (
	"github.com/milosgajdos/kraph/pkg/attrs"
	"github.com/milosgajdos/kraph/pkg/uuid"
)

type entity struct {
	uid   string
	attrs attrs.Attrs
}

// NewWithUID creates a new entity with given UID and returns it.
func NewWithUID(uid string, opts ...Option) (*entity, error) {
	if uid == "" {
		uid = uuid.New().String()
	}

	eopts := NewOptions()
	for _, apply := range opts {
		apply(&eopts)
	}

	return &entity{
		uid:   uid,
		attrs: eopts.Attrs,
	}, nil
}

// New creates a new entity and returns it.
func New(opts ...Option) (*entity, error) {
	uid := uuid.New().String()

	return NewWithUID(uid, opts...)
}

// UID returns entity UID.
func (e entity) UID() string {
	return e.uid
}

// Attrs returns entity attributes.
func (e *entity) Attrs() attrs.Attrs {
	return e.attrs
}
