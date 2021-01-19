package entity

import (
	"github.com/milosgajdos/kraph/pkg/attrs"
)

// Entity is an arbitrary entity.
type Entity interface {
	// UID returns unique ID.
	UID() string
	// Attrs returns attributes.
	Attrs() attrs.Attrs
}
