package star

import (
	"github.com/milosgajdos/kraph/pkg/api"
	"github.com/milosgajdos/kraph/pkg/api/gen"
	"github.com/milosgajdos/kraph/pkg/uuid"
)

type link struct {
	uid  uuid.UID
	opts api.LinkOptions
}

// Object is GitHub API starred repository object
type Object struct {
	*gen.Object
}

// NewObject returns a new GitHub starred repo API object
func NewObject(uid uuid.UID, name, ns string, res api.Resource, opts api.Options, links []link) *Object {
	obj := &Object{
		Object: gen.NewObject(uid, name, ns, res, opts),
	}

	for _, link := range links {
		obj.Link(link.uid, link.opts)
	}

	return obj
}
