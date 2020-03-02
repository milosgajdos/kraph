package k8s

import (
	"context"
	"fmt"
	"sync"

	"github.com/milosgajdos/kraph/api"
	"github.com/milosgajdos/kraph/query"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
)

type client struct {
	// disc is kubernetes discovery client
	disc discovery.DiscoveryInterface
	// dyn is kubernetes dynamic client
	dyn dynamic.Interface
	// ctx is client context
	ctx context.Context
	// m is API map /ns/kind/name/object
	m map[string]map[string]map[string]*Object
	// opts are client options
	opts Options
}

// NewClient returns new kubernetes API client
func NewClient(disc discovery.DiscoveryInterface, dyn dynamic.Interface, ctx context.Context, opts ...Option) *client {
	copts := Options{}
	for _, apply := range opts {
		apply(&copts)
	}

	return &client{
		disc: disc,
		dyn:  dyn,
		ctx:  ctx,
		m:    make(map[string]map[string]map[string]*Object),
		opts: copts,
	}
}

// Options returns client options
func (k *client) Options() Options {
	return k.opts
}

// Discover discovers kubernetes API and returns them
// It returns error if it fails to read the resources of if it fails to parse their versions
func (k *client) Discover() (api.API, error) {
	srvPrefResList, err := k.disc.ServerPreferredResources()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch API groups: %w", err)
	}

	api := &API{
		resourceMap: make(map[string][]Resource),
	}

	for _, srvPrefRes := range srvPrefResList {
		gv, err := schema.ParseGroupVersion(srvPrefRes.GroupVersion)
		if err != nil {
			return nil, fmt.Errorf("failed parsing %s into GroupVersion: %w", srvPrefRes.GroupVersion, err)
		}

		for _, apiResource := range srvPrefRes.APIResources {
			if !stringIn("list", apiResource.Verbs) {
				continue
			}

			resource := Resource{
				ar: apiResource,
				gv: gv,
			}

			api.resources = append(api.resources, resource)
			for _, path := range resource.Paths() {
				api.resourceMap[path] = append(api.resourceMap[path], resource)
			}
		}
	}

	return api, nil
}

// API discovery results
type result struct {
	api   string
	items []unstructured.Unstructured
	err   error
}

// processResults processes API call request results
// It builds undirected weighted graph from the received results
func (k *client) processResults(resChan <-chan result, doneChan chan struct{}, errChan chan<- error) {
	var err error
	for result := range resChan {
		if result.err != nil {
			err = result.err
			close(doneChan)
			break
		}

		for _, res := range result.items {
			ns := res.GetNamespace()
			if ns == "" {
				ns = "none"
			}

			if k.m[ns] == nil {
				k.m[ns] = make(map[string]map[string]*Object)
			}

			obj := &Object{
				obj: res,
			}

			kind := obj.Kind()
			name := obj.Name()

			if k.m[ns][kind] == nil {
				k.m[ns][kind] = make(map[string]*Object)
			}

			k.m[ns][kind][name] = obj
		}
	}

	errChan <- err
}

// Map builds a map of API resources in a given client namespace
// If the namespace is empty it queries API groups across all namespaces.
// It returns error if any of the API calls fails with error.
func (k *client) Map(a api.API) error {
	// TODO: we should take into account the client context
	// when firing goroutines and waiting for the results
	var wg sync.WaitGroup

	resChan := make(chan result, 250)
	doneChan := make(chan struct{})

	for _, resource := range a.Resources() {
		// if all namespaces are scanned and the API resource is namespaced, skip
		if k.opts.Namespace != "" && !resource.Namespaced() {
			continue
		}

		gvResClient := k.dyn.Resource(schema.GroupVersionResource{
			Group:    resource.Group(),
			Version:  resource.Version(),
			Resource: resource.Name(),
		})

		var client dynamic.ResourceInterface
		switch k.opts.Namespace {
		case "":
			client = gvResClient
		default:
			client = gvResClient.Namespace(k.opts.Namespace)
		}

		wg.Add(1)
		go func(r api.Resource) {
			defer wg.Done()
			var cont string
			for {
				res, err := client.List(metav1.ListOptions{
					Limit:    100,
					Continue: cont,
				})
				select {
				case resChan <- result{api: r.Name(), items: res.Items, err: err}:
				case <-doneChan:
					return
				}
				cont = res.GetContinue()
				if cont == "" {
					break
				}
			}
		}(resource)
	}

	errChan := make(chan error, 1)
	go k.processResults(resChan, doneChan, errChan)

	wg.Wait()
	close(resChan)

	err := <-errChan

	return err
}

func (k *client) getNamespaceKindObjects(ns, kind string, q query.Options) ([]api.Object, error) {
	var objects []api.Object

	for name, _ := range k.m[ns][kind] {
		if q.Name == "*" || q.Name == name {
			objects = append(objects, k.m[ns][kind][name])
		}
	}

	return objects, nil
}

func (k *client) getNamespaceObjects(ns string, q query.Options) ([]api.Object, error) {
	var objects []api.Object

	if q.Kind != "*" {
		return k.getNamespaceKindObjects(ns, q.Kind, q)
	}

	for kind, _ := range k.m[ns] {
		objs, err := k.getNamespaceKindObjects(ns, kind, q)
		if err != nil {
			return nil, err
		}
		objects = append(objects, objs...)
	}

	return objects, nil
}

func (k *client) getAllNamespaceObjects(q query.Options) ([]api.Object, error) {
	var objects []api.Object

	for ns, _ := range k.m {
		objs, err := k.getNamespaceObjects(ns, q)
		if err != nil {
			return nil, err
		}
		objects = append(objects, objs...)
	}

	return objects, nil
}

// Get queries the mapped API objects and returns them
func (k *client) Get(opts ...query.Option) ([]api.Object, error) {
	query := query.NewOptions()
	for _, apply := range opts {
		apply(&query)
	}

	var objects []api.Object

	if query.Namespace == "*" {
		return k.getAllNamespaceObjects(query)
	}

	if query.Namespace == "" {
		objs, err := k.getNamespaceObjects("none", query)
		if err != nil {
			return nil, err
		}
		objects = append(objects, objs...)
		return objects, nil
	}

	for ns, _ := range k.m {
		objs, err := k.getNamespaceObjects(ns, query)
		if err != nil {
			return nil, err
		}
		objects = append(objects, objs...)
	}

	return objects, nil
}
