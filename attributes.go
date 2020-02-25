package kraph

import (
	"errors"

	"gonum.org/v1/gonum/graph/encoding"
)

var (
	// ErrAttrKeyInvalid is returned when an invalid attribute key is given
	ErrAttrKeyInvalid = errors.New("invalid attribute key")
)

// Attrs provides graph attributes
type Attrs []encoding.Attribute

// Get gets an attribute value for a given key and returns it
// It returns empty string if the attribute was not found
func (a Attrs) Get(attr string) string {
	for _, attrKV := range a {
		if attrKV.Key == attr {
			return attrKV.Value
		}
	}

	return ""
}

// SetAttribute sets attribute to a given attribute value
// If the atttibute is not found it appends it to the existing attributes
func (a *Attrs) SetAttribute(attr encoding.Attribute) error {
	for i, attrKV := range *a {
		if attrKV.Key == attr.Key {
			(*a)[i].Value = attr.Value
			return nil
		}
	}

	*a = append(*a, attr)

	return nil
}

// Attrs returns all attributes
func (a Attrs) Attributes() []encoding.Attribute {
	return []encoding.Attribute(a)
}

// DOTAttrs returns GraphViz DOT attributes
func (a Attrs) DOTAttrs() []encoding.Attribute {
	return []encoding.Attribute(a)
}
