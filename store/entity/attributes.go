package entity

import (
	"gonum.org/v1/gonum/graph/encoding"
)

// Attributes are graph attributes
type Attributes map[string]string

// Get reads an attribute value for the given key and returns it.
// It returns an empty string if the attribute was not found.
func (a Attributes) Get(key string) string {
	return a[key]
}

// Set sets an attribute to the given value
func (a *Attributes) Set(key, val string) {
	(*a)[key] = val
}

// Attributes returns all attributes.
func (a Attributes) Attributes() []encoding.Attribute {
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
func (a Attributes) DOTAttributes() []encoding.Attribute {
	return a.Attributes()
}
