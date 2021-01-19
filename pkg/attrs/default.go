package attrs

import "gonum.org/v1/gonum/graph/encoding"

// attrs are graph attributes
type attrs map[string]string

// NewAttributes creates new attributes and returns it.
func New() *attrs {
	attrs := make(attrs)

	return &attrs
}

// NewCopyFrom copies attributes from a and returns it.
func NewCopyFrom(a Attrs) *attrs {
	attrs := make(attrs)

	if a != nil {
		for _, k := range a.Keys() {
			attrs[k] = a.Get(k)
		}
	}

	return &attrs
}

// Keys returns all attribute keys
func (a attrs) Keys() []string {
	keys := make([]string, len(a))

	i := 0
	for key := range a {
		keys[i] = key
		i++
	}

	return keys
}

// Get reads an attribute value for the given key and returns it.
// It returns an empty string if the attribute was not found.
func (a attrs) Get(key string) string {
	return a[key]
}

// Set sets an attribute to the given value
func (a *attrs) Set(key, val string) {
	(*a)[key] = val
}

// Attributes returns all attributes in a slice encoded
// as per gonum.graph.encoding requirements
func (a attrs) Attributes() []encoding.Attribute {
	return DOTAttrs(&a)
}
