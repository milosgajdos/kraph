package kraph

import (
	"fmt"

	"github.com/milosgajdos/kraph/pkg/api"
	"github.com/milosgajdos/kraph/pkg/attrs"
	"github.com/milosgajdos/kraph/pkg/query"
	"github.com/milosgajdos/kraph/pkg/store"
)

type kraph struct {
	store store.Store
}

// New creates new kraph and returns it
func New(opts ...Option) (Kraph, error) {
	o, err := NewOptions()
	if err != nil {
		return nil, err
	}

	for _, apply := range opts {
		apply(o)
	}

	return &kraph{
		store: o.Store,
	}, nil
}

// linkObject links obj to all of its neighbours and sets their relation to rel.
func (k *kraph) linkObjects(obj api.Object, rel api.Relation, neighbs []api.Object) error {
	from, err := k.store.Add(obj, store.AddOptions{})
	if err != nil {
		return err
	}

	for _, o := range neighbs {
		to, err := k.store.Add(o, store.AddOptions{})
		if err != nil {
			return err
		}

		attrs := attrs.New()
		if rel.String() != "" {
			attrs.Set("relation", rel.String())
		}
		attrs.Set("weight", fmt.Sprintf("%f", store.DefaultWeight))

		opts := store.LinkOptions{Attrs: attrs, Weight: store.DefaultWeight}
		if _, err := k.store.Link(from, to, opts); err != nil {
			return err
		}
	}

	return nil
}

// skipGraph skips adding API objects into graph based on defined filters.
func skipGraph(object api.Object, filters ...Filter) bool {
	if len(filters) == 0 {
		return false
	}

	for _, filter := range filters {
		if filter(object) {
			return false
		}
	}

	return true
}

// buildGraph builds a graph from given topology and returns it.
func (k *kraph) buildGraph(top api.Top, filters ...Filter) (store.Graph, error) {
	for _, object := range top.Objects() {
		if skipGraph(object, filters...) {
			continue
		}

		if len(object.Links()) == 0 {
			if _, err := k.store.Add(object, store.AddOptions{}); err != nil {
				return nil, fmt.Errorf("error adding node: %w", err)
			}
			continue
		}

		for _, link := range object.Links() {
			uid := link.To()

			q := query.Build().UID(uid, query.UIDEqFunc(uid))

			objs, err := top.Get(q)
			if err != nil {
				return nil, err
			}

			if err := k.linkObjects(object, link.Relation(), objs); err != nil {
				return nil, err
			}
		}
	}

	return k.store, nil
}

// Build builds a graph of API object using the client and returns it.
func (k *kraph) Build(client api.Client, filters ...Filter) (store.Graph, error) {
	// TODO: reset the graph before building
	// This will allow to run k.Build multiple times
	// each time building the graph from scratch
	api, err := client.Discover()
	if err != nil {
		return nil, fmt.Errorf("failed discovering API: %w", err)
	}

	top, err := client.Map(api)
	if err != nil {
		return nil, fmt.Errorf("failed mapping API: %w", err)
	}

	return k.buildGraph(top, filters...)
}

// Store returns kraph stor
func (k *kraph) Store() store.Store {
	return k.store
}
