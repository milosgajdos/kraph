package star

import (
	"github.com/milosgajdos/kraph/pkg/api"
	"github.com/milosgajdos/kraph/pkg/api/generic"
	"github.com/milosgajdos/kraph/pkg/uuid"
)

// Link is GitHub API link.
type Link struct {
	*generic.Link
}

// NewLink creates a new GitHub API link and returns it
func NewLink(from, to uuid.UID, opts api.LinkOptions) (*Link, error) {
	l, err := generic.NewLink(from, to, opts)
	if err != nil {
		return nil, err
	}

	return &Link{
		Link: l,
	}, nil
}
