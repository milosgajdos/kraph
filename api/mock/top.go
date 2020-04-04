package mock

import (
	"strings"

	"github.com/milosgajdos/kraph/api"
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

	for resName, meta := range Resources {
		groups := ResourceData[resName]["groups"]
		versions := ResourceData[resName]["versions"]
		for _, group := range groups {
			for _, version := range versions {
				gv := strings.Join([]string{group, version}, "/")
				if gvObject, ok := ObjectData[gv]; ok {
					ns := meta["ns"]
					if len(ns) == 0 {
						ns = api.NamespaceNan
					}

					nsKind := strings.Join([]string{ns, meta["kind"]}, "/")
					if names, ok := gvObject[nsKind]; ok {
						for _, name := range names {
							uid := strings.Join([]string{ns, meta["kind"], name}, "/")
							object := generic.NewObject(name, meta["kind"], ns, generic.NewUID(uid))
							top.Add(object)
						}
					}
				}
			}
		}
	}

	return top
}
