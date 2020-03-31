package k8s

import (
	"strings"

	"github.com/milosgajdos/kraph/api"
	"github.com/milosgajdos/kraph/query"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
)

// ObjRef is object reference used when linking API objects
type ObjRef struct {
	name string
	kind string
}

// Name of the API object this link points to
func (r ObjRef) Name() string {
	return r.name
}

// Kind of the API object this link points to
func (r ObjRef) Kind() string {
	return r.kind
}

// Relation is link relation
type Relation struct {
	r string
}

// String returns relation description
func (r *Relation) String() string {
	return r.r
}

// Link defines API object relation
type Link struct {
	objRef   *ObjRef
	relation *Relation
}

// ObjRef returns link object reference
func (l *Link) To() api.ObjRef {
	return l.objRef
}

// Relation returns the type of link relation
func (r *Link) Relation() api.Relation {
	return r.relation
}

// Object is API object i.e. instance of API resource
type Object struct {
	obj   unstructured.Unstructured
	links map[string]map[string]*ObjRef
}

// NewObject returns new kubernetes API object
func NewObject(obj unstructured.Unstructured) *Object {
	links := make(map[string]map[string]*ObjRef)

	for _, ref := range obj.GetOwnerReferences() {
		objRef := &ObjRef{
			name: ref.Name,
			kind: ref.Kind,
		}
		key := objRef.name + "/" + objRef.kind
		if links[key]["isOwned"] == nil {
			links[key] = make(map[string]*ObjRef)
		}
		links[key]["isOwned"] = objRef
	}

	return &Object{
		obj:   obj,
		links: links,
	}
}

// Name returns resource nam
func (o Object) Name() string {
	return strings.ToLower(o.obj.GetName())
}

// Kind returns object kind
func (o Object) Kind() string {
	return strings.ToLower(o.obj.GetKind())
}

// Namespace returns the namespace
func (o Object) Namespace() string {
	return strings.ToLower(o.obj.GetNamespace())
}

// UID returns object UID
func (o Object) UID() types.UID {
	kind := strings.ToLower(o.obj.GetKind())
	uid := o.obj.GetUID()
	if uid == "" {
		uid = types.UID(kind + "-" + strings.ToLower(o.obj.GetName()))
	}

	return uid
}

// Link links object to given ObjRef
func (o *Object) Link(ref api.ObjRef, rel api.Relation) error {
	objRef := &ObjRef{
		name: ref.Name(),
		kind: ref.Kind(),
	}

	key := objRef.name + "/" + objRef.kind
	if o.links[key][rel.String()] == nil {
		o.links[key] = make(map[string]*ObjRef)
	}

	if _, ok := o.links[key][rel.String()]; !ok {
		o.links[key][rel.String()] = objRef
	}

	return nil
}

// Links returns all object links
func (o Object) Links() []api.Link {
	var links []api.Link

	for _, rels := range o.links {
		for rel, obj := range rels {
			link := &Link{
				objRef: obj,
				relation: &Relation{
					r: rel,
				},
			}
			links = append(links, link)
		}
	}

	return links
}

// Raw returns the raw API bjoect
func (o Object) Raw() interface{} {
	return o.obj
}

// Resource is API resource
type Resource struct {
	ar metav1.APIResource
	gv schema.GroupVersion
}

// Name returns the name of the resource
func (r Resource) Name() string {
	return r.ar.Name
}

// Kind returns resource kind
func (r Resource) Kind() string {
	return r.ar.Kind
}

// Group returns the API group of the resource
func (r Resource) Group() string {
	return r.gv.Group
}

// Version returns the version of the resource
func (r Resource) Version() string {
	return r.gv.Version
}

// Namespaced returns true if the resource is namespaced
func (r Resource) Namespaced() bool {
	return r.ar.Namespaced
}

// Paths returns all possible variations of the resource paths
func (r Resource) Paths() []string {
	// WTF: SingularName is often an empty string!
	// TODO: figure this out; but for now let's set it to Kind
	singularName := r.ar.SingularName
	if singularName == "" {
		singularName = r.ar.Kind
	}
	resNames := []string{strings.ToLower(r.ar.Name), strings.ToLower(singularName)}
	resNames = append(resNames, r.ar.ShortNames...)

	var names []string
	for _, name := range resNames {
		names = append(names,
			name,
			strings.Join([]string{name, r.gv.Group}, "/"),
			strings.Join([]string{name, r.gv.Group, r.gv.Version}, "/"),
		)
	}

	return names
}

// API is API resource map
type API struct {
	resources []Resource
	// resourceMap serves as an index into APIs
	resourceMap map[string][]Resource
}

// Resources returns all API resources matching the given query
func (a *API) Resources(opts ...query.Option) []api.Resource {
	query := query.NewOptions()
	for _, apply := range opts {
		apply(&query)
	}

	var apiResources []api.Resource

	if a.resourceMap == nil {
		return apiResources
	}

	name := strings.ToLower(query.Name)
	group := query.Group
	version := query.Version

	// try pulling out the results from indexed entries
	if len(name) > 0 {
		if len(group) > 0 {
			if len(version) > 0 {
				return res2APIres(a.resourceMap[strings.Join([]string{name, group, version}, "/")])
			}
			return res2APIres(a.resourceMap[strings.Join([]string{name, group}, "/")])
		}
		// both group and version have 0 length; return the resources indexed by name
		if len(version) == 0 {
			return res2APIres(a.resourceMap[name])
		}
	}

	return a.lookupGV(group, version)
}

// lookupGV searches all resources matching the given name and/or version
func (a *API) lookupGV(group, version string) []api.Resource {
	var resources []api.Resource

	for _, r := range a.resources {
		if len(group) == 0 || group == r.gv.Group {
			if len(version) == 0 || version == r.gv.Version {
				resources = append(resources, r)
			}
		}
	}

	return resources
}

// res2APIres "converts" resources into API resources
func res2APIres(rx []Resource) []api.Resource {
	ax := make([]api.Resource, len(rx))

	for i, r := range rx {
		ax[i] = r
	}

	return ax
}
