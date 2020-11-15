package owner

import "github.com/milosgajdos/kraph/pkg/api/gen"

// Source is the source of Github API.
type Source struct {
	*gen.Source
}

// NewSource returns API source.
func NewSource(s string) *Source {
	return &Source{
		Source: gen.NewSource(s),
	}
}
