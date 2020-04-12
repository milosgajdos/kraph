package mock

import (
	"strings"

	"github.com/milosgajdos/kraph/api"
	"github.com/milosgajdos/kraph/api/generic"
)

type client struct{}

func NewClient() (api.Client, error) {
	return &client{}, nil
}

func (m *client) Discover() (api.API, error) {
	return NewAPI(), nil
}

func (m *client) Map(a api.API) (api.Top, error) {
	top := NewTop()

	for _, r := range a.Resources() {
		gv := strings.Join([]string{r.Group(), r.Version()}, "/")
		name := r.Name()
		if gvObject, ok := ObjectData[gv]; ok {
			kind := Resources[name]["kind"]
			ns := Resources[name]["ns"]
			if len(ns) == 0 {
				ns = api.NsNan
			}

			nsKind := strings.Join([]string{ns, kind}, "/")

			if names, ok := gvObject[nsKind]; ok {
				for _, name := range names {
					uid := strings.Join([]string{ns, kind, name}, "/")
					links := make(map[string]api.Relation)
					if rels, ok := ObjectLinks[uid]; ok {
						for obj, rel := range rels {
							links[obj] = generic.NewRelation(rel)
						}
					}
					object := generic.NewObject(name, kind, ns, generic.NewUID(uid), links)
					top.Add(object)
				}
			}
		}
	}
	return top, nil
}
