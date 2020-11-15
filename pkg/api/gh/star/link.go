package star

import (
	"github.com/milosgajdos/kraph/pkg/api"
	"github.com/milosgajdos/kraph/pkg/api/gen"
	"github.com/milosgajdos/kraph/pkg/uuid"
)

// Link is GitHub API link.
type Link struct {
	*gen.Link
}

// NewLink creates a new GitHub API link and returns it
func NewLink(from, to uuid.UID, opts api.LinkOptions) *Link {
	return &Link{
		Link: gen.NewLink(from, to, opts),
	}
}
