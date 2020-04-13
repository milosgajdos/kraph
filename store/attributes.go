package store

import (
	"gonum.org/v1/gonum/graph/encoding"
)

// attrs are graph attributes
type attrs map[string]string

// NewAttributes creates new attributes and returns it
func NewAttributes() Attrs {
	attrs := make(attrs)

	return &attrs
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

// Attributes returns all attributes.
func (a attrs) Attributes() []encoding.Attribute {
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

// DOTAttributes returns GraphViz DOT attributes
func (a attrs) DOTAttributes() []encoding.Attribute {
	return a.Attributes()
}

// Keys returns all the metadata keys
func (a attrs) Keys() []string {
	keys := make([]string, len(a))

	i := 0
	for key, _ := range a {
		keys[i] = key
	}

	return keys
}
