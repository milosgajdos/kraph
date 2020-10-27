package gen

import (
	"github.com/milosgajdos/kraph/pkg/api"
	"github.com/milosgajdos/kraph/pkg/uuid"
)

// NewMockLink returns a new mock API Link
func NewMockLink(from, to, rel string) api.Link {
	return NewLink(uuid.NewFromString(from), uuid.NewFromString(to), NewRelation(rel))
}
