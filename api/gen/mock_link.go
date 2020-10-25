package gen

import "github.com/milosgajdos/kraph/api"

// NewMockLink returns a new mock API Link
func NewMockLink(from, to, rel string) api.Link {
	return NewLink(NewUID(from), NewUID(to), NewRelation(rel))
}
