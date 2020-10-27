package query

import (
	"github.com/milosgajdos/kraph/pkg/attrs"
	"github.com/milosgajdos/kraph/pkg/metadata"
	"github.com/milosgajdos/kraph/pkg/uuid"
)

type MatchVal int

const (
	MatchAny MatchVal = iota
)

type MatchFunc func(interface{}) bool

type matcher struct {
	val   interface{}
	funcs []MatchFunc
}

func newMatcher(val interface{}, funcs ...MatchFunc) *matcher {
	return &matcher{
		val:   val,
		funcs: funcs,
	}
}

func (m matcher) Value() interface{} {
	return m.val
}

func (m matcher) Match(val interface{}) bool {
	match := true
	for _, fn := range m.funcs {
		match = match && fn(val)
	}

	return match
}

type match struct {
	q *Query
}

func (m match) matchVal(prop string, val interface{}) bool {
	matcher, ok := m.q.matchers[prop]
	if !ok {
		return true
	}

	if val, ok := matcher.val.(MatchVal); ok && val == MatchAny {
		return true
	}

	return matcher.Match(val)
}

func (m *match) UID() *matcher {
	return m.q.matchers["uid"]
}

func (m *match) UIDVal(u uuid.UID) bool {
	return m.matchVal("uid", u)
}

func (m *match) Namespace() *matcher {
	return m.q.matchers["ns"]
}

func (m *match) NamespaceVal(ns string) bool {
	return m.matchVal("ns", ns)
}

func (m *match) Kind() *matcher {
	return m.q.matchers["kind"]
}

func (m *match) KindVal(k string) bool {
	return m.matchVal("kind", k)
}

func (m *match) Name() *matcher {
	return m.q.matchers["name"]
}

func (m *match) NameVal(n string) bool {
	return m.matchVal("name", n)
}

func (m *match) Version() *matcher {
	return m.q.matchers["version"]
}

func (m *match) VersionVal(v string) bool {
	return m.matchVal("version", v)
}

func (m *match) Group() *matcher {
	return m.q.matchers["group"]
}

func (m *match) GroupVal(g string) bool {
	return m.matchVal("group", g)
}

func (m *match) Entity() *matcher {
	return m.q.matchers["entity"]
}

func (m *match) EntityVal(e Entity) bool {
	return m.matchVal("entity", e)
}

func (m *match) Weight() *matcher {
	return m.q.matchers["weight"]
}

func (m *match) WeightVal(w float64) bool {
	return m.matchVal("weight", w)
}

func (m *match) Attrs() *matcher {
	return m.q.matchers["attrs"]
}

func (m *match) AttrsVal(a attrs.Attrs) bool {
	return m.matchVal("attrs", a)
}

func (m *match) Metadata() *matcher {
	return m.q.matchers["metadata"]
}

func (m *match) MetadataVal(meta metadata.Metadata) bool {
	return m.matchVal("metadata", meta)
}
