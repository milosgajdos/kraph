package kraph

import (
	"fmt"

	"github.com/milosgajdos/kraph/api"
	"github.com/milosgajdos/kraph/query"
	"github.com/milosgajdos/kraph/store"
)

type kraph struct {
	store store.Store
}

// New creates new kraph and returns it
func New(opts ...Option) (Kraph, error) {
	o := NewOptions()
	for _, apply := range opts {
		apply(&o)
	}

	return &kraph{
		store: o.Store,
	}, nil
}

// linkObject links obj to all of its neighbours and sets their relation to rel.
func (k *kraph) linkObjects(obj api.Object, rel api.Relation, neighbs []api.Object) error {
	from, err := k.store.Add(obj)
	if err != nil {
		return err
	}

	for _, o := range neighbs {
		to, err := k.store.Add(o)
		if err != nil {
			return err
		}

		attrs := store.NewAttributes()
		attrs.Set("relation", rel.String())
		// TODO: this is set to default weight for now
		//attrs.Set("weight", fmt.Sprintf("%f", store.DefaultEdgeWeight))

		if e := k.store.Edge(from.ID(), to.ID()); e == nil {
			e, err = k.store.Link(from, to, store.EntAttrs(attrs))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// buildGraph builds a graph from given topology and returns it.
func (k *kraph) buildGraph(top api.Top) (store.Graph, error) {
	for _, object := range top.Objects() {
		if len(object.Links()) == 0 {
			_, err := k.store.Add(object)
			if err != nil {
				return nil, fmt.Errorf("error adding node: %w", err)
			}
			continue
		}
		for _, link := range object.Links() {
			query := []query.Option{
				query.UID(link.To().String()),
			}
			objs, err := top.Get(query...)
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

// Build builds a graph of API object using the client and the graph store.
func (k *kraph) Build(client api.Client) (store.Graph, error) {
	// TODO: reset the graph before building
	// This will allow to run Build multiple times
	// each time building the graph from scratch
	api, err := client.Discover()
	if err != nil {
		return nil, fmt.Errorf("failed discovering API: %w", err)
	}

	top, err := client.Map(api)
	if err != nil {
		return nil, fmt.Errorf("failed mapping API: %w", err)
	}

	return k.buildGraph(top)
}

// Store returns kraph stor
func (k *kraph) Store() store.Store {
	return k.store
}
