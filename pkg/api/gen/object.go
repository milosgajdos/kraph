package gen

import (
	"reflect"

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
	olinks map[string]struct{}
	opts   api.Options
}

// NewObject creates a new Object and returns it.
func NewObject(uid uuid.UID, name, ns string, res api.Resource, opts api.Options) *Object {
	return &Object{
		uid:    uid,
		name:   name,
		ns:     ns,
		res:    res,
		links:  make(map[string]api.Link),
		olinks: make(map[string]struct{}),
		opts:   opts,
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

// link creates a new link and adds it to the the objects links.
func (o *Object) link(to uuid.UID, opts api.LinkOptions) error {
	// link: o -> to
	link := NewLink(o.uid, to, opts)

	if _, ok := o.links[link.UID().String()]; !ok {
		o.links[link.UID().String()] = link
	}

	o.olinks[to.String()] = struct{}{}

	return nil
}

// Link links the object to another object.
// It creates a bidirectional link b/w the linked objects i.e.
//   * it links the object to object to
//   * it links the to object to to o object
// Two links are considered the same if both of these conditions are satisifed:
//   * they link to the same object
//   * their metadata are the same
func (o *Object) Link(to uuid.UID, opts api.LinkOptions) error {
	// check if the link to object "to" already exists
	if _, ok := o.olinks[to.String()]; !ok {
		return o.link(to, opts)
	}

	if opts.Multi {
		// if the metadata of the existing link dont match opts.Metadata
		// create a new link; links with the same Metadata would be redundant
		for _, l := range o.links {
			if l.To().String() == to.String() {
				if !reflect.DeepEqual(l.Metadata(), opts.Metadata) {
					if err := o.link(to, opts); err != nil {
						return err
					}
				}
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
	return o.opts.Metadata
}
