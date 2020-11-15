package gen

import (
	"github.com/milosgajdos/kraph/pkg/api"
	"github.com/milosgajdos/kraph/pkg/uuid"
)

// NewMockLink returns a new mock API Link
func NewMockLink(from, to uuid.UID, opts api.LinkOptions) api.Link {
	return NewLink(from, to, opts)
}
