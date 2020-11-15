package gen

import (
	"github.com/milosgajdos/kraph/pkg/api"
	"github.com/milosgajdos/kraph/pkg/metadata"
	"github.com/milosgajdos/kraph/pkg/uuid"
)

// Link links API object to another API object
type Link struct {
	uid  uuid.UID
	from uuid.UID
	to   uuid.UID
	opts api.LinkOptions
}

// NewLink returns a new link between API objects
func NewLink(from, to uuid.UID, opts api.LinkOptions) *Link {
	var uid uuid.UID = uuid.New()
	if opts.UID != nil {
		uid = opts.UID
	}

	return &Link{
		uid:  uid,
		from: from,
		to:   to,
		opts: opts,
	}
}

// UID returns link uid
func (l Link) UID() uuid.UID {
	return l.uid
}

// From returns the linking object
func (l Link) From() uuid.UID {
	return l.from
}

// To returns the linked object
func (l Link) To() uuid.UID {
	return l.to
}

// Metadata returns the link metadata
func (l Link) Metadata() metadata.Metadata {
	return l.opts.Metadata
}
