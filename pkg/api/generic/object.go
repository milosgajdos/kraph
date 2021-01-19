package generic

import (
	"github.com/milosgajdos/kraph/pkg/api"
	"github.com/milosgajdos/kraph/pkg/metadata"
	"github.com/milosgajdos/kraph/pkg/uuid"
)

// Object is a generic API object
type Object struct {
	uid    uuid.UID
	name   string
	ns     string
	res    api.Resource
	links  map[string]api.Link
	olinks map[string]api.Link
	md     metadata.Metadata
}

// NewObject creates a new Object and returns it.
func NewObject(uid uuid.UID, name, ns string, res api.Resource, opts api.Options) *Object {
	md := opts.Metadata
	if md == nil {
		md = metadata.New()
	}

	return &Object{
		uid:    uid,
		name:   name,
		ns:     ns,
		res:    res,
		links:  make(map[string]api.Link),
		olinks: make(map[string]api.Link),
		md:     md,
	}
}

// UID returns the object uid.
func (o Object) UID() uuid.UID {
	return o.uid
}

// Name returns object name.
func (o Object) Name() string {
	return o.name
}

// Namespace returns object namespace.
func (o Object) Namespace() string {
	return o.ns
}

// Resource returns the API resource the object is an instance of.
func (o Object) Resource() api.Resource {
	return o.res
}

// link creates a new link to object to.
func (o *Object) link(to uuid.UID, opts api.LinkOptions) error {
	link, err := NewLink(o.uid, to, opts)
	if err != nil {
		return err
	}

	if _, ok := o.links[link.UID().String()]; !ok {
		o.links[link.UID().String()] = link
	}

	o.olinks[to.String()] = link

	return nil
}

// Link links the object to object to.
// If link merging is requested, the new link will contain
// all the metadata of the existing link with addition to the metadata
/// that are not in the original link. The original metadata are updated.
func (o *Object) Link(to uuid.UID, opts api.LinkOptions) error {
	l, ok := o.olinks[to.String()]
	if !ok {
		return o.link(to, opts)
	}

	if opts.Merge {
		if opts.Metadata != nil {
			for _, k := range opts.Metadata.Keys() {
				l.Metadata().Set(k, opts.Metadata.Get(k))
			}
		}
	}

	return nil
}

// Links returns a slice of all object links.
func (o Object) Links() []api.Link {
	var links []api.Link

	for _, link := range o.links {
		links = append(links, link)
	}

	return links
}

// Metadata returns object metadata.
func (o *Object) Metadata() metadata.Metadata {
	return o.md
}
