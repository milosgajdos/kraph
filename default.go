package kraph

import (
	"fmt"

	"github.com/milosgajdos/kraph/pkg/api"
	"github.com/milosgajdos/kraph/pkg/attrs"
	"github.com/milosgajdos/kraph/pkg/metadata"
	"github.com/milosgajdos/kraph/pkg/query"
	"github.com/milosgajdos/kraph/pkg/store"
)

type kraph struct {
	store store.Store
	opts  *Options
}

// New creates a new kraph and returns it.
func New(store store.Store, opts ...Option) (Kraph, error) {
	o, err := NewOptions()
	if err != nil {
		return nil, err
	}

	for _, apply := range opts {
		apply(o)
	}

	return &kraph{
		store: store,
		opts:  o,
	}, nil
}

// linkObject links obj to all of its neighbours.
func (k *kraph) linkObjects(obj api.Object, link api.Link, neighbs []api.Object) error {
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
		attrs.Set("weight", fmt.Sprintf("%f", store.DefaultWeight))

		if rel := link.Metadata().Get("relation"); rel != nil {
			if r, ok := rel.(string); ok {
				attrs.Set("relation", r)
			}
		}

		opts := store.LinkOptions{
			Attrs:    attrs,
			Metadata: link.Metadata(),
			Weight:   store.DefaultWeight,
		}

		if _, err := k.store.Graph().Link(from, to, opts); err != nil {
			return err
		}
	}

	return nil
}

// skip skips adding API objects to graph based on defined filters.
func skip(object api.Object, filters ...Filter) bool {
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
func (k *kraph) buildGraph(top api.Top, filters ...Filter) error {
	for _, object := range top.Objects() {
		if skip(object, filters...) {
			continue
		}

		if len(object.Links()) == 0 {
			if _, err := k.store.Add(object, store.AddOptions{}); err != nil {
				return fmt.Errorf("error adding node: %w", err)
			}
			continue
		}

		for _, link := range object.Links() {
			uid := link.To()

			q := query.Build().UID(uid, query.UIDEqFunc(uid))

			objs, err := top.Get(q)
			if err != nil {
				return err
			}

			if err := k.linkObjects(object, link, objs); err != nil {
				return err
			}
		}
	}

	return nil
}

// Build builds a graph of API objects for the source using the client.
func (k *kraph) Build(client api.Client, filters ...Filter) error {
	api, err := client.Discover()
	if err != nil {
		return fmt.Errorf("failed discovering API: %w", err)
	}

	top, err := client.Map(api)
	if err != nil {
		return fmt.Errorf("failed mapping API: %w", err)
	}

	return k.buildGraph(top, filters...)
}

// Store returns kraph store.
func (k *kraph) Store() store.Store {
	return k.store
}

// Metadata returns kraph metadata
func (k *kraph) Metadata() metadata.Metadata {
	return k.opts.Metadata
}
