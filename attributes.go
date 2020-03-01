package kraph

import (
	"errors"

	"gonum.org/v1/gonum/graph/encoding"
)

var (
	// ErrAttrKeyInvalid is returned when an invalid attribute key is given
	ErrAttrKeyInvalid = errors.New("invalid attribute key")
)

// Attrs are graph attributes
// NOTE: we might want to make this unexported
type Attrs map[string]string

// Get gets an attribute value for a given key and returns it
// It returns empty string if the attribute was not found
func (a Attrs) Get(key string) string {
	return a[key]
}

// SetAttribute sets attribute to a given attribute value
func (a *Attrs) SetAttribute(key, val string) error {
	(*a)[key] = val

	return nil
}

// Attrs returns all attributes
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

// DOTAttrs returns GraphViz DOT attributes
func (a Attrs) DOTAttrs() []encoding.Attribute {
	return a.Attributes()
}
