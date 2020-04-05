package mock

import (
	"github.com/milosgajdos/kraph/api/generic"
)

var (
	ObjectData = map[string]map[string][]string{
		"fooGroup/v1": {
			"fooNs/fooKind": []string{"foo1", "foo2", "foo3"},
		},
		"fooGroup/v2": {
			"fooNs/fooKind": []string{"foo4", "foo5"},
		},
		"barGroup/v2": {
			"nan/barKind": []string{"bar5"},
		},
		"rndGroup/v2": {
			"rndNs/rndKind": []string{"rnd1", "rnd2", "rnd3"},
		},
		"rndGroup/v6": {
			"rndNs/rndKind": []string{"rnd6"},
		},
	}
	ObjectLinks = map[string]map[string]string{
		"fooNs/fooKind/foo1": {
			"fooNs/fooKind/foo4": "foo-foo",
			"fooNs/fooKind/foo5": "foo-foo",
			"nan/barKind/bar5":   "foo-bar",
		},
		"nan/barKind/bar5": {
			"rndNs/rndKind/rnd2": "bar-rnd",
		},
		"rndNs/rndKind/rnd2": {
			"rndNs/rndKind/rnd6": "rnd-rnd",
			"fooNs/fooKind/foo1": "rnd-foo",
		},
	}
)

// Top provides mock Topology
type Top struct {
	*generic.Top
}

// NewTop creates new mock Top(ology)
func NewTop() *Top {
	top := &Top{
		Top: generic.NewTop(),
	}

	return top
}
