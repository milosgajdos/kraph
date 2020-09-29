package attrs

import (
	"gonum.org/v1/gonum/graph/encoding"
)

// Attrs are graph attributes
type Attrs map[string]string

// NewAttributes creates new attributes and returns it
func New() *Attrs {
	attrs := make(Attrs)

	return &attrs
}

// Keys returns all attribute keys
func (a Attrs) Keys() []string {
	keys := make([]string, len(a))

	i := 0
	for key := range a {
		keys[i] = key
	}

	return keys
}

// Get reads an attribute value for the given key and returns it.
// It returns an empty string if the attribute was not found.
func (a Attrs) Get(key string) string {
	return a[key]
}

// Set sets an attribute to the given value
func (a *Attrs) Set(key, val string) {
	(*a)[key] = val
}

// Attributes returns all attributes in a slice encoded
// as per gonum.graph.encoding requirements
func (a Attrs) Attributes() []encoding.Attribute {
	attrs := make([]encoding.Attribute, len(a))

	i := 0
	for k, v := range a {
		attrs[i] = encoding.Attribute{
			Key:   k,
			Value: v,
		}
		i++
	}

	return attrs
}
