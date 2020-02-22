package kraph

import (
	"errors"

	"gonum.org/v1/gonum/graph/encoding"
)

var (
	// ErrAttrNotFound is returned when an attribute could not be found
	ErrAttrNotFound = errors.New("attribute not found")
)

// Attributes provides graph attributes
type Attributes []encoding.Attribute

// Get gets am attribute value for a given key and returns it
// It returns empty string if the attribute was not found
func (a Attributes) Get(attr string) string {
	for _, attrKV := range a {
		if attrKV.Key == attr {
			return attrKV.Value
		}
	}

	return ""
}

// SetAttribute sets attribute to a given value
// If the atttibute is not found it appends it to the existing attributes
func (a *Attributes) SetAttribute(attr encoding.Attribute) error {
	for i, attrKV := range *a {
		if attrKV.Key == attr.Key {
			(*a)[i].Value = attr.Value
			return nil
		}
	}

	*a = append(*a, attr)

	return nil
}

// Attributes returns all attributes
func (a Attributes) Attributes() []encoding.Attribute {
	return []encoding.Attribute(a)
}

// DOTAttributes returns GraphViz DOT attributes
func (a Attributes) DOTAttributes() []encoding.Attribute {
	return []encoding.Attribute(a)
}
